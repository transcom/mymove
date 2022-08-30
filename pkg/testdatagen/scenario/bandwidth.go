package scenario

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
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
	serviceMember := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Affiliation: &affiliation,
		},
	})

	return serviceMember
}

func makeAmendedOrders(appCtx appcontext.AppContext, order models.Order, userUploader *uploader.UserUploader, fileNames *[]string) models.Order {
	document := testdatagen.MakeDocument(appCtx.DB(), testdatagen.Assertions{
		Document: models.Document{
			ServiceMemberID: order.ServiceMemberID,
			ServiceMember:   order.ServiceMember,
		},
	})

	// Creates order upload documents from the files in this directory:
	// pkg/testdatagen/testdata/bandwidth_test_docs

	files := filesInBandwidthTestDirectory(fileNames)

	for _, file := range files {
		filePath := fmt.Sprintf("bandwidth_test_docs/%s", file)
		fixture := testdatagen.Fixture(filePath)

		upload := testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
			File: fixture,
			UserUpload: models.UserUpload{
				UploaderID: order.ServiceMember.UserID,
				DocumentID: &document.ID,
				Document:   document,
			},
			UserUploader: userUploader,
		})
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
	MTOShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			SITDaysAllowance:     &daysOfSIT,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               shipmentStatus,
		},
		Move: move,
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   MTOShipment,
			MTOShipmentID: MTOShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	return MTOShipment
}

func makeShipmentForMove(appCtx appcontext.AppContext, move models.Move, shipmentStatus models.MTOShipmentStatus) models.MTOShipment {
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	billableWeight := unit.Pound(4000)
	billableWeightJustification := "heavy"

	MTOShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight:        &estimatedWeight,
			PrimeActualWeight:           &actualWeight,
			ShipmentType:                models.MTOShipmentTypeHHG,
			ApprovedDate:                swag.Time(time.Now()),
			Status:                      shipmentStatus,
			BillableWeightCap:           &billableWeight,
			BillableWeightJustification: &billableWeightJustification,
		},
		Move: move,
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   MTOShipment,
			MTOShipmentID: MTOShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	return MTOShipment
}

func makePaymentRequestForShipment(appCtx appcontext.AppContext, move models.Move, shipment models.MTOShipment, primeUploader *uploader.PrimeUploader, fileNames *[]string, paymentRequestID uuid.UUID, status models.PaymentRequestStatus) {
	paymentRequest := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            paymentRequestID,
			MoveTaskOrder: move,
			IsFinal:       false,
			Status:        status,
		},
		Move: move,
	})

	dcrtCost := unit.Cents(99999)
	mtoServiceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(appCtx.DB(), testdatagen.Assertions{
		Move:        move,
		MTOShipment: shipment,
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dcrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: mtoServiceItemDCRT,
	})

	ducrtCost := unit.Cents(99999)
	mtoServiceItemDUCRT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		Move:        move,
		MTOShipment: shipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("fc14935b-ebd3-4df3-940b-f30e71b6a56c"), // DUCRT - Domestic uncrating
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ducrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: mtoServiceItemDUCRT,
	})

	files := filesInBandwidthTestDirectory(fileNames)
	// Creates prime upload documents from the files in this directory:
	// pkg/testdatagen/testdata/bandwidth_test_docs
	for _, file := range files {
		filePath := fmt.Sprintf("bandwidth_test_docs/%s", file)
		fixture := testdatagen.Fixture(filePath)
		testdatagen.MakePrimeUpload(appCtx.DB(), testdatagen.Assertions{
			File:           fixture,
			PaymentRequest: paymentRequest,
			PrimeUploader:  primeUploader,
		})
	}
}

func createOfficeUser(appCtx appcontext.AppContext) {
	/* A user with both too and tio roles */
	email := "too_tio_role@office.mil"
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

	tooTioUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"))
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            tooTioUUID,
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tooRole, tioRole},
		},
	})
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-54212ae3078f"),
			Email:  email,
			Active: true,
			UserID: &tooTioUUID,
		},
	})
}

func filesInBandwidthTestDirectory(fileNames *[]string) []string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(fmt.Errorf("failed to get current directory: %s", err))
	}

	dirName := path.Join(cwd, "pkg/testdatagen/testdata/bandwidth_test_docs")

	files, err := ioutil.ReadDir(dirName)
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
