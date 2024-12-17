package payloads

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
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
		County:         county,
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
			County:         county,
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
		County:         *models.StringPointer("WASHOE"),
	}

	oldAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         *models.StringPointer("WASHOE"),
	}

	sitOriginalAddress := models.Address{
		StreetAddress1: "123 SIT St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89501",
		County:         *models.StringPointer("WASHOE"),
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
		County:         *models.StringPointer("WASHOE"),
	}

	backupAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         *models.StringPointer("WASHOE"),
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
		County:         *models.StringPointer("WASHOE"),
	}

	backupAddress := models.Address{
		StreetAddress1: "123 Old St",
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "89502",
		County:         *models.StringPointer("WASHOE"),
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
	spaceForce := models.AffiliationSPACEFORCE
	army := models.AffiliationARMY
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
	moveA := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &army,
			},
		},
	}, nil)
	moveUSMC.Status = models.MoveStatusNeedsServiceCounseling
	scheduledPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 20, 0, 0, 0, 0, time.UTC)
	scheduledDeliveryDate := time.Date(testdatagen.GHCTestYear, time.September, 20, 0, 0, 0, 0, time.UTC)
	sitAllowance := int(90)
	gbloc := "LKNQ"
	storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)
	mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    moveSF,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:                models.MTOShipmentStatusApproved,
				ShipmentType:          models.MTOShipmentTypeHHGIntoNTSDom,
				CounselorRemarks:      handlers.FmtString("counselor remark"),
				SITDaysAllowance:      &sitAllowance,
				ScheduledPickupDate:   &scheduledPickupDate,
				ScheduledDeliveryDate: &scheduledDeliveryDate,
			},
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
	}, nil)

	moveSF.MTOShipments = append(moveSF.MTOShipments, mtoShipment)
	moveSF.ShipmentGBLOC = append(moveSF.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

	moves := models.Moves{moveUSMC}
	moveSpaceForce := models.Moves{moveSF}
	moveArmy := models.Moves{moveA}
	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct Marine move with no shipments", func() {
		payload := SearchMoves(appCtx, moves)

		suite.IsType(payload, &ghcmessages.SearchMoves{})
		suite.NotNil(payload)
	})
	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct Non-Marine move, a shipment, and delivery/pickup time.  ", func() {
		payload := SearchMoves(appCtx, moveSpaceForce)
		suite.IsType(payload, &ghcmessages.SearchMoves{})
		suite.NotNil(payload)
		suite.NotNil(mtoShipment)

		suite.NotNil(moveA)
	})
	suite.Run("Success - Returns a ghcmessages Upload payload from Upload Struct Army move, with no shipments.  ", func() {
		payload := SearchMoves(appCtx, moveArmy)
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
		poedscReServiceCode := models.ReServiceCodePODFSC
		poefscServiceName := "International POE Fuel Surcharge"
		poedscServiceName := "International POD Fuel Surcharge"
		poefscService := models.ReService{
			Code: poefscReServiceCode,
			Name: poefscServiceName,
		}
		podfscService := models.ReService{
			Code: poedscReServiceCode,
			Name: poedscServiceName,
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
}
