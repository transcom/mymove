package services

import "github.com/transcom/mymove/pkg/models"

type OfficeUserFetcher interface {
	FetchOfficeUser(field string, value string) (models.OfficeUser, error)
}
