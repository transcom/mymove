package scenario

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/random"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

// NamedScenario is a data generation scenario that has a name
type NamedScenario struct {
	Name         string
	SubScenarios map[string]func()
}

type sceneOptionsNTS struct {
	shipmentMoveCode   string
	moveStatus         models.MoveStatus
	usesExternalVendor bool
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

var estimatedWeight = unit.Pound(1400)
var actualWeight = unit.Pound(2000)
var ntsRecordedWeight = unit.Pound(2000)
var tioRemarks = "New billable weight set"

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

func makeOrdersForServiceMember(appCtx appcontext.AppContext, serviceMember models.ServiceMember, userUploader *uploader.UserUploader, fileNames *[]string) models.Order {
	document := factory.BuildDocumentLinkServiceMember(appCtx.DB(), serviceMember)

	// Creates order upload documents from the files in this directory:
	// pkg/testdatagen/testdata/bandwidth_test_docs

	files := filesInBandwidthTestDirectory(fileNames)

	for _, file := range files {
		filePath := fmt.Sprintf("bandwidth_test_docs/%s", file)
		fixture := testdatagen.Fixture(filePath)

		upload := factory.BuildUserUpload(appCtx.DB(), []factory.Customization{
			{
				Model:    document,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   appCtx,
					File:         fixture,
				},
			},
		}, nil)
		document.UserUploads = append(document.UserUploads, upload)
	}

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    document,
			LinkOnly: true,
			Type:     &factory.Documents.UploadedOrders,
		},
	}, nil)

	return orders
}

func makeMoveForOrders(appCtx appcontext.AppContext, orders models.Order, moveCode string, moveStatus models.MoveStatus,
	moveOptConfigs ...func(move *models.Move)) models.Move {

	var availableToPrimeAt *time.Time
	if moveStatus == models.MoveStatusAPPROVED || moveStatus == models.MoveStatusAPPROVALSREQUESTED {
		now := time.Now()
		availableToPrimeAt = &now
	}

	move := models.Move{
		Status:             moveStatus,
		OrdersID:           orders.ID,
		Orders:             orders,
		Locator:            moveCode,
		AvailableToPrimeAt: availableToPrimeAt,
	}

	// run configurations on move struct
	// this is to make any updates to the move struct before it gets created
	for _, config := range moveOptConfigs {
		config(&move)
	}

	move = factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             moveStatus,
				Locator:            moveCode,
				AvailableToPrimeAt: availableToPrimeAt,
			},
		},
	}, nil)

	return move
}

func createServiceMemberWithOrdersButNoMoveType(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/*
	 * A service member with orders and a move, but no move type selected
	 */
	email := "sm_no_move_type@example.com"
	uuidStr := "9ceb8321-6a82-4f6d-8bb3-a1d85922a202"
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildMove(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("7554e347-2215-484f-9240-c61bae050220"),
				FirstName:     models.StringPointer("LandingTest1"),
				LastName:      models.StringPointer("UserPerson2"),
				Edipi:         models.StringPointer("6833908164"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("b2ecbbe5-36ad-49fc-86c8-66e55e0697a7"),
				Locator: "ZPGVED",
			},
		},
	}, nil)
}

func createServiceMemberWithNoUploadedOrders(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/*
	 * Service member with no uploaded orders
	 */
	email := "needs@orde.rs"
	uuidStr := "feac0e92-66ec-4cab-ad29-538129bf918e"
	oktaID := uuid.Must(uuid.NewV4())
	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
	}, nil)

	factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("c52a9f13-ccc7-4c1b-b5ef-e1132a4f4db9"),
				FirstName:     models.StringPointer("NEEDS"),
				LastName:      models.StringPointer("ORDERS"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
}

// Create a move with both HHG and PPM shipments, used to test partial PPM orders.
func CreateMoveWithHHGAndPPM(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveInfo MoveCreatorInfo, branch models.ServiceMemberAffiliation, readyForCloseout bool) models.Move {
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        moveInfo.UserID,
				OktaID:    oktaID.String(),
				OktaEmail: moveInfo.Email,
				Active:    true,
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
				CacValidated:  true,
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

	// If not ready for closeout, just submit the shipment
	moveStatus := models.MoveStatusNeedsServiceCounseling
	if readyForCloseout {
		moveStatus = models.MoveStatusAPPROVED
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
				Status:      moveStatus,
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

	// If not ready for closeout, just submit the shipment
	ppmShipmentStatus := models.PPMShipmentStatusSubmitted
	if readyForCloseout {
		ppmShipmentStatus = models.PPMShipmentStatusNeedsCloseout
	}

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
				Status:      ppmShipmentStatus,
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

func createMoveWithPPMAndHHG(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	/*
	 * A service member with orders and a submitted move with a ppm and hhg
	 */
	email := "combo@ppm.hhg"
	uuidStr := "6016e423-f8d5-44ca-98a8-af03c8445c94"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
	}, nil)

	smIDCombo := "f6bd793f-7042-4523-aa30-34946e7339c9"
	smWithCombo := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smIDCombo),
				FirstName:     models.StringPointer("Submitted"),
				LastName:      models.StringPointer("Ppmhhg"),
				Edipi:         models.StringPointer("6833908165"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    smWithCombo,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("7024c8c5-52ca-4639-bf69-dd8238308c98"),
				Locator: "COMBOS",
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

	if *smWithCombo.Affiliation == models.AffiliationARMY || *smWithCombo.Affiliation == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(2000)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("8689afc7-84d6-4c60-a739-8cf96ede2606"),
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("8689afc7-84d6-4c60-a739-333333333333"),
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
				CounselorRemarks:     models.StringPointer("Please handle with care"),
			},
		},
	}, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusRejected,
				RejectionReason:      models.StringPointer("No longer necessary, included in other shipment"),
			},
		},
	}, nil)

	factory.BuildPPMShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				ID: uuid.FromStringOrNil("d733fe2f-b08d-434a-ad8d-551f4d597b03"),
			},
		},
	}, nil)

	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
}

func createGenericPPMRelatedMove(appCtx appcontext.AppContext, moveInfo MoveCreatorInfo, userUploader *uploader.UserUploader, moveTemplate *models.Move) models.Move {
	if moveInfo.UserID.IsNil() || moveInfo.Email == "" || moveInfo.SmID.IsNil() || moveInfo.FirstName == "" || moveInfo.LastName == "" || moveInfo.MoveID.IsNil() || moveInfo.MoveLocator == "" {
		log.Panic("All moveInfo fields must have non-zero values.")
	}

	userModel := models.User{
		ID:        moveInfo.UserID,
		OktaID:    models.UUIDPointer(uuid.Must(uuid.NewV4())).String(),
		OktaEmail: moveInfo.Email,
		Active:    true,
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
				CacValidated:  true,
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

func createUnSubmittedMoveWithMinimumPPMShipment(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a minimal PPM Shipment. This means the PPM only has required fields.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("bbb469f3-f4bc-420d-9755-b9569f81715e"),
		Email:       "dates_and_locations@ppm.unsubmitted",
		SmID:        testdatagen.ConvertUUIDStringToUUID("635e4c37-63b8-4860-9239-0e743ec383b0"),
		FirstName:   "Minimal",
		LastName:    "PPM",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("16cb4b73-cc0e-48c5-8cc7-b2a2ac52c342"),
		MoveLocator: "PPMMIN",
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID: testdatagen.ConvertUUIDStringToUUID("ffc95935-6781-4f95-9f35-16a5994cab56"),
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, true, userUploader, nil, nil, assertions.PPMShipment)
}

func createUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {

	/*
	 * A service member with orders and a PPM shipment updated with an estimated weight value and estimated incentive
	 */
	moveInfo := MoveCreatorInfo{
		UserID:           testdatagen.ConvertUUIDStringToUUID("4512dc8c-c777-444e-b6dc-7971e398f2dc"),
		Email:            "estimated_weights@ppm.unsubmitted",
		SmID:             testdatagen.ConvertUUIDStringToUUID("81b772ab-86ff-4bda-b0fa-21b14dfe14d5"),
		FirstName:        "EstimatedWeights",
		LastName:         "PPM",
		MoveID:           testdatagen.ConvertUUIDStringToUUID("e89a7018-be76-449a-b99b-e30a09c485dc"),
		MoveLocator:      "PPMEWH",
		CloseoutOfficeID: &DefaultCloseoutOfficeID,
	}
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                 testdatagen.ConvertUUIDStringToUUID("65eea403-89ac-4c2d-9b1c-0dcc8805258f"),
			EstimatedWeight:    models.PoundPointer(unit.Pound(4000)),
			HasProGear:         models.BoolPointer(false),
			EstimatedIncentive: models.CentPointer(unit.Cents(1000000)),
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, true, userUploader, nil, nil, assertions.PPMShipment)
}

func createUnSubmittedMoveWithPPMShipmentThroughAdvanceRequested(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a minimal PPM Shipment. This means the PPM only has required fields.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("dd1a3982-1ec4-4e34-a7bd-73cba4f3376a"),
		Email:       "advance_requested@ppm.unsubmitted",
		SmID:        testdatagen.ConvertUUIDStringToUUID("7a402a11-92a0-4334-b297-551be2bc44ef"),
		FirstName:   "HasAdvanceRequested",
		LastName:    "PPM",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("fe322fae-c13e-4961-9956-69fb7a491ad4"),
		MoveLocator: "PPMADV",
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                     testdatagen.ConvertUUIDStringToUUID("9160a396-9b60-41c2-af7a-aa03d5002c71"),
			EstimatedWeight:        models.PoundPointer(unit.Pound(4000)),
			HasProGear:             models.BoolPointer(false),
			EstimatedIncentive:     models.CentPointer(unit.Cents(10000000)),
			HasRequestedAdvance:    models.BoolPointer(true),
			AdvanceAmountRequested: models.CentPointer(unit.Cents(30000)),
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, true, userUploader, nil, nil, assertions.PPMShipment)
}

func createUnSubmittedMoveWithFullPPMShipment1(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a full PPM Shipment.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("1b16773e-995b-4efe-ad1c-bef2ae1253f8"),
		Email:       "full@ppm.unsubmitted",
		SmID:        testdatagen.ConvertUUIDStringToUUID("1b400031-2b78-44ce-976c-cd2e854947f8"),
		FirstName:   "Full",
		LastName:    "PPM",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("3e0b6cb9-3409-4089-83a0-0fbc3fb0b493"),
		MoveLocator: "FULLPP",
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		MTOShipment: models.MTOShipment{
			ID: testdatagen.ConvertUUIDStringToUUID("0e17d5de-b212-404d-9249-e5a160bd0c51"),
		},
		PPMShipment: models.PPMShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("11978e1c-95d3-47e6-9d3f-d1e0d8c3d11a"),
			Status: models.PPMShipmentStatusDraft,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, nil, assertions.PPMShipment)
}

func createUnSubmittedMoveWithFullPPMShipment2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a full PPM Shipment.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:           testdatagen.ConvertUUIDStringToUUID("b54d5368-a633-4e3e-a8df-22133b9f8c7c"),
		Email:            "happyPathWithEdits@ppm.unsubmitted",
		SmID:             testdatagen.ConvertUUIDStringToUUID("f7bd4d55-c245-4f58-b638-e44f98ab2f32"),
		FirstName:        "Finished",
		LastName:         "PPM",
		MoveID:           testdatagen.ConvertUUIDStringToUUID("b122621c-8577-4b3f-a392-4ade43169fe9"),
		MoveLocator:      "PPMHPE",
		CloseoutOfficeID: &DefaultCloseoutOfficeID,
	}

	departureDate := time.Date(2022, time.February, 01, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                    testdatagen.ConvertUUIDStringToUUID("d328333e-e6c8-47d7-8cdf-5864a16abf45"),
			Status:                models.PPMShipmentStatusDraft,
			EstimatedWeight:       models.PoundPointer(unit.Pound(4000)),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			ExpectedDepartureDate: departureDate,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)
}

func createUnSubmittedMoveWithFullPPMShipment3(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a full PPM Shipment.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:           testdatagen.ConvertUUIDStringToUUID("9365990e-5813-4031-aa42-170886150912"),
		Email:            "happyPathWithEditsMobile@ppm.unsubmitted",
		SmID:             testdatagen.ConvertUUIDStringToUUID("70d7372a-7e91-4b8f-927d-624cfe29ab6d"),
		FirstName:        "Finished",
		LastName:         "PPM",
		MoveID:           testdatagen.ConvertUUIDStringToUUID("4d0aa509-e6ee-4757-ad14-368e334fc51f"),
		MoveLocator:      "PPMHPM",
		CloseoutOfficeID: &DefaultCloseoutOfficeID,
	}

	departureDate := time.Date(2022, time.February, 01, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                    testdatagen.ConvertUUIDStringToUUID("6f7d6ac2-a38b-4df6-a82a-1ea9b352de89"),
			Status:                models.PPMShipmentStatusDraft,
			EstimatedWeight:       models.PoundPointer(unit.Pound(4000)),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			ExpectedDepartureDate: departureDate,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)
}

func createUnSubmittedMoveWithFullPPMShipment4(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a full PPM Shipment.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("781cf194-4eb2-4def-9da6-01abdc62333d"),
		Email:       "deleteShipmentMobile@ppm.unsubmitted",
		SmID:        testdatagen.ConvertUUIDStringToUUID("fc9264ae-4290-4445-987d-f6950b97c865"),
		FirstName:   "Delete",
		LastName:    "PPM",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("a11cae72-56f0-45a3-a546-3af43a1d50ea"),
		MoveLocator: "PPMDEL",
	}

	departureDate := time.Date(2022, time.February, 01, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                    testdatagen.ConvertUUIDStringToUUID("47f6cb98-fbd1-4b95-a91b-2d394d555d21"),
			Status:                models.PPMShipmentStatusDraft,
			EstimatedWeight:       models.PoundPointer(unit.Pound(4000)),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			ExpectedDepartureDate: departureDate,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)
}

func createUnSubmittedMoveWithFullPPMShipment5(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a full PPM Shipment.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("57d58062-93ac-4eb7-b1da-21dd137e4f65"),
		Email:       "deleteShipmentMobile@ppm.unsubmitted",
		SmID:        testdatagen.ConvertUUIDStringToUUID("d5778927-7366-44c2-8dbf-1bce14906adc"),
		FirstName:   "Delete",
		LastName:    "PPM",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("ae5e7087-8e1e-49ae-98cc-0727a5cd11eb"),
		MoveLocator: "DELPPM",
	}

	departureDate := time.Date(2022, time.February, 01, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                    testdatagen.ConvertUUIDStringToUUID("0a62f7c6-72d2-4f4f-9889-202f3c0222a6"),
			Status:                models.PPMShipmentStatusDraft,
			EstimatedWeight:       models.PoundPointer(unit.Pound(4000)),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			ExpectedDepartureDate: departureDate,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)
}

func createApprovedMoveWithPPM(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("cde987a1-a717-4a61-98b5-1f05e2e0844d"),
		Email:       "readyToFinish@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("dfbba0fc-2a70-485e-9eb2-ac80f3861032"),
		FirstName:   "Ready",
		LastName:    "Finish",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("26b960d8-a96d-4450-a441-673ccd7cc3c7"),
		MoveLocator: "PPMRF1",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("2ed2998e-ae36-46cd-af83-c3ecee55fe3e"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         testdatagen.ConvertUUIDStringToUUID("b9ae4c25-1376-4b9b-8781-106b5ae7ecab"),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPM2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {

	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("c28b2eb1-975f-49f7-b8a3-c7377c0da908"),
		Email:       "readyToFinish2@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("6456ffbb-d114-4ec5-a736-6cb63a65bfd7"),
		FirstName:   "Ready2",
		LastName:    "Finish2",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("0e33adbc-20b4-4a93-9ce5-7ee4695a0307"),
		MoveLocator: "PPMRF2",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status:           models.MoveStatusAPPROVED,
			CloseoutOfficeID: &DefaultCloseoutOfficeID,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("ef256d30-a6e7-4be8-8a60-b4ffb7dc7a7f"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         testdatagen.ConvertUUIDStringToUUID("1ce52409-009d-4d9c-a48c-b12013fa2d2b"),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPM3(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("539af373-9474-49f3-b06b-bc4b4d4111de"),
		Email:       "readyToFinish3@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("1b543655-6e5a-4ea0-b4e0-48fe4e107ef5"),
		FirstName:   "Ready3",
		LastName:    "Finish3",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("3cf2a0eb-08e6-404d-81ad-022e1aaf26aa"),
		MoveLocator: "PPMRF3",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("1f452b86-4488-46f5-98c0-b696e1410522"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         testdatagen.ConvertUUIDStringToUUID("7d8f77c3-9829-4241-b0a7-b2897f1d6822"),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPM4(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("c48998dc-8f93-437a-bd0c-2c0b187b12cb"),
		Email:       "readyToFinish4@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("16d13649-f246-456f-8093-da3a769a1247"),
		FirstName:   "Ready4",
		LastName:    "Finish4",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("9061587a-5b31-4deb-9947-703a40857fa8"),
		MoveLocator: "PPMRF4",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("ae873226-67a4-452f-b92d-924307ff2d9a"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         testdatagen.ConvertUUIDStringToUUID("881f1084-d5a8-4210-9854-fa5f01c8da81"),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPM5(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("62e20f62-638f-4390-bbc0-c672cd7fd2e3"),
		Email:       "readyToFinish5@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("55643c43-f48b-471d-8b99-b1e2a0ce5215"),
		FirstName:   "Ready5",
		LastName:    "Finish5",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("7dcbf7ef-9a74-4efa-b536-c334b2093bc0"),
		MoveLocator: "PPMRF5",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("38a9ff5a-76c5-4126-9dc8-649a1f35e847"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         testdatagen.ConvertUUIDStringToUUID("bcbd9762-2041-42e5-9b91-ba5b1ecb3487"),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPM6(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("1dca189a-ca7e-4e70-b98e-be3829e4b6cc"),
		Email:       "readyForCloseout@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("6b4ce016-9b76-44a8-a870-f378313aa1a8"),
		FirstName:   "Ready",
		LastName:    "Closeout",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("5b718175-8bc5-4ca9-a1f0-8b70d064ee92"),
		MoveLocator: "PPMCL0",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("04592d80-f67f-443e-b9d6-967a9befcc3a"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         testdatagen.ConvertUUIDStringToUUID("cdc68d38-21d9-4bd8-bd56-5b5c224ab2ab"),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPM7(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("fe825617-a53a-49bf-bf2e-c271afee344d"),
		Email:       "readyForCloseout2@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("c1ba0a4b-4873-479a-a3d8-4158afdbe7b0"),
		FirstName:   "Ready",
		LastName:    "Closeout",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("dabe45ab-aeab-4f83-b446-f1f70e265beb"),
		MoveLocator: "PPMRC2",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("0097e9d1-7579-4f6f-a71e-1b63aea0d4c7"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         testdatagen.ConvertUUIDStringToUUID("7276375a-932f-4b93-b706-3da2774dfd92"),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWeightTicket(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("33f39cca-3908-4cf5-b7d9-839741f51911"),
		Email:       "weightTicketPPM@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("a30fd609-6dcf-4dd0-a7e6-2892a31ae641"),
		FirstName:   "ActualPPM",
		LastName:    "WeightTicketComplete",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("2fdb02a5-dd80-4ec4-a9f0-f4eefb434568"),
		MoveLocator: "W3TT1K",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("3e0c9457-9010-473a-a274-fc1620d5ee16"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("7a5e932d-f1f6-435e-9518-3ee33f74bc88"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
}

func createApprovedMoveWithPPMExcessWeight(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveInfo MoveCreatorInfo) (models.Move, models.PPMShipment) {
	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          uuid.Must(uuid.NewV4()),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			AdvanceStatus:               (*models.PPMAdvanceStatus)(models.StringPointer(string(models.PPMAdvanceStatusApproved))),
			W2Address:                   &address,
			AllowableWeight:             models.PoundPointer(19000),
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.WeightTicket{
				EmptyWeight: models.PoundPointer(unit.Pound(1000)),
				FullWeight:  models.PoundPointer(unit.Pound(20000)),
			},
		},
	}, nil)

	return move, shipment
}

func createApprovedMoveWithPPMExcessWeightsAnd2WeightTickets(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	move, shipment := createApprovedMoveWithPPMExcessWeight(appCtx, userUploader,
		MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "excessweights2WTs@ppm.approved",
			SmID:        uuid.Must(uuid.NewV4()),
			FirstName:   "Two Weight Tickets",
			LastName:    "ExcessWeights",
			MoveID:      uuid.Must(uuid.NewV4()),
			MoveLocator: "XSWT02",
		})
	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.WeightTicket{
				EmptyWeight: models.PoundPointer(unit.Pound(1000)),
				FullWeight:  models.PoundPointer(unit.Pound(20000)),
			},
		},
	}, nil)
}

func createApprovedMoveWith2PPMShipmentsAndExcessWeights(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	move, _ := createApprovedMoveWithPPMExcessWeight(appCtx, userUploader,
		MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "excessweights2PPMs@ppm.approved",
			SmID:        uuid.Must(uuid.NewV4()),
			FirstName:   "Two PPMs",
			LastName:    "ExcessWeights",
			MoveID:      uuid.Must(uuid.NewV4()),
			MoveLocator: "XSWT03",
		})
	secondPPMShipment := factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)
	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    secondPPMShipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.WeightTicket{
				EmptyWeight: models.PoundPointer(unit.Pound(1000)),
				FullWeight:  models.PoundPointer(unit.Pound(20000)),
			},
		},
	}, nil)
}

func createApprovedMoveWithPPMAndHHGShipmentsAndExcessWeights(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	move, _ := createApprovedMoveWithPPMExcessWeight(appCtx, userUploader,
		MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "excessweightsPPMandHHG@ppm.approved",
			SmID:        uuid.Must(uuid.NewV4()),
			FirstName:   "PPM & HHG",
			LastName:    "ExcessWeights",
			MoveID:      uuid.Must(uuid.NewV4()),
			MoveLocator: "XSWT04",
		})
	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}
func createApprovedMoveWithAllShipmentTypesAndExcessWeights(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	move, _ := createApprovedMoveWithPPMExcessWeight(appCtx, userUploader,
		MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "excessweightsPPMandHHG@ppm.approved",
			SmID:        uuid.Must(uuid.NewV4()),
			FirstName:   "PPM & HHG",
			LastName:    "ExcessWeights",
			MoveID:      uuid.Must(uuid.NewV4()),
			MoveLocator: "XSWT05",
		})
	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createApprovedMoveWithPPMCloseoutComplete(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("f8af6fb0-101e-489c-9d9c-051931c52cf7"),
		Email:       "weightTicketPPM+closeout@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("cd4d7838-d8c1-441f-b7ce-af30b6257c3a"),
		FirstName:   "PPMCloseout",
		LastName:    "WeightTicket",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("eb6f09b4-0856-466c-b5e1-854310ccf486"),
		MoveLocator: "CLOSE0",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	approvedAdvanceStatus := models.PPMAdvanceStatusApproved
	allowableWeight := unit.Pound(4000)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("c0791087-9798-44e9-99df-59ae3ea9a71e"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("defb263e-bf01-4c67-85f5-b64ab54fd4fe"),
			ApprovedAt:                  &approvedAt,
			SubmittedAt:                 models.TimePointer(approvedAt.Add(7 * time.Hour * 24)),
			Status:                      models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			AdvanceStatus:               &approvedAdvanceStatus,
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
			AllowableWeight:             &allowableWeight,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func createApprovedMoveWithPPMCloseoutCompleteMultipleWeightTickets(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("385bb8f6-ee86-4948-b69d-615417bf71f9"),
		Email:       "weightTicketsPPM+closeout@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("7ad9b8d5-db20-4d00-946b-53531a24a9e1"),
		FirstName:   "PPMCloseout",
		LastName:    "WeightTickets",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("6c121a40-7037-46ba-9e94-1b63c598bcd9"),
		MoveLocator: "CLOSE1",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	approvedAdvanceStatus := models.PPMAdvanceStatusApproved
	allowableWeight := unit.Pound(8000)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("443750dc-def6-40ae-a60a-b6a5a4742c6b"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("08ab7a25-ef97-4134-bbb5-5be0e0de4734"),
			ApprovedAt:                  &approvedAt,
			SubmittedAt:                 models.TimePointer(approvedAt.Add(7 * time.Hour * 24)),
			Status:                      models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			AdvanceStatus:               &approvedAdvanceStatus,
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
			AllowableWeight:             &allowableWeight,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	factory.BuildWeightTicketWithConstructedWeight(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func createApprovedMoveWithPPMCloseoutCompleteWithExpenses(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("d0b0fafc-cedf-4821-8914-34a9fdea506d"),
		Email:       "expenses+closeout@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("df29b0c4-87e6-463a-a6ef-ac7d1523d9c8"),
		FirstName:   "PPMCloseout",
		LastName:    "Expenses",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("296d0c11-7b9d-4285-afd3-19c179b59508"),
		MoveLocator: "CLOSE3",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	approvedAdvanceStatus := models.PPMAdvanceStatusApproved
	allowableWeight := unit.Pound(4000)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("fce729a4-edce-45be-b91e-a70aa3cf09eb"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("645f9cd3-1aa2-4912-89fe-d0aa327226f6"),
			ApprovedAt:                  &approvedAt,
			SubmittedAt:                 models.TimePointer(approvedAt.Add(7 * time.Hour * 24)),
			Status:                      models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			AdvanceStatus:               &approvedAdvanceStatus,
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
			AllowableWeight:             &allowableWeight,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	storageType := models.MovingExpenseReceiptTypeStorage
	sitLocation := models.SITLocationTypeOrigin
	weightStored := 2000
	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
		{
			Model: models.MovingExpense{
				MovingExpenseType: &storageType,
				Description:       models.StringPointer("Storage R Us monthly rental unit"),
				SITStartDate:      models.TimePointer(time.Now()),
				SITEndDate:        models.TimePointer(time.Now().Add(30 * 24 * time.Hour)),
				SITLocation:       &sitLocation,
				WeightStored:      (*unit.Pound)(&weightStored),
			},
		},
	}, nil)
}

func createApprovedMoveWithPPMCloseoutCompleteWithAllDocTypes(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("d916c309-944b-4be6-b8ec-1ea59cffaf75"),
		Email:       "allPPMDocs+closeout@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("2da204ec-ada5-44c5-a1e1-39db1b027bdb"),
		FirstName:   "PPMCloseout",
		LastName:    "AllDocs",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("6abd318b-2eff-45d6-b282-73da0b65765d"),
		MoveLocator: "CLOSE2",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	approvedAdvanceStatus := models.PPMAdvanceStatusApproved
	allowableWeight := unit.Pound(4000)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("eb5a9e2b-cd16-4d84-8471-ccd869a589af"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("1a719536-02ba-44cd-b97d-5a0548237dc5"),
			ApprovedAt:                  &approvedAt,
			SubmittedAt:                 models.TimePointer(approvedAt.Add(7 * time.Hour * 24)),
			Status:                      models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			AdvanceStatus:               &approvedAdvanceStatus,
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
			AllowableWeight:             &allowableWeight,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicketWithConstructedWeight(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
				File:         factory.FixtureOpen("test.png"),
			},
		},
	}, nil)

	factory.BuildProgearWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    shipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
				File:         factory.FixtureOpen("test.jpg"),
			},
		},
	}, nil)

}

func createApprovedMoveWithPPMWithAboutFormComplete(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("88007896-6ae7-4600-866a-873d3bc67fd3"),
		Email:       "actualPPMDateZIPAdvanceDone@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("9d9f0509-b2fb-42a2-aab7-58dd4d79c4e7"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("acaa57ac-96f7-4411-aa07-c4bbe39e46bc"),
		MoveLocator: "ABTPPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("a742c4e9-24e3-4a97-995b-f355c6a14c04"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("f093a13b-4ab8-4545-b24c-eb44bf52e605"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWithAboutFormComplete2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("22dba194-3d9a-49c6-8328-718dd945292f"),
		Email:       "actualPPMDateZIPAdvanceDone2@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("c285f911-e432-42be-890a-965f9726b3e7"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("c20a62cb-ad19-405c-b230-dfadbd9a6eba"),
		MoveLocator: "AB2PPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("ef049132-204a-417a-a2c5-bfe2ac86e7a0"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("38f0b15a-efb9-411e-bd3d-c90514607fce"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWithAboutFormComplete3(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("9ec731d8-f347-4d34-8b54-4ce9e6ea3282"),
		Email:       "actualPPMDateZIPAdvanceDone3@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("5329c0c2-15f9-433e-9f99-7501eb68c6c1"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("a8dae89d-305a-49ae-996d-843dd7508aff"),
		MoveLocator: "AB3PPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("9cb9e177-c95c-49bf-80be-7f1b2ce41fe3"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("03d46a0d-6151-48dc-a8de-7abebd22916b"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWithAboutFormComplete4(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("2a0146c4-ec9a-4efc-a94c-6c2849c3e167"),
		Email:       "actualPPMDateZIPAdvanceDone4@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("98d28256-60e1-4792-86f1-c4e35cdef104"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("3bb2341a-9133-4d8e-abdf-0c0b18827756"),
		MoveLocator: "AB4PPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("b97b1852-aa00-4530-9ff3-a8bbcb35d928"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("379fb8f9-b210-4374-8f14-b8763be800ef"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWithAboutFormComplete5(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("bab42ae8-fe0d-4165-87be-dc1317ae0099"),
		Email:       "actualPPMDateZIPAdvanceDone5@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("71086cbf-89ee-4ca2-b063-739f3f33dab4"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("5e2916b2-dbba-4ca4-b558-d56842631757"),
		MoveLocator: "AB5PPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("91e9f99b-5923-47d3-bb80-211919ec27ce"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("c2fd7a80-afbe-425f-b7a9-bd26bd8cc965"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWithAboutFormComplete6(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("2c4eaae3-5226-456a-94d5-177c679b0656"),
		Email:       "actualPPMDateZIPAdvanceDone6@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("119f1167-fca9-4ca3-a2e9-57a033ba9dfb"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("59b67e7c-21a0-48c4-8630-c9afa206b3f2"),
		MoveLocator: "AB6PPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("5cb628ff-5670-4f16-933f-289e0c27ed25"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("2d0c3cb2-2c54-4ec0-b417-e81ab2ebd3c4"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWithAboutFormComplete7(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("c7cd77e8-74e8-4d7f-975c-d4ca18735561"),
		Email:       "actualPPMDateZIPAdvanceDone7@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("60cb3c60-68ef-47fa-b5f4-26d0e3d80e2a"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("8f451ef6-663f-49a9-b8ae-d3ecdca561d0"),
		MoveLocator: "AB7PPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("c46899a2-4c58-41b3-863c-347471ee26fc"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("59daf278-abf9-4ef1-9809-876df589890f"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMWithAboutFormComplete8(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("e5a06330-3f5c-4f50-82a6-46f1bd7dd3a6"),
		Email:       "actualPPMDateZIPAdvanceDone8@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("3719a811-83ce-4de2-b357-eb46181f0d80"),
		FirstName:   "ActualPPM",
		LastName:    "DateZIPAdvanceDone",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("6676c3cb-ad7a-4fa7-b6b2-c11c7754cad3"),
		MoveLocator: "AB8PPM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("1161ce38-441c-44ec-86fa-9e07e456cfb8"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("3faf26db-ddc4-4116-ab86-90a5e27106fd"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
}

func createApprovedMoveWithPPMMovingExpense(appCtx appcontext.AppContext, info *MoveCreatorInfo, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("146c2665-5b8a-4653-8434-9a4460de30b5"),
		Email:       "movingExpensePPM@ppm.approved",
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   "Expense",
		LastName:    "Complete",
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: "EXP3NS",
	}

	if info != nil {
		testdatagen.MergeModels(&moveInfo, *info)
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          uuid.Must(uuid.NewV4()),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)

	storageExpenseType := models.MovingExpenseReceiptTypeStorage
	sitLocation := models.SITLocationTypeOrigin
	weightStored := 2000
	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.MovingExpense{
				MovingExpenseType: &storageExpenseType,
				Description:       models.StringPointer("Storage R Us monthly rental unit"),
				SITStartDate:      models.TimePointer(time.Now()),
				SITEndDate:        models.TimePointer(time.Now().Add(30 * 24 * time.Hour)),
				SITLocation:       &sitLocation,
				WeightStored:      (*unit.Pound)(&weightStored),
			},
		},
	}, nil)
}

func createApprovedMoveWithPPMProgearWeightTicket(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("33eabbb6-416d-4d91-ba5b-bfd7d35e3037"),
		Email:       "progearWeightTicket@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("9240b1f4-352f-46b9-959a-4112ad4ae1a8"),
		FirstName:   "Progear",
		LastName:    "Complete",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("d933b7f2-41e9-4e9f-9b22-7afed753572b"),
		MoveLocator: "PR0G3R",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("bf119998-785a-4357-a3f1-5e71ee5bc757"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("9e671495-bf5a-48cf-b892-f4f3c4f1a18f"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildProgearWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
}

func createApprovedMoveWithPPMProgearWeightTicket2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID: testdatagen.ConvertUUIDStringToUUID("7d4dbc69-2973-4c8b-bf75-6fb582d7a5f6"),
		Email:  "progearWeightTicket2@ppm.approved",
		SmID:   testdatagen.ConvertUUIDStringToUUID("818f3076-78ef-4afe-abf8-62c490a9f6c4"),

		FirstName:   "Progear",
		LastName:    "Complete",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("d753eb23-b09f-4c53-b16d-fc71a56e5efd"),
		MoveLocator: "PR0G4R",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("22c401e6-91c8-48be-be8b-327326c71da4"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("24fd941f-8f27-43ad-ba68-9f6e3c181abe"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildProgearWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWithPPMShipmentReadyForFinalCloseout(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("1c842b03-fc2d-4e92-ade8-bd3e579196e0"),
		Email:       "readyForFinalComplete@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("5a21a8ed-52f5-446c-9d3e-5d8080765820"),
		FirstName:   "ReadyFor",
		LastName:    "PPMFinalCloseout",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("0b2e4341-583d-4793-b4a4-bd266534d17c"),
		MoveLocator: "PPMRFC",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	// Since we don't truncate the transportation_office table in our dev data generation workflow,
	// we need to generate an ID here instead of using a string to prevent duplicate entries.

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
			CloseoutOffice: &models.TransportationOffice{
				ID:   uuid.Must(uuid.NewV4()),
				Name: "Awesome base",
			},
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("226b81a7-9e56-4de2-b8ec-2cb5e8f72a35"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("6d1d9d00-2e5e-4830-a3c1-5c21c951e9c1"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
			FinalIncentive:              models.CentPointer(50000000),
		},
	}

	// This one is a little hairy because the move contains a
	// CloseoutOffice model
	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.WeightTicket{
				EmptyWeight: models.PoundPointer(14000),
				FullWeight:  models.PoundPointer(18000),
			},
		},
	}, nil)

	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.MovingExpense{
				Amount: models.CentPointer(45000),
			},
		},
	}, nil)

	factory.BuildProgearWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.ProgearWeightTicket{
				Weight: models.PoundPointer(1500),
			},
		},
	}, nil)
}

func createMoveWithPPMShipmentReadyForFinalCloseout2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("6f48be45-8ee0-4792-a961-ec6856e5435d"),
		Email:       "closeoutHappyPathWithEdits@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("0b17c7fe-24ae-4feb-a37a-154aa720867e"),
		FirstName:   "ReadyFor",
		LastName:    "PPMFinalCloseout",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("02c71fd2-a0dc-4975-bcd2-2b7edde22be1"),
		MoveLocator: "PPMCHE",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("f2bb3b05-e858-4717-966f-95e3f7054152"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("5d05071f-2042-40b0-a765-a17e95ec7959"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.WeightTicket{
				EmptyWeight: models.PoundPointer(14000),
				FullWeight:  models.PoundPointer(18000),
			},
		},
	}, nil)

	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.MovingExpense{
				Amount: models.CentPointer(45000),
			},
		},
	}, nil)

	factory.BuildProgearWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.ProgearWeightTicket{
				Weight: models.PoundPointer(1500),
			},
		},
	}, nil)
}

func createMoveWithPPMShipmentReadyForFinalCloseout3(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("917da44e-7e44-41be-b912-1486a72b69d8"),
		Email:       "closeoutHappyPathWithEditsMobile@ppm.approved",
		SmID:        testdatagen.ConvertUUIDStringToUUID("15fee9c1-626a-4e0e-a3fa-5409312ff955"),
		FirstName:   "ReadyFor",
		LastName:    "PPMFinalCloseout",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("15d39793-0ff1-4546-a48e-3de1fe157d95"),
		MoveLocator: "PPMCEM",
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("4da930b2-d227-4a0f-87b9-c09357e105d9"),
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          testdatagen.ConvertUUIDStringToUUID("15b3355f-8c7d-4a22-ac30-85aad77185ca"),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	factory.BuildWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.WeightTicket{
				EmptyWeight: models.PoundPointer(14000),
				FullWeight:  models.PoundPointer(18000),
			},
		},
	}, nil)

	factory.BuildMovingExpense(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.MovingExpense{
				Amount: models.CentPointer(45000),
			},
		},
	}, nil)

	factory.BuildProgearWeightTicket(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move.Orders.ServiceMember,
			LinkOnly: true,
		},
		{
			Model: models.ProgearWeightTicket{
				Weight: models.PoundPointer(1500),
			},
		},
	}, nil)
}

func createSubmittedMoveWithPPMShipment(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	/*
	 * A service member with orders and a full PPM Shipment.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("2d6a16ec-c031-42e2-aa55-90a1e29b961a"),
		Email:       "new@ppm.submitted",
		SmID:        testdatagen.ConvertUUIDStringToUUID("f1817ad8-dfd5-44c0-97eb-f634d22e147b"),
		FirstName:   "NewlySubmitted",
		LastName:    "PPM",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("5f30c363-07c0-4290-899c-3418e8472b44"),
		MoveLocator: "PPMSB1",
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		MTOShipment: models.MTOShipment{
			ID: testdatagen.ConvertUUIDStringToUUID("90ce4453-f836-4a76-a959-5f3271009f58"),
		},
		PPMShipment: models.PPMShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("2cc9fdcb-e1a6-4621-80dc-7cfa3956f2ea"),
			Status: models.PPMShipmentStatusDraft,
		},
	}

	move, _ := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, nil, assertions.PPMShipment)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)

	if err != nil {
		log.Panic(err)
	}

	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &move)

	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
}

func CreateMoveWithCloseOut(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveInfo MoveCreatorInfo, branch models.ServiceMemberAffiliation) models.Move {
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        moveInfo.UserID,
				OktaID:    oktaID.String(),
				OktaEmail: moveInfo.Email,
				Active:    true,
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
				CacValidated:  true,
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
				Status:      models.PPMShipmentStatusNeedsCloseout,
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

func createMoveWithCloseOutandNonCloseOut(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, locator string, branch models.ServiceMemberAffiliation) {
	userID := uuid.Must(uuid.NewV4())
	email := "1needscloseout@ppm.closeout"
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)
	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:     models.StringPointer("PPMSC"),
				LastName:      models.StringPointer("Submitted"),
				PersonalEmail: models.StringPointer(email),
				Affiliation:   &branch,
				CacValidated:  true,
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
				Locator:     locator,
				Status:      models.MoveStatusAPPROVED,
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
	if branch == models.AffiliationARMY || branch == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}

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

	mtoShipment2 := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
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
				Status: models.PPMShipmentStatusNeedsCloseout,
			},
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment2,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusWaitingOnCustomer,
			},
		},
	}, nil)

	factory.BuildSignedCertification(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWith2CloseOuts(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, locator string, branch models.ServiceMemberAffiliation) {
	userID := uuid.Must(uuid.NewV4())
	email := "2needcloseout@ppm.closeout"
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:     models.StringPointer("PPMSC"),
				LastName:      models.StringPointer("Submitted"),
				PersonalEmail: models.StringPointer(email),
				Affiliation:   &branch,
				CacValidated:  true,
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
				Locator:     locator,
				Status:      models.MoveStatusAPPROVED,
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

	if branch == models.AffiliationARMY || branch == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}

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

	mtoShipment2 := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
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
				Status: models.PPMShipmentStatusNeedsCloseout,
			},
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment2,
			LinkOnly: true,
		},
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsCloseout,
			},
		},
	}, nil)

	factory.BuildSignedCertification(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWithCloseOutandHHG(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, locator string, branch models.ServiceMemberAffiliation) {
	userID := uuid.Must(uuid.NewV4())
	email := "needscloseout@ppmHHG.closeout"
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:     models.StringPointer("PPMSC"),
				LastName:      models.StringPointer("Submitted"),
				PersonalEmail: models.StringPointer(email),
				Affiliation:   &branch,
				CacValidated:  true,
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
				Locator:     locator,
				Status:      models.MoveStatusAPPROVED,
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

	if branch == models.AffiliationARMY || branch == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}

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

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
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
				Status: models.PPMShipmentStatusNeedsCloseout,
			},
		},
	}, nil)

	factory.BuildSignedCertification(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func CreateMoveWithCloseoutOffice(appCtx appcontext.AppContext, moveInfo MoveCreatorInfo, userUploader *uploader.UserUploader) models.Move {
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Date(2020, time.December, 11, 12, 0, 0, 0, time.UTC)

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        moveInfo.UserID,
				OktaID:    oktaID.String(),
				OktaEmail: moveInfo.Email,
				Active:    true,
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
				CacValidated:  true,
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
				Status: models.PPMShipmentStatusNeedsCloseout,
			},
		},
	}, nil)

	return move
}

func createMovesForEachBranch(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	// Create a move for each branch
	branches := []models.ServiceMemberAffiliation{models.AffiliationARMY, models.AffiliationNAVY, models.AffiliationMARINES, models.AffiliationAIRFORCE, models.AffiliationCOASTGUARD}
	for _, branch := range branches {
		branchCode := strings.ToUpper(branch.String())[:3]
		moveInfo := MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "needscloseout@ppm.closeout",
			SmID:        uuid.Must(uuid.NewV4()),
			FirstName:   "PPMSC",
			LastName:    "Submitted",
			MoveLocator: "CO1" + branchCode,
			MoveID:      uuid.Must(uuid.NewV4()),
		}
		CreateMoveWithCloseOut(appCtx, userUploader, moveInfo, branch)
		locator := "CO2" + branchCode
		createMoveWithCloseOutandNonCloseOut(appCtx, userUploader, locator, branch)
		locator = "CO3" + branchCode
		createMoveWith2CloseOuts(appCtx, userUploader, locator, branch)
		locator = "CO4" + branchCode
		createMoveWithCloseOutandHHG(appCtx, userUploader, locator, branch)
	}
}

func CreateSubmittedMoveWithPPMShipmentForSC(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, _ services.MoveRouter, moveInfo MoveCreatorInfo) models.Move {
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        moveInfo.UserID,
				OktaID:    oktaID.String(),
				OktaEmail: moveInfo.Email,
				Active:    true,
			}},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            moveInfo.SmID,
				FirstName:     models.StringPointer(moveInfo.FirstName),
				LastName:      models.StringPointer(moveInfo.LastName),
				PersonalEmail: models.StringPointer(moveInfo.Email),
				CacValidated:  true,
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

func createSubmittedMoveWithPPMShipmentForSCWithSIT(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, _ services.MoveRouter, locator string) {
	userID := uuid.Must(uuid.NewV4())
	email := "completeWithSIT@ppm.submitted"
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()
	sitLocationType := models.SITLocationTypeOrigin

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:     models.StringPointer("PPMSC"),
				LastName:      models.StringPointer("Submitted with SIT"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
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
				Locator:          locator,
				Status:           models.MoveStatusNeedsServiceCounseling,
				SubmittedAt:      &submittedAt,
				CloseoutOfficeID: &DefaultCloseoutOfficeID,
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
				ID:                        testdatagen.ConvertUUIDStringToUUID("8158f06c-3cfa-4852-8984-c12de39da48f"),
				Status:                    models.PPMShipmentStatusSubmitted,
				SITExpected:               models.BoolPointer(true),
				SITEstimatedEntryDate:     models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
				SITEstimatedDepartureDate: models.TimePointer(time.Date(testdatagen.GHCTestYear, time.April, 16, 0, 0, 0, 0, time.UTC)),
				SITEstimatedWeight:        models.PoundPointer(unit.Pound(1234)),
				SITEstimatedCost:          models.CentPointer(unit.Cents(12345600)),
				SITLocation:               &sitLocationType,
			},
		},
	}, nil)

	factory.BuildSignedCertification(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

}

func createUnsubmittedMoveWithMultipleFullPPMShipmentComplete1(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and two full PPM Shipments.
	 */
	moveInfo := MoveCreatorInfo{
		UserID:           testdatagen.ConvertUUIDStringToUUID("afcc7029-4810-4f19-999a-2b254c659e19"),
		Email:            "multiComplete@ppm.unsubmitted",
		SmID:             testdatagen.ConvertUUIDStringToUUID("2dba3c65-1e69-429d-b797-0565014d0384"),
		FirstName:        "Multiple",
		LastName:         "Complete",
		MoveID:           testdatagen.ConvertUUIDStringToUUID("d94789bb-f8f7-4b5f-b86e-48503af70bfc"),
		MoveLocator:      "MULTI1",
		CloseoutOfficeID: &DefaultCloseoutOfficeID,
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("f5768bde-07c5-4765-b210-bcaf5f416009"),
			Status: models.PPMShipmentStatusDraft,
		},
	}

	move, _ := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createUnsubmittedMoveWithMultipleFullPPMShipmentComplete2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and two full PPM Shipments.
	 */

	moveInfo := MoveCreatorInfo{
		UserID:           testdatagen.ConvertUUIDStringToUUID("836d8363-1a5a-45b7-aee0-996a97724c24"),
		Email:            "multiComplete2@ppm.unsubmitted",
		SmID:             testdatagen.ConvertUUIDStringToUUID("bde2125f-63cf-4a4b-aff4-162a02120d89"),
		FirstName:        "Multiple2",
		LastName:         "Complete2",
		MoveID:           testdatagen.ConvertUUIDStringToUUID("839f893c-1c72-44e9-8544-298a19f1229a"),
		MoveLocator:      "MULTI2",
		CloseoutOfficeID: &DefaultCloseoutOfficeID,
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("aa677470-c7a5-4b97-b915-1b2d6a0ff58f"),
			Status: models.PPMShipmentStatusDraft,
		},
	}

	move, _ := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createSubmittedMoveWithFullPPMShipmentComplete(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {

	userID := uuid.Must(uuid.FromString("04f2a1c6-eb40-463d-8544-1909141fdedf"))
	email := "complete@ppm.submitted"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smWithPPM := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:     models.StringPointer("PPM"),
				LastName:      models.StringPointer("Submitted"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
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
				Locator: "PPMSUB",
				Status:  models.MoveStatusSUBMITTED,
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
}

func createMoveWithPPM(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	/*
	 * A service member with orders and a submitted move with a ppm
	 */
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("28837508-1942-4188-a7ef-a7b544309ea6"),
		Email:       "user@ppm",
		SmID:        testdatagen.ConvertUUIDStringToUUID("c29418e5-5d69-498d-9709-b493d5bbc814"),
		FirstName:   "Submitted",
		LastName:    "PPM",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("5174fd6c-3cab-4304-b4b3-89bd0f59b00b"),
		MoveLocator: "PPM001",
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		MTOShipment: models.MTOShipment{
			ID: testdatagen.ConvertUUIDStringToUUID("933d1c2b-5b90-4dfd-b363-5ff9a7e2b43a"),
		},
		PPMShipment: models.PPMShipment{
			ID:     testdatagen.ConvertUUIDStringToUUID("0914dfa2-6988-4a12-82b1-2586fb4aa8c7"),
			Status: models.PPMShipmentStatusDraft,
		},
	}

	move, _ := CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, nil, assertions.PPMShipment)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
}

func createMoveWithHHGMissingOrdersInfo(appCtx appcontext.AppContext, moveRouter services.MoveRouter, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	move := factory.BuildMoveWithShipment(db, []factory.Customization{
		{
			Model: models.Move{
				Locator: "REQINF",
				Status:  models.MoveStatusDRAFT,
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

	order := move.Orders
	order.TAC = nil
	order.OrdersNumber = nil
	order.DepartmentIndicator = nil
	order.OrdersTypeDetail = nil
	testdatagen.MustSave(db, &order)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	testdatagen.MustSave(db, &move)
}

func createUnsubmittedHHGMove(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/*
	 * A service member with an hhg only, unsubmitted move
	 */
	email := "hhg@only.unsubmitted"
	uuidStr := "f08146cf-4d6b-43d5-9ca5-c8d239d37b3e"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smWithHHGID := "1d06ab96-cb72-4013-b159-321d6d29c6eb"
	smWithHHG := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smWithHHGID),
				FirstName:     models.StringPointer("Unsubmitted"),
				LastName:      models.StringPointer("Hhg"),
				Edipi:         models.StringPointer("5833908165"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    smWithHHG,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("3a8c9f4f-7344-4f18-9ab5-0de3ef57b901"),
				Locator: "ONEHHG",
			},
		},
	}, nil)

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(2000)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("b67157bd-d2eb-47e2-94b6-3bc90f6fb8fe"),
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
}

func createUnsubmittedHHGMoveMultipleDestinations(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/*
		A service member with an un-submitted move that has an HHG shipment going to multiple destination addresses
	*/
	email := "multple-destinations@unsubmitted.hhg"
	userID := "81fe79a1-faaa-4735-8426-fd159e641002"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(userID)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smID := "af8f37bc-d29a-4a8a-90ac-5336a2a912b3"
	smWithHHG := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smID),
				FirstName:     models.StringPointer("Unsubmitted"),
				LastName:      models.StringPointer("Hhg"),
				Edipi:         models.StringPointer("5833908165"),
				PersonalEmail: &email,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    smWithHHG,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("c799098d-10f6-4e5a-9c88-a0de961e35b3"),
				Locator: "HHGSMA",
			},
		},
	}, nil)

	destinationAddress1 := factory.BuildAddress(db, nil, []factory.Trait{factory.GetTraitAddress3})
	destinationAddress2 := factory.BuildAddress(db, nil, []factory.Trait{factory.GetTraitAddress4})

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:           uuid.FromStringOrNil("fee1181f-22eb-452d-9252-292066e7b0a5"),
				ShipmentType: models.MTOShipmentTypeHHG,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    destinationAddress1,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:           uuid.FromStringOrNil("05361631-0e51-4a87-a8bc-82b3af120ce2"),
				ShipmentType: models.MTOShipmentTypeHHG,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    destinationAddress1,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
		{
			Model:    destinationAddress2,
			Type:     &factory.Addresses.SecondaryDeliveryAddress,
			LinkOnly: true,
		},
	}, nil)
}

func createUnsubmittedHHGMoveMultiplePickup(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/*
	 * A service member with an hhg only, unsubmitted move
	 */
	email := "hhg@multiple.pickup"
	uuidStr := "47fb0e80-6675-4ceb-b4eb-4f8e164c0f6e"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smWithHHGID := "92927bbd-5271-4a8c-b06b-fea07df84691"
	smWithHHG := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smWithHHGID),
				FirstName:     models.StringPointer("MultiplePickup"),
				LastName:      models.StringPointer("Hhg"),
				Edipi:         models.StringPointer("5833908165"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    smWithHHG,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("390341ca-2b76-4655-9555-161f4a0c9817"),
				Locator: "TWOPIC",
			},
		},
	}, nil)

	pickupAddress1 := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "1 First St",
				StreetAddress2: models.StringPointer("Apt 1"),
				StreetAddress3: models.StringPointer("Suite A"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
			},
		},
	}, nil)

	pickupAddress2 := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
			},
		},
	}, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:           uuid.FromStringOrNil("a35b1247-b4c2-48f6-9846-8e96050fbc95"),
				ShipmentType: models.MTOShipmentTypeHHG,
				ApprovedDate: models.TimePointer(time.Now()),
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    pickupAddress1,
			Type:     &factory.Addresses.PickupAddress,
			LinkOnly: true,
		},
		{
			Model:    pickupAddress2,
			Type:     &factory.Addresses.SecondaryPickupAddress,
			LinkOnly: true,
		},
	}, nil)
}

func createSubmittedHHGMoveMultiplePickupAmendedOrders(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	/*
	 * A service member with an hhg only, submitted move, with multiple addresses and amended orders
	 */
	email := "hhg@multiple.pickup.amendedOrders.submitted"
	uuidStr := "c5f202b3-90d3-46aa-8e3b-83e937fcca99"
	oktaID := uuid.Must(uuid.NewV4())

	smWithHHGID := "cfb9024b-39f3-47ca-b14b-a4e78a41e9db"

	orders := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smWithHHGID),
				FirstName:     models.StringPointer("MultiplePickup"),
				LastName:      models.StringPointer("Hhg"),
				Edipi:         models.StringPointer("5833908165"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
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

	orders = makeAmendedOrders(appCtx, orders, userUploader, &[]string{"medium.jpg", "small.pdf"})

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("e0463784-d5ea-4974-b526-f2a58c79ed07"),
				Locator: "AMENDO",
				Status:  models.MoveStatusSUBMITTED,
			},
		},
	}, nil)
	pickupAddress1 := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "1 First St",
				StreetAddress2: models.StringPointer("Apt 1"),
				StreetAddress3: models.StringPointer("Suite A"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
			},
		},
	}, nil)

	pickupAddress2 := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
			},
		},
	}, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:           uuid.FromStringOrNil("3c207b2a-d946-11eb-b8bc-0242ac130003"),
				ShipmentType: models.MTOShipmentTypeHHG,
				ApprovedDate: models.TimePointer(time.Now()),
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    pickupAddress1,
			Type:     &factory.Addresses.PickupAddress,
			LinkOnly: true,
		},
		{
			Model:    pickupAddress2,
			Type:     &factory.Addresses.SecondaryPickupAddress,
			LinkOnly: true,
		},
	}, nil)

}

func createMoveWithNTSAndNTSR(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter, opts sceneOptionsNTS) {
	db := appCtx.DB()

	email := fmt.Sprintf("nts.%s@nstr.%s", opts.shipmentMoveCode, opts.moveStatus)
	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)
	smWithNTS := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:     models.StringPointer(strings.ToTitle(string(opts.moveStatus))),
				LastName:      models.StringPointer("Nts&Nts-r"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	filterFile := &[]string{"150Kb.png"}
	orders := makeOrdersForServiceMember(appCtx, smWithNTS, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, opts.shipmentMoveCode, models.MoveStatusDRAFT)

	estimatedNTSWeight := unit.Pound(1400)
	actualNTSWeight := unit.Pound(2000)
	ntsShipment := factory.BuildNTSShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedNTSWeight,
				PrimeActualWeight:    &actualNTSWeight,
				Status:               models.MTOShipmentStatusSubmitted,
				UsesExternalVendor:   opts.usesExternalVendor,
			},
		},
	}, nil)

	factory.BuildMTOAgent(db, []factory.Customization{
		{
			Model:    ntsShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	ntsrShipment := factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedNTSWeight,
				PrimeActualWeight:    &actualNTSWeight,
				Status:               models.MTOShipmentStatusSubmitted,
				UsesExternalVendor:   opts.usesExternalVendor,
			},
		},
	}, nil)
	factory.BuildMTOAgent(db, []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				MTOAgentType: models.MTOAgentReceiving,
			},
		},
	}, nil)
	if opts.moveStatus == models.MoveStatusSUBMITTED {
		newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
		if err != nil {
			log.Panic(err)
		}

		verrs, err := models.SaveMoveDependencies(db, &move)
		if err != nil || verrs.HasAny() {
			log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
		}
	}
}

func createNTSMove(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	factory.BuildMoveWithShipment(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:    models.StringPointer("Spaceman"),
				LastName:     models.StringPointer("NTS"),
				CacValidated: true,
			},
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
			},
		},
	}, nil)
}

func createNTSRMove(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	factory.BuildMoveWithShipment(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:    models.StringPointer("Spaceman"),
				LastName:     models.StringPointer("NTS-release"),
				CacValidated: true,
			},
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
		},
	}, nil)
}

func createDefaultHHGMoveWithPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, affiliation models.ServiceMemberAffiliation) {
	createHHGMoveWithPaymentRequest(appCtx, userUploader, affiliation,
		models.Move{}, models.MTOShipment{})
}

// Creates an HHG Shipment with SIT at Origin and a payment request for first day and additional day SIT service items.
// This is to compare to calculating the cost for SIT with a PPM which excludes delivery/pickup costs because the
// address is not changing. 30 days of additional days in SIT are invoiced.
func createHHGWithOriginSITServiceItems(
	appCtx appcontext.AppContext,
	primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter,
	shipmentFetcher services.MTOShipmentFetcher,
) {
	db := appCtx.DB()
	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(db, []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)
	logger := appCtx.Logger()

	issueDate := time.Date(testdatagen.GHCTestYear, 3, 15, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(testdatagen.GHCTestYear, 8, 1, 0, 0, 0, 0, time.UTC)

	SITAllowance := 90
	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				RequestedPickupDate:  &issueDate,
				ActualPickupDate:     &issueDate,
				SITDaysAllowance:     &SITAllowance,
			},
		},
		{
			Model: models.Move{
				Locator: "ORGSIT",
			},
		},
		{
			Model: models.Order{
				IssueDate:    issueDate,
				ReportByDate: reportByDate,
			},
		},
		{
			Model: factory.BuildAddress(db, []factory.Customization{
				{
					Model: models.Address{
						City:       "Harlem",
						State:      "GA",
						PostalCode: "30813",
					},
				},
			}, nil),
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	move := shipment.MoveTaskOrder
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	submissionErr := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if submissionErr != nil {
		logger.Fatal(fmt.Sprintf("Error submitting move: %s", submissionErr))
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		logger.Fatal(fmt.Sprintf("Failed to save move and dependencies: %s", err))
	}
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	queryBuilder := query.NewQueryBuilder()
	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	handlerConfig := handlers.Config{}
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	mtoUpdater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, serviceItemCreator, moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator)
	_, approveErr := mtoUpdater.MakeAvailableToPrime(appCtx, move.ID, etag.GenerateEtag(move.UpdatedAt), true, true)

	if approveErr != nil {
		logger.Fatal("Error approving move")
	}

	// AvailableToPrimeAt is set to the current time when a move is approved, we need to update it to fall within the
	// same contract as the rest of the timestamps on our move for pricing to work.
	err = appCtx.DB().Find(&move, move.ID)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to fetch move: %s", err))
	}
	move.AvailableToPrimeAt = &May14GHCTestYear
	testdatagen.MustSave(appCtx.DB(), &move)

	// called for zip 3 domestic linehaul service item
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
		"90210", "30813").Return(2361, nil)

	shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)
	_, updateErr := shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, etag.GenerateEtag(shipment.UpdatedAt))
	if updateErr != nil {
		logger.Fatal("Error updating shipment status", zap.Error(updateErr))
	}

	// The SIT actual address will update the HHG shipment's pickup address, here we're providing the same value because
	// the prime API requires it to be specified.
	originSITAddress := shipment.PickupAddress
	originSITAddress.ID = uuid.Nil
	originSITAddress.Country = nil

	originSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		},
		{
			Model: *originSITAddress,
			Type:  &factory.Addresses.SITOriginHHGActualAddress,
		},
		{
			Model: models.MTOServiceItem{
				Reason:        models.StringPointer("Holiday break"),
				SITEntryDate:  &issueDate,
				SITPostalCode: &originSITAddress.PostalCode,
			},
		},
	}, nil)

	createdOriginServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &originSIT)
	if validErrs.HasAny() || createErr != nil {
		logger.Fatal(fmt.Sprintf("error while creating origin sit service item: %v", verrs.Errors), zap.Error(createErr))
	}
	addressCreator := address.NewAddressCreator()
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	serviceItemUpdator := mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer(), ghcrateengine.NewDomesticDestinationSITDeliveryPricer(), ghcrateengine.NewDomesticOriginSITFuelSurchargePricer())

	var originFirstDaySIT models.MTOServiceItem
	var originAdditionalDaySIT models.MTOServiceItem
	var originPickupSIT models.MTOServiceItem
	var originSITFSC models.MTOServiceItem
	for _, createdServiceItem := range *createdOriginServiceItems {
		switch createdServiceItem.ReService.Code {
		case models.ReServiceCodeDOFSIT:
			originFirstDaySIT = createdServiceItem
		case models.ReServiceCodeDOASIT:
			originAdditionalDaySIT = createdServiceItem
		case models.ReServiceCodeDOPSIT:
			originPickupSIT = createdServiceItem
		case models.ReServiceCodeDOSFSC:
			originSITFSC = createdServiceItem
		}
	}

	for _, createdServiceItem := range []models.MTOServiceItem{originFirstDaySIT, originAdditionalDaySIT, originPickupSIT, originSITFSC} {
		_, updateErr := serviceItemUpdator.ApproveOrRejectServiceItem(appCtx, createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error approving SIT service item", zap.Error(updateErr))
		}
	}

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		planner,
		ghcrateengine.NewServiceItemPricer(),
	)

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: move.ID,
	}

	var serviceItems []models.MTOServiceItem
	err = db.Eager("ReService").Where("move_id = ? AND id != ?", move.ID, originPickupSIT.ID).All(&serviceItems)
	if err != nil {
		log.Panic(err)
	}

	// additional days of SIT should exclude the initial entry day which excludes the first day of SIT
	// the prime can bill against the same addtional day SIT service item in 30 day increments per payment request
	doasitPaymentParams := []models.PaymentServiceItemParam{
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
			Value:       issueDate.Add(time.Hour * 24).Format("2006-01-02"),
		},
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
			Value:       issueDate.Add(time.Hour * 24 * 30).Format("2006-01-02"),
		}}

	paymentServiceItems := []models.PaymentServiceItem{}
	for _, serviceItem := range serviceItems {
		paymentItem := models.PaymentServiceItem{
			MTOServiceItemID: serviceItem.ID,
			MTOServiceItem:   serviceItem,
		}
		if serviceItem.ReService.Code == models.ReServiceCodeDOASIT {
			paymentItem.PaymentServiceItemParams = doasitPaymentParams
		}
		paymentServiceItems = append(paymentServiceItems, paymentItem)
	}

	paymentRequest.PaymentServiceItems = paymentServiceItems
	newPaymentRequest, createErr := paymentRequestCreator.CreatePaymentRequestCheck(appCtx, &paymentRequest)

	if createErr != nil {
		logger.Fatal("Error creating payment request", zap.Error(createErr))
	}

	factory.BuildPrimeUpload(db, []factory.Customization{
		{
			Model:    *newPaymentRequest,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{},
			ExtendedParams: &factory.PrimeUploadExtendedParams{
				PrimeUploader: primeUploader,
				AppContext:    appCtx,
			},
		},
	}, nil)
	posImage := factory.BuildProofOfServiceDoc(db, []factory.Customization{
		{
			Model:    *newPaymentRequest,
			LinkOnly: true,
		},
	}, nil)
	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("New payment request with service item params created with locator %s", move.Locator))
}

// Creates an HHG Shipment with SIT at Origin and a payment request for first day and additional day SIT service items.
// This is to compare to calculating the cost for SIT with a PPM which excludes delivery/pickup costs because the
// address is not changing. 30 days of additional days in SIT are invoiced.
func createHHGWithDestinationSITServiceItems(appCtx appcontext.AppContext, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter, shipmentFetcher services.MTOShipmentFetcher) {
	db := appCtx.DB()

	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(db, []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)

	logger := appCtx.Logger()

	issueDate := time.Date(testdatagen.GHCTestYear, 3, 15, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(testdatagen.GHCTestYear, 8, 1, 0, 0, 0, 0, time.UTC)
	SITAllowance := 90
	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				RequestedPickupDate:  &issueDate,
				ActualPickupDate:     &issueDate,
				SITDaysAllowance:     &SITAllowance,
			},
		},
		{
			Model: models.Move{
				Locator: "DSTSIT",
			},
		},
		{
			Model: models.Order{
				IssueDate:    issueDate,
				ReportByDate: reportByDate,
			},
		},
		{
			Model: factory.BuildAddress(db, []factory.Customization{
				{
					Model: models.Address{
						City:       "Harlem",
						State:      "GA",
						PostalCode: "30813",
					},
				},
			}, nil),
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	move := shipment.MoveTaskOrder
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	submissionErr := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if submissionErr != nil {
		logger.Fatal(fmt.Sprintf("Error submitting move: %s", submissionErr))
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		logger.Fatal(fmt.Sprintf("Failed to save move and dependencies: %s", err))
	}

	queryBuilder := query.NewQueryBuilder()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	//////////////////////////////////////////////////
	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	handlerConfig := handlers.Config{}
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	mtoUpdater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, serviceItemCreator, moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator)
	_, approveErr := mtoUpdater.MakeAvailableToPrime(appCtx, move.ID, etag.GenerateEtag(move.UpdatedAt), true, true)

	// AvailableToPrimeAt is set to the current time when a move is approved, we need to update it to fall within the
	// same contract as the rest of the timestamps on our move for pricing to work.
	err = appCtx.DB().Find(&move, move.ID)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to fetch move: %s", err))
	}
	move.AvailableToPrimeAt = &May14GHCTestYear
	testdatagen.MustSave(appCtx.DB(), &move)

	if approveErr != nil {
		logger.Fatal("Error approving move")
	}

	// called for zip 3 domestic linehaul service item
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
		"90210", "30813").Return(2361, nil)

	shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)
	_, updateErr := shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, etag.GenerateEtag(shipment.UpdatedAt))
	if updateErr != nil {
		logger.Fatal("Error updating shipment status", zap.Error(updateErr))
	}

	// The SIT actual address will update the HHG shipment's pickup address, here we're providing the same value because
	// the prime API requires it to be specified.
	originSITAddress := shipment.PickupAddress
	originSITAddress.ID = uuid.Nil

	destinationSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		},
		{
			Model: models.MTOServiceItem{
				Reason:       models.StringPointer("Holiday break"),
				SITEntryDate: &issueDate,
			},
		},
	}, nil)

	createdOriginServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &destinationSIT)
	if validErrs.HasAny() || createErr != nil {
		logger.Fatal(fmt.Sprintf("error while creating origin sit service item: %v", verrs.Errors), zap.Error(createErr))
	}

	addressCreator := address.NewAddressCreator()
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	serviceItemUpdator := mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer(), ghcrateengine.NewDomesticDestinationSITDeliveryPricer(), ghcrateengine.NewDomesticOriginSITFuelSurchargePricer())

	var destinationFirstDaySIT models.MTOServiceItem
	var destinationAdditionalDaySIT models.MTOServiceItem
	var destinationDeliverySIT models.MTOServiceItem
	var destinationSITFSC models.MTOServiceItem
	for _, createdServiceItem := range *createdOriginServiceItems {
		switch createdServiceItem.ReService.Code {
		case models.ReServiceCodeDDFSIT:
			destinationFirstDaySIT = createdServiceItem
		case models.ReServiceCodeDDASIT:
			destinationAdditionalDaySIT = createdServiceItem
		case models.ReServiceCodeDDDSIT:
			destinationDeliverySIT = createdServiceItem
		case models.ReServiceCodeDDSFSC:
			destinationSITFSC = createdServiceItem
		}
	}

	for _, createdServiceItem := range []models.MTOServiceItem{destinationFirstDaySIT, destinationAdditionalDaySIT, destinationDeliverySIT, destinationSITFSC} {
		_, updateErr := serviceItemUpdator.ApproveOrRejectServiceItem(appCtx, createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error approving SIT service item", zap.Error(updateErr))
		}
	}

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		planner,
		ghcrateengine.NewServiceItemPricer(),
	)

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: move.ID,
	}

	var serviceItems []models.MTOServiceItem
	err = db.Eager("ReService").Where("move_id = ? AND id != ?", move.ID, destinationDeliverySIT.ID).All(&serviceItems)
	if err != nil {
		log.Panic(err)
	}

	// additional days of SIT should exclude the initial entry day which excludes the first day of SIT
	// the prime can bill against the same addtional day SIT service item in 30 day increments

	ddasitPaymentParams := []models.PaymentServiceItemParam{
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
			Value:       issueDate.Add(time.Hour * 24).Format("2006-01-02"),
		},
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
			Value:       issueDate.Add(time.Hour * 24 * 30).Format("2006-01-02"),
		}}

	paymentServiceItems := []models.PaymentServiceItem{}
	for _, serviceItem := range serviceItems {
		paymentItem := models.PaymentServiceItem{
			MTOServiceItemID: serviceItem.ID,
			MTOServiceItem:   serviceItem,
		}
		if serviceItem.ReService.Code == models.ReServiceCodeDDASIT {
			paymentItem.PaymentServiceItemParams = ddasitPaymentParams
		}
		paymentServiceItems = append(paymentServiceItems, paymentItem)
	}

	paymentRequest.PaymentServiceItems = paymentServiceItems
	newPaymentRequest, createErr := paymentRequestCreator.CreatePaymentRequestCheck(appCtx, &paymentRequest)

	if createErr != nil {
		logger.Fatal("Error creating payment request", zap.Error(createErr))
	}

	factory.BuildPrimeUpload(db, []factory.Customization{
		{
			Model:    *newPaymentRequest,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{},
			ExtendedParams: &factory.PrimeUploadExtendedParams{
				PrimeUploader: primeUploader,
				AppContext:    appCtx,
			},
		},
	}, nil)

	posImage := factory.BuildProofOfServiceDoc(db, []factory.Customization{
		{
			Model:    *newPaymentRequest,
			LinkOnly: true,
		},
	}, nil)
	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("New payment request with service item params created with locator %s", move.Locator))
}

// Creates a payment request with domestic hhg and shorthaul shipments with
// service item pricing params for displaying cost calculations
func createHHGWithPaymentServiceItems(
	appCtx appcontext.AppContext,
	primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter,
	shipmentFetcher services.MTOShipmentFetcher,
) {
	db := appCtx.DB()
	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(db, []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)
	logger := appCtx.Logger()

	issueDate := time.Date(testdatagen.GHCTestYear, 3, 15, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(testdatagen.GHCTestYear, 8, 1, 0, 0, 0, 0, time.UTC)
	actualPickupDate := issueDate.Add(31 * 24 * time.Hour)
	SITAllowance := 90
	longhaulShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ActualPickupDate:     &actualPickupDate,
				SITDaysAllowance:     &SITAllowance,
			},
		},
		{
			Model: models.Move{
				Locator: "PARAMS",
			},
		},
		{
			Model: models.Order{
				IssueDate:    issueDate,
				ReportByDate: reportByDate,
			},
		},
	}, nil)

	move := longhaulShipment.MoveTaskOrder

	shorthaulDestinationAddress := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				PostalCode: "90211",
			},
		},
	}, nil)
	shorthaulShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				SITDaysAllowance:     &SITAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shorthaulDestinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	shipmentWithOriginalWeight := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shorthaulDestinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	shipmentWithOriginalAndReweighWeight := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shorthaulDestinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	reweighWeight := unit.Pound(100000)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: shipmentWithOriginalAndReweighWeight,
		Reweigh: models.Reweigh{
			Weight: &reweighWeight,
		},
	})

	shipmentWithOriginalAndReweighWeightReweihBolded := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shorthaulDestinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	// Make the reweigh weight and the estimated weight (original weight) be the same to create devseed
	// data where we can check that the reweigh weight is bolded.
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: shipmentWithOriginalAndReweighWeightReweihBolded,
		Reweigh: models.Reweigh{
			Weight: &estimatedWeight,
		},
	})

	billableWeightCap := unit.Pound(2000)
	billableWeightJustification := "Capped shipment"
	shipmentWithOriginalReweighAndAdjustedWeight := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:                      models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight:        &estimatedWeight,
				PrimeActualWeight:           &actualWeight,
				ShipmentType:                models.MTOShipmentTypeHHG,
				BillableWeightCap:           &billableWeightCap,
				BillableWeightJustification: &billableWeightJustification,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shorthaulDestinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: shipmentWithOriginalReweighAndAdjustedWeight,
		Reweigh: models.Reweigh{
			Weight: &reweighWeight,
		},
	})

	shipmentWithOriginalAndAdjustedWeight := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:                      models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight:        &estimatedWeight,
				PrimeActualWeight:           &actualWeight,
				ShipmentType:                models.MTOShipmentTypeHHG,
				BillableWeightCap:           &billableWeightCap,
				BillableWeightJustification: &billableWeightJustification,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shorthaulDestinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	submissionErr := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if submissionErr != nil {
		logger.Fatal(fmt.Sprintf("Error submitting move: %s", submissionErr))
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		logger.Fatal(fmt.Sprintf("Failed to save move and dependencies: %s", err))
	}

	queryBuilder := query.NewQueryBuilder()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(123, nil).Once()

	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	//////////////////////////////////////////////////
	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	handlerConfig := handlers.Config{}
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	mtoUpdater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, serviceItemCreator, moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator)
	_, approveErr := mtoUpdater.MakeAvailableToPrime(appCtx, move.ID, etag.GenerateEtag(move.UpdatedAt), true, true)

	// AvailableToPrimeAt is set to the current time when a move is approved, we need to update it to fall within the
	// same contract as the rest of the timestamps on our move for pricing to work.
	err = appCtx.DB().Find(&move, move.ID)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to fetch move: %s", err))
	}
	move.AvailableToPrimeAt = &May14GHCTestYear
	testdatagen.MustSave(appCtx.DB(), &move)

	if approveErr != nil {
		logger.Fatal("Error approving move")
	}
	// called using the addresses with origin zip of 90210 and destination zip of 94535
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(348, nil).Times(2)

	// called using the addresses with origin zip of 90210 and destination zip of 90211
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(3, nil).Times(5)

	// called for zip 3 domestic linehaul service item
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
		"94535", "94535").Return(348, nil).Times(2)

	// called for zip 5 domestic linehaul service item
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "94535", "94535").Return(348, nil).Once()

	// called for domestic shorthaul service item
	planner.On("Zip5TransitDistance", mock.AnythingOfType("*appcontext.appContext"),
		"90210", "90211").Return(3, nil).Times(7)

	// called for domestic shorthaul service item
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "90210", "90211").Return(348, nil).Times(10)

	// called for domestic origin SIT pickup service item
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "90210", "94535").Return(348, nil).Once()

	// called for domestic destination SIT delivery service item
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "94535", "90210").Return(348, nil).Times(2)

	// called for DLH, DSH, FSC service item estimated price calculations
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(400, nil).Times(3)

	for _, shipment := range []models.MTOShipment{longhaulShipment, shorthaulShipment, shipmentWithOriginalWeight, shipmentWithOriginalAndReweighWeight, shipmentWithOriginalAndReweighWeightReweihBolded, shipmentWithOriginalReweighAndAdjustedWeight, shipmentWithOriginalAndAdjustedWeight} {
		shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)
		_, updateErr := shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, etag.GenerateEtag(shipment.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error updating shipment status", zap.Error(updateErr))
		}
	}

	// There is a minimum of 29 days period for a sit service item that doesn't
	// have a departure date for the payment request param lookup to not encounter an error
	originEntryDate := actualPickupDate

	// Prep country with a real db
	country := factory.FetchOrBuildCountry(db, nil, nil)
	originSITAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
	// Manually set Country ID. Customizations will not work because DB is nil
	originSITAddress.CountryId = &country.ID
	originSITAddress.Country = nil
	originSITAddress.ID = uuid.Nil

	originSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    longhaulShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			}},
		{
			Model: originSITAddress,
			Type:  &factory.Addresses.SITOriginHHGActualAddress,
		},
		{
			Model: models.MTOServiceItem{
				Reason:        models.StringPointer("Holiday break"),
				SITEntryDate:  &originEntryDate,
				SITPostalCode: &originSITAddress.PostalCode,
				Status:        models.MTOServiceItemStatusRejected,
			},
		},
	}, nil)

	createdOriginServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &originSIT)
	if validErrs.HasAny() || createErr != nil {
		logger.Fatal(fmt.Sprintf("error while creating origin sit service item: %v", verrs.Errors), zap.Error(createErr))
	}

	destEntryDate := actualPickupDate
	destDepDate := actualPickupDate
	destSITAddress := factory.BuildAddress(db, nil, nil)
	destSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    longhaulShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			}},
		{
			Model:    destSITAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.SITDestinationFinalAddress,
		},
		{
			Model: models.MTOServiceItem{
				SITEntryDate:     &destEntryDate,
				SITDepartureDate: &destDepDate,
				SITPostalCode:    models.StringPointer("90210"),
			},
		},
	}, nil)

	createdDestServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &destSIT)
	if validErrs.HasAny() || createErr != nil {
		logger.Fatal(fmt.Sprintf("error while creating destination sit service item: %v", verrs.Errors), zap.Error(createErr))
	}

	addressCreator := address.NewAddressCreator()
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer(), ghcrateengine.NewDomesticDestinationSITDeliveryPricer(), ghcrateengine.NewDomesticOriginSITFuelSurchargePricer())

	var originFirstDaySIT models.MTOServiceItem
	var originAdditionalDaySIT models.MTOServiceItem
	var originPickupSIT models.MTOServiceItem
	var originSITFSC models.MTOServiceItem
	for _, createdServiceItem := range *createdOriginServiceItems {
		switch createdServiceItem.ReService.Code {
		case models.ReServiceCodeDOFSIT:
			originFirstDaySIT = createdServiceItem
		case models.ReServiceCodeDOASIT:
			originAdditionalDaySIT = createdServiceItem
		case models.ReServiceCodeDOPSIT:
			originPickupSIT = createdServiceItem
		case models.ReServiceCodeDOSFSC:
			originSITFSC = createdServiceItem
		}
	}

	originDepartureDate := originEntryDate.Add(15 * 24 * time.Hour)
	originPickupSIT.SITDepartureDate = &originDepartureDate

	for _, createdServiceItem := range []models.MTOServiceItem{originFirstDaySIT, originAdditionalDaySIT, originPickupSIT, originSITFSC} {
		_, updateErr := serviceItemUpdater.ApproveOrRejectServiceItem(appCtx, createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error approving SIT service item", zap.Error(updateErr))
		}
	}

	var serviceItemDDFSIT models.MTOServiceItem
	var serviceItemDDASIT models.MTOServiceItem
	var serviceItemDDDSIT models.MTOServiceItem
	var serviceItemDDSFSC models.MTOServiceItem
	for _, createdDestServiceItem := range *createdDestServiceItems {
		switch createdDestServiceItem.ReService.Code {
		case models.ReServiceCodeDDFSIT:
			serviceItemDDFSIT = createdDestServiceItem
		case models.ReServiceCodeDDASIT:
			serviceItemDDASIT = createdDestServiceItem
		case models.ReServiceCodeDDDSIT:
			serviceItemDDDSIT = createdDestServiceItem
		case models.ReServiceCodeDDSFSC:
			serviceItemDDSFSC = createdDestServiceItem
		}
	}

	destDepartureDate := destEntryDate.Add(15 * 24 * time.Hour)

	for _, createdServiceItem := range []models.MTOServiceItem{serviceItemDDASIT, serviceItemDDDSIT, serviceItemDDFSIT, serviceItemDDSFSC} {
		_, updateErr := serviceItemUpdater.ApproveOrRejectServiceItem(appCtx, createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error approving SIT service item", zap.Error(updateErr))
		}
	}

	description := "leg lamp"
	reason := "family heirloom extremely fragile"
	approvedAt := time.Now()
	itemDimension := models.MTOServiceItemDimension{
		Type:   models.DimensionTypeItem,
		Length: unit.ThousandthInches(2500),
		Height: unit.ThousandthInches(5000),
		Width:  unit.ThousandthInches(7500),
	}
	crateDimension := models.MTOServiceItemDimension{
		Type:   models.DimensionTypeCrate,
		Length: unit.ThousandthInches(30000),
		Height: unit.ThousandthInches(60000),
		Width:  unit.ThousandthInches(10000),
	}
	// cannot convert yet, has MTOServiceItemDimensions
	crating := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDCRT,
		},
		MTOServiceItem: models.MTOServiceItem{
			Status:      models.MTOServiceItemStatusApproved,
			Description: &description,
			Reason:      &reason,
			Dimensions: models.MTOServiceItemDimensions{
				itemDimension,
				crateDimension,
			},
			ApprovedAt: &approvedAt,
		},
		Move:        move,
		MTOShipment: longhaulShipment,
		Stub:        true,
	})

	// cannot convert yet, has MTOServiceItemDimensions
	uncrating := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDUCRT,
		},
		MTOServiceItem: models.MTOServiceItem{
			Description: &description,
			Reason:      &reason,
			Dimensions: models.MTOServiceItemDimensions{
				itemDimension,
				crateDimension,
			},
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: &approvedAt,
		},
		Move:        move,
		MTOShipment: longhaulShipment,
		Stub:        true,
	})

	cratingServiceItems := []models.MTOServiceItem{crating, uncrating}
	for index := range cratingServiceItems {
		_, _, cratingErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &cratingServiceItems[index])
		if cratingErr != nil {
			logger.Fatal("Error creating crating service item", zap.Error(cratingErr))
		}
	}

	shuttleDesc := "our smallest capacity shuttle vehicle"
	shuttleReason := "the bridge clearance was too low"
	estimatedShuttleWeigtht := unit.Pound(1000)
	actualShuttleWeight := unit.Pound(1500)
	originShuttle := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOSHUT,
			},
		},
		{
			Model: models.MTOServiceItem{
				Description:     &shuttleDesc,
				Reason:          &shuttleReason,
				EstimatedWeight: &estimatedShuttleWeigtht,
				ActualWeight:    &actualShuttleWeight,
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAt,
			},
		},
		{
			Model: move,

			LinkOnly: true,
		},
		{
			Model:    longhaulShipment,
			LinkOnly: true,
		},
	}, nil)

	destShuttle := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDSHUT,
			},
		},
		{
			Model: models.MTOServiceItem{
				Description:     &shuttleDesc,
				Reason:          &shuttleReason,
				EstimatedWeight: &estimatedShuttleWeigtht,
				ActualWeight:    &actualShuttleWeight,
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAt,
			},
		},
		{
			Model: move,

			LinkOnly: true,
		},
		{
			Model:    longhaulShipment,
			LinkOnly: true,
		},
	}, nil)

	shuttleServiceItems := []models.MTOServiceItem{originShuttle, destShuttle}
	for index := range shuttleServiceItems {
		_, _, shuttlingErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &shuttleServiceItems[index])
		if shuttlingErr != nil {
			logger.Fatal("Error creating shuttle service item", zap.Error(shuttlingErr))
		}
	}

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		planner,
		ghcrateengine.NewServiceItemPricer(),
	)

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: move.ID,
	}

	var serviceItems []models.MTOServiceItem
	err = db.Eager("ReService").Where("move_id = ?", move.ID).All(&serviceItems)
	if err != nil {
		log.Panic(err)
	}

	// An origin and destination SIT would normally not be on the same payment request so the TIO totals will appear
	// off.  Refer to the PARSIT move to see a reviewed and pending payment request with origin and destination SIT.
	doaPaymentStartDate := originEntryDate.Add(15 * 24 * time.Hour)
	doaPaymentEndDate := originDepartureDate.Add(15 * 24 * time.Hour)

	ddaPaymentStartDate := destEntryDate.Add(15 * 24 * time.Hour)
	daaPaymentEndDate := destDepartureDate.Add(15 * 24 * time.Hour)

	doasitPaymentParams := []models.PaymentServiceItemParam{
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
			Value:       doaPaymentStartDate.Format("2006-01-02"),
		},
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
			Value:       doaPaymentEndDate.Format("2006-01-02"),
		}}

	ddasitPaymentParams := []models.PaymentServiceItemParam{
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
			Value:       ddaPaymentStartDate.Format("2006-01-02"),
		},
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
			Value:       daaPaymentEndDate.Format("2006-01-02"),
		}}

	// Ordering the service items based on approved date to ensure the DDFSIT is after the DOASIT.
	// This avoids a flaky error when we create the service item parameters.
	sort.SliceStable(serviceItems, func(i, j int) bool {
		return serviceItems[i].ApprovedAt.String() < serviceItems[j].ApprovedAt.String()
	})
	paymentServiceItems := []models.PaymentServiceItem{}
	var serviceItemOrderString string
	for _, serviceItem := range serviceItems {
		serviceItemOrderString += serviceItem.ReService.Code.String()
		serviceItemOrderString += ", "
		paymentItem := models.PaymentServiceItem{
			MTOServiceItemID: serviceItem.ID,
			MTOServiceItem:   serviceItem,
		}
		if serviceItem.ReService.Code == models.ReServiceCodeDOASIT {
			paymentItem.PaymentServiceItemParams = doasitPaymentParams
		} else if serviceItem.ReService.Code == models.ReServiceCodeDDASIT {
			paymentItem.PaymentServiceItemParams = ddasitPaymentParams
		} // TODO: remove check once DOSFSC pricer is merged
		if serviceItem.ReService.Code != models.ReServiceCodeDOSFSC {
			paymentServiceItems = append(paymentServiceItems, paymentItem)
		}
	}

	logger.Debug(serviceItemOrderString)
	paymentRequest.PaymentServiceItems = paymentServiceItems
	newPaymentRequest, createErr := paymentRequestCreator.CreatePaymentRequestCheck(appCtx, &paymentRequest)

	if createErr != nil {
		logger.Fatal("Error creating payment request", zap.Error(createErr))
	}

	factory.BuildPrimeUpload(db, []factory.Customization{
		{
			Model:    *newPaymentRequest,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{},
			ExtendedParams: &factory.PrimeUploadExtendedParams{
				PrimeUploader: primeUploader,
				AppContext:    appCtx,
			},
		},
	}, nil)
	posImage := factory.BuildProofOfServiceDoc(db, []factory.Customization{
		{
			Model:    *newPaymentRequest,
			LinkOnly: true,
		},
	}, nil)

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("New payment request with service item params created with locator %s", move.Locator))
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
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: servicesCounseling,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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

func createHHGMoveWithPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, affiliation models.ServiceMemberAffiliation, moveTemplate models.Move, mtoShipmentTemplate models.MTOShipment) {
	db := appCtx.DB()
	logger := appCtx.Logger()
	serviceMember := models.ServiceMember{
		Affiliation:  &affiliation,
		CacValidated: true,
	}
	customer := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: serviceMember,
		},
	}, nil)

	orders := factory.BuildOrder(db, []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	moveTemplate.Status = models.MoveStatusAPPROVED
	moveTemplate.AvailableToPrimeAt = models.TimePointer(time.Now())
	mto := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveTemplate,
		},
		{
			Model:    orders,
			LinkOnly: true,
		},
	}, nil)

	shipmentPickupAddress := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				// This is a postal code that maps to the default office user gbloc LKNQ in the PostalCodeToGBLOC table
				PostalCode: "85325",
			},
		},
	}, nil)

	mtoShipmentTemplate.PrimeEstimatedWeight = &estimatedWeight
	mtoShipmentTemplate.PrimeActualWeight = &actualWeight
	mtoShipmentTemplate.ShipmentType = models.MTOShipmentTypeHHG
	mtoShipmentTemplate.ApprovedDate = models.TimePointer(time.Now())
	mtoShipmentTemplate.Status = models.MTOShipmentStatusSubmitted
	MTOShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: mtoShipmentTemplate,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    shipmentPickupAddress,
			Type:     &factory.Addresses.PickupAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(db, []factory.Customization{
		{
			Model: models.MTOAgent{
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	// setup service item
	serviceItem := testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mto,
		MTOShipment: MTOShipment,
	})

	planner := &routemocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(90210, nil)
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(910, nil)

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		planner,
		ghcrateengine.NewServiceItemPricer(),
	)

	paymentRequest := &models.PaymentRequest{
		IsFinal:         false,
		MoveTaskOrderID: mto.ID,
		PaymentServiceItems: []models.PaymentServiceItem{
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
				Status: models.PaymentServiceItemStatusRequested,
			},
		},
	}

	paymentRequest, paymentRequestErr := paymentRequestCreator.CreatePaymentRequestCheck(appCtx, paymentRequest)

	if paymentRequestErr != nil {
		logger.Fatal("error while creating payment request:", zap.Error(paymentRequestErr))
	}
	logger.Debug("create payment request ok: ", zap.Any("", paymentRequest))
}

func createHHGMoveWith10ServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	msCost := unit.Cents(10000)

	orders8 := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:           uuid.FromStringOrNil("9e8da3c7-ffe5-4f7f-b45a-8f01ccc56591"),
				CacValidated: true,
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("1d49bb07-d9dd-4308-934d-baad94f2de9b"),
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

	move8 := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders8,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("d4d95b22-2d9d-428b-9a11-284455aa87ba"),
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)
	mtoShipment8 := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
	}, nil)

	paymentRequest8 := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:      uuid.FromStringOrNil("154c9ebb-972f-4711-acb2-5911f52aced4"),
				IsFinal: false,
				Status:  models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
	}, nil)

	approvedAt := time.Now()
	serviceItemMS := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:         uuid.FromStringOrNil("4fba4249-b5aa-4c29-8448-66aa07ac8560"),
				Status:     models.MTOServiceItemStatusApproved,
				ApprovedAt: &approvedAt,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	csCost := unit.Cents(25000)
	serviceItemCS := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:         uuid.FromStringOrNil("e43c0df3-0dcd-4b70-adaa-46d669e094ad"),
				Status:     models.MTOServiceItemStatusApproved,
				ApprovedAt: &approvedAt,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	dlhCost := unit.Cents(99999)
	serviceItemDLH := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("9db1bf43-0964-44ff-8384-3297951f6781"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	fscCost := unit.Cents(55555)
	serviceItemFSC := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("b380f732-2fb2-49a0-8260-7a52ce223c59"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)

	dopCost := unit.Cents(3456)
	rejectionReason := "Customer no longer required this service"
	serviceItemDOP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("d886431c-c357-46b7-a084-a0c85dd496d4"),
				Status:          models.MTOServiceItemStatusRejected,
				RejectionReason: &rejectionReason,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dopCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOP,
			LinkOnly: true,
		},
	}, nil)

	ddpCost := unit.Cents(7890)
	serviceItemDDP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("551caa30-72fe-469a-b463-ad1f14780432"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddpCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDP,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("616dfdb5-52ec-436d-a570-a464c9dbd47a"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dpkCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDPK,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("1baeee0e-00d6-4d90-b22c-654c11d50d0f"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dupkCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDUPK,
			LinkOnly: true,
		},
	}, nil)

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	serviceItemDOFSIT := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.FromStringOrNil("61ce8a9b-5fcf-4d98-b192-a35f17819ae6"),
				PickupPostalCode: &dofsitPostal,
				Reason:           &dofsitReason,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("998beda7-e390-4a83-b15e-578a24326937"), // DOFSIT - Domestic Origin 1st Day SIT
			},
		},
	}, nil)

	dofsitCost := unit.Cents(8544)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dofsitCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOFSIT,
			LinkOnly: true,
		},
	}, nil)

	firstDeliveryDate := models.TimePointer(time.Now())
	dateOfContact := models.TimePointer(time.Now())
	customerContact1 := testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              dateOfContact.Add(time.Hour * 24),
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: *firstDeliveryDate,
		},
	})

	customerContact2 := testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			Type:                       models.CustomerContactTypeSecond,
			DateOfContact:              dateOfContact.Add(time.Hour * 48),
			TimeMilitary:               "1200Z",
			FirstAvailableDeliveryDate: firstDeliveryDate.Add(time.Hour * 24),
		},
	})
	serviceItemDDFSIT := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.FromStringOrNil("b2c770ab-db6f-465c-87f1-164ecd2f36a4"),
				CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact1, customerContact2},
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d0561c49-e1a9-40b8-a739-3e639a9d77af"), // DDFSIT - Domestic Destination 1st Day SIT
			},
		},
	}, nil)

	ddfsitCost := unit.Cents(8544)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddfsitCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDFSIT,
			LinkOnly: true,
		},
	}, nil)

	doshutCost := unit.Cents(623)
	serviceItemDOSHUT := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24034894b"),
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAt,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
			},
		},
	}, nil)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &doshutCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOSHUT,
			LinkOnly: true,
		},
	}, nil)

	serviceItemDDSHUT := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("2b0ce635-d71b-4000-a22a-7c098a3b6ae9"),
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAt,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("556663e3-675a-4b06-8da3-e4f1e9a9d3cd"), // DDSHUT - Dom Dest Shuttling
			},
		},
	}, nil)

	ddshutCost := unit.Cents(852)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddshutCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDSHUT,
			LinkOnly: true,
		},
	}, nil)

	testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("9b2b7cae-e8fa-4447-9a00-dcfc4ffc9b6f"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
	})
}

func createHHGMoveWith2PaymentRequests(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(db, []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)

	/* Customer with two payment requests */
	orders7 := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:           uuid.FromStringOrNil("4e6e4023-b089-4614-a65a-cac48027ffc2"),
				CacValidated: true,
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("f52f851e-91b8-4cb7-9f8a-6b0b8477ae2a"),
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

	mto7 := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders7,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("99783f4d-ee83-4fc9-8e0c-d32496bef32b"),
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Status:             models.MoveStatusAPPROVED,
			},
		},
	}, nil)

	mtoShipmentHHG7 := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("baa00811-2381-433e-8a96-2ced58e37a14"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(db, []factory.Customization{
		{
			Model:    mtoShipmentHHG7,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("82036387-a113-4b45-a172-94e49e4600d2"),
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	paymentRequest7 := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.FromStringOrNil("ea945ab7-099a-4819-82de-6968efe131dc"),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	// for soft deleted proof of service docs
	factory.BuildPrimeUpload(db, []factory.Customization{
		{
			Model:    paymentRequest7,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{
				ID: uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-1b557a4e425b"),
			},
		},
	}, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

	serviceItemMS7 := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-19edf785df69"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	msCost := unit.Cents(10000)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS7,
			LinkOnly: true,
		},
	}, nil)

	serviceItemDLH7 := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("aab8df9a-bbc9-4f26-a3ab-d5dcf1c8c40f"),
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	dlhCost := unit.Cents(99999)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH7,
			LinkOnly: true,
		},
	}, nil)

	additionalPaymentRequest7 := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.FromStringOrNil("540e2268-6899-4b67-828d-bb3b0331ecf2"),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	serviceItemCS7 := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-f9e646083cec"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	csCost := unit.Cents(25000)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    additionalPaymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS7,
			LinkOnly: true,
		},
	}, nil)

	MTOShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff9b5"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	serviceItemFSC7 := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-1d1c3ae66b15"),
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	fscCost := unit.Cents(55555)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    additionalPaymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC7,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWithHHGAndNTSRPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	msCost := unit.Cents(10000)

	hhgTAC := "1111"
	ntsTAC := "2222"
	hhgSAC := "3333"
	ntsSAC := "4444"

	orders := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.Order{
				TAC:    &hhgTAC,
				NtsTAC: &ntsTAC,
				SAC:    &hhgSAC,
				NtsSAC: &ntsSAC,
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

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.Must(uuid.NewV4()),
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Locator:            "HGNTSR",
			},
		},
	}, nil)
	// Create an HHG MTO Shipment
	pickupAddress := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
			},
		},
	}, nil)

	destinationAddress := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Princeton",
				State:          "NJ",
				PostalCode:     "08540",
			},
		},
	}, nil)

	hhgShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    destinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	lotNumber := "654321"

	storageFacility := factory.BuildStorageFacility(db, []factory.Customization{
		{
			Model: models.StorageFacility{
				Email:        models.StringPointer("old@email.com"),
				FacilityName: "Storage R Us",
				LotNumber:    &lotNumber,
			},
		},
		{
			Model: models.Address{
				StreetAddress1: "1234 Over Here Street",
				City:           "Houston",
				State:          "TX",
				PostalCode:     "77083",
			},
		},
	}, nil)

	tacType := models.LOATypeNTS
	sacType := models.LOATypeNTS

	serviceOrderNumber := "1234"

	// Create an NTSR MTO Shipment
	ntsrShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
				ApprovedDate:         models.TimePointer(time.Now()),
				ActualPickupDate:     models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
				TACType:              &tacType,
				SACType:              &sacType,
				ServiceOrderNumber:   &serviceOrderNumber,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(db, []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("e338e05c-6f5d-11ec-90d6-0242ac120003"),
				FirstName:    models.StringPointer("Receiving"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReceiving,
			},
		},
	}, nil)
	ntsrShipment.PickupAddressID = &pickupAddress.ID
	ntsrShipment.PickupAddress = &pickupAddress
	saveErr := db.Save(&ntsrShipment)
	if saveErr != nil {
		log.Panic("error saving NTSR shipment pickup address")
	}

	paymentRequest := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:      uuid.FromStringOrNil("3806be8d-ec39-43a2-a0ff-83b80bc4ba46"),
				IsFinal: false,
				Status:  models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	serviceItemMS := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:         uuid.Must(uuid.NewV4()),
				Status:     models.MTOServiceItemStatusApproved,
				ApprovedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	csCost := unit.Cents(25000)
	serviceItemCS := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:         uuid.Must(uuid.NewV4()),
				Status:     models.MTOServiceItemStatusApproved,
				ApprovedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	dlhCost := unit.Cents(99999)
	serviceItemDLH := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	serviceItemFSC := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	fscCost := unit.Cents(55555)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)

	serviceItemDOP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
			},
		},
	}, nil)

	dopCost := unit.Cents(3456)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dopCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOP,
			LinkOnly: true,
		},
	}, nil)

	ddpCost := unit.Cents(7890)
	serviceItemDDP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddpCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDP,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dpkCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDPK,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dupkCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDUPK,
			LinkOnly: true,
		},
	}, nil)

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	serviceItemDOFSIT := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.Must(uuid.NewV4()),
				Status:           models.MTOServiceItemStatusApproved,
				PickupPostalCode: &dofsitPostal,
				Reason:           &dofsitReason,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("998beda7-e390-4a83-b15e-578a24326937"), // DOFSIT - Domestic Origin 1st Day SIT
			},
		},
	}, nil)

	dofsitCost := unit.Cents(8544)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dofsitCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOFSIT,
			LinkOnly: true,
		},
	}, nil)

	customerContact1 := testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              time.Now().Add(time.Hour * 24),
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: time.Now(),
		},
	})

	customerContact2 := testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			Type:                       models.CustomerContactTypeSecond,
			DateOfContact:              time.Now().Add(time.Hour * 48),
			TimeMilitary:               "1200Z",
			FirstAvailableDeliveryDate: time.Now().Add(time.Hour * 24),
		},
	})

	serviceItemDDFSIT := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.Must(uuid.NewV4()),
				Status:           models.MTOServiceItemStatusApproved,
				CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact1, customerContact2},
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d0561c49-e1a9-40b8-a739-3e639a9d77af"), // DDFSIT - Domestic Destination 1st Day SIT
			},
		},
	}, nil)

	ddfsitCost := unit.Cents(8544)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddfsitCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDFSIT,
			LinkOnly: true,
		},
	}, nil)

	serviceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
	})

	dcrtCost := unit.Cents(55555)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemDLH := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemFSC := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemFSC,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemDOP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dopCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemDOP,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemDDP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddpCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemDDP,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemDUPK := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dupkCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemDUPK,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWithHHGAndNTSRMissingInfo(appCtx appcontext.AppContext, moveRouter services.MoveRouter, _ services.MTOShipmentFetcher) {
	db := appCtx.DB()
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.Move{
				Locator: "HNRMIS",
			},
		},
	}, nil)
	// original shipment that was previously approved and is now diverted
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	// new diverted shipment created by the Prime
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
}

func createMoveWithHHGAndNTSMissingInfo(appCtx appcontext.AppContext, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.Move{
				Locator: "HNTMIS",
			},
		},
	}, nil)
	// original shipment that was previously approved and is now diverted
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	// new diverted shipment created by the Prime
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
}

func createMoveWith2MinimalShipments(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.Move{
				Status:  models.MoveStatusSUBMITTED,
				Locator: "NOADDR",
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
	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	factory.BuildMTOShipmentMinimal(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:              models.MTOShipmentStatusSubmitted,
				RequestedPickupDate: &requestedPickupDate,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipmentMinimal(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:              models.MTOShipmentStatusSubmitted,
				RequestedPickupDate: &requestedPickupDate,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createApprovedMoveWithMinimalShipment(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	now := time.Now()
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				Locator:            "MISHIP",
				AvailableToPrimeAt: &now,
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
	factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		},
	}, nil)

	// requestedPickupDate := time.Now().AddDate(0, 3, 0)
	// requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	pickupAddress := factory.BuildAddress(db, nil, nil)

	shipmentFields := models.MTOShipment{
		Status: models.MTOShipmentStatusApproved,
		// RequestedPickupDate:   &requestedPickupDate,
		// RequestedDeliveryDate: &requestedDeliveryDate,
	}

	// Uncomment to create the shipment with an actual weight
	/*
		actualWeight := unit.Pound(999)
		shipmentFields.PrimeActualWeight = &actualWeight
	*/

	shipmentCustomizations := []factory.Customization{
		{
			Model: shipmentFields,
		},
		{
			Model:    pickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}

	// Uncomment to create the shipment with a destination address
	/*
		shipmentCustomizations = append(shipmentCustomizations, factory.Customization{
			Model:    factory.BuildAddress(appCtx.DB(), nil, []factory.Trait{factory.GetTraitAddress2}),
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		})
	*/

	firstShipment := factory.BuildMTOShipmentMinimal(db, shipmentCustomizations, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeFSC,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOP,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDP,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDPK,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDUPK,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWith2ShipmentsAndPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	msCost := unit.Cents(10000)

	orders := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
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
				ID:                 uuid.Must(uuid.NewV4()),
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Locator:            "REQSRV",
			},
		},
	}, nil)
	// Create an HHG MTO Shipment
	pickupAddress := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
			},
		},
	}, nil)

	destinationAddress := factory.BuildAddress(db, []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Princeton",
				State:          "NJ",
				PostalCode:     "08540",
			},
		},
	}, nil)

	hhgShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    pickupAddress,
			Type:     &factory.Addresses.PickupAddress,
			LinkOnly: true,
		},
		{
			Model:    destinationAddress,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	// Create an NTSR MTO Shipment
	ntsrShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.Must(uuid.NewV4()),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	ntsrShipment.PickupAddressID = &pickupAddress.ID
	ntsrShipment.PickupAddress = &pickupAddress
	saveErr := db.Save(&ntsrShipment)
	if saveErr != nil {
		log.Panic("error saving NTSR shipment pickup address")
	}

	paymentRequest := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:      uuid.FromStringOrNil("207216bf-0d60-4d91-957b-f0ddaeeb2dff"),
				IsFinal: false,
				Status:  models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	serviceItemMS := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:         uuid.Must(uuid.NewV4()),
				Status:     models.MTOServiceItemStatusApproved,
				ApprovedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	csCost := unit.Cents(25000)
	serviceItemCS := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:         uuid.Must(uuid.NewV4()),
				Status:     models.MTOServiceItemStatusApproved,
				ApprovedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	dlhCost := unit.Cents(99999)
	serviceItemDLH := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	serviceItemFSC := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	fscCost := unit.Cents(55555)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)

	serviceItemDOP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
			},
		},
	}, nil)

	dopCost := unit.Cents(3456)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dopCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOP,
			LinkOnly: true,
		},
	}, nil)

	ddpCost := unit.Cents(7890)
	serviceItemDDP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddpCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDP,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dpkCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDPK,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dupkCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDUPK,
			LinkOnly: true,
		},
	}, nil)

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.Must(uuid.NewV4()),
				Status:           models.MTOServiceItemStatusSubmitted,
				PickupPostalCode: &dofsitPostal,
				Reason:           &dofsitReason,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("998beda7-e390-4a83-b15e-578a24326937"), // DOFSIT - Domestic Origin 1st Day SIT
			},
		},
	}, nil)

	customerContact1 := testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              time.Now().Add(time.Hour * 24),
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: time.Now(),
		},
	})

	customerContact2 := testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeSecond,
			DateOfContact:              time.Now().Add(time.Hour * 48),
			TimeMilitary:               "1200Z",
			FirstAvailableDeliveryDate: time.Now().Add(time.Hour * 24),
		},
	})

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.Must(uuid.NewV4()),
				Status:           models.MTOServiceItemStatusSubmitted,
				CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact1, customerContact2},
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    hhgShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d0561c49-e1a9-40b8-a739-3e639a9d77af"), // DDFSIT - Domestic Destination 1st Day SIT
			},
		},
	}, nil)

	serviceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
	})

	dcrtCost := unit.Cents(55555)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemDLH := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemFSC := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemFSC,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemDOP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dopCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemDOP,
			LinkOnly: true,
		},
	}, nil)

	ntsrServiceItemDDP := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddpCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    ntsrServiceItemDDP,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.Must(uuid.NewV4()),
				Status: models.MTOServiceItemStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
			},
		},
	}, nil)
}

func createHHGMoveWith2PaymentRequestsReviewedAllRejectedServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(appCtx.DB(), []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)
	/* Customer with two payment requests */
	orders7 := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:           uuid.FromStringOrNil("4e6e4023-b089-4614-a65a-ffffffffffff"),
				CacValidated: true,
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("f52f851e-91b8-4cb7-9f8a-ffffffffffff"),
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

	locatorID := "PAYREJ"
	mto7 := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders7,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("99783f4d-ee83-4fc9-8e0c-ffffffffffff"),
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Status:             models.MoveStatusAPPROVED,
				Locator:            locatorID,
			},
		},
	}, nil)
	mtoShipmentHHG7 := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("baa00811-2381-433e-8a96-ffffffffffff"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(db, []factory.Customization{
		{
			Model:    mtoShipmentHHG7,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("82036387-a113-4b45-a172-ffffffffffff"),
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	reviewedDate := time.Now()
	paymentRequest7 := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:         uuid.FromStringOrNil("ea945ab7-099a-4819-82de-ffffffffffff"),
				IsFinal:    false,
				Status:     models.PaymentRequestStatusReviewedAllRejected,
				ReviewedAt: &reviewedDate,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	// for soft deleted proof of service docs
	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequest7,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{
				ID: uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-ffffffffffff"),
			},
		},
	}, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

	serviceItemMS7 := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-ffffffffffff"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	rejectionReason := "Just because."
	msCost := unit.Cents(10000)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents:      &msCost,
				Status:          models.PaymentServiceItemStatusDenied,
				RejectionReason: &rejectionReason,
			},
		}, {
			Model:    paymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS7,
			LinkOnly: true,
		},
	}, nil)

	serviceItemDLH7 := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("aab8df9a-bbc9-4f26-a3ab-ffffffffffff"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	dlhCost := unit.Cents(99999)
	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents:      &dlhCost,
				Status:          models.PaymentServiceItemStatusDenied,
				RejectionReason: &rejectionReason,
			},
		}, {
			Model:    paymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH7,
			LinkOnly: true,
		},
	}, nil)

	additionalPaymentRequest7 := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.FromStringOrNil("540e2268-6899-4b67-828d-ffffffffffff"),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				ReviewedAt:      &reviewedDate,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	serviceItemCS7 := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-ffffffffffff"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	csCost := unit.Cents(25000)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents:      &csCost,
				Status:          models.PaymentServiceItemStatusDenied,
				RejectionReason: &rejectionReason,
			},
		}, {
			Model:    additionalPaymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS7,
			LinkOnly: true,
		},
	}, nil)

	MTOShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-ffffffffffff"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG, // same as HHG for now
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
	}, nil)

	serviceItemFSC7 := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-ffffffffffff"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto7,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	fscCost := unit.Cents(55555)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents:      &fscCost,
				Status:          models.PaymentServiceItemStatusDenied,
				RejectionReason: &rejectionReason,
			},
		}, {
			Model:    additionalPaymentRequest7,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC7,
			LinkOnly: true,
		},
	}, nil)
}

func createTOO(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "too_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with too role */
	tooRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tooUUID := uuid.Must(uuid.FromString("dcf86235-53d3-43dd-8ee8-54212ae3078f"))
	oktaID := uuid.Must(uuid.NewV4())
	factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        tooUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole},
			}},
	}, nil)
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("144503a6-485c-463e-b943-d3c3bad11b09"),
				Email:  email,
				Active: true,
				UserID: &tooUUID,
			},
		},
	}, nil)
}

func createTIO(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "tio_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with tio role */
	tioRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	tioUUID := uuid.Must(uuid.FromString("3b2cc1b0-31a2-4d1b-874f-0591f9127374"))
	oktaID := uuid.Must(uuid.NewV4())
	factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        tioUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tioRole},
			}},
	}, nil)
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("f1828a35-43fd-42be-8b23-af4d9d51f0f3"),
				Email:  email,
				Active: true,
				UserID: &tioUUID,
			},
		},
	}, nil)
}

func createServicesCounselor(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "services_counselor_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with services counselor role */
	servicesCounselorRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeServicesCounselor).First(&servicesCounselorRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeServicesCounselor in the DB: %w", err))
	}

	servicesCounselorUUID := uuid.Must(uuid.FromString("a6c8663f-998f-4626-a978-ad60da2476ec"))
	oktaID := uuid.Must(uuid.NewV4())
	factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        servicesCounselorUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{servicesCounselorRole},
			}},
	}, nil)
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("c70d9a38-4bff-4d37-8dcc-456f317d7935"),
				Email:  email,
				Active: true,
				UserID: &servicesCounselorUUID,
			},
		},
	}, nil)
}

func createQae(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "qae_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with tio role */
	qaeRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeQae).First(&qaeRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeQae in the DB: %w", err))
	}

	qaeUUID := uuid.Must(uuid.FromString("8dbf1648-7527-4a92-b4eb-524edb703982"))
	oktaID := uuid.Must(uuid.NewV4())
	factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        qaeUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{qaeRole},
			}},
	}, nil)
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("ef4f6d1f-4ac3-4159-a364-5403e7d958ff"),
				Email:  email,
				Active: true,
				UserID: &qaeUUID,
			},
		},
	}, nil)
}

func createCustomerServiceRepresentative(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "customer_service_representative_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with RoleTypeCustomerServiceRepresentative role */
	customerServiceRepresentativeRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeCustomerServiceRepresentative).First(&customerServiceRepresentativeRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeCustomerServiceRepresentative in the DB: %w", err))
	}

	csrUUID := uuid.Must(uuid.FromString("72432922-BF2E-45DE-8837-1A458F5D1011"))
	oktaID := uuid.Must(uuid.NewV4())
	factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        csrUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{customerServiceRepresentativeRole},
			}},
	}, nil)
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("4B8C0AD8-337A-407A-9E49-074D466F837A"),
				Email:  email,
				Active: true,
				UserID: &csrUUID,
			},
		},
	}, nil)
}

func createTXO(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/* A user with both too and tio roles */
	email := "too_tio_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	tooRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	tooTioUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"))
	oktaID := uuid.Must(uuid.NewV4())
	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        tooTioUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole, tioRole},
			}},
	}, nil)
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-54212ae3078f"),
				Email:  email,
				Active: true,
				UserID: &tooTioUUID,
			},
		},
	}, nil)
	factory.BuildServiceMember(db, []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
}

func createTXOUSMC(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	emailUSMC := "too_tio_role_usmc@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", emailUSMC).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	// Makes user with both too and tio role with USMC gbloc
	tooRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	transportationOfficeUSMC := models.TransportationOffice{}
	err = db.Where("id = $1", "ccf50409-9d03-4cac-a931-580649f1647a").First(&transportationOfficeUSMC)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find transportation office USMC in the DB: %w", err))
	}

	// Makes user with both too and tio role with USMC gbloc

	tooTioWithUsmcUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-bbbbbbbbbbbb"))
	oktaWithUsmcID := uuid.Must(uuid.NewV4())
	factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        tooTioWithUsmcUUID,
				OktaID:    oktaWithUsmcID.String(),
				OktaEmail: emailUSMC,
				Active:    true,
				Roles:     []roles.Role{tooRole, tioRole},
			}},
	}, nil)
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-bbbbbbbbbbbb"),
				Email:  emailUSMC,
				Active: true,
				UserID: &tooTioWithUsmcUUID,
			},
		},
		{
			Model:    transportationOfficeUSMC,
			LinkOnly: true,
		},
	}, nil)
}

func createTXOServicesCounselor(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/* A user with both too, tio, and services counselor roles */
	email := "too_tio_services_counselor_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	officeUserRoleTypes := []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor}
	var userRoles roles.Roles
	err = db.Where("role_type IN (?)", officeUserRoleTypes).All(&userRoles)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find office user RoleType in the DB: %w", err))
	}

	tooTioServicesUUID := uuid.Must(uuid.FromString("8d78c849-0853-4eb8-a7a7-73055db7a6a8"))
	oktaID := uuid.Must(uuid.NewV4())

	// Make a user
	factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        tooTioServicesUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     userRoles,
			}},
	}, nil)

	// Make an office user associated with the previously created user
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("f3503012-e17a-4136-aa3c-508ee3b1962f"),
				Email:  email,
				Active: true,
				UserID: &tooTioServicesUUID,
			},
		},
	}, nil)
}

func createTXOServicesUSMCCounselor(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	emailUSMC := "too_tio_services_counselor_role_usmc@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", emailUSMC).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with both too, tio, and services counselor roles */
	officeUserRoleTypes := []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor}
	var userRoles roles.Roles
	err = db.Where("role_type IN (?)", officeUserRoleTypes).All(&userRoles)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find office user RoleType in the DB: %w", err))
	}

	// Makes user with too, tio, services counselor role with USMC gbloc
	transportationOfficeUSMC := models.TransportationOffice{}
	err = db.Where("id = $1", "ccf50409-9d03-4cac-a931-580649f1647a").First(&transportationOfficeUSMC)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find transportation office USMC in the DB: %w", err))
	}
	tooTioServicesWithUsmcUUID := uuid.Must(uuid.FromString("9aae1a83-6515-4c1d-84e8-f7b53dc3d5fc"))
	oktaWithUsmcID := uuid.Must(uuid.NewV4())

	// Makes a user with all office roles that is associated with USMC
	factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        tooTioServicesWithUsmcUUID,
				OktaID:    oktaWithUsmcID.String(),
				OktaEmail: emailUSMC,
				Active:    true,
				Roles:     userRoles,
			}},
	}, nil)

	// Makes an office user with the previously created user
	factory.BuildOfficeUser(db, []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("b23005d6-60ea-469f-91ab-a7daf4c686f5"),
				Email:  emailUSMC,
				Active: true,
				UserID: &tooTioServicesWithUsmcUUID,
			},
		},
		{
			Model:    transportationOfficeUSMC,
			LinkOnly: true,
		},
	}, nil)
}

func createServicesCounselorForCloseoutWithGbloc(appCtx appcontext.AppContext, userID uuid.UUID, email string, gbloc string) {
	db := appCtx.DB()
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	servicesCounselorRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeServicesCounselor).First(&servicesCounselorRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeServicesCounselor in the DB: %w", err))
	}

	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildOfficeUserWithRoles(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Email:  email,
				Active: true,
			},
		},
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.TransportationOffice{
				Gbloc: gbloc,
			},
		},
	}, []roles.RoleType{roles.RoleTypeServicesCounselor})
}

func createPrimeUser(appCtx appcontext.AppContext) models.User {
	db := appCtx.DB()
	/* A user with the prime role */

	var userRole roles.Role
	err := db.Where("role_type = (?)", roles.RoleTypePrime).First(&userRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find prime user RoleType in the DB: %w", err))
	}

	userUUID := uuid.Must(uuid.FromString("3ce06fa9-590a-48e5-9e30-6ad1e82b528c"))
	oktaID := uuid.Must(uuid.NewV4())
	email := "prime_role@office.mil"

	// Make a user
	primeUser := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        userUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
			},
		},
	}, []factory.Trait{factory.GetTraitPrimeUser})
	return primeUser
}

func createDevClientCertForUser(appCtx appcontext.AppContext, user models.User) {
	devlocalCert := factory.FetchOrBuildDevlocalClientCert(appCtx.DB())
	devlocalCert.UserID = user.ID
	testdatagen.MustSave(appCtx.DB(), &devlocalCert)
}

func createHHGMoveWithReweigh(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "REWAYD", models.MoveStatusAPPROVALSREQUESTED)
	move.TIORemarks = &tioRemarks
	testdatagen.MustSave(db, &move)
	reweighedWeight := unit.Pound(800)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		UserUploader: userUploader,
		MTOShipment: models.MTOShipment{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
		},
		Reweigh: models.Reweigh{
			Weight: &reweighedWeight,
		},
	})
	testdatagen.MakeReweigh(db, testdatagen.Assertions{UserUploader: userUploader})
}

func createHHGMoveWithBillableWeights(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	makeAmendedOrders(appCtx, orders, userUploader, &[]string{"small.pdf"})
	move := makeMoveForOrders(appCtx, orders, "BILWEI", models.MoveStatusAPPROVALSREQUESTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("6cd95b06-fef3-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

// creates a mix of shipments statuses with estimated, actual, and reweigh weights for testing the MTO page
func createReweighWithMixedShipmentStatuses(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "WTSTAT", models.MoveStatusAPPROVALSREQUESTED)

	// shipment is not yet approved so will be excluded from MTO weight calculations
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	divertedEstimated := unit.Pound(5000)
	divertedActual := unit.Pound(6000)
	// shipment was diverted so will have weights values already
	divertedShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				Diversion:            true,
				PrimeEstimatedWeight: &divertedEstimated,
				PrimeActualWeight:    &divertedActual,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	diveretedReweigh := unit.Pound(5500)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: divertedShipment,
		Reweigh: models.Reweigh{
			Weight: &diveretedReweigh,
		},
	})

	canceledEstimated := unit.Pound(5000)
	canceledActual := unit.Pound(6000)
	canceledReweigh := unit.Pound(5500)
	// cancelled shipment will still appear on MTO page but will not be included in weight calculations
	canceledShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusCanceled,
				PrimeEstimatedWeight: &canceledEstimated,
				PrimeActualWeight:    &canceledActual,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: canceledShipment,
		Reweigh: models.Reweigh{
			Weight: &canceledReweigh,
		},
	})

	approvedEstimated := unit.Pound(1000)
	approvedActual := unit.Pound(1500)
	approvedReweigh := unit.Pound(1250)
	approvedShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusApproved,
				PrimeEstimatedWeight: &approvedEstimated,
				PrimeActualWeight:    &approvedActual,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: approvedShipment,
		Reweigh: models.Reweigh{
			Weight: &approvedReweigh,
		},
	})

	approvedReweighRequestedEstimated := unit.Pound(1000)
	approvedReweighRequestedActual := unit.Pound(1500)
	approvedReweighRequestedShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusApproved,
				PrimeEstimatedWeight: &approvedReweighRequestedEstimated,
				PrimeActualWeight:    &approvedReweighRequestedActual,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: approvedReweighRequestedShipment,
	})

	divRequestedEstimated := unit.Pound(1000)
	divRequestedActual := unit.Pound(1500)
	divRequestedReweigh := unit.Pound(1750)
	divRequestedShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusDiversionRequested,
				PrimeEstimatedWeight: &divRequestedEstimated,
				PrimeActualWeight:    &divRequestedActual,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: divRequestedShipment,
		Reweigh: models.Reweigh{
			Weight: &divRequestedReweigh,
		},
	})

	cancellationRequestedEstimated := unit.Pound(1000)
	cancellationRequestedActual := unit.Pound(1500)
	cancellationRequestedReweigh := unit.Pound(1250)
	cancellationRequestedShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusCancellationRequested,
				PrimeEstimatedWeight: &cancellationRequestedEstimated,
				PrimeActualWeight:    &cancellationRequestedActual,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: cancellationRequestedShipment,
		Reweigh: models.Reweigh{
			Weight: &cancellationRequestedReweigh,
		},
	})
}

func createReweighWithMultipleShipments(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	serviceMember.FirstName = models.StringPointer("MultipleShipments")
	serviceMember.LastName = models.StringPointer("Reweighs")
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MULTRW", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(3000)
	now := time.Now()
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("5b72c64e-ffad-11eb-9a03-0242ac130003"),
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         &now,
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	shipmentWithMissingReweigh := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("6192766e-ffad-11eb-9a03-0242ac130003"),
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         &now,
				Status:               models.MTOShipmentStatusApproved,
				CounselorRemarks:     models.StringPointer("Please handle with care"),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweighWithNoWeightForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipmentWithMissingReweigh)

	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         &now,
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))

	shipmentForReweigh := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         &now,
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipmentForReweigh, unit.Pound(1541))
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("78a475d6-ffb8-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

func createReweighWithShipmentMissingReweigh(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MISHRW", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(6000)
	now := time.Now()
	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         &now,
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("4a1b0048-ffe7-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighWithNoWeightForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment)
}

func createReweighWithShipmentMaxBillableWeightExceeded(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MAXCED", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(8900)
	now := time.Now()
	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         &now,
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("096496b0-ffea-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(123456))
}

func createReweighWithShipmentNoEstimatedWeight(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "NOESTW", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks

	actualHHGWeight := unit.Pound(6000)
	now := time.Now()
	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeActualWeight: &actualHHGWeight,
				ShipmentType:      models.MTOShipmentTypeHHG,
				ApprovedDate:      &now,
				Status:            models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("c5c32296-0147-11ec-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

func createReweighWithShipmentDeprecatedPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	email := "deprecatedPaymentRequest@hhg.hhg"
	uuidStr := "6995a480-2e90-4d9b-90df-0f9b42277653"
	oktaID := uuid.Must(uuid.NewV4())
	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smID := "6c4074fe-ba11-471f-89f2-cf4f8c075377"
	sm := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smID),
				FirstName:     models.StringPointer("Deprecated"),
				LastName:      models.StringPointer("PaymentRequest"),
				Edipi:         models.StringPointer("6833908165"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    sm,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:         uuid.FromStringOrNil("bb0c2329-e225-41cc-a931-823c6026425b"),
				Locator:    "DEPPRQ",
				TIORemarks: &tioRemarks,
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
	actualHHGWeight := unit.Pound(6000)
	now := time.Now()
	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeActualWeight: &actualHHGWeight,
				ShipmentType:      models.MTOShipmentTypeHHG,
				ApprovedDate:      &now,
				Status:            models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	filterFile := &[]string{"150Kb.png"}
	paymentRequestID := uuid.Must(uuid.FromString("f80a07d3-0dcf-431f-b72c-dfd77e0483f6"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusDeprecated)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

func createReweighWithShipmentEDIErrorPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	email := "errrorPaymentRequest@hhg.hhg"
	uuidStr := "91252539-e8d0-4b9c-9722-d57c3b30bfb9"
	oktaID := uuid.Must(uuid.NewV4())
	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	smID := "8edb2121-3f7f-46f8-b8be-33ee60371369"
	sm := factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smID),
				FirstName:     models.StringPointer("Error"),
				LastName:      models.StringPointer("PaymentRequest"),
				Edipi:         models.StringPointer("6833908166"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model:    sm,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:         uuid.FromStringOrNil("18175273-1274-459e-b419-96450e49dafc"),
				Locator:    "ERRPRQ",
				TIORemarks: &tioRemarks,
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
	actualHHGWeight := unit.Pound(6000)
	now := time.Now()
	shipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeActualWeight: &actualHHGWeight,
				ShipmentType:      models.MTOShipmentTypeHHG,
				ApprovedDate:      &now,
				Status:            models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	err := moveRouter.Submit(appCtx, &move, &newSignedCertification)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	filterFile := &[]string{"150Kb.png"}
	paymentRequestID := uuid.Must(uuid.FromString("cc967c33-674e-4987-b4fc-b48624191c43"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusEDIError)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

func createHHGMoveWithTaskOrderServices(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {

	db := appCtx.DB()
	mtoWithTaskOrderServices := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("9c7b255c-2981-4bf8-839f-61c7458e2b4d"),
				Locator:            "RDY4PY",
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Status:             models.MoveStatusAPPROVED,
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
	estimated := unit.Pound(1400)
	actual := unit.Pound(1349)
	mtoShipment4 := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("c3a9e368-188b-4828-a64a-204da9b988c2"),
				RequestedPickupDate:  models.TimePointer(time.Now()),
				ScheduledPickupDate:  models.TimePointer(time.Now().AddDate(0, 0, -1)),
				PrimeEstimatedWeight: &estimated, // so we can price Dom. Destination Price
				PrimeActualWeight:    &actual,    // so we can price DLH
				Status:               models.MTOShipmentStatusApproved,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
	}, nil)
	mtoShipment5 := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("01b9671e-b268-4906-967b-ba661a1d3933"),
				RequestedPickupDate:  models.TimePointer(time.Now()),
				ScheduledPickupDate:  models.TimePointer(time.Now().AddDate(0, 0, -1)),
				PrimeEstimatedWeight: &estimated,
				PrimeActualWeight:    &actual,
				Status:               models.MTOShipmentStatusApproved,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("94bc8b44-fefe-469f-83a0-39b1e31116fb"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment4,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // Dom. Destination Price
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("fd6741a5-a92c-44d5-8303-1d7f5e60afbf"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment4,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("6431e3e2-4ee4-41b5-b226-393f9133eb6c"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment4,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("eee4b555-2475-4e67-a5b8-102f28d950f8"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment5,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("a6e5debc-9e73-421b-8f68-92936ce34737"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment5,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("999504a9-45b0-477f-a00b-3ede8ffde379"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment5,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK
			},
		},
	}, nil)

	factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("ca9aeb58-e5a9-44b0-abe8-81d233dbdebf"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("722a6f4e-b438-4655-88c7-51609056550d"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)
}

func createWebhookSubscriptionForPaymentRequestUpdate(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	// Create one webhook subscription for PaymentRequestUpdate
	testdatagen.MakeWebhookSubscription(db, testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			CallbackURL: "https://primelocal:9443/support/v1/webhook-notify",
		},
	})
}

func createMoveWithServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	orders9 := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("796a0acd-1ccb-4a2f-a9b3-e44906ced698"),
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

	move9 := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders9,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:     uuid.FromStringOrNil("7cbe57ba-fd3a-45a7-aa9a-1970f1908ae7"),
				Status: models.MoveStatusAPPROVED,
			},
		},
	}, nil)
	mtoShipment9 := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("ec3f4edf-1463-43fb-98c4-272d3acb204a"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move9,
			LinkOnly: true,
		},
	}, nil)

	paymentRequest9 := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:         uuid.FromStringOrNil("cfd110d4-1f62-401c-a92c-39987a0b4228"),
				IsFinal:    false,
				Status:     models.PaymentRequestStatusReviewed,
				ReviewedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move9,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
		}, {
			Model:    paymentRequest9,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusDenied,
			},
		}, {
			Model:    paymentRequest9,
			LinkOnly: true,
		},
	}, nil)

	customizations9 := []factory.Customization{
		{
			Model:    move9,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment9,
			LinkOnly: true,
		},
		{
			Model:    paymentRequest9,
			LinkOnly: true,
		},
	}

	currentTime := time.Now()
	const testDateFormat = "060102"

	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(testDateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "24246",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		db,
		models.ReServiceCodeDLH,
		basicPaymentServiceItemParams,
		customizations9, nil,
	)
}

func createMoveWithBasicServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	orders10 := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("796a0acd-1ccb-4a2f-a9b3-e44906ced699"),
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

	move10 := factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders10,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("7cbe57ba-fd3a-45a7-aa9a-1970f1908ae8"),
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move10,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	paymentRequest10 := factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:         uuid.FromStringOrNil("cfd110d4-1f62-401c-a92c-39987a0b4229"),
				Status:     models.PaymentRequestStatusReviewed,
				ReviewedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move10,
			LinkOnly: true,
		},
	}, nil)

	serviceItemA := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	serviceItemB := factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
		}, {
			Model:    serviceItemA,
			LinkOnly: true,
		}, {
			Model:    paymentRequest10,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentServiceItem(db, []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusDenied,
			},
		}, {
			Model:    serviceItemB,
			LinkOnly: true,
		}, {
			Model:    paymentRequest10,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWithUniqueDestinationAddress(appCtx appcontext.AppContext) {
	db := appCtx.DB()

	order := factory.BuildOrder(db, []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Apt 2"),
				StreetAddress3: models.StringPointer("Suite B"),
				City:           "Columbia",
				State:          "SC",
				PostalCode:     "29212",
			},
			Type: &factory.Addresses.DutyLocationAddress,
		},
		{
			Model: models.Order{
				OrdersNumber: models.StringPointer("ORDER3"),
				TAC:          models.StringPointer("F8E1"),
			},
		},
	}, nil)

	factory.BuildMove(db, []factory.Customization{
		{
			Model:    order,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("ecbc2e6a-1b45-403b-9bd4-ea315d4d3d93"),
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Status:             models.MoveStatusAPPROVED,
			},
		},
	}, nil)
}

/*
Generic helper function that lets you create a move with any staus and with any shipment type
*/
func CreateMoveWithMTOShipment(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType, shipmentType models.MTOShipmentType, destinationType *models.DestinationType, locator string, moveStatus models.MoveStatus) models.Move {
	if shipmentType == models.MTOShipmentTypeBoatHaulAway || shipmentType == models.MTOShipmentTypeBoatTowAway { // Add boat specific fields in relevant PR
		log.Panic(fmt.Errorf("Unable to generate random integer for submitted move date"), zap.Error(errors.New("Not yet implemented")))
	}

	db := appCtx.DB()
	submittedAt := time.Now()
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	ordersNumber := "8675309"
	departmentIndicator := "ARMY"
	tac := "E19A"
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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
				Status:      moveStatus,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := factory.BuildAddress(db, nil, nil)

	if destinationType != nil { // Destination type is only used for retirement moves
		retirementMTOShipment := factory.BuildMTOShipment(db, []factory.Customization{
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

		if shipmentType == models.MTOShipmentTypeMobileHome {
			factory.BuildMobileHomeShipment(appCtx.DB(), []factory.Customization{
				{
					Model: models.MobileHome{
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(300),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: models.IntPointer(72),
					},
				},
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    retirementMTOShipment,
					LinkOnly: true,
				},
			}, nil)
		}
	}

	requestedPickupDate = submittedAt.Add(30 * 24 * time.Hour)
	requestedDeliveryDate = requestedPickupDate.Add(7 * 24 * time.Hour)
	regularMTOShipment := factory.BuildMTOShipment(db, []factory.Customization{
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

	if shipmentType == models.MTOShipmentTypeMobileHome {
		factory.BuildMobileHomeShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.MobileHome{
					Year:           models.IntPointer(2000),
					Make:           models.StringPointer("Boat Make"),
					Model:          models.StringPointer("Boat Model"),
					LengthInInches: models.IntPointer(300),
					WidthInInches:  models.IntPointer(108),
					HeightInInches: models.IntPointer(72),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    regularMTOShipment,
				LinkOnly: true,
			},
		}, nil)
	}

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
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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
	retirementMTOShipment := factory.BuildMTOShipment(db, []factory.Customization{
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

	if shipmentType == models.MTOShipmentTypeMobileHome {
		factory.BuildMobileHomeShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.MobileHome{
					Year:           models.IntPointer(2000),
					Make:           models.StringPointer("Boat Make"),
					Model:          models.StringPointer("Boat Model"),
					LengthInInches: models.IntPointer(300),
					WidthInInches:  models.IntPointer(108),
					HeightInInches: models.IntPointer(72),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    retirementMTOShipment,
				LinkOnly: true,
			},
		}, nil)
	}

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

func CreateNeedsServicesCounselingWithAmendedOrders(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, ordersType internalmessages.OrdersType, shipmentType models.MTOShipmentType, destinationType *models.DestinationType, locator string) models.Move {
	db := appCtx.DB()
	submittedAt := time.Now()
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	ordersNumber := "8675309"
	departmentIndicator := "ARMY"
	tac := "E19A"
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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
	orders = makeAmendedOrders(appCtx, orders, userUploader, &[]string{"medium.jpg", "small.pdf"})
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
	return move
}

/*
Create Needs Service Counseling without all required order information
*/
func createNeedsServicesCounselingWithoutCompletedOrders(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType, shipmentType models.MTOShipmentType, destinationType *models.DestinationType, locator string) {
	db := appCtx.DB()
	submittedAt := time.Now()
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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
	if *orders.ServiceMember.Affiliation == models.AffiliationARMY || *orders.ServiceMember.Affiliation == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}

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
}

func createUserWithLocatorAndDODID(appCtx appcontext.AppContext, locator string, dodID string) {
	db := appCtx.DB()
	submittedAt := time.Now()
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.ServiceMember{
				Edipi:        models.StringPointer(dodID),
				FirstName:    models.StringPointer("QAETestFirst"),
				CacValidated: true,
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

}

func createNeedsServicesCounselingSingleHHG(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType, locator string) {
	db := appCtx.DB()
	submittedAt := time.Now()
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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

}

func CreateNeedsServicesCounselingMinimalNTSR(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType, locator string) models.Move {
	db := appCtx.DB()
	submittedAt := time.Now()
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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

func createHHGNeedsServicesCounselingUSMC(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	marineCorps := models.AffiliationMARINES
	submittedAt := time.Now()

	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Move{
				Locator:     "USMCSS",
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
		{
			Model: models.ServiceMember{
				Affiliation:  &marineCorps,
				LastName:     models.StringPointer("Marine"),
				FirstName:    models.StringPointer("Ted"),
				CacValidated: true,
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
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
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
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
	}, nil)
}

func createHHGNeedsServicesCounselingUSMC2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	marineCorps := models.AffiliationMARINES
	submittedAt := time.Now()

	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Move{
				Locator:     "USMCSC",
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
		{
			Model: models.ServiceMember{
				Affiliation:  &marineCorps,
				LastName:     models.StringPointer("Marine"),
				FirstName:    models.StringPointer("Barbara"),
				CacValidated: true,
			},
		},
		{
			Model: models.TransportationOffice{
				Gbloc: "ZANY",
			},
			Type: &factory.TransportationOffices.CloseoutOffice,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
	requestedPickupDate := submittedAt.Add(20 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(14 * 24 * time.Hour)
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
	}, nil)

}

func CreateHHGNeedsServicesCounselingUSMC3(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, locator string) models.Move {
	db := appCtx.DB()

	marineCorps := models.AffiliationMARINES
	submittedAt := time.Now()

	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Move{
				Locator:     locator,
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
		{
			Model: models.ServiceMember{
				Affiliation: &marineCorps,
				LastName:    models.StringPointer("Marine"),
				FirstName:   models.StringPointer("Ted"),
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
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
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
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
	}, nil)
	return move
}

func createHHGServicesCounselingCompleted(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	servicesCounselingCompletedAt := time.Now()
	submittedAt := servicesCounselingCompletedAt.Add(-7 * 24 * time.Hour)
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Move{
				Locator:                      "CSLCMP",
				Status:                       models.MoveStatusServiceCounselingCompleted,
				SubmittedAt:                  &submittedAt,
				ServiceCounselingCompletedAt: &servicesCounselingCompletedAt,
			},
		},
	}, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
}

func createHHGNoShipments(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	submittedAt := time.Now()
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	factory.BuildMove(db, []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:     "NOSHIP",
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)
}

func createHHGMoveWithMultipleOrdersFiles(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"2mb.png", "150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MULTOR", models.MoveStatusSUBMITTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("aca5cc9c-c266-4a7d-895d-dc3c9c0d9894"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createHHGMoveWithAmendedOrders(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"2mb.png", "150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	makeAmendedOrders(appCtx, orders, userUploader, &[]string{"medium.jpg", "small.pdf"})
	move := makeMoveForOrders(appCtx, orders, "AMDORD", models.MoveStatusAPPROVALSREQUESTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("c47999c4-afa8-4c87-8a0e-7763b4e5d4c5"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createHHGMoveWithRiskOfExcess(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"2mb.png", "150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:                  models.MoveStatusAPPROVALSREQUESTED,
				Locator:                 "RISKEX",
				AvailableToPrimeAt:      &now,
				ExcessWeightQualifiedAt: &now,
			},
		},
	}, nil)
	shipment := makeRiskOfExcessShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("50b35add-705a-468b-8bad-056f5d9ef7e1"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createMoveWithDivertedShipments(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				Locator:            "DVRS0N",
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)
	// original shipment that was previously approved and is now diverted
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:       models.MTOShipmentStatusSubmitted,
				ApprovedDate: models.TimePointer(time.Now()),
				Diversion:    true,
			},
		},
	}, nil)
	// new diverted shipment created by the Prime
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:    models.MTOShipmentStatusSubmitted,
				Diversion: true,
			},
		},
	}, nil)
}

func createMoveWithSITExtensionHistory(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "SITEXT", models.MoveStatusAPPROVALSREQUESTED)

	// manually calculated SIT days including SIT extension approved days
	sitDaysAllowance := 270
	mtoShipmentSIT := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &sitDaysAllowance,
			},
		},
	}, nil)

	year, month, day := time.Now().Add(time.Hour * 24 * -60).Date()

	threeMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	twoMonthsAgo := threeMonthsAgo.Add(time.Hour * 24 * 30)
	postalCode := "90210"
	reason := "peak season all trucks in use"

	// This will in practice not exist without DOFSIT and DOASIT
	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &threeMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &threeMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:           models.MTOServiceItemStatusApproved,
				SITEntryDate:     &threeMonthsAgo,
				SITDepartureDate: &twoMonthsAgo,
				SITPostalCode:    &postalCode,
				Reason:           &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &twoMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &twoMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDASIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &twoMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	MakeSITExtensionsForShipment(appCtx, mtoShipmentSIT)

	factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:         uuid.Must(uuid.NewV4()),
				Status:     models.PaymentRequestStatusReviewed,
				ReviewedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

}

func createMoveWithFutureSIT(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "SITFUT", models.MoveStatusAPPROVALSREQUESTED)

	// manually calculated SIT days including SIT extension approved days
	sitDaysAllowance := 270
	mtoShipmentSIT := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &sitDaysAllowance,
			},
		},
	}, nil)

	year, month, day := time.Now().Add(time.Hour * 24 * 90).Date()

	threeMonthsFromNow := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	postalCode := "90210"
	reason := "peak season all trucks in use"

	// This will in practice not exist without DOFSIT and DOASIT
	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &threeMonthsFromNow,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &threeMonthsFromNow,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &threeMonthsFromNow,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		},
		{
			Model:    mtoShipmentSIT,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentRequest(db, []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:         uuid.Must(uuid.NewV4()),
				Status:     models.PaymentRequestStatusReviewed,
				ReviewedAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

}

func createMoveWithOriginAndDestinationSIT(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveLocator string) models.MTOServiceItem {
	db := appCtx.DB()

	sitDaysAllowance := 90
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.Move{
				ID:                 uuid.Must(uuid.NewV4()),
				Locator:            moveLocator,
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.Entitlement{
				StorageInTransit: &sitDaysAllowance,
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
	factory.BuildMTOServiceItemBasic(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		},
	}, nil)

	mtoShipment := factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &sitDaysAllowance,
			},
		},
	}, nil)

	year, month, day := time.Now().Add(time.Hour * 24 * -60).Date()
	twoMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	oneMonthAgo := twoMonthsAgo.Add(time.Hour * 24 * 30)
	postalCode := "90210"
	reason := "peak season all trucks in use"
	// This will in practice not exist without DOFSIT and DOASIT
	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:           models.MTOServiceItemStatusApproved,
				SITEntryDate:     &twoMonthsAgo,
				SITDepartureDate: &oneMonthAgo,
				SITPostalCode:    &postalCode,
				Reason:           &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	approvedAt := time.Now()
	oneWeekAgo := oneMonthAgo.Add(time.Hour * 24 * 23)
	dddsit := factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &oneWeekAgo,
				Reason:       &reason,
				ApprovedAt:   &approvedAt,
			},
		},
		{
			Model: models.Address{},
			Type:  &factory.Addresses.SITDestinationOriginalAddress,
		},
		{
			Model: models.Address{},
			Type:  &factory.Addresses.SITDestinationFinalAddress,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: dddsit,
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			Type: models.CustomerContactTypeSecond,
		},
		MTOServiceItem: dddsit,
	})

	return dddsit
}

func createPaymentRequestsWithPartialSITInvoice(appCtx appcontext.AppContext, primeUploader *uploader.PrimeUploader) {
	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(appCtx.DB(), []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)

	// Move available to the prime with 3 shipments (control, 2 w/ SITS)
	availableToPrimeAt := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.Move{
				Locator:            "PARSIT",
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		},
	}, nil)
	oneHundredAndTwentyDays := 120
	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &oneHundredAndTwentyDays,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	firstPrimeUpload := factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.PaymentRequest{
				Status: models.PaymentRequestStatusReviewed,
			},
		},
		{
			Model: models.PrimeUpload{},
			ExtendedParams: &factory.PrimeUploadExtendedParams{
				PrimeUploader: primeUploader,
				AppContext:    appCtx,
			},
		},
	}, nil)

	firstPaymentRequest := firstPrimeUpload.ProofOfServiceDoc.PaymentRequest

	secondPrimeUpload := factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.PaymentRequest{
				SequenceNumber: 2,
			},
		},
		{
			Model: models.PrimeUpload{},
			ExtendedParams: &factory.PrimeUploadExtendedParams{
				PrimeUploader: primeUploader,
				AppContext:    appCtx,
			},
		},
	}, nil)

	secondPaymentRequest := secondPrimeUpload.ProofOfServiceDoc.PaymentRequest

	year, month, day := time.Now().Date()
	originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
	originDepartureDate := originEntryDate.Add(time.Hour * 24 * 30)

	destinationEntryDate := time.Date(year, month, day-89, 0, 0, 0, 0, time.UTC)
	destinationDepartureDate := destinationEntryDate.Add(time.Hour * 24 * 60)

	// First reviewed payment request with 30 days billed for origin SIT
	dofsit := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
		}, {
			Model:    firstPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    dofsit,
			LinkOnly: true,
		},
	}, nil)

	doasit := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// Creates the approved payment service item for DOASIT w/ SIT start date param
	doasitParam := factory.BuildPaymentServiceItemParam(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItemParam{
				Value: originEntryDate.Format("2006-01-02"),
			},
		},
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
		},
		{
			Model: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
		},
		{
			Model:    firstPaymentRequest,
			LinkOnly: true,
		},
		{
			Model:    doasit,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// Creates the SIT end date param for existing DOASIT payment request service item
	factory.BuildPaymentServiceItemParam(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItemParam{
				Value: originDepartureDate.Format("2006-01-02"),
			},
		},
		{
			Model: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
		},
		{
			Model:    doasitParam.PaymentServiceItem,
			LinkOnly: true,
		},
		{
			Model:    firstPaymentRequest,
			LinkOnly: true,
		},
		{
			Model:    doasit,
			LinkOnly: true,
		},
	}, nil)

	// Creates the NumberDaysSIT param for existing DOASIT payment request service item
	factory.BuildPaymentServiceItemParam(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItemParam{
				Value: "30",
			},
		},
		{
			Model: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
		},
		{
			Model:    doasitParam.PaymentServiceItem,
			LinkOnly: true,
		},
		{
			Model:    firstPaymentRequest,
			LinkOnly: true,
		},
		{
			Model:    doasit,
			LinkOnly: true,
		},
	}, nil)

	dopsit := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:           models.MTOServiceItemStatusApproved,
				SITEntryDate:     &originEntryDate,
				SITDepartureDate: &originDepartureDate,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
		}, {
			Model:    dopsit,
			LinkOnly: true,
		}, {
			Model:    firstPaymentRequest,
			LinkOnly: true,
		},
	}, nil)

	// Destination SIT service items for the second payment request
	ddfsit := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &destinationEntryDate,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    ddfsit,
			LinkOnly: true,
		},
		{
			Model:    secondPaymentRequest,
			LinkOnly: true,
		},
	}, nil)

	ddasit := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &destinationEntryDate,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDASIT,
			},
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	ddasitParam := factory.BuildPaymentServiceItemParam(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItemParam{
				Value: destinationEntryDate.Format("2006-01-02"),
			},
		},
		{
			Model: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
		},
		{
			Model:    secondPaymentRequest,
			LinkOnly: true,
		},
		{
			Model:    ddasit,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPaymentServiceItemParam(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItemParam{
				Value: destinationDepartureDate.Format("2006-01-02"),
			},
		},
		{
			Model: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
		},
		{
			Model:    ddasitParam.PaymentServiceItem,
			LinkOnly: true,
		},
		{
			Model:    secondPaymentRequest,
			LinkOnly: true,
		},
		{
			Model:    ddasit,
			LinkOnly: true,
		},
	}, nil)

	// Creates the NumberDaysSIT param for existing DOASIT payment request service item
	factory.BuildPaymentServiceItemParam(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItemParam{
				Value: "60",
			},
		},
		{
			Model: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
		},
		{
			Model:    ddasitParam.PaymentServiceItem,
			LinkOnly: true,
		},
		{
			Model:    secondPaymentRequest,
			LinkOnly: true,
		},
		{
			Model:    ddasit,
			LinkOnly: true,
		},
	}, nil)

	// Will leave the departure date blank with 30 days left in SIT Days authorized
	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &destinationEntryDate,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWithAllPendingTOOActions(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	makeAmendedOrders(appCtx, orders, userUploader, &[]string{"medium.jpg", "small.pdf"})
	move := makeMoveForOrders(appCtx, orders, "PENDNG", models.MoveStatusAPPROVALSREQUESTED)
	now := time.Now()
	move.ExcessWeightQualifiedAt = &now
	testdatagen.MustSave(db, &move)
	shipment := makeRiskOfExcessShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	MakeSITWithPendingSITExtensionsForShipment(appCtx, shipment)
	paymentRequestID := uuid.Must(uuid.FromString("70b35add-605a-289d-8dad-056f5d9ef7e1"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func MakeSITWithPendingSITExtensionsForShipment(appCtx appcontext.AppContext, shipment models.MTOShipment) {
	db := appCtx.DB()

	year, month, day := time.Now().Date()
	thirtyDaysAgo := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
	factory.BuildMTOServiceItem(db, []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOServiceItem{
				SITEntryDate: &thirtyDaysAgo,
				Status:       models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		},
	}, nil)

	for i := 0; i < 2; i++ {
		factory.BuildSITDurationUpdate(db, []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)
	}
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
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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

func createMoveWithNTSRShipment(appCtx appcontext.AppContext, locator string, usesExternalVendor bool) {
	db := appCtx.DB()
	submittedAt := time.Now()
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	orders := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
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
}

// createRandomMove creates a random move with fake data that has been approved for usage
func createRandomMove(
	appCtx appcontext.AppContext,
	possibleStatuses []models.MoveStatus,
	allDutyLocations []models.DutyLocation,
	dutyLocationsInGBLOC []models.DutyLocation,
	withFullOrder bool,
	_ *uploader.UserUploader,
	moveTemplate models.Move,
	mtoShipmentTemplate models.MTOShipment,
	orderTemplate models.Order,
	serviceMemberTemplate models.ServiceMember,
) models.Move {
	db := appCtx.DB()
	randDays, err := random.GetRandomInt(366)
	if err != nil {
		log.Panic(fmt.Errorf("Unable to generate random integer for submitted move date"), zap.Error(err))
	}
	submittedAt := time.Now().AddDate(0, 0, randDays*-1)

	if serviceMemberTemplate.Affiliation == nil {
		randomAffiliation, err := random.GetRandomInt(5)
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for affiliation"), zap.Error(err))
		}
		serviceMemberTemplate.Affiliation = &[]models.ServiceMemberAffiliation{
			models.AffiliationARMY,
			models.AffiliationAIRFORCE,
			models.AffiliationNAVY,
			models.AffiliationCOASTGUARD,
			models.AffiliationMARINES}[randomAffiliation]
	}

	customs := []factory.Customization{
		{
			Model: serviceMemberTemplate,
		},
		{
			Model: orderTemplate,
		},
	}

	dutyLocationCount := len(allDutyLocations)
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	if orderTemplate.OriginDutyLocationID == nil {
		// We can pick any origin duty location not only one in the office user's GBLOC
		if *serviceMemberTemplate.Affiliation == models.AffiliationMARINES {
			randDutyStationIndex, err := random.GetRandomInt(dutyLocationCount)
			if err != nil {
				log.Panic(fmt.Errorf("Unable to generate random integer for duty location"), zap.Error(err))
			}
			customs = append(customs, factory.Customization{
				Model:    allDutyLocations[randDutyStationIndex],
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			})
		} else {
			randDutyStationIndex, err := random.GetRandomInt(len(dutyLocationsInGBLOC))
			if err != nil {
				log.Panic(fmt.Errorf("Unable to generate random integer for duty location"), zap.Error(err))
			}
			customs = append(customs, factory.Customization{
				Model:    dutyLocationsInGBLOC[randDutyStationIndex],
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			})
		}
	}

	if orderTemplate.NewDutyLocationID == uuid.Nil {
		randDutyStationIndex, err := random.GetRandomInt(dutyLocationCount)
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for duty location"), zap.Error(err))
		}
		customs = append(customs, factory.Customization{
			Model:    allDutyLocations[randDutyStationIndex],
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		})
	}

	randomFirst, randomLast := fakedata.RandomName()
	serviceMemberTemplate.FirstName = &randomFirst
	serviceMemberTemplate.LastName = &randomLast
	serviceMemberTemplate.CacValidated = true

	// assertions passed in means we cannot yet convert to BuildOrder
	var order models.Order
	if withFullOrder {
		order = factory.BuildOrder(db, customs, nil)
	} else {
		order = factory.BuildOrderWithoutDefaults(db, customs, nil)
	}

	if moveTemplate.SubmittedAt == nil {
		moveTemplate.SubmittedAt = &submittedAt
	}

	if moveTemplate.Status == "" {
		randStatusIndex, err := random.GetRandomInt(len(possibleStatuses))
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for move status"), zap.Error(err))
		}
		moveTemplate.Status = possibleStatuses[randStatusIndex]

		if moveTemplate.Status == models.MoveStatusServiceCounselingCompleted {
			counseledAt := submittedAt.Add(3 * 24 * time.Hour)
			moveTemplate.ServiceCounselingCompletedAt = &counseledAt
		}
	}

	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveTemplate,
		},
		{
			Model:    order,
			LinkOnly: true,
		},
	}, nil)

	shipmentStatus := models.MTOShipmentStatusSubmitted
	if mtoShipmentTemplate.Status != "" {
		shipmentStatus = mtoShipmentTemplate.Status
	}

	laterRequestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	laterRequestedDeliveryDate := laterRequestedPickupDate.Add(7 * 24 * time.Hour)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                shipmentStatus,
				RequestedPickupDate:   &laterRequestedPickupDate,
				RequestedDeliveryDate: &laterRequestedDeliveryDate,
				ApprovedDate:          mtoShipmentTemplate.ApprovedDate,
				Diversion:             mtoShipmentTemplate.Diversion,
			},
		},
	}, nil)

	earlierRequestedPickupDate := submittedAt.Add(30 * 24 * time.Hour)
	earlierRequestedDeliveryDate := earlierRequestedPickupDate.Add(7 * 24 * time.Hour)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                shipmentStatus,
				RequestedPickupDate:   &earlierRequestedPickupDate,
				RequestedDeliveryDate: &earlierRequestedDeliveryDate,
				ApprovedDate:          mtoShipmentTemplate.ApprovedDate,
				Diversion:             mtoShipmentTemplate.Diversion,
			},
		},
	}, nil)

	return move
}

func createMultipleMovesTwoMovesHHGAndPPMShipments(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	userID := uuid.Must(uuid.NewV4())
	oktaID := uuid.Must(uuid.NewV4())
	email := "multiplemoves@HHG_PPM.com"

	originDutyLocation := factory.BuildDutyLocation(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				Name:                       "Woodbine, KY 40356",
				ProvidesServicesCounseling: true,
			},
		},
	}, nil)
	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	affiliation := models.AffiliationARMY
	serviceMember := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation:   &affiliation,
				FirstName:     models.StringPointer("Lola"),
				LastName:      models.StringPointer("Smith"),
				Edipi:         models.StringPointer("8362534853"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	// Move A
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	ordersNumber := "8675309"
	departmentIndicator := "ARMY"
	tac := "E19A"
	ordersA := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveADetails models.Move
	moveADetails.Status = models.MoveStatusAPPROVED
	moveADetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveADetails.Locator = "MMOVEA"
	moveADetails.CreatedAt = time.Now().Add(time.Hour * 24)

	moveA := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveADetails,
		},
		{
			Model:    ordersA,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	requestedPickupDateMoveA := moveA.CreatedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDateMoveA := requestedPickupDateMoveA.Add(7 * 24 * time.Hour)
	destinationAddressMoveA := factory.BuildAddress(db, nil, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    moveA,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDateMoveA,
				RequestedDeliveryDate: &requestedDeliveryDateMoveA,
			},
		},
		{
			Model:    destinationAddressMoveA,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	// Move B
	ordersB := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveBDetails models.Move
	moveBDetails.Status = models.MoveStatusAPPROVED
	moveBDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveBDetails.Locator = "MMOVEB"
	moveBDetails.CreatedAt = time.Now().Add(time.Hour * -24)

	moveB := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveBDetails,
		},
		{
			Model:    ordersB,
			LinkOnly: true,
		},
	}, nil)

	requestedPickupDateMoveB := moveB.CreatedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDateMoveB := requestedPickupDateMoveB.Add(7 * 24 * time.Hour)
	destinationAddressMoveB := factory.BuildAddress(db, nil, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDateMoveB,
				RequestedDeliveryDate: &requestedDeliveryDateMoveB,
			},
		},
		{
			Model:    destinationAddressMoveB,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveB,
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
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Move C
	ordersC := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveCDetails models.Move
	moveCDetails.Status = models.MoveStatusAPPROVED
	moveCDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveCDetails.Locator = "MMOVEC"
	moveCDetails.CreatedAt = time.Now().Add(time.Hour * -48)

	moveC := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveCDetails,
		},
		{
			Model:    ordersC,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveC,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.NTSRaw,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
}

func createMultipleMovesThreeMovesHHGPPMNTSShipments(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	userID := uuid.Must(uuid.NewV4())
	oktaID := uuid.Must(uuid.NewV4())
	email := "multiplemoves@HHG_PPM_NTS.com"

	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"

	originDutyLocation := factory.BuildDutyLocation(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				Name:                       "Hope, KY 40356",
				ProvidesServicesCounseling: true,
			},
		},
	}, nil)
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	affiliation := models.AffiliationAIRFORCE
	serviceMember := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation:   &affiliation,
				FirstName:     models.StringPointer("Hannah"),
				LastName:      models.StringPointer("James"),
				Edipi:         models.StringPointer("8362534857"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	// Move A
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	ordersNumber := "8675309"
	departmentIndicator := "AIR FORCE"
	tac := "E13V"
	ordersA := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveADetails models.Move
	moveADetails.Status = models.MoveStatusNeedsServiceCounseling
	moveADetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveADetails.Locator = "MMOVED"
	moveADetails.CreatedAt = time.Now().Add(time.Hour * 24)

	moveA := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveADetails,
		},
		{
			Model:    ordersA,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	requestedPickupDateMoveA := moveA.CreatedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDateMoveA := requestedPickupDateMoveA.Add(7 * 24 * time.Hour)
	destinationAddressMoveA := factory.BuildAddress(db, nil, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    moveA,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDateMoveA,
				RequestedDeliveryDate: &requestedDeliveryDateMoveA,
			},
		},
		{
			Model:    destinationAddressMoveA,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	// Move B
	ordersB := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveBDetails models.Move
	moveBDetails.Status = models.MoveStatusAPPROVED
	moveBDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveBDetails.Locator = "MMOVEE"
	moveBDetails.CreatedAt = time.Now().Add(time.Hour * -24)

	moveB := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveBDetails,
		},
		{
			Model:    ordersB,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveB,
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
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Move C
	ordersC := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveCDetails models.Move
	moveCDetails.Status = models.MoveStatusAPPROVED
	moveCDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveCDetails.Locator = "MMOVEF"
	moveCDetails.CreatedAt = time.Now().Add(time.Hour * -48)

	moveC := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveCDetails,
		},
		{
			Model:    ordersC,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveC,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.NTSRaw,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
}

func createMultipleMovesThreeMovesNTSHHGShipments(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	userID := uuid.Must(uuid.NewV4())
	oktaID := uuid.Must(uuid.NewV4())
	email := "multiplemoves@NTS_HHG.com"

	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"

	originDutyLocation := factory.BuildDutyLocation(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				Name:                       "Centre, KY 40356",
				ProvidesServicesCounseling: true,
			},
		},
	}, nil)
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	affiliation := models.AffiliationNAVY
	serviceMember := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation:   &affiliation,
				FirstName:     models.StringPointer("Jenna"),
				LastName:      models.StringPointer("Ken"),
				Edipi:         models.StringPointer("8362534854"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	// Move A
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	ordersNumber := "8675309"
	departmentIndicator := "NAVY"
	tac := "E13V"
	ordersA := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveADetails models.Move
	moveADetails.Status = models.MoveStatusAPPROVALSREQUESTED
	moveADetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveADetails.Locator = "MMOVEG"
	moveADetails.CreatedAt = time.Now().Add(time.Hour * 24)

	moveA := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveADetails,
		},
		{
			Model:    ordersA,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	requestedPickupDateMoveA := moveA.CreatedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDateMoveA := requestedPickupDateMoveA.Add(7 * 24 * time.Hour)
	destinationAddressMoveA := factory.BuildAddress(db, nil, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    moveA,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDateMoveA,
				RequestedDeliveryDate: &requestedDeliveryDateMoveA,
			},
		},
		{
			Model:    destinationAddressMoveA,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveA,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.NTSRaw,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Move B
	ordersB := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveBDetails models.Move
	moveBDetails.Status = models.MoveStatusAPPROVED
	moveBDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveBDetails.Locator = "MMOVEH"
	moveBDetails.CreatedAt = time.Now().Add(time.Hour * -24)

	moveB := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveBDetails,
		},
		{
			Model:    ordersB,
			LinkOnly: true,
		},
	}, nil)

	requestedPickupDateMoveB := moveA.CreatedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDateMoveB := requestedPickupDateMoveA.Add(7 * 24 * time.Hour)
	destinationAddressMoveB := factory.BuildAddress(db, nil, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDateMoveB,
				RequestedDeliveryDate: &requestedDeliveryDateMoveB,
			},
		},
		{
			Model:    destinationAddressMoveB,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.NTSRaw,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.NTSRaw,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Move C
	ordersC := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveCDetails models.Move
	moveCDetails.Status = models.MoveStatusAPPROVED
	moveCDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveCDetails.Locator = "MMOVEI"
	moveCDetails.CreatedAt = time.Now().Add(time.Hour * -48)

	moveC := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveCDetails,
		},
		{
			Model:    ordersC,
			LinkOnly: true,
		},
	}, nil)

	requestedPickupDateMoveC := moveA.CreatedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDateMoveC := requestedPickupDateMoveA.Add(7 * 24 * time.Hour)
	destinationAddressMoveC := factory.BuildAddress(db, nil, nil)

	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    moveC,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDateMoveC,
				RequestedDeliveryDate: &requestedDeliveryDateMoveC,
			},
		},
		{
			Model:    destinationAddressMoveC,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveC,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.NTSRaw,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
}

func createMultipleMovesThreeMovesPPMShipments(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	userID := uuid.Must(uuid.NewV4())
	oktaID := uuid.Must(uuid.NewV4())
	email := "multiplemoves@PPM.com"

	newDutyLocation := factory.FetchOrBuildCurrentDutyLocation(db)
	newDutyLocation.Address.PostalCode = "52549"

	originDutyLocation := factory.BuildDutyLocation(db, []factory.Customization{
		{
			Model: models.DutyLocation{
				Name:                       "Marion, KY 40356",
				ProvidesServicesCounseling: true,
			},
		},
	}, nil)
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			}},
	}, nil)

	affiliation := models.AffiliationCOASTGUARD
	serviceMember := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation:   &affiliation,
				FirstName:     models.StringPointer("Tori"),
				LastName:      models.StringPointer("Ross"),
				Edipi:         models.StringPointer("8362534852"),
				PersonalEmail: models.StringPointer(email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	// Move A
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hhgPermitted := internalmessages.OrdersTypeDetailHHGPERMITTED
	ordersNumber := "8675309"
	departmentIndicator := "COAST GUARD"
	tac := "E19A"
	ordersA := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveADetails models.Move
	moveADetails.Status = models.MoveStatusAPPROVED
	moveADetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveADetails.Locator = "MMOVEJ"
	moveADetails.CreatedAt = time.Now().Add(time.Hour * 24)

	moveA := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveADetails,
		},
		{
			Model:    ordersA,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveA,
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
			Model:    moveA,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Move B
	ordersNumber = "8675302"
	ordersB := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveBDetails models.Move
	moveBDetails.Status = models.MoveStatusAPPROVED
	moveBDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveBDetails.Locator = "MMOVEK"
	moveBDetails.CreatedAt = time.Now().Add(time.Hour * 24)

	moveB := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveBDetails,
		},
		{
			Model:    ordersB,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveB,
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
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Move C
	ordersNumber = "8675301"
	ordersC := factory.BuildOrderWithoutDefaults(db, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				OrdersType:          pcos,
				OrdersTypeDetail:    &hhgPermitted,
				OrdersNumber:        &ordersNumber,
				DepartmentIndicator: &departmentIndicator,
				TAC:                 &tac,
			},
		},
	}, nil)

	var moveCDetails models.Move
	moveCDetails.Status = models.MoveStatusAPPROVED
	moveCDetails.AvailableToPrimeAt = models.TimePointer(time.Now())
	moveCDetails.Locator = "MMOVEL"
	moveCDetails.CreatedAt = time.Now().Add(time.Hour * 24)

	moveC := factory.BuildMove(db, []factory.Customization{
		{
			Model: moveCDetails,
		},
		{
			Model:    ordersC,
			LinkOnly: true,
		},
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	requestedPickupDateMoveB := moveB.CreatedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDateMoveB := requestedPickupDateMoveB.Add(7 * 24 * time.Hour)
	destinationAddressMoveB := factory.BuildAddress(db, nil, nil)
	factory.BuildMTOShipment(db, []factory.Customization{
		{
			Model:    moveB,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusSubmitted,
				RequestedPickupDate:   &requestedPickupDateMoveB,
				RequestedDeliveryDate: &requestedDeliveryDateMoveB,
			},
		},
		{
			Model:    destinationAddressMoveB,
			Type:     &factory.Addresses.DeliveryAddress,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    moveC,
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
			Model:    moveC,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
}

func CreateBoatHaulAwayMoveForSC(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, _ services.MoveRouter, moveInfo MoveCreatorInfo) models.Move {
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        moveInfo.UserID,
				OktaID:    oktaID.String(),
				OktaEmail: moveInfo.Email,
				Active:    true,
			}},
	}, nil)

	smWithBoat := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            moveInfo.SmID,
				FirstName:     models.StringPointer(moveInfo.FirstName),
				LastName:      models.StringPointer(moveInfo.LastName),
				PersonalEmail: models.StringPointer(moveInfo.Email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    smWithBoat,
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

	if *smWithBoat.Affiliation == models.AffiliationARMY || *smWithBoat.Affiliation == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}
	mtoShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeBoatHaulAway,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildBoatShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.BoatShipment{
				Type:           models.BoatShipmentTypeHaulAway,
				Year:           models.IntPointer(2000),
				Make:           models.StringPointer("Boat Make"),
				Model:          models.StringPointer("Boat Model"),
				LengthInInches: models.IntPointer(300),
				WidthInInches:  models.IntPointer(108),
				HeightInInches: models.IntPointer(72),
				HasTrailer:     models.BoolPointer(true),
				IsRoadworthy:   models.BoolPointer(false),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
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

func CreateBoatHaulAwayMoveForTOO(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, _ services.MoveRouter, moveInfo MoveCreatorInfo) models.Move {
	oktaID := uuid.Must(uuid.NewV4())
	submittedAt := time.Now()

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        moveInfo.UserID,
				OktaID:    oktaID.String(),
				OktaEmail: moveInfo.Email,
				Active:    true,
			}},
	}, nil)

	smWithBoat := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            moveInfo.SmID,
				FirstName:     models.StringPointer(moveInfo.FirstName),
				LastName:      models.StringPointer(moveInfo.LastName),
				PersonalEmail: models.StringPointer(moveInfo.Email),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    smWithBoat,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:          moveInfo.MoveID,
				Locator:     moveInfo.MoveLocator,
				Status:      models.MoveStatusSUBMITTED,
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

	if *smWithBoat.Affiliation == models.AffiliationARMY || *smWithBoat.Affiliation == models.AffiliationAIRFORCE {
		move.CloseoutOfficeID = &DefaultCloseoutOfficeID
		testdatagen.MustSave(appCtx.DB(), &move)
	}
	mtoShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeBoatHaulAway,
				Status:       models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildBoatShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.BoatShipment{
				Type:           models.BoatShipmentTypeHaulAway,
				Year:           models.IntPointer(2000),
				Make:           models.StringPointer("Boat Make"),
				Model:          models.StringPointer("Boat Model"),
				LengthInInches: models.IntPointer(300),
				WidthInInches:  models.IntPointer(108),
				HeightInInches: models.IntPointer(72),
				HasTrailer:     models.BoolPointer(true),
				IsRoadworthy:   models.BoolPointer(false),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
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
