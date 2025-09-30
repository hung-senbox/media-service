package helper

import (
	"media-service/internal/department/dto/request"
	"media-service/internal/gateway/dto"
	"media-service/pkg/constants"
	"strconv"
	"strings"
)

func BuildDepartmentMessagesUpload(departmentID string, req request.UploadDepartmentRequest) dto.UploadMessageLanguagesRequest {
	return dto.UploadMessageLanguagesRequest{
		MessageLanguages: []dto.UploadMessageRequest{
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentMessageKey),
				Value:      req.Message,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentNoteKey),
				Value:      req.Note,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentNameKey),
				Value:      req.Name,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentDescKey),
				Value:      req.Description,
				LanguageID: req.LanguageID,
			},
		},
	}
}

func BuildDepartmentMessagesUpdate(departmentID string, req request.UpdateDepartmentRequest) dto.UploadMessageLanguagesRequest {
	return dto.UploadMessageLanguagesRequest{
		MessageLanguages: []dto.UploadMessageRequest{
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentMessageKey),
				Value:      req.Message,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentNoteKey),
				Value:      req.Note,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentNameKey),
				Value:      req.Name,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     departmentID,
				Type:       "department",
				Key:        string(constants.DepartmentDescKey),
				Value:      req.Description,
				LanguageID: req.LanguageID,
			},
		},
	}
}

func ParseAppLanguage(header string, defaultVal uint) uint {
	header = strings.TrimSpace(strings.Trim(header, "\""))
	if val, err := strconv.Atoi(header); err == nil {
		return uint(val)
	}
	return defaultVal
}
