package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildPPMShipment() {
	suite.Run("Successful creation of default PPMShipment", func() {
		// Under test:      BuildPPMShipment
		// Mocked:          None
		// Set up:          Create a PPM shipment with no customizations or traits
		// Expected outcome:PPMShipment should be created with default values
		defaultPPM := models.PPMShipment{
			ExpectedDepartureDate:  time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
			SITExpected:            models.BoolPointer(false),
			Status:                 models.PPMShipmentStatusSubmitted,
			EstimatedWeight:        models.PoundPointer(unit.Pound(4000)),
			HasProGear:             models.BoolPointer(true),
			ProGearWeight:          models.PoundPointer(unit.Pound(1987)),
			SpouseProGearWeight:    models.PoundPointer(unit.Pound(498)),
			EstimatedIncentive:     models.CentPointer(unit.Cents(1000000)),
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(598700)),
			PickupAddress: &models.Address{
				StreetAddress1: "987 New Street",
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50309",
				County:         "POLK",
			},
			SecondaryPickupAddress: &models.Address{
				StreetAddress1: "123 Main Street",
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50309",
				County:         "POLK",
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "123 New Street",
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30813",
				County:         "COLUMBIA",
			},
			SecondaryDestinationAddress: &models.Address{
				StreetAddress1: "1234 Main Street",
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30813",
				County:         "COLUMBIA",
			},
		}

		// SETUP
		ppmShipment := BuildPPMShipment(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultPPM.PickupAddress.StreetAddress1, ppmShipment.PickupAddress.StreetAddress1)
		suite.Equal(defaultPPM.DestinationAddress.StreetAddress1, ppmShipment.DestinationAddress.StreetAddress1)
		suite.Equal(defaultPPM.PickupAddress.City, ppmShipment.PickupAddress.City)
		suite.Equal(defaultPPM.DestinationAddress.City, ppmShipment.DestinationAddress.City)
		suite.Equal(defaultPPM.PickupAddress.State, ppmShipment.PickupAddress.State)
		suite.Equal(defaultPPM.DestinationAddress.State, ppmShipment.DestinationAddress.State)
		suite.Equal(defaultPPM.PickupAddress.PostalCode, ppmShipment.PickupAddress.PostalCode)
		suite.Equal(defaultPPM.DestinationAddress.PostalCode, ppmShipment.DestinationAddress.PostalCode)
		suite.Equal(defaultPPM.ExpectedDepartureDate, ppmShipment.ExpectedDepartureDate)
		suite.Equal(defaultPPM.SITExpected, ppmShipment.SITExpected)
		suite.Equal(defaultPPM.Status, ppmShipment.Status)
		suite.Equal(defaultPPM.EstimatedWeight, ppmShipment.EstimatedWeight)
		suite.Equal(defaultPPM.HasProGear, ppmShipment.HasProGear)
		suite.Equal(defaultPPM.ProGearWeight, ppmShipment.ProGearWeight)
		suite.Equal(defaultPPM.SpouseProGearWeight, ppmShipment.SpouseProGearWeight)
		suite.Equal(defaultPPM.EstimatedIncentive, ppmShipment.EstimatedIncentive)
		suite.Equal(defaultPPM.HasRequestedAdvance, ppmShipment.HasRequestedAdvance)
		suite.Equal(defaultPPM.AdvanceAmountRequested, ppmShipment.AdvanceAmountRequested)
	})

	suite.Run("Successful creation of minimal PPMShipment", func() {
		// Under test:      BuildMinimalPPMShipment
		// Mocked:          None
		// Set up:          Create a Minimal PPM shipment with no customizations or traits
		// Expected outcome:PPMShipment should be created with default values
		defaultPPM := models.PPMShipment{
			Status:                models.PPMShipmentStatusDraft,
			ExpectedDepartureDate: time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
			SITExpected:           models.BoolPointer(false),
		}

		// SETUP
		ppmShipment := BuildMinimalPPMShipment(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultPPM.Status, ppmShipment.Status)
		suite.Equal(defaultPPM.ExpectedDepartureDate, ppmShipment.ExpectedDepartureDate)
		suite.Equal(defaultPPM.SITExpected, ppmShipment.SITExpected)
		suite.Nil(ppmShipment.EstimatedWeight)
		suite.Nil(ppmShipment.HasProGear)
		suite.Nil(ppmShipment.ProGearWeight)
		suite.Nil(ppmShipment.SpouseProGearWeight)
		suite.Nil(ppmShipment.EstimatedIncentive)
		suite.Nil(ppmShipment.HasRequestedAdvance)
		suite.Nil(ppmShipment.AdvanceAmountRequested)
	})

	suite.Run("Successful creation of customized PPMShipment", func() {
		// Under test:      BuildPPMShipment
		// Set up:          Create a PPM shipment and pass custom fields
		// Expected outcome:PPMShipment should be created with custom fields
		// SETUP
		sitLocation := models.SITLocationTypeDestination
		customPPM := models.PPMShipment{
			ID:                     uuid.Must(uuid.NewV4()),
			Status:                 models.PPMShipmentStatusWaitingOnCustomer,
			ExpectedDepartureDate:  time.Now(),
			HasProGear:             models.BoolPointer(true),
			ProGearWeight:          models.PoundPointer(unit.Pound(1989)),
			EstimatedWeight:        models.PoundPointer(unit.Pound(3000)),
			SpouseProGearWeight:    models.PoundPointer(unit.Pound(123)),
			EstimatedIncentive:     models.CentPointer(unit.Cents(1005000)),
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(600000)),
			SITExpected:            models.BoolPointer(true),
			SITLocation:            &sitLocation,
			SITEstimatedWeight:     models.PoundPointer(unit.Pound(2000)),
		}
		customAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}

		// CALL FUNCTION UNDER TEST
		ppmShipment := BuildPPMShipment(suite.DB(), []Customization{
			{Model: customPPM},
			{
				Model: customAddress,
				Type:  &Addresses.W2Address,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customPPM.ExpectedDepartureDate, ppmShipment.ExpectedDepartureDate)
		suite.Equal(customPPM.SITExpected, ppmShipment.SITExpected)
		suite.Equal(customPPM.Status, ppmShipment.Status)
		suite.Equal(customPPM.EstimatedWeight, ppmShipment.EstimatedWeight)
		suite.Equal(customPPM.HasProGear, ppmShipment.HasProGear)
		suite.Equal(customPPM.ProGearWeight, ppmShipment.ProGearWeight)
		suite.Equal(customPPM.SpouseProGearWeight, ppmShipment.SpouseProGearWeight)
		suite.Equal(customPPM.EstimatedIncentive, ppmShipment.EstimatedIncentive)
		suite.Equal(customPPM.HasRequestedAdvance, ppmShipment.HasRequestedAdvance)
		suite.Equal(customPPM.AdvanceAmountRequested, ppmShipment.AdvanceAmountRequested)
		// Check that the address and phoneline were customized
		suite.Equal(customAddress.StreetAddress1, ppmShipment.W2Address.StreetAddress1)
	})

	suite.Run("Successful creation of PPMShipment with trait", func() {
		// Under test:       BuildPPMShipment
		// Set up:           Pass in a trait
		// Expected outcome: PPMShipment should be created.

		expectedPPM := models.PPMShipment{
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualPickupPostalCode:      models.StringPointer("30813"),
			ActualDestinationPostalCode: models.StringPointer("50309"),
			HasRequestedAdvance:         models.BoolPointer(true),
			AdvanceAmountRequested:      models.CentPointer(unit.Cents(598700)),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(598700)),
		}

		w2Address := models.Address{
			StreetAddress1: "987 New Street",
		}

		ppmShipment := BuildPPMShipment(suite.DB(), nil, []Trait{GetTraitApprovedPPMWithActualInfo})

		suite.Equal(expectedPPM.Status, ppmShipment.Status)
		suite.Equal(expectedPPM.ActualPickupPostalCode, ppmShipment.ActualPickupPostalCode)
		suite.Equal(expectedPPM.ActualDestinationPostalCode, ppmShipment.ActualDestinationPostalCode)
		suite.Equal(expectedPPM.HasRequestedAdvance, ppmShipment.HasRequestedAdvance)
		suite.Equal(expectedPPM.AdvanceAmountRequested, ppmShipment.AdvanceAmountRequested)
		suite.Equal(expectedPPM.HasReceivedAdvance, ppmShipment.HasReceivedAdvance)
		suite.Equal(expectedPPM.AdvanceAmountReceived, ppmShipment.AdvanceAmountReceived)
		suite.NotNil(ppmShipment.SubmittedAt)
		suite.NotNil(ppmShipment.ApprovedAt)
		suite.NotNil(ppmShipment.ExpectedDepartureDate)
		suite.NotNil(ppmShipment.ActualMoveDate)
		suite.Equal(w2Address.StreetAddress1, ppmShipment.W2Address.StreetAddress1)
		suite.Equal(models.MoveStatusAPPROVED, ppmShipment.Shipment.MoveTaskOrder.Status)
		suite.Equal(models.MTOShipmentStatusApproved, ppmShipment.Shipment.Status)
	})

	suite.Run("Successful return of linkOnly PPMShipment", func() {
		// Under test:       BuildPPMShipment
		// Set up:           Pass in a linkOnly PPMShipment
		// Expected outcome: No new PPMShipment should be created.

		// Check num PPMShipment records
		precount, err := suite.DB().Count(&models.PPMShipment{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		ppmShipment := BuildPPMShipment(suite.DB(), []Customization{
			{
				Model: models.PPMShipment{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.PPMShipment{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, ppmShipment.ID)
	})

	suite.Run("Successful return of stubbed PPMShipment", func() {
		// Under test:       BuildPPMShipment
		// Set up:           Pass in a linkOnly PPMShipment
		// Expected outcome: No new PPMShipment should be created.

		// Check num PPMShipment records
		precount, err := suite.DB().Count(&models.PPMShipment{})
		suite.NoError(err)

		// Nil passed in as db
		ppmShipment := BuildPPMShipment(nil, []Customization{
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusWaitingOnCustomer,
				},
			},
		}, nil)
		count, err := suite.DB().Count(&models.PPMShipment{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status)

	})

	suite.Run("PPM Shipment ready for final customer CloseOut", func() {
		// Under test:       BuildPPMShipmentReadyForFinalCustomerCloseOut
		// Set up:           build without custom user uploader
		// Expected outcome: New PPMShipment should be created with
		// Weight Ticket

		ppmShipment := BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		suite.NotNil(ppmShipment.ActualPickupPostalCode)
		suite.NotNil(ppmShipment.ActualDestinationPostalCode)
		suite.NotNil(ppmShipment.AOAPacket)
		suite.NotNil(ppmShipment.AOAPacketID)
		suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status)

		suite.NotEmpty(ppmShipment.WeightTickets)
		suite.Equal(1, len(ppmShipment.WeightTickets))
		suite.False(ppmShipment.WeightTickets[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.WeightTickets[0].PPMShipmentID)
	})

	suite.Run("Stubbed PPM Shipment ready for final customer CloseOut", func() {
		// Under test:       BuildPPMShipmentReadyForFinalCustomerCloseOut
		// Set up:           build without db
		// Expected outcome: New PPMShipment should be created with
		// Weight Ticket

		ppmShipment := BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)

		suite.False(ppmShipment.ID.IsNil())
		suite.NotNil(ppmShipment.ActualPickupPostalCode)
		suite.NotNil(ppmShipment.ActualDestinationPostalCode)
		suite.NotNil(ppmShipment.AOAPacket)
		suite.NotNil(ppmShipment.AOAPacketID)
		suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status)

		suite.NotEmpty(ppmShipment.WeightTickets)
		suite.Equal(1, len(ppmShipment.WeightTickets))
		suite.Equal(ppmShipment.ID, ppmShipment.WeightTickets[0].PPMShipmentID)
	})

	suite.Run("PPM Shipment ready for final customer CloseOut with all doc types", func() {
		// Under test:       BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes
		// Set up:           build without custom user uploader
		// Expected outcome: New PPMShipment should be created with
		// Weight Ticket, Progear Weight Ticket, and Moving Expense

		ppmShipment := BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(suite.DB(), nil)

		suite.NotNil(ppmShipment.ActualPickupPostalCode)
		suite.NotNil(ppmShipment.ActualDestinationPostalCode)
		suite.NotNil(ppmShipment.AOAPacket)
		suite.NotNil(ppmShipment.AOAPacketID)
		suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status)

		suite.NotEmpty(ppmShipment.WeightTickets)
		suite.Equal(1, len(ppmShipment.WeightTickets))
		suite.False(ppmShipment.WeightTickets[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.WeightTickets[0].PPMShipmentID)

		suite.NotEmpty(ppmShipment.ProgearWeightTickets)
		suite.Equal(1, len(ppmShipment.ProgearWeightTickets))
		suite.False(ppmShipment.ProgearWeightTickets[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.ProgearWeightTickets[0].PPMShipmentID)

		suite.NotEmpty(ppmShipment.MovingExpenses)
		suite.Equal(1, len(ppmShipment.MovingExpenses))
		suite.False(ppmShipment.MovingExpenses[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.MovingExpenses[0].PPMShipmentID)
	})

	suite.Run("PPM Shipment that needs approval with all doc types", func() {
		// Under test:       BuildPPMShipmentThatNeedsCloseoutWithAllDocTypes
		// Set up:           build without custom user uploader
		// Expected outcome: New PPMShipment should be created with
		// Weight Ticket, Progear Weight Ticket, and Moving Expense

		ppmShipment := BuildPPMShipmentThatNeedsCloseoutWithAllDocTypes(suite.DB(), nil)

		suite.NotNil(ppmShipment.ActualPickupPostalCode)
		suite.NotNil(ppmShipment.ActualDestinationPostalCode)
		suite.NotNil(ppmShipment.AOAPacket)
		suite.NotNil(ppmShipment.AOAPacketID)
		suite.Equal(models.PPMShipmentStatusNeedsCloseout, ppmShipment.Status)

		suite.NotEmpty(ppmShipment.WeightTickets)
		suite.Equal(1, len(ppmShipment.WeightTickets))
		suite.False(ppmShipment.WeightTickets[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.WeightTickets[0].PPMShipmentID)

		suite.NotEmpty(ppmShipment.ProgearWeightTickets)
		suite.Equal(1, len(ppmShipment.ProgearWeightTickets))
		suite.False(ppmShipment.ProgearWeightTickets[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.ProgearWeightTickets[0].PPMShipmentID)

		suite.NotEmpty(ppmShipment.MovingExpenses)
		suite.Equal(1, len(ppmShipment.MovingExpenses))
		suite.False(ppmShipment.MovingExpenses[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.MovingExpenses[0].PPMShipmentID)
	})

	suite.Run("PPM Shipment that is missing payment packet", func() {
		// Under test:       BuildPPMShipmentThatNeedsCloseoutWithAllDocTypes
		// Set up:           build without custom user uploader
		// Expected outcome: New PPMShipment should be created with
		// Weight Ticket, Progear Weight Ticket, and Moving Expense

		ppmShipment := BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(suite.DB(), nil, nil)

		suite.NotNil(ppmShipment.ActualPickupPostalCode)
		suite.NotNil(ppmShipment.ActualDestinationPostalCode)
		suite.NotNil(ppmShipment.AOAPacket)
		suite.NotNil(ppmShipment.AOAPacketID)
		suite.Equal(models.PPMShipmentStatusCloseoutComplete, ppmShipment.Status)

		suite.NotEmpty(ppmShipment.WeightTickets)
		suite.Equal(1, len(ppmShipment.WeightTickets))
		suite.False(ppmShipment.WeightTickets[0].ID.IsNil())
		suite.Equal(ppmShipment.ID, ppmShipment.WeightTickets[0].PPMShipmentID)

		suite.Empty(ppmShipment.ProgearWeightTickets)
		suite.Empty(ppmShipment.MovingExpenses)
	})
}
