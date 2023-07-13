package testharness

import (
	"log"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

// NamedScenario is a data generation scenario that has a name
type NamedScenario struct {
	Name         string
	SubScenarios map[string]func()
}

// May15TestYear is a May 15 of TestYear
var May15TestYear = time.Date(testdatagen.TestYear, time.May, 15, 0, 0, 0, 0, time.UTC)

// Oct1TestYear is October 1 of TestYear
var Oct1TestYear = time.Date(testdatagen.TestYear, time.October, 1, 0, 0, 0, 0, time.UTC)

// Dec31TestYear is December 31 of TestYear
var Dec31TestYear = time.Date(testdatagen.TestYear, time.December, 31, 0, 0, 0, 0, time.UTC)

// May14FollowingYear is May 14 of the year AFTER TestYear
var May14FollowingYear = time.Date(testdatagen.TestYear+1, time.May, 14, 0, 0, 0, 0, time.UTC)

var May14GHCTestYear = time.Date(testdatagen.GHCTestYear, time.May, 14, 0, 0, 0, 0, time.UTC)

// Closeout offices populated via migrations, this is the ID of one within the GBLOC 'KKFA' with the name 'Creech AFB'
var DefaultCloseoutOfficeID = uuid.FromStringOrNil("5de30a80-a8e5-458c-9b54-edfae7b8cdb9")

// fully public to facilitate reuse outside of this package
type MoveCreatorInfo struct {
	UserID           uuid.UUID
	Email            string
	SmID             uuid.UUID
	FirstName        string
	LastName         string
	MoveID           uuid.UUID
	MoveLocator      string
	CloseoutOfficeID *uuid.UUID
}

func createGenericPPMRelatedMove(appCtx appcontext.AppContext, moveInfo MoveCreatorInfo, userUploader *uploader.UserUploader, moveTemplate *models.Move) models.Move {
	if moveInfo.UserID.IsNil() || moveInfo.Email == "" || moveInfo.SmID.IsNil() || moveInfo.FirstName == "" || moveInfo.LastName == "" || moveInfo.MoveID.IsNil() || moveInfo.MoveLocator == "" {
		log.Panic("All moveInfo fields must have non-zero values.")
	}

	userModel := models.User{
		ID:            moveInfo.UserID,
		LoginGovUUID:  models.UUIDPointer(uuid.Must(uuid.NewV4())),
		LoginGovEmail: moveInfo.Email,
		Active:        true,
	}

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: userModel,
		},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
		{
			Model: models.ServiceMember{
				ID:            moveInfo.SmID,
				FirstName:     models.StringPointer(moveInfo.FirstName),
				LastName:      models.StringPointer(moveInfo.LastName),
				Edipi:         models.StringPointer(factory.RandomEdipi()),
				PersonalEmail: models.StringPointer(moveInfo.Email),
			},
		},
	}, nil)

	if moveInfo.CloseoutOfficeID == nil && (*smWithPPM.Affiliation == models.AffiliationARMY || *smWithPPM.Affiliation == models.AffiliationAIRFORCE) {
		moveInfo.CloseoutOfficeID = &DefaultCloseoutOfficeID
	}

	var customMove models.Move
	if moveTemplate != nil {
		customMove = *moveTemplate
	}
	customMove.ID = moveInfo.MoveID
	customMove.Locator = moveInfo.MoveLocator

	customs := []factory.Customization{
		{
			Model:    smWithPPM,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}
	// this is slightly hacky, but it makes the transformation from
	// using testdatagen.Assertions easier
	if customMove.CloseoutOffice != nil {
		customCloseoutOffice := *customMove.CloseoutOffice
		customMove.CloseoutOffice = nil
		customs = append(customs, factory.Customization{
			Model: customCloseoutOffice,
			Type:  &factory.TransportationOffices.CloseoutOffice,
		})
	} else if moveInfo.CloseoutOfficeID != nil {
		var closeoutOffice models.TransportationOffice
		err := appCtx.DB().Find(&closeoutOffice, *moveInfo.CloseoutOfficeID)
		if err != nil {
			log.Panicf("Cannot load closeout office with ID '%s' from DB: %s",
				moveInfo.CloseoutOfficeID, err)
		}
		customs = append(customs, factory.Customization{
			Model:    closeoutOffice,
			LinkOnly: true,
			Type:     &factory.TransportationOffices.CloseoutOffice,
		})
	}

	customs = append(customs, factory.Customization{
		Model: customMove,
	})

	move := factory.BuildMove(appCtx.DB(), customs, nil)

	return move
}

func CreateGenericMoveWithPPMShipment(appCtx appcontext.AppContext, moveInfo MoveCreatorInfo, useMinimalPPMShipment bool, userUploader *uploader.UserUploader, mtoShipmentTemplate *models.MTOShipment, moveTemplate *models.Move, ppmShipmentTemplate models.PPMShipment) (models.Move, models.PPMShipment) {

	if ppmShipmentTemplate.ID.IsNil() {
		log.Panic("PPMShipment ID cannot be nil.")
	}

	move := createGenericPPMRelatedMove(appCtx, moveInfo, userUploader, moveTemplate)

	customs := []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}

	// This is slightly hacky, but when converting from
	// testdatagen.Assertions, this makes the changes a bit less
	// invasive
	if ppmShipmentTemplate.W2Address != nil {
		customs = append(customs, factory.Customization{
			Model:    *ppmShipmentTemplate.W2Address,
			LinkOnly: true,
			Type:     &factory.Addresses.W2Address,
		})
		ppmShipmentTemplate.W2Address = nil
	}
	customs = append(customs, factory.Customization{
		Model: ppmShipmentTemplate,
	})

	if mtoShipmentTemplate != nil {
		customs = append(customs, factory.Customization{
			Model: *mtoShipmentTemplate,
		})
	}
	if useMinimalPPMShipment {
		return move, factory.BuildMinimalPPMShipment(appCtx.DB(), customs, nil)
	}

	// assertions passed in means we cannot yet convert to BuildPPMShipment
	return move, factory.BuildPPMShipment(appCtx.DB(), customs, nil)
}

func CreateMoveWithCloseOut(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveInfo MoveCreatorInfo, branch models.ServiceMemberAffiliation) models.Move {
	loginGovUUID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:            moveInfo.UserID,
				LoginGovUUID:  &loginGovUUID,
				LoginGovEmail: moveInfo.Email,
				Active:        true,
			}},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            moveInfo.SmID,
				FirstName:     models.StringPointer(moveInfo.FirstName),
				LastName:      models.StringPointer(moveInfo.LastName),
				PersonalEmail: models.StringPointer(moveInfo.Email),
				Affiliation:   &branch,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	address := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
				Country:        models.StringPointer("US"),
			},
		},
	}, nil)

	newDutyLocation := factory.BuildDutyLocation(appCtx.DB(), []factory.Customization{
		{
			Model:    address,
			LinkOnly: true,
		},
	}, nil)

	var closeoutOffice models.TransportationOffice
	if moveInfo.CloseoutOfficeID != nil {
		err := appCtx.DB().Q().Where(`id=$1`, moveInfo.CloseoutOfficeID).First(&closeoutOffice)
		if err != nil {
			log.Panic(err)
		}
	} else if branch == models.AffiliationARMY || branch == models.AffiliationAIRFORCE {
		err := appCtx.DB().Q().Where(`id=$1`, DefaultCloseoutOfficeID).First(&closeoutOffice)
		if err != nil {
			log.Panic(err)
		}
	}

	customs := []factory.Customization{
		{
			Model:    smWithPPM,
			LinkOnly: true,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Move{
				ID:          moveInfo.MoveID,
				Locator:     moveInfo.MoveLocator,
				Status:      models.MoveStatusAPPROVED,
				SubmittedAt: &submittedAt,
				PPMType:     models.StringPointer("FULL"),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}

	if !closeoutOffice.ID.IsNil() {
		customs = append(customs, factory.Customization{
			Model:    closeoutOffice,
			LinkOnly: true,
			Type:     &factory.TransportationOffices.CloseoutOffice,
		})
	}

	move := factory.BuildMove(appCtx.DB(), customs, nil)

	mtoShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				Status:      models.PPMShipmentStatusNeedsPaymentApproval,
				SubmittedAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)

	factory.BuildSignedCertification(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	return move
}

func CreateMoveWithCloseoutOffice(appCtx appcontext.AppContext, moveInfo MoveCreatorInfo, userUploader *uploader.UserUploader) models.Move {
	loginGovUUID := uuid.Must(uuid.NewV4())
	submittedAt := time.Date(2020, time.December, 11, 12, 0, 0, 0, time.UTC)

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:            moveInfo.UserID,
				LoginGovUUID:  &loginGovUUID,
				LoginGovEmail: moveInfo.Email,
				Active:        true,
			}},
	}, nil)

	branch := models.AffiliationAIRFORCE
	serviceMember := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            moveInfo.SmID,
				FirstName:     models.StringPointer(moveInfo.FirstName),
				LastName:      models.StringPointer(moveInfo.LastName),
				PersonalEmail: models.StringPointer(moveInfo.Email),
				Affiliation:   &branch,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Name: "Los Angeles AFB"},
		},
	}, nil)

	// Make a move with the closeout office
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    closeoutOffice,
			LinkOnly: true,
			Type:     &factory.TransportationOffices.CloseoutOffice,
		},
		{
			Model: models.Move{
				ID:          moveInfo.MoveID,
				Locator:     moveInfo.MoveLocator,
				SubmittedAt: &submittedAt,
				Status:      models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	mtoShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsPaymentApproval,
			},
		},
	}, nil)

	return move
}

func CreateSubmittedMoveWithPPMShipmentForSC(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, _ services.MoveRouter, moveInfo MoveCreatorInfo) models.Move {
	loginGovUUID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:            moveInfo.UserID,
				LoginGovUUID:  &loginGovUUID,
				LoginGovEmail: moveInfo.Email,
				Active:        true,
			}},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            moveInfo.SmID,
				FirstName:     models.StringPointer(moveInfo.FirstName),
				LastName:      models.StringPointer(moveInfo.LastName),
				PersonalEmail: models.StringPointer(moveInfo.Email),
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    smWithPPM,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:          moveInfo.MoveID,
				Locator:     moveInfo.MoveLocator,
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	if *smWithPPM.Affiliation == models.AffiliationARMY || *smWithPPM.Affiliation == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}
	mtoShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildSignedCertification(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	return move
}

// A generic method
func CreateMoveWithOptions(appCtx appcontext.AppContext, assertions testdatagen.Assertions) models.Move {

	ordersType := assertions.Order.OrdersType
	shipmentType := assertions.MTOShipment.ShipmentType
	destinationType := assertions.MTOShipment.DestinationType
	locator := assertions.Move.Locator
	status := assertions.Move.Status
	servicesCounseling := assertions.DutyLocation.ProvidesServicesCounseling
	usesExternalVendor := assertions.MTOShipment.UsesExternalVendor

	db := appCtx.DB()
	submittedAt := time.Now()
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: servicesCounseling,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType: ordersType,
			},
		},
	}, nil)
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:     locator,
				Status:      status,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := factory.BuildAddress(db, nil, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          shipmentType,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				DestinationType:       destinationType,
				UsesExternalVendor:    usesExternalVendor,
			},
		},
		{
			Model:    destinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	requestedPickupDate = submittedAt.Add(30 * 24 * time.Hour)
	requestedDeliveryDate = requestedPickupDate.Add(7 * 24 * time.Hour)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          shipmentType,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
	}, nil)

	return move
}

/*
Create Needs Service Counseling - pass in orders with all required information, shipment type, destination type, locator
*/
func CreateNeedsServicesCounseling(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType, shipmentType models.MTOShipmentType, destinationType *models.DestinationType, locator string) models.Move {
	db := appCtx.DB()
	submittedAt := time.Now()
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	ordersNumber := "8675309"
	departmentIndicator := "ARMY"
	tac := "E19A"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          ordersType,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:     locator,
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := factory.BuildAddress(db, nil, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          shipmentType,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				DestinationType:       destinationType,
			},
		},
		{
			Model:    destinationAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	requestedPickupDate = submittedAt.Add(30 * 24 * time.Hour)
	requestedDeliveryDate = requestedPickupDate.Add(7 * 24 * time.Hour)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          shipmentType,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
	}, nil)
	officeUser := factory.BuildOfficeUserWithRoles(db, nil, []roles.RoleType{roles.RoleTypeTOO})
	factory.BuildCustomerSupportRemark(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    officeUser,
			LinkOnly: true,
		},
		{
			Model: models.CustomerSupportRemark{
				Content: "The customer mentioned that they need to provide some more complex instructions for pickup and drop off.",
			},
		},
	}, nil)

	return move
}

func CreateNeedsServicesCounselingMinimalNTSR(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType, locator string) models.Move {
	db := appCtx.DB()
	submittedAt := time.Now()
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType: ordersType,
			},
		},
	}, nil)
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:     locator,
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
	// Makes a basic NTS-R shipment with minimal info.
	requestedDeliveryDate := time.Now().AddDate(0, 0, 14)
	destinationAddress := factory.BuildAddress(db, nil, nil)
	factory.BuildMTOShipmentMinimal(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHGOutOfNTSDom,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
		{
			Model:    destinationAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	return move
}

// MakeSITExtensionsForShipment helper function
func MakeSITExtensionsForShipment(appCtx appcontext.AppContext, shipment models.MTOShipment) {
	db := appCtx.DB()
	sitContractorRemarks1 := "The customer requested an extension."
	sitOfficeRemarks1 := "The service member is unable to move into their new home at the expected time."
	approvedDays := 90

	factory.BuildSITDurationUpdate(db, []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model: models.SITDurationUpdate{
				ContractorRemarks: &sitContractorRemarks1,
				OfficeRemarks:     &sitOfficeRemarks1,
				ApprovedDays:      &approvedDays,
			},
		},
	}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})

	factory.BuildSITDurationUpdate(db, []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model: models.SITDurationUpdate{
				ApprovedDays: &approvedDays,
			},
		},
	}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})
}

func CreateMoveWithHHGAndNTSShipments(appCtx appcontext.AppContext, locator string, usesExternalVendor bool) models.Move {
	db := appCtx.DB()
	submittedAt := time.Now()
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
			},
		},
	}, nil)
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:     locator,
				Status:      models.MoveStatusSUBMITTED,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
	// Makes a basic HHG shipment to reflect likely real scenario
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := factory.BuildAddress(db, nil, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
		{
			Model:    destinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildNTSShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{

				Status:             models.MTOShipmentStatusSubmitted,
				UsesExternalVendor: usesExternalVendor,
			},
		},
	}, nil)

	return move
}

func CreateMoveWithHHGAndNTSRShipments(appCtx appcontext.AppContext, locator string, usesExternalVendor bool) models.Move {
	db := appCtx.DB()
	submittedAt := time.Now()
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
			},
		},
	}, nil)
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:     locator,
				Status:      models.MoveStatusSUBMITTED,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
	// Makes a basic HHG shipment to reflect likely real scenario
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := factory.BuildAddress(db, nil, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
		{
			Model:    destinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:             models.MTOShipmentStatusSubmitted,
				UsesExternalVendor: usesExternalVendor,
			},
		},
	}, nil)

	return move
}

func CreateMoveWithNTSShipment(appCtx appcontext.AppContext, locator string, usesExternalVendor bool) models.Move {
	db := appCtx.DB()
	submittedAt := time.Now()
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
			},
		},
	}, nil)
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:     locator,
				Status:      models.MoveStatusSUBMITTED,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
	factory.BuildNTSShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:             models.MTOShipmentStatusSubmitted,
				UsesExternalVendor: usesExternalVendor,
			},
		},
	}, nil)

	return move
}
