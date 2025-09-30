package repository

import (
	"context"
	"errors"
	"fmt"
	"media-service/internal/department/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DepartmentRepository interface {
	UploadDepartment(ctx context.Context, department *model.Department) error
	UpdateDepartment(ctx context.Context, department *model.Department) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Department, error)
	GetByIDAndOrgID(ctx context.Context, id primitive.ObjectID, orgID string) (*model.Department, error)
	GetDepartmentsByOrgID(ctx context.Context, orgID string) ([]*model.Department, error)
	AssignLeader(ctx context.Context, deptID primitive.ObjectID, leader model.Leader) (*model.Leader, error)
	AssignStaff(ctx context.Context, deptID primitive.ObjectID, staff model.Staff) (*model.Staff, error)
	RemoveStaffByIndex(ctx context.Context, deptID primitive.ObjectID, index int) error
	GetByOwnerID(ctx context.Context, ownerID string) ([]*model.Department, error)
	GetDepartmentsByRegionID(ctx context.Context, regionID string) ([]*model.Department, error)
	RemoveLeader(ctx context.Context, deptID string) error
}

type departmentRepository struct {
	collection *mongo.Collection
}

func NewDepartmentRepository(collection *mongo.Collection) DepartmentRepository {
	return &departmentRepository{collection}
}

func (r *departmentRepository) UploadDepartment(ctx context.Context, department *model.Department) error {
	// Nếu chưa có _id thì tự sinh
	if department.ID.IsZero() {
		department.ID = primitive.NewObjectID()
	}

	department.CreatedAt = time.Now()
	department.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, department)
	return err
}

func (r *departmentRepository) UpdateDepartment(ctx context.Context, dept *model.Department) error {
	filter := bson.M{"_id": dept.ID}
	update := bson.M{
		"$set": bson.M{
			"location_id":          dept.LocationID,
			"region_id":            dept.RegionID,
			"icon":                 dept.Icon,
			"updated_at":           time.Now().Unix(),
			"url":                  dept.Url,
			"is_published_message": dept.IsPublishedMessage,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *departmentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Department, error) {
	var dept model.Department
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&dept)
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) GetByIDAndOrgID(ctx context.Context, id primitive.ObjectID, orgID string) (*model.Department, error) {
	var dept model.Department
	err := r.collection.FindOne(ctx, bson.M{
		"_id":             id,
		"organization_id": orgID,
	}).Decode(&dept)
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) GetDepartmentsByOrgID(ctx context.Context, orgID string) ([]*model.Department, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"organization_id": orgID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var departments []*model.Department
	if err = cursor.All(ctx, &departments); err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *departmentRepository) AssignLeader(ctx context.Context, deptID primitive.ObjectID, leader model.Leader) (*model.Leader, error) {

	filter := bson.M{"_id": deptID}
	update := bson.M{
		"$set": bson.M{
			"leader":     leader,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("department not found")
	}

	return &leader, nil
}

func (r *departmentRepository) AssignStaff(ctx context.Context, deptID primitive.ObjectID, staff model.Staff) (*model.Staff, error) {
	// 1. Check staff tồn tại trong department
	filter := bson.M{
		"_id": deptID,
		"staffs": bson.M{
			"$elemMatch": bson.M{
				"owner_id":   staff.OwnerID,
				"owner_role": staff.OwnerRole,
			},
		},
	}

	var dept model.Department
	err := r.collection.FindOne(ctx, filter).Decode(&dept)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("failed to check staff existence: %w", err)
	}

	// kiem tra da gan trong leader chua
	if dept.Leader.OwnerID == staff.OwnerID && dept.Leader.OwnerRole == staff.OwnerRole {
		return nil, fmt.Errorf("the staff already exists as leader")
	}

	if err == nil {
		for _, st := range dept.Staffs {
			if st.OwnerID == staff.OwnerID && st.OwnerRole == staff.OwnerRole {
				if st.Index != staff.Index {
					return nil, fmt.Errorf("the staff was assigned")
				}
			}
		}
	}

	// 2. Update staff nếu index đã có
	filter = bson.M{
		"_id":          deptID,
		"staffs.index": staff.Index,
	}
	update := bson.M{
		"$set": bson.M{
			"staffs.$.owner_id":   staff.OwnerID,
			"staffs.$.owner_role": staff.OwnerRole,
			"staffs.$.index":      staff.Index,
			"updated_at":          primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// 3. Nếu không match → push staff mới
	if res.MatchedCount == 0 {
		filter = bson.M{"_id": deptID}
		update = bson.M{
			"$push": bson.M{
				"staffs": staff,
			},
			"$set": bson.M{
				"updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		}
		_, err = r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return &staff, nil
}

func (r *departmentRepository) RemoveStaffByIndex(ctx context.Context, deptID primitive.ObjectID, index int) error {
	filter := bson.M{"_id": deptID}

	update := bson.M{
		"$pull": bson.M{
			"staffs": bson.M{"index": index},
		},
		"$set": bson.M{
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("department not found")
	}

	return nil
}

func (r *departmentRepository) GetByOwnerID(ctx context.Context, ownerID string) ([]*model.Department, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"leader.owner_id": ownerID},
			{"staffs.owner_id": ownerID},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var departments []*model.Department
	if err := cursor.All(ctx, &departments); err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *departmentRepository) GetDepartmentsByRegionID(ctx context.Context, regionID string) ([]*model.Department, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"region_id": regionID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var departments []*model.Department
	if err = cursor.All(ctx, &departments); err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *departmentRepository) RemoveLeader(ctx context.Context, deptID string) error {
	objectID, err := primitive.ObjectIDFromHex(deptID)
	if err != nil {
		return fmt.Errorf("invalid department ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$unset": bson.M{
			"leader": 1,
		},
		"$set": bson.M{
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no department found with id %s", deptID)
	}

	return nil
}
