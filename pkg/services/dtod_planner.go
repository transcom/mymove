package services

import "github.com/transcom/mymove/pkg/models"

// PaymentRequestListFetcher is the exported interface for fetching a list of payment requests
//go:generate mockery -name PaymentRequestListFetcher
type DTODPlannerMileageFetcher interface {
	FetchPaymentRequestList(officeUserID uuid.UUID, params *FetchPaymentRequestListParams) (*models.PaymentRequests, int, error)
	FetchPaymentRequestListByMove(officeUserID uuid.UUID, locator string) (*models.PaymentRequests, error)
}

