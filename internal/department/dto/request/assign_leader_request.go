package request

type AssignLeaderRequest struct {
	DepartmentID string `json:"department_id" binding:"required"`
	OwnerRole    string `json:"owner_role" binding:"required"`
	OwnerID      string `json:"owner_id" binding:"required"`
}
