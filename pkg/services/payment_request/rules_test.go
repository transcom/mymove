package paymentrequest

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PaymentRequestServiceSuite) TestValidationRules() {

	suite.Run("checkMTOIDField", func() {

		suite.Run("success", func() {

			move := factory.BuildMove(suite.DB(), nil, nil)
			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: move.ID,
			}

			err := checkMTOIDField().Validate(suite.AppContextForTest(), paymentRequest, nil)
			suite.NoError(err)
		})

		suite.Run("failure", func() {

			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")),
			}

			err := checkMTOIDField().Validate(suite.AppContextForTest(), paymentRequest, nil)
			switch err.(type) {
			case apperror.InvalidCreateInputError:
				suite.Equal(err.Error(), "Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create")
			default:
				suite.Failf("expected *apperror.InvalidCreateInputError", "%v", err)
			}
		})

	})

	suite.Run("checkMTOIDMatchesServiceItemMTOID", func() {

		suite.Run("success", func() {

			move := factory.BuildMove(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
			testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EndDate: time.Now().Add(time.Hour * 24),
				},
			})
			estimatedWeight := unit.Pound(2048)
			serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
					},
				},
				{
					Model: models.MTOShipment{
						PrimeEstimatedWeight: &estimatedWeight,
					},
				},
				{
					Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
				},
			}, nil)

			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: move.ID,
				IsFinal:         false,
				PaymentServiceItems: models.PaymentServiceItems{
					{
						MTOServiceItemID: serviceItem.ID,
						MTOServiceItem:   serviceItem,
						PaymentServiceItemParams: models.PaymentServiceItemParams{
							{
								IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
								Value:       "3254",
							},
							{
								IncomingKey: models.ServiceItemParamNameRequestedPickupDate.String(),
								Value:       "2019-12-16",
							},
						},
					},
				},
			}

			err := checkMTOIDMatchesServiceItemMTOID().Validate(suite.AppContextForTest(), paymentRequest, nil)
			suite.NoError(err)
		})

		suite.Run("failure", func() {

			move := factory.BuildMove(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

			testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EndDate: time.Now().Add(time.Hour * 24),
				},
			})
			estimatedWeight := unit.Pound(2048)
			serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
					},
				},
				{
					Model: models.MTOShipment{
						PrimeEstimatedWeight: &estimatedWeight,
					},
				},
				{
					Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
				},
			}, nil)

			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				PaymentServiceItems: models.PaymentServiceItems{
					{
						MTOServiceItemID: serviceItem.ID,
						MTOServiceItem:   serviceItem,
						PaymentServiceItemParams: models.PaymentServiceItemParams{
							{
								IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
								Value:       "3254",
							},
							{
								IncomingKey: models.ServiceItemParamNameRequestedPickupDate.String(),
								Value:       "2019-12-16",
							},
						},
					},
				},
			}

			err := checkMTOIDMatchesServiceItemMTOID().Validate(suite.AppContextForTest(), paymentRequest, nil)
			suite.Error(err)
			suite.Contains(err.Error(), "Conflict Error: Payment Request MoveTaskOrderID does not match Service Item MoveTaskOrderID")
		})

	})

	suite.Run("checkStatusOfExistingPaymentRequest", func() {

		suite.Run("success", func() {

			move := factory.BuildMove(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
			testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EndDate: time.Now().Add(time.Hour * 24),
				},
			})
			estimatedWeight := unit.Pound(2048)
			serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
					},
				},
				{
					Model: models.MTOShipment{
						PrimeEstimatedWeight: &estimatedWeight,
					},
				},
				{
					Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
				},
			}, nil)

			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: move.ID,
				IsFinal:         false,
				PaymentServiceItems: models.PaymentServiceItems{
					{
						MTOServiceItemID: serviceItem.ID,
						MTOServiceItem:   serviceItem,
						PaymentServiceItemParams: models.PaymentServiceItemParams{
							{
								IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
								Value:       "3254",
							},
							{
								IncomingKey: models.ServiceItemParamNameRequestedPickupDate.String(),
								Value:       "2022-03-16",
							},
						},
						PaymentRequest: models.PaymentRequest{
							ReviewedAt: models.TimePointer(time.Now()),
						},
					},
				},
			}

			err := checkStatusOfExistingPaymentRequest().Validate(suite.AppContextForTest(), paymentRequest, nil)
			suite.NoError(err)
		})

		suite.Run("failure", func() {

			move := factory.BuildMove(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

			testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EndDate: time.Now().Add(time.Hour * 24),
				},
			})
			estimatedWeight := unit.Pound(2048)
			serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
					},
				},
				{
					Model: models.MTOShipment{
						PrimeEstimatedWeight: &estimatedWeight,
					},
				},
				{
					Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
				},
			}, nil)

			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: move.ID,
				IsFinal:         false,
				ReviewedAt:      nil,
				PaymentServiceItems: models.PaymentServiceItems{
					{
						MTOServiceItemID: serviceItem.ID,
						MTOServiceItem:   serviceItem,
						PaymentServiceItemParams: models.PaymentServiceItemParams{
							{
								IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
								Value:       "3254",
							},
							{
								IncomingKey: models.ServiceItemParamNameRequestedPickupDate.String(),
								Value:       "2023-12-16",
							},
						},
					},
				},
			}

			err := checkStatusOfExistingPaymentRequest().Validate(suite.AppContextForTest(), paymentRequest, nil)
			suite.Error(err)
			suite.Contains(err.Error(), "Conflict Error: Payment Request for Service Item is already paid or requested")
		})

	})

}
