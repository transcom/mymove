package ppmshipment

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmShipmentFetcher is the concrete struct implementing the PPMShipmentFetcher interface
type ppmShipmentFetcher struct{}

// NewPPMShipmentFetcher creates a new PPMShipmentFetcher
func NewPPMShipmentFetcher() services.PPMShipmentFetcher {
	return &ppmShipmentFetcher{}
}

// These are helper constants for requesting eager preload associations
const (
	// EagerPreloadAssociationShipment is the name of the association for the shipment
	EagerPreloadAssociationShipment = "Shipment"
	// EagerPreloadAssociationServiceMember is the name of the association for the service member
	EagerPreloadAssociationServiceMember = "Shipment.MoveTaskOrder.Orders.ServiceMember"
	// EagerPreloadAssociationWeightTickets is the name of the association for the weight tickets
	EagerPreloadAssociationWeightTickets = "WeightTickets"
	// EagerPreloadAssociationProgearWeightTickets is the name of the association for the pro-gear weight tickets
	EagerPreloadAssociationProgearWeightTickets = "ProgearWeightTickets"
	// EagerPreloadAssociationMovingExpenses is the name of the association for the moving expenses
	EagerPreloadAssociationMovingExpenses = "MovingExpenses"
	// EagerPreloadAssociationW2Address is the name of the association for the W2 address
	EagerPreloadAssociationW2Address = "W2Address"
	// EagerPreloadAssociationAOAPacket is the name of the association for the AOA packet
	EagerPreloadAssociationAOAPacket = "AOAPacket"
	// EagerPreloadAssociationPaymentPacket is the name of the association for the payment packet
	EagerPreloadAssociationPaymentPacket = "PaymentPacket"
	// EagerPreloadAssociationPickupAddress is the name of the association for the Pickup address
	EagerPreloadAssociationPickupAddress = "PickupAddress"
	// EagerPreloadAssociationSecondaryPickupAddress is the name of the association for the Secondary Pickup address
	EagerPreloadAssociationSecondaryPickupAddress = "SecondaryPickupAddress"
	// EagerPreloadAssociationSecondaryPickupAddress is the name of the association for the Tertiary Pickup address
	EagerPreloadAssociationTertiaryPickupAddress = "TertiaryPickupAddress"
	// EagerPreloadAssociationDestinationAddress is the name of the association for the Destination address
	EagerPreloadAssociationDestinationAddress = "DestinationAddress"
	// EagerPreloadAssociationSecondaryDestinationAddress is the name of the association for the Secondary Destination address
	EagerPreloadAssociationSecondaryDestinationAddress = "SecondaryDestinationAddress"
	// EagerPreloadAssociationSecondaryDestinationAddress is the name of the association for the Tertiary Destination address
	EagerPreloadAssociationTertiaryDestinationAddress = "TertiaryDestinationAddress"
)

// These are helper constants for requesting post load associations, meaning associations that can't be eager pre-loaded
// due to bugs in pop
const (
	// PostLoadAssociationSignedCertification is the name of the association for the signed certification
	PostLoadAssociationSignedCertification = "SignedCertification"
	// PostLoadAssociationWeightTicketUploads is the name of the association for the weight ticket uploads
	PostLoadAssociationWeightTicketUploads = "WeightTicketUploads"
	// PostLoadAssociationProgearWeightTicketUploads is the name of the association for the pro-gear weight ticket uploads
	PostLoadAssociationProgearWeightTicketUploads = "ProgearWeightTicketUploads"
	// PostLoadAssociationMovingExpenseUploads is the name of the association for the moving expense uploads
	PostLoadAssociationMovingExpenseUploads = "MovingExpenseUploads"
	// PostLoadAssociationUploadedOrders is the name of the association for the orders uploaded by the service member
	PostLoadAssociationUploadedOrders = "UploadedOrders"
)

// GetListOfAllPreloadAssociations returns all associations for a PPMShipment that can be eagerly preloaded for ease of use.
func GetListOfAllPreloadAssociations() []string {
	return []string{
		EagerPreloadAssociationShipment,
		EagerPreloadAssociationServiceMember,
		EagerPreloadAssociationWeightTickets,
		EagerPreloadAssociationProgearWeightTickets,
		EagerPreloadAssociationMovingExpenses,
		EagerPreloadAssociationW2Address,
		EagerPreloadAssociationAOAPacket,
		EagerPreloadAssociationPaymentPacket,
		EagerPreloadAssociationPickupAddress,
		EagerPreloadAssociationDestinationAddress,
		EagerPreloadAssociationSecondaryPickupAddress,
		EagerPreloadAssociationSecondaryDestinationAddress,
		EagerPreloadAssociationTertiaryPickupAddress,
		EagerPreloadAssociationTertiaryDestinationAddress,
	}
}

// GetListOfAllPostloadAssociations returns all associations for a PPMShipment that can't be eagerly preloaded due to bugs in pop
func GetListOfAllPostloadAssociations() []string {
	return []string{
		PostLoadAssociationSignedCertification,
		PostLoadAssociationWeightTicketUploads,
		PostLoadAssociationProgearWeightTicketUploads,
		PostLoadAssociationMovingExpenseUploads,
		PostLoadAssociationUploadedOrders,
	}
}

// GetPPMShipment returns a PPMShipment with any desired associations by ID
func (f ppmShipmentFetcher) GetPPMShipment(
	appCtx appcontext.AppContext,
	ppmShipmentID uuid.UUID,
	eagerPreloadAssociations []string,
	postloadAssociations []string,
) (*models.PPMShipment, error) {
	if eagerPreloadAssociations != nil {
		validPreloadAssociations := make(map[string]bool)
		for _, v := range GetListOfAllPreloadAssociations() {
			validPreloadAssociations[v] = true
		}

		for _, association := range eagerPreloadAssociations {
			if !validPreloadAssociations[association] {
				msg := fmt.Sprintf("Requested eager preload association %s is not implemented", association)

				return nil, apperror.NewNotImplementedError(msg)
			}
		}
	}

	var ppmShipment models.PPMShipment

	q := appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope(models.PPMShipment{}))

	if eagerPreloadAssociations != nil {
		q.EagerPreload(eagerPreloadAssociations...)
	}

	if appCtx.Session() != nil && appCtx.Session().IsMilApp() {
		q.
			InnerJoin("mto_shipments", "mto_shipments.id = ppm_shipments.shipment_id").
			InnerJoin("moves", "moves.id = mto_shipments.move_id").
			InnerJoin("orders", "orders.id = moves.orders_id").
			Where("orders.service_member_id = ?", appCtx.Session().ServiceMemberID)
	}

	err := q.Find(&ppmShipment, ppmShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
		}
	}

	ppmShipment.WeightTickets = ppmShipment.WeightTickets.FilterDeleted()
	ppmShipment.ProgearWeightTickets = ppmShipment.ProgearWeightTickets.FilterDeleted()
	ppmShipment.MovingExpenses = ppmShipment.MovingExpenses.FilterDeleted()

	if postloadAssociations != nil {
		postloadErr := f.PostloadAssociations(appCtx, &ppmShipment, postloadAssociations)

		if postloadErr != nil {
			return nil, postloadErr
		}
	}

	return &ppmShipment, nil
}

// PostloadAssociations loads associations that can't be eager preloaded due to bugs in pop
func (f ppmShipmentFetcher) PostloadAssociations(
	appCtx appcontext.AppContext,
	ppmShipment *models.PPMShipment,
	postloadAssociations []string,
) error {
	for _, association := range postloadAssociations {
		switch association {
		case PostLoadAssociationSignedCertification:
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
		case PostLoadAssociationWeightTicketUploads:
			for i := range ppmShipment.WeightTickets {
				weightTicket := &ppmShipment.WeightTickets[i]

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
		case PostLoadAssociationProgearWeightTicketUploads:
			for i := range ppmShipment.ProgearWeightTickets {
				progearWeightTicket := &ppmShipment.ProgearWeightTickets[i]
				err := appCtx.DB().Load(progearWeightTicket,
					"Document.UserUploads.Upload")

				if err != nil {
					return apperror.NewQueryError("ProgearWeightTickets", err, "failed to load ProgearWeightTickets document uploads")
				}

				progearWeightTicket.Document.UserUploads = progearWeightTicket.Document.UserUploads.FilterDeleted()
			}

		case PostLoadAssociationMovingExpenseUploads:
			for i := range ppmShipment.MovingExpenses {
				movingExpense := &ppmShipment.MovingExpenses[i]
				err := appCtx.DB().Load(movingExpense,
					"Document.UserUploads.Upload")

				if err != nil {
					return apperror.NewQueryError("MovingExpenses", err, "failed to load MovingExpenses document uploads")
				}

				movingExpense.Document.UserUploads = movingExpense.Document.UserUploads.FilterDeleted()
			}
		case PostLoadAssociationUploadedOrders:
			err := appCtx.DB().Load(ppmShipment, "Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads.Upload")

			if err != nil {
				return apperror.NewQueryError("PPMShipment", err, "failed to load PPMShipment uploaded orders")
			}
		default:
			return apperror.NewNotImplementedError(fmt.Sprintf("Requested post load association %s is not implemented", association))
		}
	}

	return nil
}

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
			"ProgearWeightTickets",
			"W2Address.Country",
			"PickupAddress.Country",
			"SecondaryPickupAddress.Country",
			"TertiaryPickupAddress.Country",
			"DestinationAddress.Country",
			"SecondaryDestinationAddress.Country",
			"TertiaryDestinationAddress.Country",
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
			"ProgearWeightTickets",
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
	ppmShipment.WeightTickets = ppmShipment.WeightTickets.FilterDeleted()
	for i := range ppmShipment.WeightTickets {
		weightTicket := &ppmShipment.WeightTickets[i]
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

	ppmShipment.ProgearWeightTickets = ppmShipment.ProgearWeightTickets.FilterDeleted()
	for i := range ppmShipment.ProgearWeightTickets {
		progearWeightTicket := &ppmShipment.ProgearWeightTickets[i]
		err := appCtx.DB().Load(progearWeightTicket,
			"Document.UserUploads.Upload")

		if err != nil {
			return apperror.NewQueryError("ProgearWeightTickets", err, "failed to load ProgearWeightTickets document uploads")
		}

		progearWeightTicket.Document.UserUploads = progearWeightTicket.Document.UserUploads.FilterDeleted()
	}

	ppmShipment.MovingExpenses = ppmShipment.MovingExpenses.FilterDeleted()
	for i := range ppmShipment.MovingExpenses {
		movingExpense := &ppmShipment.MovingExpenses[i]
		err := appCtx.DB().Load(movingExpense,
			"Document.UserUploads.Upload")

		if err != nil {
			return apperror.NewQueryError("MovingExpenses", err, "failed to load  MovingExpenses document uploads")
		}

		movingExpense.Document.UserUploads = movingExpense.Document.UserUploads.FilterDeleted()
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

// returns true if moves orders are from a location that does not provide service counseling
func IsPrimeCounseledPPM(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (bool, error) {
	var ppmDutyLocation models.DutyLocation

	err := appCtx.DB().Q().
		Join("orders", "duty_locations.id = orders.origin_duty_location_id").
		Join("moves", "orders.id = moves.orders_id ").
		Join("mto_shipments", "moves.id = mto_shipments.move_id").
		Join("ppm_shipments", "mto_shipments.id = ppm_shipments.shipment_id").
		Where("ppm_shipments.shipment_id = ?", mtoShipmentID).
		First(&ppmDutyLocation)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, apperror.NewNotFoundError(mtoShipmentID, "while looking for PPMShipment")
		default:
			return false, apperror.NewQueryError("PPMShipment", err, "")
		}
	}

	return !ppmDutyLocation.ProvidesServicesCounseling, err
}
