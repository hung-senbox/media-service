package cache

func UserCacheKey(userID string) string {
	return "user:" + userID
}

func StudentCacheKey(studentID string) string {
	return "student:" + studentID
}

func TeacherCacheKey(teacherID string) string {
	return "teacher:" + teacherID
}
