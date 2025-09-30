package request

type RemoveStaffRequest struct {
	DepartmentID string `json:"department_id"`
	Index        int    `json:"index" binding:"gte=0"`
}
