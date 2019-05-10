package services

import "github.com/transcom/mymove/pkg/models"

type OfficeUserFetcher interface {
	FetchOfficeUser(field string, value interface{}) (models.OfficeUser, error)
}

type OfficeUserListFetcher interface {
	FetchOfficeUserList(filters map[string]interface{}) (models.OfficeUsers, error)
}
