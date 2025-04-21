package paymentrequest

import (
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
			testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate: testdatagen.ContractStartDate,
					EndDate:   testdatagen.ContractEndDate,
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

			testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate: testdatagen.ContractStartDate,
					EndDate:   testdatagen.ContractEndDate,
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

	suite.Run("checkValidSitAddlDates", func() {

		suite.Run("success", func() {

			move := factory.BuildMove(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
			testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate: testdatagen.ContractStartDate,
					EndDate:   testdatagen.ContractEndDate,
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
								IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
								Value:       "2024-02-22",
							},
							{
								IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
								Value:       "2024-02-23",
							},
						},
					},
				},
			}

			err := checkValidSitAddlDates().Validate(suite.AppContextForTest(), paymentRequest, nil)
			suite.NoError(err)
		})

		suite.Run("failure", func() {
			testCases := []struct {
				sitStartDate string
				sitEndDate   string
				errorString  string
			}{
				{"01-01-2024", "2024-02-21", "Invalid Create Input Error: SITPaymentRequestStart must be a valid date value of YYYY-MM-DD"},
				{"2024-02-21", "01-01-2024", "Invalid Create Input Error: SITPaymentRequestEnd must be a valid date value of YYYY-MM-DD"},
				{"2024-02-22", "2024-02-21", "Invalid Create Input Error: SITPaymentRequestStart must be a date that comes before SITPaymentRequestEnd"},
			}

			move := factory.BuildMove(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
			testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate: testdatagen.ContractStartDate,
					EndDate:   testdatagen.ContractEndDate,
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
			for _, testCase := range testCases {

				paymentRequest := models.PaymentRequest{
					MoveTaskOrderID: move.ID,
					IsFinal:         false,
					PaymentServiceItems: models.PaymentServiceItems{
						{
							MTOServiceItemID: serviceItem.ID,
							MTOServiceItem:   serviceItem,
							PaymentServiceItemParams: models.PaymentServiceItemParams{
								{
									IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
									Value:       testCase.sitStartDate,
								},
								{
									IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
									Value:       testCase.sitEndDate,
								},
							},
						},
					},
				}

				err := checkValidSitAddlDates().Validate(suite.AppContextForTest(), paymentRequest, nil)
				suite.Error(err)
				suite.Contains(err.Error(), testCase.errorString)
			}

		})

	})

}
