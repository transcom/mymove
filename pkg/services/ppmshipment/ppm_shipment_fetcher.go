package ppmshipment

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
)

func FindPPMShipment(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, eagerAssociations ...string) (*models.PPMShipment, error) {
	var shipment models.PPMShipment
	findShipmentQuery := appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope())

	if len(eagerAssociations) > 0 {
		findShipmentQuery.Eager(eagerAssociations...)
	}

	err := findShipmentQuery.Find(&shipment, ppmShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for shipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "")
		}
	}

	return &shipment, nil
}

func FindPPMShipmentWithDocument(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, documentID uuid.UUID) error {
	var weightTicket models.WeightTicket
	var proGear models.ProgearWeightTicket
	var movingExpense models.MovingExpense

	err := appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope()).
		Where("ppm_shipment_id = ? AND (empty_document_id = ? OR full_document_id = ?)", ppmShipmentID, documentID, documentID).
		First(&weightTicket)

	if err != nil {
		switch err {
		case sql.ErrNoRows: // not ready to return an error unless the document is also not part of pro gear or expenses
		default:
			return apperror.NewQueryError("PPMShipment", err, "")
		}
	} else {
		return nil
	}

	err = appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope()).
		Where("ppm_shipment_id = ? AND (empty_document_id = ? OR full_document_id = ? OR constructed_weight_document_id = ?)", ppmShipmentID, documentID, documentID, documentID).
		First(&proGear)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
		default:
			return apperror.NewQueryError("PPMShipment", err, "")
		}
	} else {
		return nil
	}

	err = appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope()).
		Where("ppm_shipment_id = ? AND document_id = ?", ppmShipmentID, documentID).
		First(&movingExpense)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(documentID, "document does not exist for the given shipment")
		default:
			return apperror.NewQueryError("PPMShipment", err, "")
		}
	}

	return nil
}
