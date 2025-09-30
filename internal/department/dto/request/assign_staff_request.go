package request

type AssignStaffRequest struct {
	DepartmentID string `json:"department_id" binding:"required"`
	OwnerRole    string `json:"owner_role" binding:"required"`
	OwnerID      string `json:"owner_id" binding:"required"`
	Index        int    `json:"index" binding:"gte=0"`
}
