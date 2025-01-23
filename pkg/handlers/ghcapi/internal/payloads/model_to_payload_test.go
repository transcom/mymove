package payloads

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/unit"
)

func TestOrder(_ *testing.T) {
	order := &models.Order{}
	Order(order)
}

func (suite *PayloadsSuite) TestOrderWithMove() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	moves := models.Moves{}
	moves = append(moves, move)
	order := factory.BuildOrder(nil, []factory.Customization{
		{
			Model: models.Order{
				ID:            uuid.Must(uuid.NewV4()),
				HasDependents: *models.BoolPointer(true),
				Moves:         moves,
			},
		},
	}, nil)
	Order(&order)
}

func (suite *PayloadsSuite) TestBoatShipment() {
	boat := factory.BuildBoatShipment(suite.DB(), nil, nil)
	boatShipment := BoatShipment(nil, &boat)
	suite.NotNil(boatShipment)

}

func (suite *PayloadsSuite) TestMobileHomeShipment() {
	mobileHome := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)
	mobileHomeShipment := MobileHomeShipment(nil, &mobileHome)
	suite.NotNil(mobileHomeShipment)
}

func (suite *PayloadsSuite) TestMovingExpense() {
	contractExpense := models.MovingExpenseReceiptTypeContractedExpense
	weightStored := 2000
	sitLocation := models.SITLocationTypeDestination
	sitReimburseableAmount := 500

	movingExpense := models.MovingExpense{
		PPMShipmentID:          uuid.Must(uuid.NewV4()),
		DocumentID:             uuid.Must(uuid.NewV4()),
		MovingExpenseType:      &contractExpense,
		Reason:                 models.StringPointer("no good"),
		SITStartDate:           models.TimePointer(time.Now()),
		SITEndDate:             models.TimePointer(time.Now()),
		WeightStored:           (*unit.Pound)(&weightStored),
		SITLocation:            &sitLocation,
		SITReimburseableAmount: (*unit.Cents)(&sitReimburseableAmount),
	}
	movingExpenseValues := MovingExpense(nil, &movingExpense)
	suite.NotNil(movingExpenseValues)
}

func (suite *PayloadsSuite) TestMovingExpenses() {
	contractExpense := models.MovingExpenseReceiptTypeContractedExpense
	weightStored := 2000
	sitLocation := models.SITLocationTypeDestination
	sitReimburseableAmount := 500
	movingExpenses := models.MovingExpenses{}

	movingExpense := models.MovingExpense{
		PPMShipmentID:          uuid.Must(uuid.NewV4()),
		DocumentID:             uuid.Must(uuid.NewV4()),
		MovingExpenseType:      &contractExpense,
		Reason:                 models.StringPointer("no good"),
		SITStartDate:           models.TimePointer(time.Now()),
		SITEndDate:             models.TimePointer(time.Now()),
		WeightStored:           (*unit.Pound)(&weightStored),
		SITLocation:            &sitLocation,
		SITReimburseableAmount: (*unit.Cents)(&sitReimburseableAmount),
	}
	movingExpenseTwo := models.MovingExpense{
		PPMShipmentID:          uuid.Must(uuid.NewV4()),
		DocumentID:             uuid.Must(uuid.NewV4()),
		MovingExpenseType:      &contractExpense,
		Reason:                 models.StringPointer("no good"),
		SITStartDate:           models.TimePointer(time.Now()),
		SITEndDate:             models.TimePointer(time.Now()),
		WeightStored:           (*unit.Pound)(&weightStored),
		SITLocation:            &sitLocation,
		SITReimburseableAmount: (*unit.Cents)(&sitReimburseableAmount),
	}
	movingExpenses = append(movingExpenses, movingExpense, movingExpenseTwo)
	movingExpensesValue := MovingExpenses(nil, movingExpenses)
	suite.NotNil(movingExpensesValue)
}

// TestMove makes sure zero values/optional fields are handled
func TestMove(t *testing.T) {
	_, err := Move(&models.Move{}, &test.FakeS3Storage{})
	if err != nil {
		t.Fail()
	}
}

func (suite *PayloadsSuite) TestExcessWeightInMovePayload() {
	now := time.Now()

	suite.Run("successfully converts excess weight in model to payload", func() {
		move := models.Move{

			ExcessWeightQualifiedAt:                        &now,
			ExcessUnaccompaniedBaggageWeightQualifiedAt:    &now,
			ExcessUnaccompaniedBaggageWeightAcknowledgedAt: &now,
			ExcessWeightAcknowledgedAt:                     &now,
		}

		payload, err := Move(&move, &test.FakeS3Storage{})
		suite.NoError(err)
		suite.Equal(handlers.FmtDateTimePtr(move.ExcessWeightQualifiedAt), payload.ExcessWeightQualifiedAt)
		suite.Equal(handlers.FmtDateTimePtr(move.ExcessUnaccompaniedBaggageWeightQualifiedAt), payload.ExcessUnaccompaniedBaggageWeightQualifiedAt)
		suite.Equal(handlers.FmtDateTimePtr(move.ExcessUnaccompaniedBaggageWeightAcknowledgedAt), payload.ExcessUnaccompaniedBaggageWeightAcknowledgedAt)
		suite.Equal(handlers.FmtDateTimePtr(move.ExcessWeightAcknowledgedAt), payload.ExcessWeightAcknowledgedAt)
	})
}

func (suite *PayloadsSuite) TestPaymentRequestQueue() {
	officeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Email: "officeuser1@example.com",
			},
		},
		{
			Model: models.User{
				Privileges: []models.Privilege{
					{
						PrivilegeType: models.PrivilegeTypeSupervisor,
					},
				},
				Roles: []roles.Role{
					{
						RoleType: roles.RoleTypeTIO,
					},
				},
			},
		},
	}, nil)
	officeUserTIO := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})

	gbloc := "LKNQ"

	approvedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	approvedMove.ShipmentGBLOC = append(approvedMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

	pr2 := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model:    approvedMove,
			LinkOnly: true,
		},
		{
			Model: models.TransportationOffice{
				Gbloc: "LKNQ",
			},
			Type: &factory.TransportationOffices.OriginDutyLocation,
		},
		{
			Model: models.DutyLocation{
				Name: "KJKJKJKJKJK",
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
	}, nil)

	paymentRequests := models.PaymentRequests{pr2}
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "PPSO",
				ProvidesCloseout: true,
			},
		},
	}, nil)
	var officeUsers models.OfficeUsers
	var officeUsersSafety models.OfficeUsers
	officeUsers = append(officeUsers, officeUser)
	var paymentRequestsQueue = QueuePaymentRequests(&paymentRequests, officeUsers, officeUser, officeUsersSafety)

	suite.Run("Test Payment request is assignable due to not being assigend", func() {
		paymentRequestCopy := *paymentRequestsQueue
		suite.NotNil(paymentRequestsQueue)
		suite.IsType(paymentRequestsQueue, &ghcmessages.QueuePaymentRequests{})
		suite.Nil(paymentRequestCopy[0].AssignedTo)
	})

	suite.Run("Test Payment request has no counseling office", func() {
		paymentRequestCopy := *paymentRequestsQueue
		suite.NotNil(paymentRequestsQueue)
		suite.IsType(paymentRequestsQueue, &ghcmessages.QueuePaymentRequests{})
		suite.Nil(paymentRequestCopy[0].CounselingOffice)
	})

	paymentRequests[0].MoveTaskOrder.TIOAssignedUser = &officeUserTIO
	paymentRequests[0].MoveTaskOrder.CounselingOffice = &transportationOffice

	paymentRequestsQueue = QueuePaymentRequests(&paymentRequests, officeUsers, officeUser, officeUsersSafety)

	suite.Run("Test PaymentRequest has both Counseling Office and TIO AssignedUser ", func() {
		PaymentRequestsCopy := *paymentRequestsQueue

		suite.NotNil(PaymentRequests)
		suite.IsType(&ghcmessages.QueuePaymentRequests{}, paymentRequestsQueue)
		suite.IsType(&ghcmessages.QueuePaymentRequest{}, PaymentRequestsCopy[0])
		suite.Equal(PaymentRequestsCopy[0].AssignedTo.FirstName, officeUserTIO.FirstName)
		suite.Equal(PaymentRequestsCopy[0].AssignedTo.LastName, officeUserTIO.LastName)
		suite.Equal(*PaymentRequestsCopy[0].CounselingOffice, transportationOffice.Name)
	})

	suite.Run("Test PaymentRequest is assignable due to user Supervisor role", func() {
		paymentRequests := QueuePaymentRequests(&paymentRequests, officeUsers, officeUser, officeUsersSafety)
		paymentRequestCopy := *paymentRequests
		suite.Equal(paymentRequestCopy[0].Assignable, true)
	})

	officeUserHQ := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeHQ})
	suite.Run("Test PaymentRequest is not assignable due to user HQ role", func() {
		paymentRequests := QueuePaymentRequests(&paymentRequests, officeUsers, officeUserHQ, officeUsersSafety)
		paymentRequestCopy := *paymentRequests
		suite.Equal(paymentRequestCopy[0].Assignable, false)
	})
}

func (suite *PayloadsSuite) TestFetchPPMShipment() {

	ppmShipmentID, _ := uuid.NewV4()
	streetAddress1 := "MacDill AFB"
	streetAddress2, streetAddress3 := "", ""
	city := "Tampa"
	state := "FL"
	postalcode := "33621"
	county := "HILLSBOROUGH"

	country := models.Country{
		Country: "US",
	}

	expectedAddress := models.Address{
		StreetAddress1: streetAddress1,
		StreetAddress2: &streetAddress2,
		StreetAddress3: &streetAddress3,
		City:           city,
		State:          state,
		PostalCode:     postalcode,
		Country:        &country,
		County:         &county,
	}

	isActualExpenseReimbursement := true

	expectedPPMShipment := models.PPMShipment{
		ID:                           ppmShipmentID,
		PickupAddress:                &expectedAddress,
		DestinationAddress:           &expectedAddress,
		IsActualExpenseReimbursement: &isActualExpenseReimbursement,
	}

	suite.Run("Success -", func() {
		returnedPPMShipment := PPMShipment(nil, &expectedPPMShipment)

		suite.IsType(returnedPPMShipment, &ghcmessages.PPMShipment{})
		suite.Equal(&streetAddress1, returnedPPMShipment.PickupAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress2, returnedPPMShipment.PickupAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress3, returnedPPMShipment.PickupAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.PickupAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.PickupAddress.City)
		suite.Equal(&state, returnedPPMShipment.PickupAddress.State)
		suite.Equal(&country.Country, returnedPPMShipment.PickupAddress.Country)
		suite.Equal(&county, returnedPPMShipment.PickupAddress.County)

		suite.Equal(&streetAddress1, returnedPPMShipment.DestinationAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress2, returnedPPMShipment.DestinationAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress3, returnedPPMShipment.DestinationAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.DestinationAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.DestinationAddress.City)
		suite.Equal(&state, returnedPPMShipment.DestinationAddress.State)
		suite.Equal(&country.Country, returnedPPMShipment.DestinationAddress.Country)
		suite.Equal(&county, returnedPPMShipment.DestinationAddress.County)
		suite.True(*returnedPPMShipment.IsActualExpenseReimbursement)
	})

	suite.Run("Destination street address 1 returns empty string to convey OPTIONAL state ", func() {
		expected_street_address_1 := ""
		expectedAddress2 := models.Address{
			StreetAddress1: expected_street_address_1,
			StreetAddress2: &streetAddress2,
			StreetAddress3: &streetAddress3,
			City:           city,
			State:          state,
			PostalCode:     postalcode,
			Country:        &country,
			County:         &county,
		}

		expectedPPMShipment2 := models.PPMShipment{
			ID:                 ppmShipmentID,
			PickupAddress:      &expectedAddress,
			DestinationAddress: &expectedAddress2,
		}
		returnedPPMShipment := PPMShipment(nil, &expectedPPMShipment2)

		suite.IsType(returnedPPMShipment, &ghcmessages.PPMShipment{})
		suite.Equal(&streetAddress1, returnedPPMShipment.PickupAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress2, returnedPPMShipment.PickupAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress3, returnedPPMShipment.PickupAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.PickupAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.PickupAddress.City)
		suite.Equal(&state, returnedPPMShipment.PickupAddress.State)
		suite.Equal(&county, returnedPPMShipment.PickupAddress.County)

		suite.Equal(&expected_street_address_1, returnedPPMShipment.DestinationAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress2, returnedPPMShipment.DestinationAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress3, returnedPPMShipment.DestinationAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.DestinationAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.DestinationAddress.City)
		suite.Equal(&state, returnedPPMShipment.DestinationAddress.State)
		suite.Equal(&county, returnedPPMShipment.DestinationAddress.County)
	})
}

func (suite *PayloadsSuite) TestUpload() {
	uploadID, _ := uuid.NewV4()
	testURL := "https://testurl.com"

	basicUpload := models.Upload{
		ID:          uploadID,
		Filename:    "fileName",
		ContentType: "image/png",
		Bytes:       1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedUpload := Upload(suite.storer, basicUpload, testURL)

		suite.IsType(returnedUpload, &ghcmessages.Upload{})
		expectedID := handlers.FmtUUIDValue(basicUpload.ID)
		suite.Equal(expectedID, returnedUpload.ID)
		suite.Equal(basicUpload.Filename, returnedUpload.Filename)
		suite.Equal(basicUpload.ContentType, returnedUpload.ContentType)
		suite.Equal(basicUpload.Bytes, returnedUpload.Bytes)
		suite.Equal(testURL, returnedUpload.URL.String())
	})
}

func (suite *PayloadsSuite) TestShipmentAddressUpdate() {
	id, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()

	newAddress := models.Address{
		StreetAddress1: "123 New St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89503",
		County:         models.StringPointer("WASHOE"),
	}

	oldAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         models.StringPointer("WASHOE"),
	}

	sitOriginalAddress := models.Address{
		StreetAddress1: "123 SIT St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89501",
		County:         models.StringPointer("WASHOE"),
	}
	officeRemarks := "some office remarks"
	newSitDistanceBetween := 0
	oldSitDistanceBetween := 0

	shipmentAddressUpdate := models.ShipmentAddressUpdate{
		ID:                    id,
		ShipmentID:            id2,
		NewAddress:            newAddress,
		OriginalAddress:       oldAddress,
		SitOriginalAddress:    &sitOriginalAddress,
		ContractorRemarks:     "some remarks",
		OfficeRemarks:         &officeRemarks,
		Status:                models.ShipmentAddressUpdateStatusRequested,
		NewSitDistanceBetween: &newSitDistanceBetween,
		OldSitDistanceBetween: &oldSitDistanceBetween,
	}

	emptyShipmentAddressUpdate := models.ShipmentAddressUpdate{ID: uuid.Nil}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedShipmentAddressUpdate := ShipmentAddressUpdate(&shipmentAddressUpdate)

		suite.IsType(returnedShipmentAddressUpdate, &ghcmessages.ShipmentAddressUpdate{})
	})
	suite.Run("Failure - Returns nil", func() {
		returnedShipmentAddressUpdate := ShipmentAddressUpdate(&emptyShipmentAddressUpdate)

		suite.Nil(returnedShipmentAddressUpdate)
	})
}

func (suite *PayloadsSuite) TestMoveWithGBLOC() {
	defaultOrdersNumber := "ORDER3"
	defaultTACNumber := "F8E1"
	defaultDepartmentIndicator := "AIR_AND_SPACE_FORCE"
	defaultGrade := "E_1"
	defaultHasDependents := false
	defaultSpouseHasProGear := false
	defaultOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	defaultOrdersTypeDetail := internalmessages.OrdersTypeDetail("HHG_PERMITTED")
	defaultStatus := models.OrderStatusDRAFT
	testYear := 2018
	defaultIssueDate := time.Date(testYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	defaultReportByDate := time.Date(testYear, time.August, 1, 0, 0, 0, 0, time.UTC)
	defaultGBLOC := "KKFA"

	originDutyLocation := models.DutyLocation{
		Name: "Custom Origin",
	}
	originDutyLocationTOName := "origin duty location transportation office"
	firstName := "customFirst"
	lastName := "customLast"
	serviceMember := models.ServiceMember{
		FirstName: &firstName,
		LastName:  &lastName,
	}
	uploadedOrders := models.Document{
		ID: uuid.Must(uuid.NewV4()),
	}
	dependents := 7
	entitlement := models.Entitlement{
		TotalDependents: &dependents,
	}
	amendedOrders := models.Document{
		ID: uuid.Must(uuid.NewV4()),
	}
	// Create order
	order := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model: originDutyLocation,
			Type:  &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.TransportationOffice{
				Name: originDutyLocationTOName,
			},
			Type: &factory.TransportationOffices.OriginDutyLocation,
		},
		{
			Model: serviceMember,
		},
		{
			Model: uploadedOrders,
			Type:  &factory.Documents.UploadedOrders,
		},
		{
			Model: entitlement,
		},
		{
			Model: amendedOrders,
			Type:  &factory.Documents.UploadedAmendedOrders,
		},
	}, nil)

	suite.Equal(defaultOrdersNumber, *order.OrdersNumber)
	suite.Equal(defaultTACNumber, *order.TAC)
	suite.Equal(defaultDepartmentIndicator, *order.DepartmentIndicator)
	suite.Equal(defaultGrade, string(*order.Grade))
	suite.Equal(defaultHasDependents, order.HasDependents)
	suite.Equal(defaultSpouseHasProGear, order.SpouseHasProGear)
	suite.Equal(defaultOrdersType, order.OrdersType)
	suite.Equal(defaultOrdersTypeDetail, *order.OrdersTypeDetail)
	suite.Equal(defaultStatus, order.Status)
	suite.Equal(defaultIssueDate, order.IssueDate)
	suite.Equal(defaultReportByDate, order.ReportByDate)
	suite.Equal(defaultGBLOC, *order.OriginDutyLocationGBLOC)

	suite.Equal(originDutyLocation.Name, order.OriginDutyLocation.Name)
	suite.Equal(originDutyLocationTOName, order.OriginDutyLocation.TransportationOffice.Name)
	suite.Equal(*serviceMember.FirstName, *order.ServiceMember.FirstName)
	suite.Equal(*serviceMember.LastName, *order.ServiceMember.LastName)
	suite.Equal(uploadedOrders.ID, order.UploadedOrdersID)
	suite.Equal(uploadedOrders.ID, order.UploadedOrders.ID)
	suite.Equal(*entitlement.TotalDependents, *order.Entitlement.TotalDependents)
	suite.Equal(amendedOrders.ID, *order.UploadedAmendedOrdersID)
	suite.Equal(amendedOrders.ID, order.UploadedAmendedOrders.ID)
}

func (suite *PayloadsSuite) TestWeightTicketUpload() {
	uploadID, _ := uuid.NewV4()
	testURL := "https://testurl.com"
	isWeightTicket := true

	basicUpload := models.Upload{
		ID:          uploadID,
		Filename:    "fileName",
		ContentType: "image/png",
		Bytes:       1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedUpload := WeightTicketUpload(suite.storer, basicUpload, testURL, isWeightTicket)

		suite.IsType(returnedUpload, &ghcmessages.Upload{})
		expectedID := handlers.FmtUUIDValue(basicUpload.ID)
		suite.Equal(expectedID, returnedUpload.ID)
		suite.Equal(basicUpload.Filename, returnedUpload.Filename)
		suite.Equal(basicUpload.ContentType, returnedUpload.ContentType)
		suite.Equal(basicUpload.Bytes, returnedUpload.Bytes)
		suite.Equal(testURL, returnedUpload.URL.String())
		suite.Equal(isWeightTicket, returnedUpload.IsWeightTicket)
	})
}

func (suite *PayloadsSuite) TestProofOfServiceDoc() {
	uploadID1, _ := uuid.NewV4()
	uploadID2, _ := uuid.NewV4()
	isWeightTicket := true

	// Create sample ProofOfServiceDoc
	proofOfServiceDoc := models.ProofOfServiceDoc{
		ID:               uuid.Must(uuid.NewV4()),
		PaymentRequestID: uuid.Must(uuid.NewV4()),
		IsWeightTicket:   isWeightTicket,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Create sample PrimeUploads
	primeUpload1 := models.PrimeUpload{
		ID:                  uuid.Must(uuid.NewV4()),
		ProofOfServiceDocID: uuid.Must(uuid.NewV4()),
		ContractorID:        uuid.Must(uuid.NewV4()),
		UploadID:            uploadID1,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	primeUpload2 := models.PrimeUpload{
		ID:                  uuid.Must(uuid.NewV4()),
		ProofOfServiceDocID: uuid.Must(uuid.NewV4()),
		ContractorID:        uuid.Must(uuid.NewV4()),
		UploadID:            uploadID2,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	proofOfServiceDoc.PrimeUploads = []models.PrimeUpload{primeUpload1, primeUpload2}

	suite.Run("Success - Returns a ghcmessages Proof of Service payload from a Struct", func() {
		returnedProofOfServiceDoc, _ := ProofOfServiceDoc(proofOfServiceDoc, suite.storer)

		suite.IsType(returnedProofOfServiceDoc, &ghcmessages.ProofOfServiceDoc{})
	})
}

func (suite *PayloadsSuite) TestCustomer() {
	id, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()

	residentialAddress := models.Address{
		StreetAddress1: "123 New St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89503",
		County:         models.StringPointer("WASHOE"),
	}

	backupAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         models.StringPointer("WASHOE"),
	}

	phone := "444-555-6677"

	firstName := "First"
	lastName := "Last"
	affiliation := models.AffiliationARMY
	email := "dontEmailMe@gmail.com"
	cacValidated := true
	customer := models.ServiceMember{
		ID:                   id,
		UserID:               id2,
		FirstName:            &firstName,
		LastName:             &lastName,
		Affiliation:          &affiliation,
		PersonalEmail:        &email,
		Telephone:            &phone,
		ResidentialAddress:   &residentialAddress,
		BackupMailingAddress: &backupAddress,
		CacValidated:         cacValidated,
	}

	suite.Run("Success - Returns a ghcmessages Customer payload from Customer Struct", func() {
		customer := Customer(&customer)

		suite.IsType(customer, &ghcmessages.Customer{})
	})
}

func (suite *PayloadsSuite) TestEntitlement() {
	entitlementID, _ := uuid.NewV4()
	dependentsAuthorized := true
	nonTemporaryStorage := true
	privatelyOwnedVehicle := true
	proGearWeight := 1000
	proGearWeightSpouse := 500
	storageInTransit := 90
	totalDependents := 2
	requiredMedicalEquipmentWeight := 200
	accompaniedTour := true
	dependentsUnderTwelve := 1
	dependentsTwelveAndOver := 1
	authorizedWeight := 8000
	ubAllowance := 300

	entitlement := &models.Entitlement{
		ID:                             entitlementID,
		DBAuthorizedWeight:             &authorizedWeight,
		DependentsAuthorized:           &dependentsAuthorized,
		NonTemporaryStorage:            &nonTemporaryStorage,
		PrivatelyOwnedVehicle:          &privatelyOwnedVehicle,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		StorageInTransit:               &storageInTransit,
		TotalDependents:                &totalDependents,
		RequiredMedicalEquipmentWeight: requiredMedicalEquipmentWeight,
		AccompaniedTour:                &accompaniedTour,
		DependentsUnderTwelve:          &dependentsUnderTwelve,
		DependentsTwelveAndOver:        &dependentsTwelveAndOver,
		UpdatedAt:                      time.Now(),
		UBAllowance:                    &ubAllowance,
	}

	returnedEntitlement := Entitlement(entitlement)
	returnedUBAllowance := entitlement.UBAllowance

	suite.IsType(&ghcmessages.Entitlements{}, returnedEntitlement)

	suite.Equal(strfmt.UUID(entitlementID.String()), returnedEntitlement.ID)
	suite.Equal(authorizedWeight, int(*returnedEntitlement.AuthorizedWeight))
	suite.Equal(entitlement.DependentsAuthorized, returnedEntitlement.DependentsAuthorized)
	suite.Equal(entitlement.NonTemporaryStorage, returnedEntitlement.NonTemporaryStorage)
	suite.Equal(entitlement.PrivatelyOwnedVehicle, returnedEntitlement.PrivatelyOwnedVehicle)
	suite.Equal(int(*returnedUBAllowance), int(*returnedEntitlement.UnaccompaniedBaggageAllowance))
	suite.Equal(int64(proGearWeight), returnedEntitlement.ProGearWeight)
	suite.Equal(int64(proGearWeightSpouse), returnedEntitlement.ProGearWeightSpouse)
	suite.Equal(storageInTransit, int(*returnedEntitlement.StorageInTransit))
	suite.Equal(totalDependents, int(returnedEntitlement.TotalDependents))
	suite.Equal(int64(requiredMedicalEquipmentWeight), returnedEntitlement.RequiredMedicalEquipmentWeight)
	suite.Equal(models.BoolPointer(accompaniedTour), returnedEntitlement.AccompaniedTour)
	suite.Equal(dependentsUnderTwelve, int(*returnedEntitlement.DependentsUnderTwelve))
	suite.Equal(dependentsTwelveAndOver, int(*returnedEntitlement.DependentsTwelveAndOver))
}

func (suite *PayloadsSuite) TestCreateCustomer() {
	id, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	oktaID := "thisIsNotARealID"

	oktaUser := models.CreatedOktaUser{
		ID: oktaID,
		Profile: struct {
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			MobilePhone string `json:"mobilePhone"`
			SecondEmail string `json:"secondEmail"`
			Login       string `json:"login"`
			Email       string `json:"email"`
		}{
			Email: "john.doe@example.com",
		},
	}

	residentialAddress := models.Address{
		StreetAddress1: "123 New St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89503",
		County:         models.StringPointer("WASHOE"),
	}

	backupAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         models.StringPointer("WASHOE"),
	}

	phone := "444-555-6677"
	backupContact := models.BackupContact{
		Name:  "Billy Bob",
		Email: "billBob@mail.mil",
		Phone: &phone,
	}

	firstName := "First"
	lastName := "Last"
	affiliation := models.AffiliationARMY
	email := "dontEmailMe@gmail.com"
	sm := models.ServiceMember{
		ID:                   id,
		UserID:               id2,
		FirstName:            &firstName,
		LastName:             &lastName,
		Affiliation:          &affiliation,
		PersonalEmail:        &email,
		Telephone:            &phone,
		ResidentialAddress:   &residentialAddress,
		BackupMailingAddress: &backupAddress,
	}

	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct", func() {
		returnedShipmentAddressUpdate := CreatedCustomer(&sm, &oktaUser, &backupContact)

		suite.IsType(returnedShipmentAddressUpdate, &ghcmessages.CreatedCustomer{})
	})
}

func (suite *PayloadsSuite) TestMoveTaskOrder() {
	move := factory.BuildMove(suite.DB(), nil, nil)
	moveTaskOrder := MoveTaskOrder(&move)
	suite.NotNil(moveTaskOrder)
}

func (suite *PayloadsSuite) TestTransportationOffice() {
	office := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				ID: uuid.Must(uuid.NewV4()),
			},
		}}, nil)
	transportationOffice := TransportationOffice(&office)
	suite.NotNil(transportationOffice)
}
func (suite *PayloadsSuite) TestTransportationOffices() {
	office := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				ID: uuid.Must(uuid.NewV4()),
			},
		}}, nil)
	officeTwo := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				ID: uuid.Must(uuid.NewV4()),
			},
		}}, nil)
	transportationOfficeList := models.TransportationOffices{}
	transportationOfficeList = append(transportationOfficeList, office, officeTwo)
	value := TransportationOffices(transportationOfficeList)
	suite.NotNil(value)
}
func (suite *PayloadsSuite) TestListMove() {

	marines := models.AffiliationMARINES
	listMove := ListMove(nil)

	suite.Nil(listMove)
	moveUSMC := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &marines,
			},
		},
	}, nil)

	listMove = ListMove(&moveUSMC)
	suite.NotNil(listMove)
}

func (suite *PayloadsSuite) TestListMoves() {
	list := models.Moves{}

	marines := models.AffiliationMARINES
	spaceForce := models.AffiliationSPACEFORCE
	moveUSMC := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &marines,
			},
		},
	}, nil)
	moveSF := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &spaceForce,
			},
		},
	}, nil)
	list = append(list, moveUSMC, moveSF)
	value := ListMoves(&list)
	suite.NotNil(value)
}
func (suite *PayloadsSuite) TestSearchMoves() {
	appCtx := suite.AppContextForTest()

	marines := models.AffiliationMARINES
	moveUSMC := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &marines,
			},
		},
	}, nil)

	moves := models.Moves{moveUSMC}
	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct Marine move with no shipments", func() {
		payload := SearchMoves(appCtx, moves)

		suite.IsType(payload, &ghcmessages.SearchMoves{})
		suite.NotNil(payload)
	})
}

func (suite *PayloadsSuite) TestMarketCode() {
	suite.Run("returns nil when marketCode is nil", func() {
		var marketCode *models.MarketCode = nil
		result := MarketCode(marketCode)
		suite.Equal(result, "")
	})

	suite.Run("returns string when marketCode is not nil", func() {
		marketCodeDomestic := models.MarketCodeDomestic
		result := MarketCode(&marketCodeDomestic)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("d", result, "Expected result to be 'd' for domestic market code")
	})

	suite.Run("returns string when marketCode is international", func() {
		marketCodeInternational := models.MarketCodeInternational
		result := MarketCode(&marketCodeInternational)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("i", result, "Expected result to be 'i' for international market code")
	})
}

func (suite *PayloadsSuite) TestReServiceItem() {
	suite.Run("returns nil when reServiceItem is nil", func() {
		var reServiceItem *models.ReServiceItem = nil
		result := ReServiceItem(reServiceItem)
		suite.Nil(result, "Expected result to be nil when reServiceItem is nil")
	})

	suite.Run("correctly maps ReServiceItem with all fields populated", func() {
		isAutoApproved := true
		marketCodeInternational := models.MarketCodeInternational
		reServiceCode := models.ReServiceCodePOEFSC
		poefscServiceName := "International POE Fuel Surcharge"
		reService := models.ReService{
			Code: reServiceCode,
			Name: poefscServiceName,
		}
		ubShipmentType := models.MTOShipmentTypeUnaccompaniedBaggage
		reServiceItem := &models.ReServiceItem{
			IsAutoApproved: isAutoApproved,
			MarketCode:     marketCodeInternational,
			ReService:      reService,
			ShipmentType:   ubShipmentType,
		}
		result := ReServiceItem(reServiceItem)

		suite.NotNil(result, "Expected result to not be nil when reServiceItem has values")
		suite.Equal(isAutoApproved, result.IsAutoApproved, "Expected IsAutoApproved to match")
		suite.True(result.IsAutoApproved, "Expected IsAutoApproved to be true")
		suite.Equal(string(marketCodeInternational), result.MarketCode, "Expected MarketCode to match")
		suite.Equal(string(reServiceItem.ReService.Code), result.ServiceCode, "Expected ServiceCode to match")
		suite.Equal(string(reServiceItem.ReService.Name), result.ServiceName, "Expected ServiceName to match")
		suite.Equal(string(ubShipmentType), result.ShipmentType, "Expected ShipmentType to match")
	})
}

func (suite *PayloadsSuite) TestReServiceItems() {
	suite.Run("Correctly maps ReServiceItems with all fields populated", func() {
		isAutoApprovedTrue := true
		isAutoApprovedFalse := false
		marketCodeInternational := models.MarketCodeInternational
		marketCodeDomestic := models.MarketCodeDomestic
		poefscReServiceCode := models.ReServiceCodePOEFSC
		podfscReServiceCode := models.ReServiceCodePODFSC
		poefscServiceName := "International POE Fuel Surcharge"
		podfscServiceName := "International POD Fuel Surcharge"
		poefscService := models.ReService{
			Code: poefscReServiceCode,
			Name: poefscServiceName,
		}
		podfscService := models.ReService{
			Code: podfscReServiceCode,
			Name: podfscServiceName,
		}
		hhgShipmentType := models.MTOShipmentTypeHHG
		ubShipmentType := models.MTOShipmentTypeUnaccompaniedBaggage
		poefscServiceItem := models.ReServiceItem{
			IsAutoApproved: isAutoApprovedTrue,
			MarketCode:     marketCodeInternational,
			ReService:      poefscService,
			ShipmentType:   ubShipmentType,
		}
		podfscServiceItem := models.ReServiceItem{
			IsAutoApproved: isAutoApprovedFalse,
			MarketCode:     marketCodeDomestic,
			ReService:      podfscService,
			ShipmentType:   hhgShipmentType,
		}
		reServiceItems := make(models.ReServiceItems, 2)
		reServiceItems[0] = poefscServiceItem
		reServiceItems[1] = podfscServiceItem
		result := ReServiceItems(reServiceItems)

		suite.NotNil(result, "Expected result to not be nil when reServiceItems has values")
		suite.Equal(poefscServiceItem.IsAutoApproved, result[0].IsAutoApproved, "Expected IsAutoApproved to match")
		suite.True(result[0].IsAutoApproved, "Expected IsAutoApproved to be true")
		suite.Equal(string(marketCodeInternational), result[0].MarketCode, "Expected MarketCode to match")
		suite.Equal(string(poefscServiceItem.ReService.Code), result[0].ServiceCode, "Expected ServiceCode to match")
		suite.Equal(string(poefscServiceItem.ReService.Name), result[0].ServiceName, "Expected ServiceName to match")
		suite.Equal(string(ubShipmentType), result[0].ShipmentType, "Expected ShipmentType to match")
		suite.Equal(podfscServiceItem.IsAutoApproved, result[1].IsAutoApproved, "Expected IsAutoApproved to match")
		suite.False(result[1].IsAutoApproved, "Expected IsAutoApproved to be false")
		suite.Equal(string(marketCodeDomestic), result[1].MarketCode, "Expected MarketCode to match")
		suite.Equal(string(podfscServiceItem.ReService.Code), result[1].ServiceCode, "Expected ServiceCode to match")
		suite.Equal(string(podfscServiceItem.ReService.Name), result[1].ServiceName, "Expected ServiceName to match")
		suite.Equal(string(hhgShipmentType), result[1].ShipmentType, "Expected ShipmentType to match")
	})
}

func (suite *PayloadsSuite) TestGsrAppeal() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

	suite.Run("returns nil when gsrAppeal is nil", func() {
		var gsrAppeal *models.GsrAppeal = nil
		result := GsrAppeal(gsrAppeal)
		suite.Nil(result, "Expected result to be nil when gsrAppeal is nil")
	})

	suite.Run("correctly maps GsrAppeal with all fields populated", func() {
		gsrAppealID := uuid.Must(uuid.NewV4())
		reportViolationID := uuid.Must(uuid.NewV4())
		evaluationReportID := uuid.Must(uuid.NewV4())
		appealStatus := models.AppealStatusSustained
		isSeriousIncident := true
		remarks := "Sample remarks"
		createdAt := time.Now()

		gsrAppeal := &models.GsrAppeal{
			ID:                      gsrAppealID,
			ReportViolationID:       &reportViolationID,
			EvaluationReportID:      evaluationReportID,
			OfficeUser:              &officeUser,
			OfficeUserID:            officeUser.ID,
			IsSeriousIncidentAppeal: &isSeriousIncident,
			AppealStatus:            appealStatus,
			Remarks:                 remarks,
			CreatedAt:               createdAt,
		}

		result := GsrAppeal(gsrAppeal)

		suite.NotNil(result, "Expected result to not be nil when gsrAppeal has values")
		suite.Equal(handlers.FmtUUID(gsrAppealID), &result.ID, "Expected ID to match")
		suite.Equal(handlers.FmtUUID(reportViolationID), &result.ViolationID, "Expected ViolationID to match")
		suite.Equal(handlers.FmtUUID(evaluationReportID), &result.ReportID, "Expected ReportID to match")
		suite.Equal(handlers.FmtUUID(officeUser.ID), &result.OfficeUserID, "Expected OfficeUserID to match")
		suite.Equal(ghcmessages.GSRAppealStatusType(appealStatus), result.AppealStatus, "Expected AppealStatus to match")
		suite.Equal(remarks, result.Remarks, "Expected Remarks to match")
		suite.Equal(strfmt.DateTime(createdAt), result.CreatedAt, "Expected CreatedAt to match")
		suite.True(result.IsSeriousIncident, "Expected IsSeriousIncident to be true")
	})

	suite.Run("handles nil ReportViolationID without panic", func() {
		gsrAppealID := uuid.Must(uuid.NewV4())
		evaluationReportID := uuid.Must(uuid.NewV4())
		isSeriousIncident := false
		appealStatus := models.AppealStatusRejected
		remarks := "Sample remarks"
		createdAt := time.Now()

		gsrAppeal := &models.GsrAppeal{
			ID:                      gsrAppealID,
			ReportViolationID:       nil,
			EvaluationReportID:      evaluationReportID,
			OfficeUser:              &officeUser,
			OfficeUserID:            officeUser.ID,
			IsSeriousIncidentAppeal: &isSeriousIncident,
			AppealStatus:            appealStatus,
			Remarks:                 remarks,
			CreatedAt:               createdAt,
		}

		result := GsrAppeal(gsrAppeal)

		suite.NotNil(result, "Expected result to not be nil when gsrAppeal has values")
		suite.Equal(handlers.FmtUUID(gsrAppealID), &result.ID, "Expected ID to match")
		suite.Equal(strfmt.UUID(""), result.ViolationID, "Expected ViolationID to be nil when ReportViolationID is nil")
		suite.Equal(handlers.FmtUUID(evaluationReportID), &result.ReportID, "Expected ReportID to match")
		suite.Equal(handlers.FmtUUID(officeUser.ID), &result.OfficeUserID, "Expected OfficeUserID to match")
		suite.Equal(ghcmessages.GSRAppealStatusType(appealStatus), result.AppealStatus, "Expected AppealStatus to match")
		suite.Equal(remarks, result.Remarks, "Expected Remarks to match")
		suite.Equal(strfmt.DateTime(createdAt), result.CreatedAt, "Expected CreatedAt to match")
		suite.False(result.IsSeriousIncident, "Expected IsSeriousIncident to be false")
	})
}

func (suite *PayloadsSuite) TestVLocation() {
	suite.Run("correctly maps VLocation with all fields populated", func() {
		city := "LOS ANGELES"
		state := "CA"
		postalCode := "90210"
		county := "LOS ANGELES"
		usPostRegionCityID := uuid.Must(uuid.NewV4())

		vLocation := &models.VLocation{
			CityName:             city,
			StateName:            state,
			UsprZipID:            postalCode,
			UsprcCountyNm:        county,
			UsPostRegionCitiesID: &usPostRegionCityID,
		}

		payload := VLocation(vLocation)

		suite.IsType(payload, &ghcmessages.VLocation{})
		suite.Equal(handlers.FmtUUID(usPostRegionCityID), &payload.UsPostRegionCitiesID, "Expected UsPostRegionCitiesID to match")
		suite.Equal(city, payload.City, "Expected City to match")
		suite.Equal(state, payload.State, "Expected State to match")
		suite.Equal(postalCode, payload.PostalCode, "Expected PostalCode to match")
		suite.Equal(county, *(payload.County), "Expected County to match")
	})
}

func (suite *PayloadsSuite) TestMTOServiceItemModel() {
	suite.Run("returns nil when MTOServiceItem is nil", func() {
		var serviceItem *models.MTOServiceItem = nil
		result := MTOServiceItemModel(serviceItem, suite.storer)
		suite.Nil(result, "Expected result to be nil when MTOServiceItem is nil")
	})

	suite.Run("successfully converts MTOServiceItem to payload", func() {
		serviceID := uuid.Must(uuid.NewV4())
		moveID := uuid.Must(uuid.NewV4())
		shipID := uuid.Must(uuid.NewV4())
		reServiceID := uuid.Must(uuid.NewV4())
		now := time.Now()

		mockReService := models.ReService{
			ID:   reServiceID,
			Code: models.ReServiceCodeICRT,
			Name: "Some ReService",
		}

		mockPickupAddress := models.Address{
			ID:        uuid.Must(uuid.NewV4()),
			IsOconus:  models.BoolPointer(false),
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockMTOShipment := models.MTOShipment{
			ID:            shipID,
			PickupAddress: &mockPickupAddress,
		}

		mockServiceItem := models.MTOServiceItem{
			ID:              serviceID,
			MoveTaskOrderID: moveID,
			MTOShipmentID:   &shipID,
			MTOShipment:     mockMTOShipment,
			ReServiceID:     reServiceID,
			ReService:       mockReService,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		result := MTOServiceItemModel(&mockServiceItem, suite.storer)
		suite.NotNil(result, "Expected result to not be nil when MTOServiceItem is valid")
		suite.Equal(handlers.FmtUUID(serviceID), result.ID, "Expected ID to match")
		suite.Equal(handlers.FmtUUID(moveID), result.MoveTaskOrderID, "Expected MoveTaskOrderID to match")
		suite.Equal(handlers.FmtUUIDPtr(&shipID), result.MtoShipmentID, "Expected MtoShipmentID to match")
		suite.Equal(handlers.FmtString(models.MarketConus.FullString()), result.Market, "Expected Market to be CONUS")
	})

	suite.Run("sets Market to OCONUS when PickupAddress.IsOconus is true for ICRT", func() {
		reServiceID := uuid.Must(uuid.NewV4())

		mockReService := models.ReService{
			ID:   reServiceID,
			Code: models.ReServiceCodeICRT,
			Name: "Test ReService",
		}

		mockPickupAddress := models.Address{
			ID:       uuid.Must(uuid.NewV4()),
			IsOconus: models.BoolPointer(true),
		}

		mockMTOShipment := models.MTOShipment{
			PickupAddress: &mockPickupAddress,
		}

		mockServiceItem := models.MTOServiceItem{
			ReService:   mockReService,
			MTOShipment: mockMTOShipment,
		}

		result := MTOServiceItemModel(&mockServiceItem, suite.storer)
		suite.NotNil(result, "Expected result to not be nil for valid MTOServiceItem")
		suite.Equal(handlers.FmtString(models.MarketOconus.FullString()), result.Market, "Expected Market to be OCONUS")
	})

	suite.Run("sets Market to CONUS when PickupAddress.IsOconus is false for ICRT", func() {
		reServiceID := uuid.Must(uuid.NewV4())

		mockReService := models.ReService{
			ID:   reServiceID,
			Code: models.ReServiceCodeICRT,
			Name: "Test ReService",
		}

		mockPickupAddress := models.Address{
			ID:       uuid.Must(uuid.NewV4()),
			IsOconus: models.BoolPointer(false),
		}

		mockMTOShipment := models.MTOShipment{
			PickupAddress: &mockPickupAddress,
		}

		mockServiceItem := models.MTOServiceItem{
			ReService:   mockReService,
			MTOShipment: mockMTOShipment,
		}

		result := MTOServiceItemModel(&mockServiceItem, suite.storer)
		suite.NotNil(result, "Expected result to not be nil for valid MTOServiceItem")
		suite.Equal(handlers.FmtString(models.MarketConus.FullString()), result.Market, "Expected Market to be CONUS")
	})

	suite.Run("sets Market to CONUS when DestinationAddress.IsOconus is false for IUCRT", func() {
		reServiceID := uuid.Must(uuid.NewV4())

		mockReService := models.ReService{
			ID:   reServiceID,
			Code: models.ReServiceCodeIUCRT,
			Name: "Test ReService",
		}

		mockDestinationAddress := models.Address{
			ID:       uuid.Must(uuid.NewV4()),
			IsOconus: models.BoolPointer(false),
		}

		mockMTOShipment := models.MTOShipment{
			DestinationAddress: &mockDestinationAddress,
		}

		mockServiceItem := models.MTOServiceItem{
			ReService:   mockReService,
			MTOShipment: mockMTOShipment,
		}

		result := MTOServiceItemModel(&mockServiceItem, suite.storer)
		suite.NotNil(result, "Expected result to not be nil for valid MTOServiceItem")
		suite.Equal(handlers.FmtString(models.MarketConus.FullString()), result.Market, "Expected Market to be CONUS")
	})

	suite.Run("sets Market to OCONUS when DestinationAddress.IsOconus is true for IUCRT", func() {
		reServiceID := uuid.Must(uuid.NewV4())

		mockReService := models.ReService{
			ID:   reServiceID,
			Code: models.ReServiceCodeIUCRT,
			Name: "Test ReService",
		}

		mockDestinationAddress := models.Address{
			ID:       uuid.Must(uuid.NewV4()),
			IsOconus: models.BoolPointer(true),
		}

		mockMTOShipment := models.MTOShipment{
			DestinationAddress: &mockDestinationAddress,
		}

		mockServiceItem := models.MTOServiceItem{
			ReService:   mockReService,
			MTOShipment: mockMTOShipment,
		}

		result := MTOServiceItemModel(&mockServiceItem, suite.storer)
		suite.NotNil(result, "Expected result to not be nil for valid MTOServiceItem")
		suite.Equal(handlers.FmtString(models.MarketOconus.FullString()), result.Market, "Expected Market to be OCONUS")
	})

	suite.Run("sets Sort from correct serviceItem", func() {
		reServiceID := uuid.Must(uuid.NewV4())

		reServiceItems := make(models.ReServiceItems, 3)
		mockReService := models.ReService{
			ID:             reServiceID,
			Code:           models.ReServiceCodeUBP,
			Name:           "Test ReService",
			ReServiceItems: &reServiceItems,
		}

		mockMTOShipment := models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeUnaccompaniedBaggage,
			MarketCode:   models.MarketCodeInternational,
		}

		reServiceItems[0] = models.ReServiceItem{
			ReService:    mockReService,
			ShipmentType: models.MTOShipmentTypeHHG,
			MarketCode:   models.MarketCodeInternational,
			Sort:         models.StringPointer("0"),
		}
		reServiceItems[1] = models.ReServiceItem{
			ReService:    mockReService,
			ShipmentType: models.MTOShipmentTypeUnaccompaniedBaggage,
			MarketCode:   models.MarketCodeInternational,
			Sort:         models.StringPointer("1"),
		}
		reServiceItems[2] = models.ReServiceItem{
			ReService:    mockReService,
			ShipmentType: models.MTOShipmentTypeUnaccompaniedBaggage,
			MarketCode:   models.MarketCodeDomestic,
			Sort:         models.StringPointer("2"),
		}

		mockMtoServiceItem := models.MTOServiceItem{
			ReService:   mockReService,
			MTOShipment: mockMTOShipment,
		}

		result := MTOServiceItemModel(&mockMtoServiceItem, suite.storer)
		suite.NotNil(result, "Expected result to not be nil for valid MTOServiceItem")
		suite.Equal("1", *result.Sort, "Expected to get the Sort value by matching the correct ReServiceItem using ShipmentType and MarketCode.")
	})
}

func (suite *PayloadsSuite) TestPort() {

	suite.Run("returns nil when PortLocation is nil", func() {
		var mtoServiceItems models.MTOServiceItems = nil
		result := Port(mtoServiceItems, "POE")
		suite.Nil(result, "Expected result to be nil when Port Location is nil")
	})

	suite.Run("Success - Maps PortLocation to Port payload", func() {
		// Use the factory to create a port location
		portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)

		mtoServiceItem := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
			{
				Model:    portLocation,
				LinkOnly: true,
				Type:     &factory.PortLocations.PortOfEmbarkation,
			},
		}, nil)

		// Actual
		mtoServiceItems := models.MTOServiceItems{mtoServiceItem}
		result := Port(mtoServiceItems, "POE")

		// Assert
		suite.IsType(&ghcmessages.Port{}, result)
		suite.Equal(strfmt.UUID(portLocation.ID.String()), result.ID)
		suite.Equal(portLocation.Port.PortType.String(), result.PortType)
		suite.Equal(portLocation.Port.PortCode, result.PortCode)
		suite.Equal(portLocation.Port.PortName, result.PortName)
		suite.Equal(portLocation.City.CityName, result.City)
		suite.Equal(portLocation.UsPostRegionCity.UsprcCountyNm, result.County)
		suite.Equal(portLocation.UsPostRegionCity.UsPostRegion.State.StateName, result.State)
		suite.Equal(portLocation.UsPostRegionCity.UsprZipID, result.Zip)
		suite.Equal(portLocation.Country.CountryName, result.Country)
	})
}

func (suite *PayloadsSuite) TestMTOShipment_POE_POD_Locations() {
	suite.Run("Only POE Location is set", func() {
		// Create mock data for MTOServiceItems with POE and POD
		poePortLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)

		poefscServiceItem := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.ReService{
					Code:     models.ReServiceCodePOEFSC,
					Priority: 1,
				},
			},
			{
				Model:    poePortLocation,
				LinkOnly: true,
				Type:     &factory.PortLocations.PortOfEmbarkation,
			},
		}, nil)

		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MTOServiceItems: models.MTOServiceItems{poefscServiceItem},
				},
			},
		}, nil)

		payload := MTOShipment(nil, &mtoShipment, nil)

		// Assertions
		suite.NotNil(payload, "Expected payload to not be nil")
		suite.NotNil(payload.PoeLocation, "Expected POELocation to not be nil")
		suite.Equal("PDX", payload.PoeLocation.PortCode, "Expected POE Port Code to match")
		suite.Equal("PORTLAND INTL", payload.PoeLocation.PortName, "Expected POE Port Name to match")
		suite.Nil(payload.PodLocation, "Expected PODLocation to be nil when POELocation is set")
	})

	suite.Run("Only POD Location is set", func() {
		// Create mock data for MTOServiceItems with POE and POD
		podPortLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)

		podfscServiceItem := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.ReService{
					Code:     models.ReServiceCodePODFSC,
					Priority: 1,
				},
			},
			{
				Model:    podPortLocation,
				LinkOnly: true,
				Type:     &factory.PortLocations.PortOfDebarkation,
			},
		}, nil)

		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MTOServiceItems: models.MTOServiceItems{podfscServiceItem},
				},
			},
		}, nil)

		payload := MTOShipment(nil, &mtoShipment, nil)

		// Assertions
		suite.NotNil(payload, "Expected payload to not be nil")
		suite.NotNil(payload.PodLocation, "Expected PODLocation to not be nil")
		suite.Equal("PDX", payload.PodLocation.PortCode, "Expected POD Port Code to match")
		suite.Equal("PORTLAND INTL", payload.PodLocation.PortName, "Expected POD Port Name to match")
		suite.Nil(payload.PoeLocation, "Expected PODLocation to be nil when PODLocation is set")
	})
}
func (suite *PayloadsSuite) TestMTOShipment() {
	suite.Run("transforms standard MTOShipment without SIT overrides", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mtoShipment.PrimeEstimatedWeight = models.PoundPointer(1000)
		mtoShipment.PrimeActualWeight = models.PoundPointer(1100)
		miles := unit.Miles(1234)
		mtoShipment.Distance = &miles

		payload := MTOShipment(suite.storer, &mtoShipment, nil)

		suite.NotNil(payload)
		suite.Equal(strfmt.UUID(mtoShipment.ID.String()), payload.ID)
		suite.Equal(handlers.FmtPoundPtr(mtoShipment.PrimeEstimatedWeight), payload.PrimeEstimatedWeight)
		suite.Equal(handlers.FmtPoundPtr(mtoShipment.PrimeActualWeight), payload.PrimeActualWeight)
		suite.Equal(handlers.FmtInt64(1234), payload.Distance)
		suite.Nil(payload.SitStatus)
	})

	suite.Run("SIT overrides total SIT days with SITStatus payload", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mtoShipment.SITDaysAllowance = models.IntPointer(90)

		sitStatusPayload := &ghcmessages.SITStatus{
			TotalSITDaysUsed:   handlers.FmtInt64(int64(10)),
			TotalDaysRemaining: handlers.FmtInt64(int64(40)),
		}

		payload := MTOShipment(suite.storer, &mtoShipment, sitStatusPayload)

		suite.NotNil(payload)
		suite.NotNil(payload.SitDaysAllowance)
		suite.Equal(int64(50), *payload.SitDaysAllowance)
		suite.Equal(sitStatusPayload, payload.SitStatus)
	})

	suite.Run("handles nil Distance", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mtoShipment.Distance = nil

		payload := MTOShipment(suite.storer, &mtoShipment, nil)
		suite.Nil(payload.Distance)
	})

	suite.Run("checks scheduled dates and actual dates set", func() {
		now := time.Now()
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ScheduledPickupDate:   &now,
					ScheduledDeliveryDate: &now,
					ActualPickupDate:      models.TimePointer(now.AddDate(0, 0, 1)),
					ActualDeliveryDate:    models.TimePointer(now.AddDate(0, 0, 2)),
				},
			},
		}, nil)
		payload := MTOShipment(suite.storer, &mtoShipment, nil)
		suite.NotNil(payload.ScheduledPickupDate)
		suite.NotNil(payload.ScheduledDeliveryDate)
		suite.NotNil(payload.ActualPickupDate)
		suite.NotNil(payload.ActualDeliveryDate)
	})
}

func (suite *PayloadsSuite) TestCounselingOffices() {
	suite.Run("correctly maps transportaion offices to counseling offices payload", func() {
		office1 := factory.BuildTransportationOffice(nil, []factory.Customization{
			{
				Model: models.TransportationOffice{
					ID:   uuid.Must(uuid.NewV4()),
					Name: "PPPO Fort Liberty",
				},
			},
		}, nil)

		office2 := factory.BuildTransportationOffice(nil, []factory.Customization{
			{
				Model: models.TransportationOffice{
					ID:   uuid.Must(uuid.NewV4()),
					Name: "PPPO Fort Walker",
				},
			},
		}, nil)

		offices := models.TransportationOffices{office1, office2}

		payload := CounselingOffices(offices)

		suite.IsType(payload, ghcmessages.CounselingOffices{})
		suite.Equal(2, len(payload))
		suite.Equal(office1.ID.String(), payload[0].ID.String())
		suite.Equal(office2.ID.String(), payload[1].ID.String())
	})
}

func (suite *PayloadsSuite) TestMTOShipments() {
	suite.Run("multiple shipments with partial SIT status map", func() {
		shipment1 := factory.BuildMTOShipment(suite.DB(), nil, nil)
		shipment2 := factory.BuildMTOShipment(suite.DB(), nil, nil)

		shipments := models.MTOShipments{shipment1, shipment2}

		// SIT status map that only has SIT info for shipment2
		sitStatusMap := map[string]*ghcmessages.SITStatus{
			shipment2.ID.String(): {
				TotalDaysRemaining: handlers.FmtInt64(30),
				TotalSITDaysUsed:   handlers.FmtInt64(10),
			},
		}

		payload := MTOShipments(suite.storer, &shipments, sitStatusMap)
		suite.NotNil(payload)
		suite.Len(*payload, 2)

		// Shipment1 has no SIT override
		suite.Nil((*payload)[0].SitStatus)

		// Shipment2 has SIT override
		suite.NotNil((*payload)[1].SitStatus)
		suite.Equal(int64(40), *(*payload)[1].SitDaysAllowance)
	})

	suite.Run("nil slice returns empty payload (or nil) gracefully", func() {
		var emptyShipments models.MTOShipments
		payload := MTOShipments(suite.storer, &emptyShipments, nil)
		suite.NotNil(payload)
		suite.Len(*payload, 0)
	})
}

func (suite *PayloadsSuite) TestMTOAgent() {
	suite.Run("transforms a single MTOAgent", func() {
		agent := factory.BuildMTOAgent(suite.DB(), nil, nil)
		payload := MTOAgent(&agent)
		suite.NotNil(payload)
		suite.Equal(strfmt.UUID(agent.ID.String()), payload.ID)
		suite.Equal(string(agent.MTOAgentType), payload.AgentType)
	})
}

func (suite *PayloadsSuite) TestMTOAgents() {
	suite.Run("transforms multiple MTOAgents", func() {
		agent1 := factory.BuildMTOAgent(suite.DB(), nil, nil)
		agent2 := factory.BuildMTOAgent(suite.DB(), nil, nil)
		agents := models.MTOAgents{agent1, agent2}

		payload := MTOAgents(&agents)
		suite.Len(*payload, 2)
		suite.Equal(strfmt.UUID(agent1.ID.String()), (*payload)[0].ID)
		suite.Equal(strfmt.UUID(agent2.ID.String()), (*payload)[1].ID)
	})

	suite.Run("empty slice yields empty payload", func() {
		agents := models.MTOAgents{}
		payload := MTOAgents(&agents)
		suite.NotNil(payload)
		suite.Len(*payload, 0)
	})
}

func (suite *PayloadsSuite) TestPaymentRequests() {
	suite.Run("transforms multiple PaymentRequests", func() {
		pr := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		prs := models.PaymentRequests{pr}

		payload, err := PaymentRequests(suite.AppContextForTest(), &prs, suite.storer)
		suite.NoError(err)
		suite.Len(*payload, 1)
		suite.Equal(strfmt.UUID(pr.ID.String()), (*payload)[0].ID)
	})
}

func (suite *PayloadsSuite) TestPaymentRequest() {
	suite.Run("single PaymentRequest with EDI error info, GEX timestamps, TPPS data", func() {
		pr := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		pr.SentToGexAt = models.TimePointer(time.Now().Add(-1 * time.Hour))
		pr.ReceivedByGexAt = models.TimePointer(time.Now())

		tppsReport := models.TPPSPaidInvoiceReportEntry{
			LineNetCharge:                   2500,
			SellerPaidDate:                  time.Now(),
			InvoiceTotalChargesInMillicents: 500000,
		}
		pr.TPPSPaidInvoiceReports = models.TPPSPaidInvoiceReportEntrys{tppsReport}

		result, err := PaymentRequest(suite.AppContextForTest(), &pr, suite.storer)
		suite.NoError(err)
		suite.NotNil(result)
		suite.NotNil(result.EdiErrorType)
	})
}

func (suite *PayloadsSuite) TestPaymentServiceItem() {
	suite.Run("transforms PaymentServiceItem including MTOServiceItem code and name", func() {
		psi := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		psi.MTOServiceItem.ReService.Code = models.ReServiceCodeDLH
		psi.MTOServiceItem.ReService.Name = "Domestic Linehaul"

		payload := PaymentServiceItem(&psi)
		suite.NotNil(payload)
		suite.Equal(string(models.ReServiceCodeDLH), payload.MtoServiceItemCode)
		suite.Equal("Domestic Linehaul", payload.MtoServiceItemName)
		suite.Equal(ghcmessages.PaymentServiceItemStatus(psi.Status), payload.Status)
	})
}

func (suite *PayloadsSuite) TestPaymentServiceItems() {
	suite.Run("transforms multiple PaymentServiceItems with TPPS data", func() {
		psi1 := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		psi2 := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		items := models.PaymentServiceItems{psi1, psi2}

		tppsReports := models.TPPSPaidInvoiceReportEntrys{
			{
				ProductDescription: string(psi1.MTOServiceItem.ReService.Code),
				LineNetCharge:      1500,
			},
		}

		payload := PaymentServiceItems(&items, &tppsReports)
		suite.NotNil(payload)
		suite.Len(*payload, 2)
		suite.Equal(int64(1500), *(*payload)[0].TppsInvoiceAmountPaidPerServiceItemMillicents)
	})
}

func (suite *PayloadsSuite) TestPaymentServiceItemParam() {
	suite.Run("transforms PaymentServiceItemParam", func() {
		paramKey := factory.BuildServiceItemParamKey(suite.DB(), nil, nil)
		param := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{Model: paramKey},
		}, nil)

		payload := PaymentServiceItemParam(param)
		suite.NotNil(payload)
	})

	suite.Run("handles minimal PaymentServiceItemParam", func() {
		param := models.PaymentServiceItemParam{}
		payload := PaymentServiceItemParam(param)
		suite.NotNil(payload)
	})
}

func (suite *PayloadsSuite) TestPaymentServiceItemParams() {
	suite.Run("transforms slice of PaymentServiceItemParams", func() {
		param1 := factory.BuildPaymentServiceItemParam(suite.DB(), nil, nil)
		param2 := factory.BuildPaymentServiceItemParam(suite.DB(), nil, nil)
		params := models.PaymentServiceItemParams{param1, param2}

		payload := PaymentServiceItemParams(&params)
		suite.NotNil(payload)
		suite.Len(*payload, 2)
	})
}

func (suite *PayloadsSuite) TestServiceRequestDoc() {
	suite.Run("transforms ServiceRequestDocument with multiple uploads", func() {
		serviceRequest := factory.BuildServiceRequestDocument(suite.DB(), nil, nil)
		payload, err := ServiceRequestDoc(serviceRequest, suite.storer)
		suite.NoError(err)
		suite.NotNil(payload)
	})

	suite.Run("handles empty list of uploads", func() {
		serviceRequest := models.ServiceRequestDocument{}
		payload, err := ServiceRequestDoc(serviceRequest, suite.storer)
		suite.NoError(err)
		suite.NotNil(payload)
		suite.Empty(payload.Uploads)
	})
}

func (suite *PayloadsSuite) TestMTOServiceItemSingleModel() {
	suite.Run("transforms basic MTOServiceItem with SIT data", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{Model: mtoShipment, LinkOnly: true},
		}, nil)
		serviceItem.SITEntryDate = models.TimePointer(time.Now().AddDate(0, 0, -2))

		payload := MTOServiceItemSingleModel(&serviceItem)
		suite.NotNil(payload)
		suite.Equal(handlers.FmtDateTimePtr(serviceItem.SITEntryDate), payload.SitEntryDate)
	})
}
