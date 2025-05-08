package acknowledgemovesshipments

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveAndShipmentAcknowledgementUpdater struct {
}

// NewMoveAndShipmentAcknowledgementUpdater returns a new moveAndShipmentAcknowledgementUpdater
func NewMoveAndShipmentAcknowledgementUpdater() services.MoveAndShipmentAcknowledgementUpdater {
	return &moveAndShipmentAcknowledgementUpdater{}
}

func (p *moveAndShipmentAcknowledgementUpdater) AcknowledgeMovesAndShipments(appCtx appcontext.AppContext, moves *models.Moves) error {

	// Convert to JSON so we can pass it to the DB procedure
	jsonData, err := customMarshalJSON(*moves)
	if err != nil {
		return fmt.Errorf("error converting moves to json for prime_acknowledge_moves_shipments procedure: %w", err)
	}

	// Call procedure to update move and shipment acknowledgement dates
	err = appCtx.DB().RawQuery("CALL prime_acknowledge_moves_shipments($1)", string(jsonData)).Exec()
	if err != nil {
		return fmt.Errorf("error executing prime_acknowledge_moves_shipments procedure: %w", err)
	}
	return nil
}

func customMarshalJSON(v interface{}) ([]byte, error) {
	switch obj := v.(type) {
	case models.Move:
		return marshalMove(obj)
	case models.Moves:
		return marshalMoves(obj)
	case models.MTOShipment:
		return marshalMTOShipment(obj)
	case models.MTOShipments:
		return marshalMTOShipments(obj)
	default:
		return nil, fmt.Errorf("unsupported type for custom JSON marshaling")
	}
}

// Only marshaling ID, MTOShipments, and Prime AcknowledgedAt fields
func marshalMove(m models.Move) ([]byte, error) {
	marshaledShipments, err := marshalMTOShipments(m.MTOShipments)
	if err != nil {
		return nil, fmt.Errorf("error marshaling MTOShipments: %w", err)
	}

	return json.Marshal(struct {
		ID                  uuid.UUID       `json:"id"`
		MTOShipments        json.RawMessage `json:"mtoShipments"`
		PrimeAcknowledgedAt *time.Time      `json:"primeAcknowledgedAt"`
	}{
		ID:                  m.ID,
		MTOShipments:        marshaledShipments,
		PrimeAcknowledgedAt: m.PrimeAcknowledgedAt,
	})
}

func marshalMoves(moves models.Moves) ([]byte, error) {
	customMoves := make([]json.RawMessage, len(moves))
	for i, move := range moves {
		marshaledMove, err := marshalMove(move)
		if err != nil {
			return nil, err
		}
		customMoves[i] = marshaledMove
	}
	return json.Marshal(customMoves)
}

// Only marshaling ID, and Prime AcknowledgedAt fields
func marshalMTOShipment(s models.MTOShipment) ([]byte, error) {
	return json.Marshal(struct {
		ID                  uuid.UUID  `json:"id"`
		PrimeAcknowledgedAt *time.Time `json:"primeAcknowledgedAt"`
	}{
		ID:                  s.ID,
		PrimeAcknowledgedAt: s.PrimeAcknowledgedAt,
	})
}

func marshalMTOShipments(shipments models.MTOShipments) ([]byte, error) {
	customShipments := make([]json.RawMessage, len(shipments))
	for i, shipment := range shipments {
		marshaledShipment, err := marshalMTOShipment(shipment)
		if err != nil {
			return nil, err
		}
		customShipments[i] = marshaledShipment
	}
	return json.Marshal(customShipments)
}
