package service

import (
	"context"
	"errors"
	"fmt"
	"media-service/helper"
	"media-service/internal/department/dto/request"
	"media-service/internal/department/dto/response"
	"media-service/internal/department/mapper"
	"media-service/internal/department/model"
	"media-service/internal/department/repository"
	"media-service/internal/gateway"
	"media-service/internal/gateway/dto"
	"media-service/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DepartmentService interface {
	UploadDepartment(ctx context.Context, req request.UploadDepartmentRequest) (*response.DepartmentResponseDTO, error)
	UpdateDepartment(ctx context.Context, req request.UpdateDepartmentRequest) (*response.DepartmentResponseDTO, error)
	UploadDepartmentMenu(ctx context.Context, req request.UploadSectionMenuDepartmentRequest) error
	GetDepartments4Web(ctx context.Context) ([]*response.GetDepartments4Web, error)
	GetDeparmentDetail4Web(ctx context.Context, departmentID string) (*response.GetDepartment4Web, error)
	AssignLeader(ctx context.Context, req request.AssignLeaderRequest) (*model.Leader, error)
	AssignStaff(ctx context.Context, req request.AssignStaffRequest) (*model.Staff, error)
	RemoveStaffByIndex(ctx context.Context, departmentID string, index int) error
	GetDepartments4App(ctx context.Context) ([]*response.GetDepartment4App, error)
	GetDepartments4Gateway(ctx context.Context) ([]*response.GetDepartment4Gateway, error)
	UploadDepartmentMenuOrganization(ctx context.Context, req request.UploadDepartmentMenuOrganizationRequest) error
	RemoveLeader(ctx context.Context, req request.RemoveLeaderRequest) error
	GetDepartmentsByOrganization4Gateway(ctx context.Context, organizationID string) ([]*response.GetDepartment4Gateway, error)
}

type departmentService struct {
	userGateway            gateway.UserGateway
	menuGateway            gateway.MenuGateway
	messageLanguageGateway gateway.MessageLanguageGateway
	classRoomGateway       gateway.ClassroomGateway
	repo                   repository.DepartmentRepository
	regionRepo             repository.RegionRepository
}

func NewDepartmentService(
	repo repository.DepartmentRepository,
	userGateway gateway.UserGateway,
	messageLanguageGateway gateway.MessageLanguageGateway,
	menuGateway gateway.MenuGateway,
	classroomGW gateway.ClassroomGateway,
	regionRepo repository.RegionRepository) DepartmentService {
	return &departmentService{
		userGateway:            userGateway,
		menuGateway:            menuGateway,
		messageLanguageGateway: messageLanguageGateway,
		classRoomGateway:       classroomGW,
		repo:                   repo,
		regionRepo:             regionRepo,
	}
}

// UploadDepartment service layer
func (s *departmentService) UploadDepartment(ctx context.Context, req request.UploadDepartmentRequest) (*response.DepartmentResponseDTO, error) {

	// get organization admin from user context
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	// check is super admin & check org admin
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied: super admin cannot perform this action")
	}

	organizationAdminID := currentUser.OrganizationAdmin.ID

	department := &model.Department{
		ID:                 primitive.NewObjectID(),
		LocationID:         req.LocationID,
		OrganizationID:     organizationAdminID,
		RegionID:           req.RegionID,
		Name:               req.Name,
		Description:        req.Description,
		Note:               req.Note,
		Icon:               req.Icon,
		Leader:             model.Leader{},
		Staffs:             []model.Staff{},
		IsPublishedMessage: req.IsPublishedMessage,
		Url:                req.Url,
	}

	// Gọi repository để insert
	err = s.repo.UploadDepartment(ctx, department)
	if err != nil {
		return nil, err
	}

	// goi messs lang gw upload message
	err = s.uploadMessages(ctx, helper.BuildDepartmentMessagesUpload(department.ID.Hex(), req))

	if err != nil {
		return nil, fmt.Errorf("upload department messages failed: %w", err)
	}
	return mapper.MapDepartmentToResponse(department, nil, ""), nil
}

func (s *departmentService) UpdateDepartment(ctx context.Context, req request.UpdateDepartmentRequest) (*response.DepartmentResponseDTO, error) {
	// get organization admin from user context
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	// check is super admin & check org admin
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied: super admin cannot perform this action")
	}

	// Convert string ID sang ObjectID
	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, errors.New("invalid department ID")
	}

	// Check tồn tại
	_, err = s.repo.GetByIDAndOrgID(ctx, objID, currentUser.OrganizationAdmin.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}

	// Overwrite dữ liệu theo request
	department := &model.Department{
		ID:                 objID,
		LocationID:         req.LocationID,
		RegionID:           req.RegionID,
		Name:               req.Name,
		Description:        req.Description,
		Note:               req.Note,
		Icon:               req.Icon,
		UpdatedAt:          time.Now(),
		IsPublishedMessage: req.IsPublishedMessage,
		Url:                req.Url,
	}

	// Gọi repository để update
	err = s.repo.UpdateDepartment(ctx, department)
	if err != nil {
		return nil, err
	}

	// goi messs lang gw upload message
	err = s.uploadMessages(ctx, helper.BuildDepartmentMessagesUpdate(department.ID.Hex(), req))

	if err != nil {
		return nil, fmt.Errorf("upload department messages failed: %w", err)
	}

	return mapper.MapDepartmentToResponse(department, nil, ""), nil
}

func (s *departmentService) UploadDepartmentMenu(ctx context.Context, req request.UploadSectionMenuDepartmentRequest) error {
	err := s.menuGateway.UploadDepartmentMenu(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (s *departmentService) GetDepartments4Web(ctx context.Context) ([]*response.GetDepartments4Web, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	// kiểm tra quyền
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, nil
	}

	// 1. Lấy tất cả regions theo OrgID
	regions, err := s.regionRepo.GetAllByOrgID(ctx, currentUser.OrganizationAdmin.ID)
	if err != nil {
		return nil, err
	}

	var result []*response.GetDepartments4Web

	// 2. Lặp qua từng region
	for _, region := range regions {
		// 2.1 Lấy danh sách departments theo regionID
		departments, err := s.repo.GetDepartmentsByRegionID(ctx, region.ID.Hex())
		if err != nil {
			return nil, err
		}

		iconUrls := make(map[string]string)
		leaders := make(map[string]response.LeaderResponseDTO)
		staffsMap := make(map[string][]response.StaffResponseDTO)
		homeMenusMap := make(map[string][]dto.MenuResponse)
		organizationMenusMap := make(map[string][]dto.MenuResponse)
		messageLanguagesMap := make(map[string][]dto.MessageLanguageResponse)
		locationNamesMap := make(map[string]string)

		for _, dept := range departments {
			// get home menus
			homeMenus, _ := s.menuGateway.GetDepartmentMenu(ctx, dept.ID.Hex())
			homeMenusMap[dept.ID.Hex()] = homeMenus

			// get organization menus
			organizationMenus, _ := s.menuGateway.GetDepartmentMenuOrganization(ctx, dept.ID.Hex(), dept.OrganizationID)
			organizationMenusMap[dept.ID.Hex()] = organizationMenus

			// get icon url
			iconUrl, _ := s.menuGateway.GetImageURL(ctx, dept.Icon, string(constants.ImageModePrivate))
			iconUrls[dept.ID.Hex()] = iconUrl

			// get message languages
			msgLangs, _ := s.messageLanguageGateway.GetMessageLanguages(ctx, dept.ID.Hex())
			messageLanguagesMap[dept.ID.Hex()] = msgLangs

			// get location name
			if dept.LocationID != "" {
				classroom, err := s.classRoomGateway.GetClassroomByID(ctx, dept.LocationID)
				if err == nil && classroom != nil {
					locationNamesMap[dept.LocationID] = classroom.Name
				}
			}

			// map leader
			var leaderDTO response.LeaderResponseDTO
			switch dept.Leader.OwnerRole {
			case string(constants.OwnerRoleTeacher):
				teacherInfo, err := s.userGateway.GetTeacherInfo(ctx, dept.Leader.OwnerID)
				if err == nil && teacherInfo != nil {
					leaderDTO = response.LeaderResponseDTO{
						OwnerID:   dept.Leader.OwnerID,
						OwnerRole: dept.Leader.OwnerRole,
						Name:      teacherInfo.Name,
						AvatarUrl: teacherInfo.Avatar.ImageUrl,
					}
				}
			case string(constants.OwnerRoleStaff):
				staffInfo, err := s.userGateway.GetStaffInfo(ctx, dept.Leader.OwnerID)
				if err == nil && staffInfo != nil {
					leaderDTO = response.LeaderResponseDTO{
						OwnerID:   dept.Leader.OwnerID,
						OwnerRole: dept.Leader.OwnerRole,
						Name:      staffInfo.Name,
						AvatarUrl: staffInfo.Avatar.ImageUrl,
					}
				}
			default:
				leaderDTO = response.LeaderResponseDTO{
					OwnerID:   dept.Leader.OwnerID,
					OwnerRole: dept.Leader.OwnerRole,
				}
			}
			leaders[dept.ID.Hex()] = leaderDTO

			// map staffs
			var staffs []response.StaffResponseDTO
			for _, st := range dept.Staffs {
				var staffDTO response.StaffResponseDTO
				switch st.OwnerRole {
				case string(constants.OwnerRoleTeacher):
					teacherInfo, err := s.userGateway.GetTeacherInfo(ctx, st.OwnerID)
					if err == nil && teacherInfo != nil {
						staffDTO = response.StaffResponseDTO{
							OwnerID:   st.OwnerID,
							OwnerRole: st.OwnerRole,
							Index:     st.Index,
							Name:      teacherInfo.Name,
							AvatarUrl: teacherInfo.Avatar.ImageUrl,
						}
					}
				case string(constants.OwnerRoleStaff):
					staffInfo, err := s.userGateway.GetStaffInfo(ctx, st.OwnerID)
					if err == nil && staffInfo != nil {
						staffDTO = response.StaffResponseDTO{
							OwnerID:   st.OwnerID,
							OwnerRole: st.OwnerRole,
							Index:     st.Index,
							Name:      staffInfo.Name,
							AvatarUrl: staffInfo.Avatar.ImageUrl,
						}
					}
				default:
					staffDTO = response.StaffResponseDTO{
						OwnerID:   st.OwnerID,
						OwnerRole: st.OwnerRole,
						Index:     st.Index,
					}
				}
				staffs = append(staffs, staffDTO)
			}
			staffsMap[dept.ID.Hex()] = staffs
		}

		// 3. Gom nhóm cho từng region
		group := &response.GetDepartments4Web{
			RegionID:    region.ID.Hex(),
			RegionName:  region.Name,
			Departments: mapper.MapDepartmentsToResponses4Web(departments, homeMenusMap, iconUrls, leaders, staffsMap, organizationMenusMap, messageLanguagesMap, locationNamesMap),
		}
		result = append(result, group)
	}

	return result, nil
}

func (s *departmentService) GetDeparmentDetail4Web(ctx context.Context, departmentID string) (*response.GetDepartment4Web, error) {
	// get organization admin from user context
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	// check is super admin & check org admin
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, nil
	}

	// convert departmentID string -> ObjectID
	objID, err := primitive.ObjectIDFromHex(departmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid department id: %w", err)
	}

	department, err := s.repo.GetByIDAndOrgID(ctx, objID, currentUser.OrganizationAdmin.ID)
	if err != nil {
		return nil, err
	}

	// gọi gateway để lấy menu của department
	homeMenus, _ := s.menuGateway.GetDepartmentMenu(ctx, departmentID)
	orgMenus, _ := s.menuGateway.GetDepartmentMenuOrganization(ctx, departmentID, department.OrganizationID)

	// get icon url
	iconUrl, _ := s.menuGateway.GetImageURL(ctx, department.Icon, string(constants.ImageModePrivate))

	// map sang DTO cho web
	return mapper.MapDepartmentToResponse4Web(department, homeMenus, iconUrl, response.LeaderResponseDTO{}, nil, orgMenus, nil, ""), nil
}

func (s *departmentService) AssignLeader(ctx context.Context, req request.AssignLeaderRequest) (*model.Leader, error) {

	leader := model.Leader{
		OwnerID:   req.OwnerID,
		OwnerRole: req.OwnerRole,
	}

	// convert departmentID string -> ObjectID
	objID, err := primitive.ObjectIDFromHex(req.DepartmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid department id: %w", err)
	}

	return s.repo.AssignLeader(ctx, objID, leader)
}

func (s *departmentService) AssignStaff(ctx context.Context, req request.AssignStaffRequest) (*model.Staff, error) {

	staff := model.Staff{
		OwnerID:   req.OwnerID,
		OwnerRole: req.OwnerRole,
		Index:     req.Index,
	}
	// convert departmentID string -> ObjectID
	objID, err := primitive.ObjectIDFromHex(req.DepartmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid department id: %w", err)
	}

	return s.repo.AssignStaff(ctx, objID, staff)
}

func (s *departmentService) RemoveStaffByIndex(ctx context.Context, departmentID string, index int) error {
	// convert departmentID string -> ObjectID
	objID, err := primitive.ObjectIDFromHex(departmentID)
	if err != nil {
		return fmt.Errorf("invalid department id: %w", err)
	}

	return s.repo.RemoveStaffByIndex(ctx, objID, index)
}

func (s *departmentService) GetDepartments4App(ctx context.Context) ([]*response.GetDepartment4App, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	var departments []*model.Department

	// Lấy danh sách teachers
	teachers, _ := s.userGateway.GetTeachersByUser(ctx, currentUser.ID)
	for _, teacher := range teachers {
		if teacher != nil && teacher.ID != "" {
			teacherDepartments, err := s.repo.GetByOwnerID(ctx, teacher.ID)
			if err != nil {
				return nil, fmt.Errorf("get teacher departments failed: %w", err)
			}
			departments = append(departments, teacherDepartments...)
		}
	}

	// Lấy danh sách staffs
	staffs, _ := s.userGateway.GetStaffsByUser(ctx, currentUser.ID)
	for _, staff := range staffs {
		if staff != nil && staff.ID != "" {
			staffDepartments, err := s.repo.GetByOwnerID(ctx, staff.ID)
			if err != nil {
				return nil, fmt.Errorf("get staff departments failed: %w", err)
			}
			departments = append(departments, staffDepartments...)
		}
	}

	// Nếu không có teacher hoặc staff thì departments sẽ rỗng
	if len(departments) == 0 {
		return nil, nil
	}

	// --- Xóa trùng lặp theo Department.ID ---
	uniqueDepartments := make([]*model.Department, 0, len(departments))
	seen := make(map[string]bool)
	for _, d := range departments {
		if d != nil {
			idStr := d.ID.Hex()
			if !seen[idStr] {
				seen[idStr] = true
				uniqueDepartments = append(uniqueDepartments, d)
			}
		}
	}

	// ---- Gọi message language ----
	responses := make([]*response.GetDepartment4App, 0, len(uniqueDepartments))
	for _, dept := range uniqueDepartments {
		msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, dept.ID.Hex())
		responses = append(responses, mapper.MapDepartmentToResponse4App(dept, msg))
	}

	return responses, nil
}

func (s *departmentService) GetDepartments4Gateway(ctx context.Context) ([]*response.GetDepartment4Gateway, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	var departments []*model.Department

	// Lấy danh sách teachers
	teachers, _ := s.userGateway.GetTeachersByUser(ctx, currentUser.ID)
	for _, teacher := range teachers {
		if teacher != nil && teacher.ID != "" {
			teacherDepartments, err := s.repo.GetByOwnerID(ctx, teacher.ID)
			if err != nil {
				return nil, fmt.Errorf("get teacher departments failed: %w", err)
			}
			departments = append(departments, teacherDepartments...)
		}
	}

	// Lấy danh sách staffs
	staffs, _ := s.userGateway.GetStaffsByUser(ctx, currentUser.ID)
	for _, staff := range staffs {
		if staff != nil && staff.ID != "" {
			staffDepartments, err := s.repo.GetByOwnerID(ctx, staff.ID)
			if err != nil {
				return nil, fmt.Errorf("get staff departments failed: %w", err)
			}
			departments = append(departments, staffDepartments...)
		}
	}

	// Nếu không có teacher hoặc staff thì departments sẽ rỗng
	if len(departments) == 0 {
		return nil, nil
	}

	// --- Xóa trùng lặp theo Department.ID ---
	uniqueDepartments := make([]*model.Department, 0, len(departments))
	seen := make(map[string]bool)
	for _, d := range departments {
		if d != nil {
			idStr := d.ID.Hex()
			if !seen[idStr] {
				seen[idStr] = true
				uniqueDepartments = append(uniqueDepartments, d)
			}
		}
	}

	// ---- Gọi message language ----
	responses := make([]*response.GetDepartment4Gateway, 0, len(uniqueDepartments))
	for _, dept := range uniqueDepartments {
		msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, dept.ID.Hex())
		responses = append(responses, mapper.MapDepartmentToResponse4Gateway(dept, msg))
	}

	return responses, nil
}
func (s *departmentService) GetDepartmentsByOrganization4Gateway(ctx context.Context, organizationID string) ([]*response.GetDepartment4Gateway, error) {
	// lấy departments
	departments, err := s.repo.GetDepartmentsByOrgID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	// map sang response
	responses := make([]*response.GetDepartment4Gateway, len(departments))
	for i, dept := range departments {
		// lấy message language theo department
		msg, err := s.messageLanguageGateway.GetMessageLanguage(ctx, dept.ID.Hex())
		if err != nil {
			msg = dto.MessageLanguageResponse{}
		}
		responses[i] = mapper.MapDepartmentToResponse4Gateway(dept, msg)
	}

	return responses, nil
}

func (s *departmentService) UploadDepartmentMenuOrganization(ctx context.Context, req request.UploadDepartmentMenuOrganizationRequest) error {
	err := s.menuGateway.UploadDepartmentMenuOrganization(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (s *departmentService) RemoveLeader(ctx context.Context, req request.RemoveLeaderRequest) error {
	return s.repo.RemoveLeader(ctx, req.DepartmentID)
}

func (s *departmentService) uploadMessages(ctx context.Context, req dto.UploadMessageLanguagesRequest) error {
	err := s.messageLanguageGateway.UploadMessages(ctx, req)

	if err != nil {
		return err
	}

	return nil
}
