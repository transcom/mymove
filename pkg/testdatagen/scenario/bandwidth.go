package scenario

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

// bandwidthScenario builds 1 move with 3 documents of varying sizes for each
// of the User Order uploads, and the Prime's Proof of Service docs
type bandwidthScenario NamedScenario

// BandwidthScenario is the thing
var BandwidthScenario = bandwidthScenario{Name: "bandwidth"}

func createHHGMove150Kb(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "S150KB", models.MoveStatusSUBMITTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("68034aa3-831c-4d2d-9fd4-b66bc0cc5130"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createHHGMove2mb(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"2mb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MED2MB", models.MoveStatusSUBMITTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("4de88d57-9723-446b-904c-cf8d0a834687"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createHHGMove25mb(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"25mb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "LG25MB", models.MoveStatusSUBMITTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("aca5cc9c-c266-4a7d-895d-dc3c9c0d9894"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func makeServiceMember(appCtx appcontext.AppContext) models.ServiceMember {
	affiliation := models.AffiliationCOASTGUARD
	serviceMember := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &affiliation,
			},
		},
	}, nil)

	return serviceMember
}

func makeAmendedOrders(appCtx appcontext.AppContext, order models.Order, userUploader *uploader.UserUploader, fileNames *[]string) models.Order {
	document := factory.BuildDocumentLinkServiceMember(appCtx.DB(), order.ServiceMember)

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

	order.UploadedAmendedOrders = &document
	order.UploadedAmendedOrdersID = &document.ID
	saveErr := appCtx.DB().Save(&order)
	if saveErr != nil {
		log.Panic("error saving amended orders upload to orders")
	}

	return order
}

func makeRiskOfExcessShipmentForMove(appCtx appcontext.AppContext, move models.Move, shipmentStatus models.MTOShipmentStatus) models.MTOShipment {
	estimatedWeight := unit.Pound(7200)
	actualWeight := unit.Pound(7400)
	daysOfSIT := 90
	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				SITDaysAllowance:     &daysOfSIT,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               shipmentStatus,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	return MTOShipment
}

func makeShipmentForMove(appCtx appcontext.AppContext, move models.Move, shipmentStatus models.MTOShipmentStatus) models.MTOShipment {
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	billableWeight := unit.Pound(4000)
	billableWeightJustification := "heavy"

	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:        &estimatedWeight,
				PrimeActualWeight:           &actualWeight,
				ShipmentType:                models.MTOShipmentTypeHHG,
				ApprovedDate:                models.TimePointer(time.Now()),
				Status:                      shipmentStatus,
				BillableWeightCap:           &billableWeight,
				BillableWeightJustification: &billableWeightJustification,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	return MTOShipment
}

func makePaymentRequestForShipment(appCtx appcontext.AppContext, move models.Move, shipment models.MTOShipment, primeUploader *uploader.PrimeUploader, fileNames *[]string, paymentRequestID uuid.UUID, status models.PaymentRequestStatus) {
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:      paymentRequestID,
				IsFinal: false,
				Status:  status,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	dcrtCost := unit.Cents(99999)
	mtoServiceItemDCRT, err := testdatagen.MakeMTOServiceItemDomesticCrating(appCtx.DB(), testdatagen.Assertions{
		Move:        move,
		MTOShipment: shipment,
	})
	if err != nil {
		log.Panic(err)
	}

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    mtoServiceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	ducrtCost := unit.Cents(99999)
	mtoServiceItemDUCRT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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
				ID: uuid.FromStringOrNil("fc14935b-ebd3-4df3-940b-f30e71b6a56c"), // DUCRT - Domestic uncrating
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ducrtCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    mtoServiceItemDUCRT,
			LinkOnly: true,
		},
	}, nil)

	files := filesInBandwidthTestDirectory(fileNames)
	// Creates prime upload documents from the files in this directory:
	// pkg/testdatagen/testdata/bandwidth_test_docs
	for _, file := range files {
		filePath := fmt.Sprintf("bandwidth_test_docs/%s", file)
		fixture := testdatagen.Fixture(filePath)
		factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model: models.PrimeUpload{},
				ExtendedParams: &factory.PrimeUploadExtendedParams{
					PrimeUploader: primeUploader,
					AppContext:    appCtx,
					File:          fixture,
				},
			},
		}, nil)
	}
}

func createOfficeUser(appCtx appcontext.AppContext) {
	/* A user with both too and tio roles */
	email := "too_tio_role@office.mil"
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	tooTioUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"))
	OktaID := uuid.Must(uuid.NewV4())
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        tooTioUUID,
				OktaID:    OktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole, tioRole},
			},
		},
	}, nil)
	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-54212ae3078f"),
				Email:  email,
				Active: true,
				UserID: &tooTioUUID,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
}

func filesInBandwidthTestDirectory(fileNames *[]string) []string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(fmt.Errorf("failed to get current directory: %s", err))
	}

	dirName := path.Join(cwd, "pkg/testdatagen/testdata/bandwidth_test_docs")

	files, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}

	var docs []string

	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		if fileNames == nil || len(*fileNames) == 0 {
			docs = append(docs, file.Name())
		} else {
			for _, fileName := range *fileNames {
				if fileName == file.Name() {
					docs = append(docs, file.Name())
				}
			}
		}
	}

	return docs
}

// Run does that data load thing
func (e bandwidthScenario) Run(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	createOfficeUser(appCtx)
	createHHGMove150Kb(appCtx, userUploader, primeUploader)
	createHHGMove2mb(appCtx, userUploader, primeUploader)
	createHHGMove25mb(appCtx, userUploader, primeUploader)
}
