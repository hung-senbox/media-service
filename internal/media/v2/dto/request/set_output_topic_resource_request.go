package request

type SetOutputTopicResourceRequest struct {
	TopicResourceID string `json:"topic_resource_id" binding:"required"`
	TargetStudentID string `json:"target_student_id" binding:"required"`
	TargetTopicID   string `json:"target_topic_id" binding:"required"`
}
