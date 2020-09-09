package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakePaymentRequest creates a single PaymentRequest and associated set relationships
func MakePaymentRequest(db *pop.Connection, assertions Assertions) models.PaymentRequest {
	// Create new PaymentRequest if not provided
	// ID is required because it must be populated for Eager saving to work.
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	paymentRequestNumber := assertions.PaymentRequest.PaymentRequestNumber
	sequenceNumber := assertions.PaymentRequest.SequenceNumber
	if paymentRequestNumber == "" {
		if sequenceNumber == 0 {
			sequenceNumber = 1
		}
		paymentRequestNumber = fmt.Sprintf("%s-%d", *moveTaskOrder.ReferenceID, sequenceNumber)
	}

	paymentRequest := models.PaymentRequest{
		MoveTaskOrder:        moveTaskOrder,
		MoveTaskOrderID:      moveTaskOrder.ID,
		IsFinal:              false,
		RejectionReason:      nil,
		Status:               models.PaymentRequestStatusPending,
		PaymentRequestNumber: paymentRequestNumber,
		SequenceNumber:       sequenceNumber,
	}

	// Overwrite values with those from assertions
	mergeModels(&paymentRequest, assertions.PaymentRequest)

	mustCreate(db, &paymentRequest)

	return paymentRequest
}

// MakePaymentRequestWithParams creates a single PaymentRequest and associated PaymentServiceItems and Params
func MakePaymentRequestWithParams(db *pop.Connection, assertions Assertions) models.PaymentRequest {
	paymentRequest := MakePaymentRequest(db, assertions)
	// add PSIs
	var params models.PaymentServiceItemParams

	paymentServiceItem := MakePaymentServiceItem(db, Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDLH,
		},
		PaymentServiceItem: models.PaymentServiceItem{
			PaymentRequestID: paymentRequest.ID,
		},
	})

	// contract code param
	contractCodeServiceItemParamKey := FetchOrMakeServiceItemParamKey(db,
		Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:  models.ServiceItemParamNameContractCode,
				Type: models.ServiceItemParamTypeString,
			},
		})

	contractCodeServiceItemParam := MakePaymentServiceItemParam(db,
		Assertions{
			PaymentServiceItem:  paymentServiceItem,
			ServiceItemParamKey: contractCodeServiceItemParamKey,
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: DefaultContractCode,
			},
		})
	params = append(params, contractCodeServiceItemParam)

	// billed actual weight param
	weightServiceItemParamKey := FetchOrMakeServiceItemParamKey(db,
		Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:  models.ServiceItemParamNameWeightBilledActual,
				Type: models.ServiceItemParamTypeInteger,
			},
		})

	weightServiceItemParam := MakePaymentServiceItemParam(db,
		Assertions{
			PaymentServiceItem:  paymentServiceItem,
			ServiceItemParamKey: weightServiceItemParamKey,
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "4242",
			},
		})
	params = append(params, weightServiceItemParam)

	// distance zip3 param
	distanceServiceItemParamKey := FetchOrMakeServiceItemParamKey(db,
		Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:  models.ServiceItemParamNameDistanceZip3,
				Type: models.ServiceItemParamTypeInteger,
			},
		})

	distanceServiceItemParam := MakePaymentServiceItemParam(db,
		Assertions{
			PaymentServiceItem:  paymentServiceItem,
			ServiceItemParamKey: distanceServiceItemParamKey,
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "2424",
			},
		})
	params = append(params, distanceServiceItemParam)

	paymentServiceItem.PaymentServiceItemParams = params
	mustSave(db, &paymentServiceItem)

	return paymentRequest
}

// MakeDefaultPaymentRequest makes an PaymentRequest with default values
func MakeDefaultPaymentRequest(db *pop.Connection) models.PaymentRequest {
	return MakePaymentRequest(db, Assertions{})
}
