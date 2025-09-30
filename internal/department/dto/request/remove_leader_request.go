package request

type RemoveLeaderRequest struct {
	DepartmentID string `json:"department_id" binding:"required"`
}
