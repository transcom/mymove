package services

import "github.com/transcom/mymove/pkg/models"

// OfficeUserFetcher is the exported interface for fetching a single office user
type OfficeUserFetcher interface {
	FetchOfficeUser(filters map[string]interface{}) (models.OfficeUser, error)
}

// OfficeUserListFetcher is the exported interface for fetching multiple office users
type OfficeUserListFetcher interface {
	FetchOfficeUserList(filters map[string]interface{}) (models.OfficeUsers, error)
}
