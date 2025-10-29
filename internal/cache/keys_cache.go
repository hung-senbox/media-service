package cache

import "media-service/pkg/constants"

func UserCacheKey(userID string) string {
	return constants.ProfileCachePrefix + "user:" + userID
}

func StudentCacheKey(studentID string) string {
	return constants.ProfileCachePrefix + "student:" + studentID
}

func TeacherCacheKey(teacherID string) string {
	return constants.ProfileCachePrefix + "teacher:" + teacherID
}

func StaffCacheKey(staffID string) string {
	return constants.ProfileCachePrefix + "staff:" + staffID
}

func ParentCacheKey(parentID string) string {
	return constants.ProfileCachePrefix + "parent:" + parentID
}

func ChildCacheKey(childID string) string {
	return constants.ProfileCachePrefix + "child:" + childID
}

func TeacherByUserAndOrgCacheKey(userID, orgID string) string {
	return constants.ProfileCachePrefix + "teacher-by:" + userID + ":" + orgID
}

func UserByTeacherCacheKey(teacherID string) string {
	return constants.ProfileCachePrefix + "user-by:" + teacherID
}
