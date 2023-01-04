package move

import (
	"errors"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type CloseoutOfficeUpdater struct {
	services.MoveFetcher
	services.TransportationOfficesFetcher
}

func NewCloseoutOfficeUpdater(moveFetcher services.MoveFetcher, transportationOfficeFetcher services.TransportationOfficesFetcher) services.MoveCloseoutOfficeUpdater {
	return &CloseoutOfficeUpdater{moveFetcher, transportationOfficeFetcher}
}

func (s CloseoutOfficeUpdater) UpdateCloseoutOffice(appCtx appcontext.AppContext, moveLocator string, closeoutOfficeID uuid.UUID, eTag string) (*models.Move, error) {
	move, err := s.MoveFetcher.FetchMove(appCtx, moveLocator, &services.MoveFetcherParams{IncludeHidden: false})
	if err != nil {
		return nil, err
	}

	if eTag != etag.GenerateEtag(move.UpdatedAt) {
		return nil, apperror.NewPreconditionFailedError(move.ID, errors.New("If-Match eTag provided did not match the move's updatedAt value"))
	}

	transportationOffice, err := s.GetTransportationOffice(appCtx, closeoutOfficeID, true)
	if err != nil {
		return nil, err
	}

	move.CloseoutOfficeID = &transportationOffice.ID
	move.CloseoutOffice = transportationOffice

	verrs, err := appCtx.DB().ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "")
	}

	return move, nil
}
