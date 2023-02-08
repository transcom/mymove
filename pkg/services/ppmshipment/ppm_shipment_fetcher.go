package ppmshipment

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
)

func FindPPMShipmentAndWeightTickets(appCtx appcontext.AppContext, id uuid.UUID) (*models.PPMShipment, error) {
	var ppmShipment models.PPMShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment",
			"WeightTickets",
		).
		Find(&ppmShipment, id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(id, "while looking for PPMShipmentAndWeightTickets")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipmentAndWeightTickets")
		}
	}
	ppmShipment.WeightTickets = ppmShipment.WeightTickets.FilterDeleted()

	return &ppmShipment, nil
}

func FindPPMShipmentByMTOID(appCtx appcontext.AppContext, mtoID uuid.UUID) (*models.PPMShipment, error) {
	var ppmShipment models.PPMShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment",
			"WeightTickets",
			"MovingExpenses",
			"ProgearExpenses",
			"W2Address",
		).
		Where("shipment_id = ?", mtoID).First(&ppmShipment)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoID, "while looking for PPMShipment by MTO ShipmentID")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
		}
	}

	err = loadPPMAssociations(appCtx, &ppmShipment)
	if err != nil {
		return nil, err
	}

	return &ppmShipment, nil
}

// FindPPMShipment returns a PPMShipment with associations by ID
func FindPPMShipment(appCtx appcontext.AppContext, id uuid.UUID) (*models.PPMShipment, error) {
	var ppmShipment models.PPMShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment",
			"WeightTickets",
			"MovingExpenses",
			"ProgearExpenses",
			"W2Address",
		).
		Find(&ppmShipment, id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(id, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
		}
	}

	err = loadPPMAssociations(appCtx, &ppmShipment)
	if err != nil {
		return nil, err
	}

	return &ppmShipment, nil
}

func loadPPMAssociations(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) error {
	for i := range ppmShipment.WeightTickets {
		if weightTicket := &ppmShipment.WeightTickets[i]; weightTicket.DeletedAt == nil {
			err := appCtx.DB().Load(weightTicket,
				"EmptyDocument.UserUploads.Upload",
				"FullDocument.UserUploads.Upload",
				"ProofOfTrailerOwnershipDocument.UserUploads.Upload")

			if err != nil {
				return apperror.NewQueryError("WeightTicket", err, "failed to load WeightTicket document uploads")
			}

			weightTicket.EmptyDocument.UserUploads = weightTicket.EmptyDocument.UserUploads.FilterDeleted()
			weightTicket.FullDocument.UserUploads = weightTicket.FullDocument.UserUploads.FilterDeleted()
			weightTicket.ProofOfTrailerOwnershipDocument.UserUploads = weightTicket.ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()
		}
	}

	for i := range ppmShipment.ProgearExpenses {
		if progearExpense := &ppmShipment.ProgearExpenses[i]; progearExpense.DeletedAt == nil {
			err := appCtx.DB().Load(progearExpense,
				"Document.UserUploads.Upload")

			if err != nil {
				return apperror.NewQueryError("ProgearExpenses", err, "failed to load ProgearExpenses document uploads")
			}

			progearExpense.Document.UserUploads = progearExpense.Document.UserUploads.FilterDeleted()
		}
	}

	for i := range ppmShipment.MovingExpenses {
		if movingExpense := &ppmShipment.MovingExpenses[i]; movingExpense.DeletedAt == nil {
			err := appCtx.DB().Load(movingExpense,
				"Document.UserUploads.Upload")

			if err != nil {
				return apperror.NewQueryError("MovingExpenses", err, "failed to load ProgearExpenses document uploads")
			}

			movingExpense.Document.UserUploads = movingExpense.Document.UserUploads.FilterDeleted()
		}
	}

	// We can't load SignedCertification with EagerPreload because of a bug in Pop, so we'll load it directly next.
	loadErr := appCtx.DB().Load(ppmShipment, "SignedCertification")
	// Pop will load an empty struct here and we'll get validation errors when attempting to save, if we don't set the
	// field to nil
	if ppmShipment.SignedCertification != nil && ppmShipment.SignedCertification.ID.IsNil() {
		ppmShipment.SignedCertification = nil
	}

	if loadErr != nil {
		return apperror.NewQueryError("PPMShipment", loadErr, "unable to load SignedCertification")
	}

	return nil
}

func FindPPMShipmentWithDocument(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, documentID uuid.UUID) error {
	var weightTicket models.WeightTicket
	var proGear models.ProgearWeightTicket
	var movingExpense models.MovingExpense

	err := appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope()).
		Where("ppm_shipment_id = ? AND (empty_document_id = ? OR full_document_id = ? OR proof_of_trailer_ownership_document_id = ?)", ppmShipmentID, documentID, documentID, documentID).
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
		Where("ppm_shipment_id = ? AND document_id = ?", ppmShipmentID, documentID).
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

func FetchPPMShipmentFromMTOShipmentID(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.PPMShipment, error) {
	var ppmShipment models.PPMShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).EagerPreload("Shipment", "W2Address", "WeightTickets").
		Where("ppm_shipments.shipment_id = ?", mtoShipmentID).
		First(&ppmShipment)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "")
		}
	}
	return &ppmShipment, nil
}
