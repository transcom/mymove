package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// AccessCodeValidator is the service object interface for ValidateAccessCode
//go:generate $GOPATH/src/github.com/transcom/mymove/bin/mockery -name AccessCodeValidator -output=$GOPATH/src/github.com/transcom/mymove/mocks
type AccessCodeValidator interface {
	ValidateAccessCode(code string, moveType models.SelectedMoveType) (*models.AccessCode, bool, error)
}

// AccessCodeClaimer is the service object interface for ValidateAccessCode
//go:generate $GOPATH/src/github.com/transcom/mymove/bin/mockery -name AccessCodeClaimer -output=$GOPATH/src/github.com/transcom/mymove/mocks
type AccessCodeClaimer interface {
	ClaimAccessCode(code string, serviceMemberID uuid.UUID) (*models.AccessCode, *validate.Errors, error)
}
