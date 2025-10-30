package cache

import "media-service/pkg/constants"

func UserCacheKey(userID string) string {
	return constants.MainCachePrefix + "user:" + userID
}

func StudentCacheKey(studentID string) string {
	return constants.MainCachePrefix + "student:" + studentID
}

func TeacherCacheKey(teacherID string) string {
	return constants.MainCachePrefix + "teacher:" + teacherID
}

func StaffCacheKey(staffID string) string {
	return constants.MainCachePrefix + "staff:" + staffID
}

func ParentCacheKey(parentID string) string {
	return constants.MainCachePrefix + "parent:" + parentID
}

func ChildCacheKey(childID string) string {
	return constants.MainCachePrefix + "child:" + childID
}

func TeacherByUserAndOrgCacheKey(userID, orgID string) string {
	return constants.MainCachePrefix + "teacher-by-user-org:" + userID + ":" + orgID
}

func StaffByUserAndOrgCacheKey(userID, orgID string) string {
	return constants.MainCachePrefix + "staff-by-user-org:" + userID + ":" + orgID
}

func UserByTeacherCacheKey(teacherID string) string {
	return constants.MainCachePrefix + "user-by-teacher:" + teacherID
}

func ParentByUserCacheKey(userID string) string {
	return constants.MainCachePrefix + "parent-by-user:" + userID
}
