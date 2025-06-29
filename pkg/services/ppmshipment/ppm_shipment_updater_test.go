package ppmshipment

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const dateOnly = "2006-01-02"

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {

	// One-time test setup

	fakeEstimatedIncentive := models.CentPointer(unit.Cents(1000000))

	type updateSubtestData struct {
		ppmShipmentUpdater services.PPMShipmentUpdater
	}

	// setUpForTests - Sets up objects/mocks that need to be set up on a per-test basis.
	setUpForTests := func(estimatedIncentiveAmount *unit.Cents, sitEstimatedCost *unit.Cents, maxIncentiveAmount *unit.Cents, estimatedIncentiveError error) (subtestData updateSubtestData) {
		ppmEstimator := mocks.PPMEstimator{}
		ppmEstimator.
			On(
				"FinalIncentiveWithDefaultChecks",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(nil, nil)

		ppmEstimator.
			On(
				"EstimateIncentiveWithDefaultChecks",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(estimatedIncentiveAmount, sitEstimatedCost, estimatedIncentiveError)

		ppmEstimator.
			On(
				"MaxIncentive",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(maxIncentiveAmount, nil)

		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		subtestData.ppmShipmentUpdater = NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		return subtestData
	}

	setUpForFinalIncentiveTests := func(finalIncentiveAmount *unit.Cents, finalIncentiveError error, estimatedIncentiveAmount *unit.Cents, sitEstimatedCost *unit.Cents, estimatedIncentiveError error) (subtestData updateSubtestData) {
		ppmEstimator := mocks.PPMEstimator{}
		ppmEstimator.
			On(
				"FinalIncentiveWithDefaultChecks",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(finalIncentiveAmount, finalIncentiveError)

		ppmEstimator.
			On(
				"EstimateIncentiveWithDefaultChecks",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(estimatedIncentiveAmount, sitEstimatedCost, estimatedIncentiveError)

		ppmEstimator.
			On(
				"MaxIncentive",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(nil, nil)

		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		subtestData.ppmShipmentUpdater = NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		return subtestData
	}
	setupPricerData := func() {
		testdatagen.FetchOrMakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
			GHCDieselFuelPrice: models.GHCDieselFuelPrice{
				FuelPriceInMillicents: unit.Millicents(281400),
				PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
			},
		})

		originDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "056",
				ServicesSchedule: 3,
				SITPDSchedule:    3,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             originDomesticServiceArea.Contract,
				ContractID:           originDomesticServiceArea.ContractID,
				StartDate:            time.Date(2019, time.June, 1, 0, 0, 0, 0, time.UTC),
				EndDate:              time.Date(2020, time.May, 31, 0, 0, 0, 0, time.UTC),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originDomesticServiceArea.Contract,
				ContractID:          originDomesticServiceArea.ContractID,
				DomesticServiceArea: originDomesticServiceArea,
				Zip3:                "902",
			},
		})

		destDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    originDomesticServiceArea.Contract,
				ContractID:  originDomesticServiceArea.ContractID,
				ServiceArea: "208",
			},
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            destDomesticServiceArea.Contract,
				ContractID:          destDomesticServiceArea.ContractID,
				DomesticServiceArea: destDomesticServiceArea,
				Zip3:                "308",
			},
		})

		testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
				Contract:              originDomesticServiceArea.Contract,
				ContractID:            originDomesticServiceArea.ContractID,
				DomesticServiceArea:   originDomesticServiceArea,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				WeightLower:           unit.Pound(500),
				WeightUpper:           unit.Pound(4999),
				MilesLower:            2001,
				MilesUpper:            2500,
				PriceMillicents:       unit.Millicents(412400),
			},
		})

		testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
				Contract:              originDomesticServiceArea.Contract,
				ContractID:            originDomesticServiceArea.ContractID,
				DomesticServiceArea:   originDomesticServiceArea,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				WeightLower:           unit.Pound(500),
				WeightUpper:           unit.Pound(4999),
				MilesLower:            2001,
				MilesUpper:            2500,
				IsPeakPeriod:          true,
				PriceMillicents:       unit.Millicents(437600),
			},
		})

		testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
				Contract:              originDomesticServiceArea.Contract,
				ContractID:            originDomesticServiceArea.ContractID,
				DomesticServiceArea:   originDomesticServiceArea,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				WeightLower:           unit.Pound(5000),
				WeightUpper:           unit.Pound(9999),
				MilesLower:            2001,
				MilesUpper:            2500,
				PriceMillicents:       unit.Millicents(606800),
			},
		})

		dopService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             dopService.ID,
				Service:               dopService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          false,
				PriceCents:            unit.Cents(404),
			},
		})

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             dopService.ID,
				Service:               dopService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            unit.Cents(465),
			},
		})

		ddpService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddpService.ID,
				Service:               ddpService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          false,
				PriceCents:            unit.Cents(832),
			},
		})

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddpService.ID,
				Service:               ddpService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            unit.Cents(957),
			},
		})

		dpkService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)

		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				ContractID:   originDomesticServiceArea.ContractID,
				Contract:     originDomesticServiceArea.Contract,
				ServiceID:    dpkService.ID,
				Service:      dpkService,
				IsPeakPeriod: false,
				Schedule:     3,
				PriceCents:   7395,
			},
		})

		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				ContractID:   originDomesticServiceArea.ContractID,
				Contract:     originDomesticServiceArea.Contract,
				ServiceID:    dpkService.ID,
				Service:      dpkService,
				IsPeakPeriod: true,
				Schedule:     3,
				PriceCents:   8000,
			},
		})

		dupkService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)

		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				ContractID:   destDomesticServiceArea.ContractID,
				Contract:     destDomesticServiceArea.Contract,
				ServiceID:    dupkService.ID,
				Service:      dupkService,
				IsPeakPeriod: false,
				Schedule:     2,
				PriceCents:   597,
			},
		})

		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				ContractID:   destDomesticServiceArea.ContractID,
				Contract:     destDomesticServiceArea.Contract,
				ServiceID:    dupkService.ID,
				Service:      dupkService,
				IsPeakPeriod: true,
				Schedule:     2,
				PriceCents:   650,
			},
		})

		dofsitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             dofsitService.ID,
				Service:               dofsitService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          false,
				PriceCents:            1153,
			},
		})

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             dofsitService.ID,
				Service:               dofsitService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            1326,
			},
		})

		doasitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOASIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             doasitService.ID,
				Service:               doasitService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          false,
				PriceCents:            46,
			},
		})

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             doasitService.ID,
				Service:               doasitService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            53,
			},
		})

		ddfsitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddfsitService.ID,
				Service:               ddfsitService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          false,
				PriceCents:            1612,
			},
		})

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddfsitService.ID,
				Service:               ddfsitService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            1854,
			},
		})

		ddasitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDASIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddasitService.ID,
				Service:               ddasitService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          false,
				PriceCents:            55,
			},
		})

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddasitService.ID,
				Service:               ddasitService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            63,
			},
		})
	}

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil, nil)

		originalPPM := factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: testdatagen.NextValidMoveDate,
					SITExpected:           models.BoolPointer(false),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "987 Other Avenue",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50309",
					County:         models.StringPointer("POLK"),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "987 Other Avenue",
					StreetAddress2: models.StringPointer("P.O. Box 12345"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Fort Eisenhower",
					State:          "GA",
					PostalCode:     "50309",
					County:         models.StringPointer("COLUMBIA"),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		newPPM := models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek),
			SITExpected:           models.BoolPointer(true),
			PickupAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50308",
				County:         models.StringPointer("POLK"),
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 12345"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30183",
				County:         models.StringPointer("COLUMBIA"),
			},
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that should now be updated
		suite.Equal(newPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(newPPM.PickupAddress.PostalCode, updatedPPM.PickupAddress.PostalCode)
		suite.Equal(newPPM.DestinationAddress.PostalCode, updatedPPM.DestinationAddress.PostalCode)
		suite.Equal(newPPM.SITExpected, updatedPPM.SITExpected)

		// Estimated Incentive shouldn't be set since we don't have all the necessary fields.
		suite.Nil(updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment and shipment market code reflects international shipment", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil, nil)

		originalPPM := factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: testdatagen.NextValidMoveDate,
					SITExpected:           models.BoolPointer(false),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "987 Other Avenue",
					StreetAddress2: models.StringPointer("P.O. Box 1234"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50309",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "987 Other Avenue",
					StreetAddress2: models.StringPointer("P.O. Box 12345"),
					StreetAddress3: models.StringPointer("c/o Another Person"),
					City:           "Fort Eisenhower",
					State:          "GA",
					PostalCode:     "50309",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		newPPM := models.PPMShipment{
			PickupAddress: &models.Address{
				StreetAddress1: "987 Cold Avenue",
				City:           "Anchorage",
				State:          "AK",
				PostalCode:     "99501",
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 12345"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30183",
			},
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)
		suite.NilOrNoVerrs(err)

		// since one of the addresses is being updated to be OCONUS, the shipment's market code should change
		updatedShipment := models.MTOShipment{}
		err = suite.DB().Find(&updatedShipment, updatedPPM.ShipmentID)
		suite.NoError(err)
		suite.Equal(updatedShipment.MarketCode, models.MarketCodeInternational)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated dates & locations - weights already set", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))
		newFakeMaxIncentive := models.CentPointer(unit.Cents(5000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, newFakeMaxIncentive, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: testdatagen.NextValidMoveDate,
					SITExpected:           models.BoolPointer(false),
					EstimatedWeight:       models.PoundPointer(4000),
					HasProGear:            models.BoolPointer(false),
					EstimatedIncentive:    fakeEstimatedIncentive,
				},
			},
		}, nil)

		newPPM := models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek),
			SITExpected:           models.BoolPointer(true),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)

		// Fields that should now be updated
		suite.Equal(newPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(newPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*newFakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
		suite.Equal(*newFakeMaxIncentive, *updatedPPM.MaxIncentive)
		suite.Equal(updatedPPM.Shipment.MarketCode, models.MarketCodeDomestic)
	})

	suite.Run("Can successfully update a PPMShipment - add estimated weights - no pro gear", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), nil, nil)

		newPPM := models.PPMShipment{
			EstimatedWeight: models.PoundPointer(4000),
			HasProGear:      models.BoolPointer(false),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.Equal(*fakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - add estimated weights - has pro gear", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), nil, nil)

		newPPM := models.PPMShipment{
			EstimatedWeight:     models.PoundPointer(4000),
			HasProGear:          models.BoolPointer(true),
			ProGearWeight:       models.PoundPointer(1000),
			SpouseProGearWeight: models.PoundPointer(0),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*newPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*newPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)

		// EstimatedIncentive should have been calculated and set
		suite.Nil(originalPPM.EstimatedIncentive)
		suite.Equal(*fakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated weights - pro gear no to yes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedWeight:    models.PoundPointer(4000),
					HasProGear:         models.BoolPointer(false),
					EstimatedIncentive: fakeEstimatedIncentive,
				},
			},
		}, nil)
		newPPM := models.PPMShipment{
			EstimatedWeight:     models.PoundPointer(4500),
			HasProGear:          models.BoolPointer(true),
			ProGearWeight:       models.PoundPointer(1000),
			SpouseProGearWeight: models.PoundPointer(0),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*newPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*newPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*newFakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit estimated weights - pro gear yes to no", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedWeight:     models.PoundPointer(4000),
					HasProGear:          models.BoolPointer(true),
					ProGearWeight:       models.PoundPointer(1000),
					SpouseProGearWeight: models.PoundPointer(0),
					EstimatedIncentive:  fakeEstimatedIncentive,
				},
			},
		}, nil)

		newPPM := models.PPMShipment{
			EstimatedWeight: models.PoundPointer(4500),
			HasProGear:      models.BoolPointer(false),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)

		// Fields that should now be updated
		suite.Equal(*newPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*newPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)
		suite.Equal(*newFakeEstimatedIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can successfully update a PPMShipment - edit just allowable weight", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedWeight: models.PoundPointer(4000),
					AllowableWeight: models.PoundPointer(3000),
				},
			},
		}, nil)

		newPPM := models.PPMShipment{
			AllowableWeight: models.PoundPointer(4545),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.EstimatedWeight, updatedPPM.EstimatedWeight)

		// Fields that should now be updated
		suite.Equal(*newPPM.AllowableWeight, *updatedPPM.AllowableWeight)
	})

	suite.Run("Can successfully update a PPMShipment - add advance info - no advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedWeight:    models.PoundPointer(4000),
					HasProGear:         models.BoolPointer(false),
					EstimatedIncentive: fakeEstimatedIncentive,
				},
			},
		}, nil)
		newPPM := models.PPMShipment{
			HasRequestedAdvance: models.BoolPointer(false),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Nil(updatedPPM.AdvanceAmountRequested)
	})

	suite.Run("Can successfully update a PPMShipment - add advance info - yes advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedWeight:    models.PoundPointer(4000),
					HasProGear:         models.BoolPointer(false),
					EstimatedIncentive: fakeEstimatedIncentive,
				},
			},
		}, nil)
		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(300000)),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)
		suite.Nil(updatedPPM.ProGearWeight)
		suite.Nil(updatedPPM.SpouseProGearWeight)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
	})

	suite.Run("Can successfully update a PPMShipment - office user edits requested advance", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})
		approved := models.PPMAdvanceStatusApproved
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedIncentive:     fakeEstimatedIncentive,
					HasRequestedAdvance:    models.BoolPointer(true),
					AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
					AdvanceStatus:          &approved,
				},
			},
		}, nil)
		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(200000)),
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(&approved, updatedPPM.AdvanceStatus)
	})

	suite.Run("Can successfully update a PPMShipment - office user approves advance request", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedIncentive:     fakeEstimatedIncentive,
					HasRequestedAdvance:    models.BoolPointer(true),
					AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
				},
			},
		}, nil)

		approved := models.PPMAdvanceStatusApproved

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
			AdvanceStatus:          &approved,
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(&approved, updatedPPM.AdvanceStatus)
	})

	suite.Run("Can successfully update a PPMShipment - office user rejects advance request", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedIncentive:     fakeEstimatedIncentive,
					HasRequestedAdvance:    models.BoolPointer(true),
					AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
				},
			},
		}, nil)
		rejected := models.PPMAdvanceStatusRejected

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
			AdvanceStatus:          &rejected,
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(&rejected, updatedPPM.AdvanceStatus)
	})

	suite.Run("Can successfully update a PPMShipment - edit advance - advance requested no to yes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					EstimatedIncentive:  fakeEstimatedIncentive,
					HasRequestedAdvance: models.BoolPointer(false),
				},
			},
		}, nil)
		approved := models.PPMAdvanceStatusApproved

		newPPM := models.PPMShipment{
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
			AdvanceStatus:          &approved,
		}

		subtestData := setUpForTests(originalPPM.EstimatedIncentive, nil, nil, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(originalPPM.SITExpected, updatedPPM.SITExpected)
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.EstimatedIncentive, *updatedPPM.EstimatedIncentive)

		// Fields that should now be updated
		suite.Equal(*newPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)
		suite.Equal(*newPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(&approved, updatedPPM.AdvanceStatus)
	})

	suite.Run("Can successfully update a PPMShipment - edit SIT - yes to no", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, nil, nil)
		sitLocation := models.SITLocationTypeOrigin

		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &sitLocation,
					SITEstimatedEntryDate:     models.TimePointer(testdatagen.NextValidMoveDate),
					SITEstimatedDepartureDate: models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek)),
					SITEstimatedWeight:        models.PoundPointer(1000),
					SITEstimatedCost:          models.CentPointer(unit.Cents(69900)),
				},
			},
		}, nil)
		newPPM := models.PPMShipment{
			SITExpected: models.BoolPointer(false),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(*originalPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)

		// Fields that should now be updated
		suite.Equal(*newPPM.SITExpected, *updatedPPM.SITExpected)
		suite.Nil(updatedPPM.SITLocation)
		suite.Nil(updatedPPM.SITEstimatedEntryDate)
		suite.Nil(updatedPPM.SITEstimatedDepartureDate)
		suite.Nil(updatedPPM.SITEstimatedWeight)
		suite.Nil(updatedPPM.SITEstimatedCost)
	})

	suite.Run("Can successfully update a PPMShipment - edit SIT - no to yes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))
		newFakeSITEstimatedCost := models.CentPointer(unit.Cents(62500))

		subtestData := setUpForTests(newFakeEstimatedIncentive, newFakeSITEstimatedCost, nil, nil)
		sitLocation := models.SITLocationTypeOrigin

		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					SITExpected: models.BoolPointer(false),
				},
			},
		}, nil)
		newPPM := models.PPMShipment{
			SITExpected:               models.BoolPointer(true),
			SITLocation:               &sitLocation,
			SITEstimatedEntryDate:     models.TimePointer(testdatagen.NextValidMoveDate),
			SITEstimatedDepartureDate: models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek)),
			SITEstimatedWeight:        models.PoundPointer(1000),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that shouldn't have changed
		suite.Equal(originalPPM.ExpectedDepartureDate.Format(dateOnly), updatedPPM.ExpectedDepartureDate.Format(dateOnly))
		suite.Equal(*originalPPM.EstimatedWeight, *updatedPPM.EstimatedWeight)
		suite.Equal(*originalPPM.HasProGear, *updatedPPM.HasProGear)
		suite.Equal(*originalPPM.ProGearWeight, *updatedPPM.ProGearWeight)
		suite.Equal(*originalPPM.SpouseProGearWeight, *updatedPPM.SpouseProGearWeight)
		suite.Equal(*originalPPM.AdvanceAmountRequested, *updatedPPM.AdvanceAmountRequested)
		suite.Equal(*originalPPM.HasRequestedAdvance, *updatedPPM.HasRequestedAdvance)

		// Fields that should now be updated
		suite.Equal(*newPPM.SITExpected, *updatedPPM.SITExpected)
		suite.Equal(*newPPM.SITLocation, *updatedPPM.SITLocation)
		suite.Equal(*newPPM.SITEstimatedEntryDate, *updatedPPM.SITEstimatedEntryDate)
		suite.Equal(*newPPM.SITEstimatedDepartureDate, *updatedPPM.SITEstimatedDepartureDate)
		suite.Equal(*newPPM.SITEstimatedWeight, *updatedPPM.SITEstimatedWeight)
		suite.Equal(*newFakeSITEstimatedCost, *updatedPPM.SITEstimatedCost)
	})

	suite.Run("Can successfully update a PPMShipment - final incentive and actual move date", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForFinalIncentiveTests(nil, nil, nil, nil, nil)

		today := time.Now()

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ActualMoveDate:  &today,
					EstimatedWeight: models.PoundPointer(unit.Pound(5000)),
				},
			},
		}, nil)

		newPPM := originalPPM

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Fields that should now be updated
		suite.Equal(newPPM.FinalIncentive, updatedPPM.FinalIncentive)
		suite.Equal(newPPM.ActualMoveDate, updatedPPM.ActualMoveDate)
	})

	suite.Run("Can't update if Shipment can't be found", func() {
		badMTOShipmentID := uuid.Must(uuid.NewV4())

		subtestData := setUpForTests(nil, nil, nil, nil)

		updatedPPMShipment, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextWithSessionForTest(&auth.Session{}), &models.PPMShipment{}, badMTOShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment by MTO ShipmentID", badMTOShipmentID.String()), err.Error())
	})

	suite.Run("Can't update if there is invalid input", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(nil, nil, nil, nil)

		originalPPMShipment := factory.BuildPPMShipment(appCtx.DB(), nil, nil)

		// Easiest invalid input to trigger is to set an invalid AdvanceAmountRequested value. The rest are harder to
		// trigger based on how the service object is set up.
		newPPMShipment := models.PPMShipment{
			AdvanceAmountRequested: models.CentPointer(unit.Cents(3000000)),
		}

		updatedPPMShipment, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPMShipment, originalPPMShipment.ShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while validating the PPM shipment.", err.Error())
	})

	suite.Run("Can't update if there is an error calculating incentive", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		fakeEstimatedIncentiveError := errors.New("failed to calculate incentive")
		subtestData := setUpForTests(nil, nil, nil, fakeEstimatedIncentiveError)

		originalPPMShipment := factory.BuildPPMShipment(appCtx.DB(), nil, nil)

		newPPMShipment := models.PPMShipment{
			HasRequestedAdvance: models.BoolPointer(false),
		}

		updatedPPMShipment, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPMShipment, originalPPMShipment.ShipmentID)

		suite.Nil(updatedPPMShipment)

		suite.Error(err)
		suite.Equal(fakeEstimatedIncentiveError, err)
	})

	suite.Run("Can successfully update a PPMShipment - add W-2 address", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), nil, nil)

		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "08004"

		newPPM := models.PPMShipment{
			W2Address: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.NotNil(updatedPPM.W2AddressID)
		suite.Equal(streetAddress1, updatedPPM.W2Address.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.W2Address.StreetAddress2)
		suite.Equal(city, updatedPPM.W2Address.City)
		suite.Equal(state, updatedPPM.W2Address.State)
		suite.Equal(postalCode, updatedPPM.W2Address.PostalCode)
	})

	suite.Run("Can successfully update a PPMShipment - add W-2 address with empty strings for optional fields", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), nil, nil)

		streetAddress1 := "1819 S Cedar Street"
		city := "Fayetteville"
		state := "NC"
		postalCode := "28314"

		newPPM := models.PPMShipment{
			W2Address: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: models.StringPointer(""),
				StreetAddress3: models.StringPointer(""),
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
		}
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.NotNil(updatedPPM.W2AddressID)
		suite.Equal(streetAddress1, updatedPPM.W2Address.StreetAddress1)
		suite.Equal(city, updatedPPM.W2Address.City)
		suite.Equal(state, updatedPPM.W2Address.State)
		suite.Equal(postalCode, updatedPPM.W2Address.PostalCode)
		suite.Nil(updatedPPM.W2Address.StreetAddress2)
		suite.Nil(updatedPPM.W2Address.StreetAddress3)
		suite.NotNil(updatedPPM.W2Address.Country)
	})

	suite.Run("Can successfully update a PPMShipment - modify W-2 address", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil, nil)

		address := factory.BuildAddress(appCtx.DB(), nil, nil)
		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    address,
				LinkOnly: true,
				Type:     &factory.Addresses.W2Address,
			},
		}, nil)
		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Cookstown"
		state := "NJ"
		postalCode := "08511"

		newPPM := models.PPMShipment{
			W2Address: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.Equal(address.ID, *updatedPPM.W2AddressID)
		suite.Equal(streetAddress1, updatedPPM.W2Address.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.W2Address.StreetAddress2)
		suite.Equal(city, updatedPPM.W2Address.City)
		suite.Equal(state, updatedPPM.W2Address.State)
		suite.Equal(postalCode, updatedPPM.W2Address.PostalCode)
		suite.Equal(*address.StreetAddress3, *updatedPPM.W2Address.StreetAddress3)
		suite.Equal(address.CountryId, updatedPPM.W2Address.CountryId)
	})

	suite.Run("Can successfully update a PPMShipment - add Pickup and Destination address", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		subtestData := setUpForTests(fakeEstimatedIncentive, nil, nil, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), nil, nil)

		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "08004"

		newPPM := models.PPMShipment{
			PickupAddress: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
			DestinationAddress: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
			SecondaryPickupAddress: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
			SecondaryDestinationAddress: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
			TertiaryPickupAddress: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
			TertiaryDestinationAddress: &models.Address{
				StreetAddress1: streetAddress1,
				StreetAddress2: &streetAddress2,
				City:           city,
				State:          state,
				PostalCode:     postalCode,
			},
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.NotNil(updatedPPM.PickupAddressID)
		suite.NotNil(updatedPPM.PickupAddress)
		suite.Equal(streetAddress1, updatedPPM.PickupAddress.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.PickupAddress.StreetAddress2)
		suite.Equal(city, updatedPPM.PickupAddress.City)
		suite.Equal(state, updatedPPM.PickupAddress.State)
		suite.Equal(postalCode, updatedPPM.PickupAddress.PostalCode)

		suite.NotNil(updatedPPM.DestinationAddressID)
		suite.NotNil(updatedPPM.DestinationAddress)
		suite.Equal(streetAddress1, updatedPPM.DestinationAddress.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.DestinationAddress.StreetAddress2)
		suite.Equal(city, updatedPPM.DestinationAddress.City)
		suite.Equal(state, updatedPPM.DestinationAddress.State)
		suite.Equal(postalCode, updatedPPM.DestinationAddress.PostalCode)

		suite.NotNil(updatedPPM.SecondaryPickupAddressID)
		suite.NotNil(updatedPPM.SecondaryPickupAddress)
		suite.Equal(streetAddress1, updatedPPM.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.SecondaryPickupAddress.StreetAddress2)
		suite.Equal(city, updatedPPM.SecondaryPickupAddress.City)
		suite.Equal(state, updatedPPM.SecondaryPickupAddress.State)
		suite.Equal(postalCode, updatedPPM.SecondaryPickupAddress.PostalCode)

		suite.NotNil(updatedPPM.SecondaryDestinationAddressID)
		suite.NotNil(updatedPPM.SecondaryDestinationAddress)
		suite.Equal(streetAddress1, updatedPPM.SecondaryDestinationAddress.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.SecondaryDestinationAddress.StreetAddress2)
		suite.Equal(city, updatedPPM.SecondaryDestinationAddress.City)
		suite.Equal(state, updatedPPM.SecondaryDestinationAddress.State)
		suite.Equal(postalCode, updatedPPM.SecondaryDestinationAddress.PostalCode)

		suite.NotNil(updatedPPM.TertiaryPickupAddressID)
		suite.NotNil(updatedPPM.TertiaryPickupAddress)
		suite.Equal(streetAddress1, updatedPPM.TertiaryPickupAddress.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.TertiaryPickupAddress.StreetAddress2)
		suite.Equal(city, updatedPPM.TertiaryPickupAddress.City)
		suite.Equal(state, updatedPPM.TertiaryPickupAddress.State)
		suite.Equal(postalCode, updatedPPM.TertiaryPickupAddress.PostalCode)

		suite.NotNil(updatedPPM.TertiaryDestinationAddressID)
		suite.NotNil(updatedPPM.TertiaryDestinationAddress)
		suite.Equal(streetAddress1, updatedPPM.TertiaryDestinationAddress.StreetAddress1)
		suite.Equal(streetAddress2, *updatedPPM.TertiaryDestinationAddress.StreetAddress2)
		suite.Equal(city, updatedPPM.TertiaryDestinationAddress.City)
		suite.Equal(state, updatedPPM.TertiaryDestinationAddress.State)
		suite.Equal(postalCode, updatedPPM.TertiaryDestinationAddress.PostalCode)
	})
	suite.Run("Can successfully update a PPM Shipment SIT estimated cost", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		setupPricerData()
		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))
		newFakeSITEstimatedCost := models.CentPointer(unit.Cents(62500))

		subtestData := setUpForTests(newFakeEstimatedIncentive, newFakeSITEstimatedCost, nil, nil)
		sitLocationDestination := models.SITLocationTypeDestination
		entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)
		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "30813"
		destinationAddress := &models.Address{
			StreetAddress1: streetAddress1,
			StreetAddress2: &streetAddress2,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate:     entryDate.Add(time.Hour * 24 * 30),
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &sitLocationDestination,
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					SITEstimatedWeight:        models.PoundPointer(1000),
					SITEstimatedCost:          newFakeSITEstimatedCost,
				},
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		originalPPM.DestinationAddress = destinationAddress
		mockedPlanner := &routemocks.Planner{}
		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"90210", "30813").Return(2294, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentSITEstimatedCost(appCtx, &originalPPM)

		suite.NilOrNoVerrs(err)
		suite.NotEqual(*updatedPPM.SITEstimatedCost, newFakeSITEstimatedCost)
	})
	suite.Run("Can't find contract for PPM SIT Estimated Cost calculation", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))
		newFakeSITEstimatedCost := models.CentPointer(unit.Cents(62500))

		subtestData := setUpForTests(newFakeEstimatedIncentive, newFakeSITEstimatedCost, nil, nil)
		sitLocationDestination := models.SITLocationTypeDestination
		entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		// we do not have a contract for this date
		invalidDate := time.Date(2017, time.March, 15, 0, 0, 0, 0, time.UTC)
		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "30813"
		destinationAddress := &models.Address{
			StreetAddress1: streetAddress1,
			StreetAddress2: &streetAddress2,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate:     invalidDate,
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &sitLocationDestination,
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					SITEstimatedWeight:        models.PoundPointer(1000),
					SITEstimatedCost:          newFakeSITEstimatedCost,
				},
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		originalPPM.DestinationAddress = destinationAddress
		mockedPlanner := &routemocks.Planner{}
		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"90210", "30813").Return(2294, nil)

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentSITEstimatedCost(appCtx, &originalPPM)

		suite.Error(err)
		suite.Nil(updatedPPM)
	})
	suite.Run("Can't successfully find the PPM Shipment to update SIT estimated cost", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(2000000))
		newFakeSITEstimatedCost := models.CentPointer(unit.Cents(62500))

		subtestData := setUpForTests(newFakeEstimatedIncentive, newFakeSITEstimatedCost, nil, nil)
		sitLocationDestination := models.SITLocationTypeDestination
		entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)
		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "30813"
		destinationAddress := &models.Address{
			StreetAddress1: streetAddress1,
			StreetAddress2: &streetAddress2,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate:     entryDate.Add(time.Hour * 24 * 30),
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &sitLocationDestination,
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					SITEstimatedWeight:        models.PoundPointer(1000),
					SITEstimatedCost:          newFakeSITEstimatedCost,
				},
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		originalPPM.DestinationAddress = destinationAddress
		originalPPM.ID = uuid.Nil
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentSITEstimatedCost(appCtx, &originalPPM)

		suite.Error(err)
		suite.Nil(updatedPPM)
	})

	suite.Run("Can successfully update a PPMShipment - cap estimated incentive to max incentive value", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		newFakeEstimatedIncentive := models.CentPointer(unit.Cents(8000000))
		newFakeMaxIncentive := models.CentPointer(unit.Cents(5000000))

		subtestData := setUpForTests(newFakeEstimatedIncentive, nil, newFakeMaxIncentive, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: testdatagen.NextValidMoveDate,
					SITExpected:           models.BoolPointer(false),
					EstimatedWeight:       models.PoundPointer(4000),
					HasProGear:            models.BoolPointer(false),
					EstimatedIncentive:    fakeEstimatedIncentive,
				},
			},
		}, nil)

		newPPM := models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek),
			SITExpected:           models.BoolPointer(true),
		}

		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.NilOrNoVerrs(err)

		// Verify estimated incentive is capped at the max
		suite.Equal(*newFakeMaxIncentive, *updatedPPM.EstimatedIncentive)
	})

	suite.Run("Can update gun safe authorized when HasGunSafe value changes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					HasGunSafe: models.BoolPointer(false),
				},
			},
		}, nil)

		newPPM := models.PPMShipment{
			HasGunSafe: models.BoolPointer(true),
		}

		subtestData := setUpForTests(nil, nil, nil, nil)
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)
		suite.NotNil(updatedPPM)
		suite.NilOrNoVerrs(err)

		var updatedEntitlement models.Entitlement
		err = appCtx.DB().Find(&updatedEntitlement, originalPPM.Shipment.MoveTaskOrder.Orders.EntitlementID)
		suite.NoError(err)

		suite.True(updatedEntitlement.GunSafe)
	})

	suite.Run("Returns error if entitlement is nil when updating gun safe", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		// Create an order and link it to the entitlement
		orders := factory.BuildOrder(suite.DB(), nil, nil)

		// Manually set EntitlementID to a fake UUID to simulate broken preload
		orders.EntitlementID = nil
		orders.Entitlement = nil
		suite.MustSave(&orders)

		// Create a move and link it to the order
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    orders,
				LinkOnly: true,
			},
		}, nil)

		// Create an MTOShipment linked to the move
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
					Status:       models.MTOShipmentStatusDraft,
				},
			},
		}, nil)

		originalPPM := factory.BuildMinimalPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model: models.PPMShipment{
					HasGunSafe: models.BoolPointer(false),
				},
			},
		}, nil)

		newPPM := models.PPMShipment{
			HasGunSafe: models.BoolPointer(true),
		}

		subtestData := setUpForTests(nil, nil, nil, nil)
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)

		suite.Nil(updatedPPM)
		suite.Error(err)
		suite.IsType(apperror.QueryError{}, err)
		suite.Contains(err.Error(), "Move is missing an associated entitlement.")

	})
	suite.Run("updating PPM with valid GCC multiplier date updates PPM - expected departure date", func() {
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		originalPPM := factory.BuildPPMShipment(appCtx.DB(), nil, nil)

		newPPM := models.PPMShipment{
			ExpectedDepartureDate: validGccMultiplierDate,
			GCCMultiplierID:       nil,
		}

		subtestData := setUpForTests(nil, nil, nil, nil)
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)
		suite.NilOrNoVerrs(err)
		suite.NotNil(updatedPPM.GCCMultiplierID)
	})
	suite.Run("updating PPM with invalid GCC multiplier date updates PPM multiplier to nil - expected departure date", func() {
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
		invalidGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-04-02")
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: validGccMultiplierDate,
				},
			},
		}, nil)

		newPPM := models.PPMShipment{
			ExpectedDepartureDate: invalidGccMultiplierDate,
		}

		subtestData := setUpForTests(nil, nil, nil, nil)
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)
		suite.NilOrNoVerrs(err)
		suite.Nil(updatedPPM.GCCMultiplierID)
	})

	suite.Run("updating PPM with valid GCC multiplier date updates PPM multiplier - actual move date", func() {
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
		invalidGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-04-02")
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		// this PPM will have a 1x multiplier (nil)
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: invalidGccMultiplierDate,
					GCCMultiplierID:       nil,
				},
			},
		}, nil)

		// this should change it to 1.3x (not nil)
		newPPM := models.PPMShipment{
			ActualMoveDate:  &validGccMultiplierDate,
			GCCMultiplierID: nil,
		}

		subtestData := setUpForTests(nil, nil, nil, nil)
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)
		suite.NilOrNoVerrs(err)
		suite.NotNil(updatedPPM.GCCMultiplierID)
	})
	suite.Run("updating PPM with invalid GCC multiplier date updates PPM multiplier - actual move date", func() {
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
		invalidGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-04-02")
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		// this PPM should have a 1.3x multiplier (not nil)
		originalPPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: validGccMultiplierDate,
				},
			},
		}, nil)

		// this should change it to 1x (nil)
		newPPM := models.PPMShipment{
			ActualMoveDate:  &invalidGccMultiplierDate,
			GCCMultiplierID: nil,
		}

		subtestData := setUpForTests(nil, nil, nil, nil)
		updatedPPM, err := subtestData.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(appCtx, &newPPM, originalPPM.ShipmentID)
		suite.NilOrNoVerrs(err)
		suite.Nil(updatedPPM.GCCMultiplierID)
	})
}
