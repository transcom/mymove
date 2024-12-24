package testharness

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func newUserUploader(appCtx appcontext.AppContext) *uploader.UserUploader {
	// initialize this directly with defaults instead of using command
	// line options. Simple for now, we can revist if we need to
	fsParams := storage.NewFilesystemParams("tmp", "storage")
	storer := storage.NewFilesystem(fsParams)

	userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate user uploader", zap.Error(err))
	}
	return userUploader
}

func newPrimeUploader(appCtx appcontext.AppContext) *uploader.PrimeUploader {
	// initialize this directly with defaults instead of using command
	// line options. Simple for now, we can revist if we need to
	fsParams := storage.NewFilesystemParams("tmp", "storage")
	storer := storage.NewFilesystem(fsParams)

	primeUploader, err := uploader.NewPrimeUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate prime uploader", zap.Error(err))
	}
	return primeUploader
}

type userInfo struct {
	email     string
	firstName string
	lastName  string
}

func newUserInfo(emailSubstring string) userInfo {
	email := strings.ToLower(fmt.Sprintf("%s"+emailSubstring+"_%s@example.com",
		testdatagen.MakeRandomString(5), testdatagen.MakeRandomString(8)))
	username := strings.Split(email, "@")[0]
	firstName := strings.Split(username, "_")[0]
	lastName := username[len(firstName)+1:]
	return userInfo{
		email:     email,
		firstName: firstName,
		lastName:  lastName,
	}
}

func MakeMoveWithOrders(db *pop.Connection) models.Move {
	userInfo := newUserInfo("customer")
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				PersonalEmail: models.StringPointer(userInfo.email),
				CacValidated:  true,
			},
		},
	}, nil)

	return move
}

func MakeSpouseProGearMove(db *pop.Connection) models.Move {
	userInfo := newUserInfo("customer")
	move := factory.BuildMove(db, []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				PersonalEmail: models.StringPointer(userInfo.email),
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model: models.Order{
				HasDependents:    true,
				SpouseHasProGear: true,
			},
		},
	}, nil)

	return move
}

func MakeWithShipmentMove(appCtx appcontext.AppContext) models.Move {
	userInfo := newUserInfo("customer")
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				PersonalEmail: models.StringPointer(userInfo.email),
				CacValidated:  true,
			},
		},
		{
			Model: models.Order{
				HasDependents:    true,
				SpouseHasProGear: true,
			},
		},
	}, nil)

	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// This is a postal code that maps to the default office user gbloc KKFA in the PostalCodeToGBLOC table
				PostalCode: "85004",
			},
		},
	}, nil)

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove

}

// MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO is a function
// that creates an HHG move with service items and payments requests with files
// from the Prime for review by thte TOO
// copied almost verbatim from e2ebasic createHHGMoveWithServiceItemsAndPaymentRequestsAndFiles
func MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	primeUploader := newPrimeUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		},
	}, nil)
	sitDaysAllowance := 270
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	actualPickupDate := time.Now().AddDate(0, 0, 1)
	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				Status:               models.MTOShipmentStatusSubmitted,
				SITDaysAllowance:     &sitDaysAllowance,
				ActualPickupDate:     &actualPickupDate,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal: false,
				Status:  models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	twoMonthsAgo := threeMonthsAgo.AddDate(0, 1, 0)
	sitCost := unit.Cents(200000)
	sitItems := factory.BuildOriginSITServiceItems(appCtx.DB(), mto, MTOShipment, &threeMonthsAgo, &twoMonthsAgo)
	sitItems = append(sitItems, factory.BuildDestSITServiceItems(appCtx.DB(), mto, MTOShipment, &twoMonthsAgo, nil)...)
	for i := range sitItems {
		factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &sitCost,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			}, {
				Model:    sitItems[i],
				LinkOnly: true,
			},
		}, nil)
	}
	scenario.MakeSITExtensionsForShipment(appCtx, MTOShipment)

	currentTimeDCRT := time.Now()

	basicPaymentServiceItemParamsDCRT := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractYearName,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameEscalationCompounded,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.FormatFloat(1.125, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "1.71",
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameCubicFeetBilled,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "4.00",
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "32210",
		},
		{
			Key:     models.ServiceItemParamNameDimensionHeight,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "10",
		},
		{
			Key:     models.ServiceItemParamNameDimensionLength,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "12",
		},
		{
			Key:     models.ServiceItemParamNameDimensionWidth,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "3",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDCRT,
		basicPaymentServiceItemParamsDCRT,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    MTOShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)

	dcrtCost := unit.Cents(99999)
	mtoServiceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(appCtx.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: MTOShipment,
	})

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
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
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

	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}, nil)
	posImage := factory.BuildProofOfServiceDoc(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}, nil)
	primeContractor := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, mto.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// MakeHHGMoveWithIntlCratingServiceItemsTOO is a function
// that creates an HHG move with international service items
// from the Prime for review by the TOO
func MakeHHGMoveWithIntlCratingServiceItemsTOO(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	primeUploader := newPrimeUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	actualPickupDate := time.Now().AddDate(0, 0, 1)

	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				Status:               models.MTOShipmentStatusSubmitted,
				ActualPickupDate:     &actualPickupDate,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal: false,
				Status:  models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	_ = factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("86203d72-7f7c-49ff-82f0-5b95e4958f60"), // ICRT - Domestic uncrating
			},
		},
	}, nil)

	_ = factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4132416b-b1aa-42e7-98f2-0ac0a03e8a31"), // IUCRT - Domestic uncrating
			},
		},
	}, nil)

	_ = factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("86203d72-7f7c-49ff-82f0-5b95e4958f60"), // ICRT - Domestic uncrating
			},
		},
	}, nil)

	_ = factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4132416b-b1aa-42e7-98f2-0ac0a03e8a31"), // IUCRT - Domestic uncrating
			},
		},
	}, nil)

	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}, nil)
	posImage := factory.BuildProofOfServiceDoc(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}, nil)
	primeContractor := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, mto.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// MakeHHGMoveForTOOAfterActualPickupDate is a function
// that creates an HHG move with an actual pickup date in the past for diversion testing
// copied almost verbatim from e2ebasic createHHGMoveWithServiceItemsAndPaymentRequestsAndFiles
func MakeHHGMoveForTOOAfterActualPickupDate(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		},
	}, nil)
	sitDaysAllowance := 270
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	twoMonthsAgo := threeMonthsAgo.AddDate(0, 1, 0)
	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				Status:               models.MTOShipmentStatusSubmitted,
				SITDaysAllowance:     &sitDaysAllowance,
				ActualPickupDate:     &twoMonthsAgo,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, mto.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// copied almost verbatim from e2ebasic
func MakePrimeSimulatorMoveNeedsShipmentUpdate(appCtx appcontext.AppContext) models.Move {
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
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

	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipmentFields := models.MTOShipment{
		Status:                models.MTOShipmentStatusSubmitted,
		RequestedPickupDate:   &requestedPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
	}

	firstShipment := factory.BuildMTOShipmentMinimal(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: shipmentFields,
		},
		{
			Model:    pickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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
	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

func MakePrimeSimulatorMoveSameBasePointCity(appCtx appcontext.AppContext) models.Move {
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status:               models.MoveStatusAPPROVED,
				AvailableToPrimeAt:   &now,
				ApprovalsRequestedAt: &now,
				SubmittedAt:          &now,
			},
		},
	}, nil)
	factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
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

	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	pickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "1 First St",
				StreetAddress2: models.StringPointer("Apt 1"),
				City:           "Miami Gardens",
				State:          "FL",
				PostalCode:     "33169",
			},
		},
	}, nil)
	destinationAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Bldg 2"),
				City:           "Key West",
				State:          "FL",
				PostalCode:     "33040",
			},
		},
	}, nil)

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	shipmentFields := models.MTOShipment{
		PrimeEstimatedWeight:  &estimatedWeight,
		PrimeActualWeight:     &actualWeight,
		Status:                models.MTOShipmentStatusApproved,
		RequestedPickupDate:   &requestedPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
	}

	firstShipment := factory.BuildMTOShipmentMinimal(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: shipmentFields,
		},
		{
			Model:    pickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    destinationAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithNTSAndNeedsSC is similar to old shared.createUserWithLocatorAndDODID
func MakeHHGMoveWithNTSAndNeedsSC(appCtx appcontext.AppContext) models.Move {

	submittedAt := time.Now()
	dodID := testdatagen.MakeRandomNumberString(10)
	userInfo := newUserInfo("customer")

	orders := factory.BuildOrderWithoutDefaults(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				Edipi:         models.StringPointer(dodID),
				CacValidated:  true,
			},
		},
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
			Type: &factory.DutyLocations.OriginDutyLocation,
		},
	}, nil)
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status:      models.MoveStatusNeedsServiceCounseling,
				SubmittedAt: &submittedAt,
			},
		},
		{
			Model:    orders,
			LinkOnly: true,
		},
	}, nil)
	// Makes a basic HHG shipment to reflect likely real scenario
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := factory.BuildAddress(appCtx.DB(), nil, nil)
	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
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
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeGoodTACAndLoaCombination builds a good TAC and LOA and returns the TAC
// so that e2e_tests can supply a "Valid" TAC that isn't expired
// or missing a LOA
func MakeGoodTACAndLoaCombination(appCtx appcontext.AppContext) models.TransportationAccountingCode {
	// Transcom Relational Database Management (TRDM) TGET data
	// Creats an active and linked together transportation accounting code and line of accounting
	// Said TAC and LOA are active within a date range of 1 year
	ordersIssueDate := time.Now()
	startDate := ordersIssueDate.AddDate(-1, 0, 0)
	endDate := ordersIssueDate.AddDate(1, 0, 0)
	tacCode := factory.MakeRandomString(4)
	loaSysID := factory.MakeRandomString(10)

	// Ensure all DFAS elements are present
	factory.BuildLineOfAccounting(appCtx.DB(), []factory.Customization{
		{
			Model: models.LineOfAccounting{
				LoaBgnDt:               &startDate,
				LoaEndDt:               &endDate,
				LoaSysID:               &loaSysID,
				LoaHsGdsCd:             models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
				LoaDptID:               models.StringPointer("1"),
				LoaTnsfrDptNm:          models.StringPointer("1"),
				LoaBafID:               models.StringPointer("1"),
				LoaTrsySfxTx:           models.StringPointer("1"),
				LoaMajClmNm:            models.StringPointer("1"),
				LoaOpAgncyID:           models.StringPointer("1"),
				LoaAlltSnID:            models.StringPointer("1"),
				LoaPgmElmntID:          models.StringPointer("1"),
				LoaTskBdgtSblnTx:       models.StringPointer("1"),
				LoaDfAgncyAlctnRcpntID: models.StringPointer("1"),
				LoaJbOrdNm:             models.StringPointer("1"),
				LoaSbaltmtRcpntID:      models.StringPointer("1"),
				LoaWkCntrRcpntNm:       models.StringPointer("1"),
				LoaMajRmbsmtSrcID:      models.StringPointer("1"),
				LoaDtlRmbsmtSrcID:      models.StringPointer("1"),
				LoaCustNm:              models.StringPointer("1"),
				LoaObjClsID:            models.StringPointer("1"),
				LoaSrvSrcID:            models.StringPointer("1"),
				LoaSpclIntrID:          models.StringPointer("1"),
				LoaBdgtAcntClsNm:       models.StringPointer("1"),
				LoaDocID:               models.StringPointer("1"),
				LoaClsRefID:            models.StringPointer("1"),
				LoaInstlAcntgActID:     models.StringPointer("1"),
				LoaLclInstlID:          models.StringPointer("1"),
				LoaFmsTrnsactnID:       models.StringPointer("1"),
				LoaTrnsnID:             models.StringPointer("1"),
				LoaUic:                 models.StringPointer("1"),
				LoaBgFyTx:              models.IntPointer(2023),
				LoaEndFyTx:             models.IntPointer(2025),
			},
		},
	}, nil)
	// Create the TAC and associate loa based on LoaSysID
	tac := factory.BuildTransportationAccountingCodeWithoutAttachedLoa(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				TAC:               tacCode,
				TrnsprtnAcntBgnDt: &startDate,
				TrnsprtnAcntEndDt: &endDate,
				TacFnBlModCd:      models.StringPointer("1"),
				LoaSysID:          &loaSysID,
			},
		},
	}, nil)
	return tac
}

// MakeNTSRMoveWithPaymentRequest is similar to old shared.createNTSRMoveWithPaymentRequest
func MakeNTSRMoveWithPaymentRequest(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	currentTime := time.Now()
	tac := "1111"

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
	}, nil)

	// Create Orders
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.Order{
				TAC: &tac,
			},
		},
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

	// Create Move
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				SubmittedAt:        &currentTime,
			},
		},
	}, nil)
	// Create Pickup Address
	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85004",
			},
		},
	}, nil)

	// Create Storage Facility
	storageFacility := factory.BuildStorageFacility(appCtx.DB(), nil, []factory.Trait{
		factory.GetTraitStorageFacilityKKFA,
	})
	// Create NTS-R Shipment
	tacType := models.LOATypeHHG
	serviceOrderNumber := testdatagen.MakeRandomNumberString(4)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	ntsRecordedWeight := unit.Pound(2000)
	ntsrShipment := factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				NTSRecordedWeight:    &ntsRecordedWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				TACType:              &tacType,
				Status:               models.MTOShipmentStatusApproved,
				ServiceOrderNumber:   &serviceOrderNumber,
				UsesExternalVendor:   true,
			},
		},
	}, nil)

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	// Create Payment Request
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// create service item
	msCostcos := unit.Cents(32400)
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeCS,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			}},
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &msCostcos,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// MakeHHGMoveWithServiceItemsandPaymentRequestsForTIO is basically
// scenario.createMoveWithServiceItemsandPaymentRequests01
func MakeHHGMoveWithServiceItemsandPaymentRequestsForTIO(appCtx appcontext.AppContext) models.Move {
	/*
		Creates a move for the TIO flow
	*/
	userUploader := newUserUploader(appCtx)

	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)
	csCost := unit.Cents(25000)
	fscCost := unit.Cents(55555)

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
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

	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)

	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// This is a postal code that maps to the default office user gbloc KKFA in the PostalCodeToGBLOC table
				PostalCode: "85004",
			},
		},
	}, nil)

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	mtoShipmentHHG := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	paymentRequestHHG := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	// for soft deleted proof of service docs
	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequestHHG,
			LinkOnly: true,
		},
	}, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

	serviceItemMS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	// Shuttling service item
	doshutCost := unit.Cents(623)
	approvedAtTime := time.Now()
	serviceItemDOSHUT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTime,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &doshutCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOSHUT,
			LinkOnly: true,
		},
	}, nil)

	currentTime := time.Now()

	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Pound(4000))),
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOSHUT,
		basicPaymentServiceItemParams,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Crating service item
	dcrtCost := unit.Cents(623)
	approvedAtTimeCRT := time.Now()
	serviceItemDCRT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTimeCRT,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Dom Crating
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	currentTimeDCRT := time.Now()

	basicPaymentServiceItemParamsDCRT := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractYearName,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameEscalationCompounded,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.FormatFloat(1.125, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "1.71",
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameCubicFeetBilled,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "4.00",
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "32210",
		},
		{
			Key:     models.ServiceItemParamNameDimensionHeight,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "10",
		},
		{
			Key:     models.ServiceItemParamNameDimensionLength,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "12",
		},
		{
			Key:     models.ServiceItemParamNameDimensionWidth,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "3",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDCRT,
		basicPaymentServiceItemParamsDCRT,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Domestic line haul service item
	serviceItemDLH := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	createdAtTime := time.Now().Add(time.Duration(time.Hour * -24))
	additionalPaymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
				SequenceNumber:  2,
				CreatedAt:       createdAtTime,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	serviceItemCS := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)
	serviceItemFSC := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
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

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)
	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, mto.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

func MakeHHGMoveWithServiceItemsandPaymentRequestReviewedForQAE(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)
	csCost := unit.Cents(25000)
	fscCost := unit.Cents(55555)

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
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

	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)

	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// This is a postal code that maps to the default office user gbloc KKFA in the PostalCodeToGBLOC table
				PostalCode: "85004",
			},
		},
	}, nil)

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	mtoShipmentHHG := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	paymentRequestHHG := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	// for soft deleted proof of service docs
	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequestHHG,
			LinkOnly: true,
		},
	}, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

	serviceItemMS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	// Shuttling service item
	doshutCost := unit.Cents(623)
	approvedAtTime := time.Now()
	serviceItemDOSHUT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTime,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &doshutCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOSHUT,
			LinkOnly: true,
		},
	}, nil)

	currentTime := time.Now()

	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Pound(4000))),
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOSHUT,
		basicPaymentServiceItemParams,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Crating service item
	dcrtCost := unit.Cents(623)
	approvedAtTimeCRT := time.Now()
	serviceItemDCRT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTimeCRT,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Dom Crating
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	currentTimeDCRT := time.Now()

	basicPaymentServiceItemParamsDCRT := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractYearName,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameEscalationCompounded,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.FormatFloat(1.125, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "1.71",
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameCubicFeetBilled,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "4.00",
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "32210",
		},
		{
			Key:     models.ServiceItemParamNameDimensionHeight,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "10",
		},
		{
			Key:     models.ServiceItemParamNameDimensionLength,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "12",
		},
		{
			Key:     models.ServiceItemParamNameDimensionWidth,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "3",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDCRT,
		basicPaymentServiceItemParamsDCRT,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Domestic line haul service item
	serviceItemDLH := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	createdAtTime := time.Now().Add(time.Duration(time.Hour * -24))
	additionalPaymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
				SequenceNumber:  2,
				CreatedAt:       createdAtTime,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	serviceItemCS := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)
	serviceItemFSC := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
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

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)
	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, mto.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

func MakeHHGMoveWithServiceItemsandPaymentRequestWithDocsReviewedForQAE(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	primeUploader := newPrimeUploader(appCtx)

	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)
	csCost := unit.Cents(25000)
	fscCost := unit.Cents(55555)

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
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

	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)

	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// This is a postal code that maps to the default office user gbloc KKFA in the PostalCodeToGBLOC table
				PostalCode: "85004",
			},
		},
	}, nil)

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	mtoShipmentHHG := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	paymentRequestHHG := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequestHHG,
			LinkOnly: true,
		},
	}, nil)
	posImage := factory.BuildProofOfServiceDoc(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequestHHG,
			LinkOnly: true,
		},
	}, nil)
	primeContractor := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	serviceItemMS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	// Shuttling service item
	doshutCost := unit.Cents(623)
	approvedAtTime := time.Now()
	serviceItemDOSHUT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTime,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &doshutCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOSHUT,
			LinkOnly: true,
		},
	}, nil)

	currentTime := time.Now()

	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Pound(4000))),
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOSHUT,
		basicPaymentServiceItemParams,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Crating service item
	dcrtCost := unit.Cents(623)
	approvedAtTimeCRT := time.Now()
	serviceItemDCRT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTimeCRT,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Dom Crating
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	currentTimeDCRT := time.Now()

	basicPaymentServiceItemParamsDCRT := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractYearName,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameEscalationCompounded,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.FormatFloat(1.125, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "1.71",
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameCubicFeetBilled,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "4.00",
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "32210",
		},
		{
			Key:     models.ServiceItemParamNameDimensionHeight,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "10",
		},
		{
			Key:     models.ServiceItemParamNameDimensionLength,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "12",
		},
		{
			Key:     models.ServiceItemParamNameDimensionWidth,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "3",
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDCRT,
		basicPaymentServiceItemParamsDCRT,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Domestic line haul service item
	serviceItemDLH := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	createdAtTime := time.Now().Add(time.Duration(time.Hour * -24))
	additionalPaymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
				SequenceNumber:  2,
				CreatedAt:       createdAtTime,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	serviceItemCS := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)
	serviceItemFSC := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model:    mto,
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

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		},
	}, nil)
	posImage2 := factory.BuildProofOfServiceDoc(appCtx.DB(), []factory.Customization{
		{
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		},
	}, nil)
	primeContractor2 := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")

	// Creates custom test.jpg prime upload
	file2 := testdatagen.Fixture("test.jpg")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage2.ID, primeContractor2, uploader.File{File: file2}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, mto.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// like scenario.createNTSRMoveWithServiceItemsAndPaymentRequest
func MakeNTSRMoveWithServiceItemsAndPaymentRequest(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	currentTime := time.Now()
	tac := "1111"
	tac2 := "2222"
	sac := "3333"
	sac2 := "4444"

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
	}, nil)

	// Create Orders
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.Order{
				TAC:    &tac,
				NtsTAC: &tac2,
				SAC:    &sac,
				NtsSAC: &sac2,
			},
		},
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

	// Create Move
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				SubmittedAt:        &currentTime,
			},
		},
	}, nil)
	// Create Pickup Address
	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85004",
			},
		},
	}, nil)

	// Create Storage Facility
	storageFacility := factory.BuildStorageFacility(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85005",
			},
		},
	}, nil)
	// Create NTS-R Shipment
	tacType := models.LOATypeHHG
	sacType := models.LOATypeNTS
	serviceOrderNumber := "1234"
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	ntsrShipment := factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				TACType:              &tacType,
				Status:               models.MTOShipmentStatusApproved,
				SACType:              &sacType,
				ServiceOrderNumber:   &serviceOrderNumber,
			},
		},
	}, nil)

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	// Create Payment Request
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// Create Domestic linehaul service item
	dlCost := unit.Cents(80000)
	dlItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DL Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDLH,
		dlItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &dlCost,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Fuel surcharge service item
	fsCost := unit.Cents(10700)
	fsItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "FS Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeFSC,
		fsItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &fsCost,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic origin price service item
	doCost := unit.Cents(15000)
	doItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DO Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOP,
		doItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &doCost,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic destination price service item
	ddpCost := unit.Cents(15000)
	ddpItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DDP Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDDP,
		ddpItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &ddpCost,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic unpacking service item
	duCost := unit.Cents(45900)
	duItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DUPK Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "5.79",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDUPK,
		duItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &duCost,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// MakeHHGMoveWithRetireeForTOO creates a retiree move for TOO
func MakeHHGMoveWithRetireeForTOO(appCtx appcontext.AppContext) models.Move {
	retirement := internalmessages.OrdersTypeRETIREMENT
	hhg := models.MTOShipmentTypeHHG
	hor := models.DestinationTypeHomeOfRecord
	originDutyLocation := factory.FetchOrBuildCurrentDutyLocation(appCtx.DB())
	move := scenario.CreateMoveWithOptions(appCtx, testdatagen.Assertions{
		Order: models.Order{
			OrdersType:         retirement,
			OriginDutyLocation: &originDutyLocation,
		},
		MTOShipment: models.MTOShipment{
			ShipmentType:    hhg,
			DestinationType: &hor,
		},
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: false,
		},
	})

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithNTSShipmentsForTOO creates an HHG Move with NTS Shipment
func MakeHHGMoveWithNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSShipments(appCtx, locator, false)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithPPMShipmentsForTOO creates an HHG Move with a PPM shipment.
func MakeHHGMoveWithPPMShipmentsForTOO(appCtx appcontext.AppContext, readyForCloseout bool) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)
	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}
	move := scenario.CreateMoveWithHHGAndPPM(appCtx, userUploader, moveInfo, models.AffiliationARMY, readyForCloseout)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithExternalNTSShipmentsForTOO creates an HHG Move with
// NTS Shipment by external vendor
func MakeHHGMoveWithExternalNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSShipments(appCtx, locator, true)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithApprovedNTSShipmentsForTOO creates an HHG Move with approved NTS
// Shipments
func MakeHHGMoveWithApprovedNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSShipments(appCtx, locator, false)

	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	err := moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic("Failed to approve move: %w", err)
	}

	err = appCtx.DB().Save(&move)
	if err != nil {
		log.Panic("Failed to save move: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	orders := newmove.Orders
	orders.SAC = models.StringPointer("4K988AS098F")
	orders.TAC = models.StringPointer("E15A")
	orders.NtsSAC = models.StringPointer("3L988AS098F")
	orders.NtsTAC = models.StringPointer("F123")
	err = appCtx.DB().Save(&orders)
	if err != nil {
		log.Panic("Failed to save orders: %w", err)
	}

	planner := &routemocks.Planner{}

	// mock any and all planner calls
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(2361, nil)

	queryBuilder := query.NewQueryBuilder()
	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
	shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)

	updatedShipments := make([]*models.MTOShipment, len(newmove.MTOShipments))
	for i := range newmove.MTOShipments {
		shipment := newmove.MTOShipments[i]
		updatedShipments[i], err = shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, etag.GenerateEtag(shipment.UpdatedAt))
		if err != nil {
			log.Panic("Error updating shipment status %w", err)
		}
	}

	storageFacility := factory.BuildStorageFacility(appCtx.DB(), nil, nil)

	updatedShipment := updatedShipments[1]

	sacType := models.LOATypeHHG
	updatedShipment.SACType = &sacType
	tacType := models.LOATypeNTS
	updatedShipment.TACType = &tacType
	updatedShipment.ServiceOrderNumber = models.StringPointer("999999")
	updatedShipment.StorageFacilityID = &storageFacility.ID
	err = appCtx.DB().Save(updatedShipment)
	if err != nil {
		log.Panic("Failed to save shipment: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err = models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeMoveWithNTSShipmentsForTOO creates an HHG Move with NTS Shipment
func MakeMoveWithNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithNTSShipment(appCtx, locator, true)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithNTSSRhipmentsForTOO creates an HHG Move with NTS-R Shipment
func MakeHHGMoveWithNTSRShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSRShipments(appCtx, locator, false)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithApprovedNTSShipmentsForTOO creates an HHG Move with approved NTS
// Shipments
func MakeHHGMoveWithApprovedNTSRShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSRShipments(appCtx, locator, false)

	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	err := moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic("Failed to approve move: %w", err)
	}

	err = appCtx.DB().Save(&move)
	if err != nil {
		log.Panic("Failed to save move: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	orders := newmove.Orders
	orders.SAC = models.StringPointer("4K988AS098F")
	orders.TAC = models.StringPointer("E15A")
	orders.NtsSAC = models.StringPointer("3L988AS098F")
	orders.NtsTAC = models.StringPointer("F123")
	err = appCtx.DB().Save(&orders)
	if err != nil {
		log.Panic("Failed to save orders: %w", err)
	}

	planner := &routemocks.Planner{}

	// mock any and all planner calls
	planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(2361, nil)

	queryBuilder := query.NewQueryBuilder()
	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
	shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)

	updatedShipments := make([]*models.MTOShipment, len(newmove.MTOShipments))
	for i := range newmove.MTOShipments {
		shipment := newmove.MTOShipments[i]
		updatedShipments[i], err = shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, nil, etag.GenerateEtag(shipment.UpdatedAt))
		if err != nil {
			log.Panic("Error updating shipment status %w", err)
		}
	}

	storageFacility := factory.BuildStorageFacility(appCtx.DB(), nil, nil)

	updatedShipment := updatedShipments[1]

	sacType := models.LOATypeHHG
	updatedShipment.SACType = &sacType
	tacType := models.LOATypeNTS
	updatedShipment.TACType = &tacType
	updatedShipment.ServiceOrderNumber = models.StringPointer("999999")
	updatedShipment.StorageFacilityID = &storageFacility.ID
	err = appCtx.DB().Save(updatedShipment)
	if err != nil {
		log.Panic("Failed to save shipment: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err = models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithExternalNTSRShipmentsForTOO creates an HHG Move with
// NTS Shipment by external vendor
func MakeHHGMoveWithExternalNTSRShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSRShipments(appCtx, locator, true)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeMoveWithMinimalNTSRNeedsSC creates an Move with
// NTS-R Shipment
func MakeMoveWithMinimalNTSRNeedsSC(appCtx appcontext.AppContext) models.Move {
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounselingMinimalNTSR(appCtx, pcos, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveNeedsSC creates an fully ready move needing SC approval
func MakeHHGMoveNeedsSC(appCtx appcontext.AppContext) models.Move {
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hhg := models.MTOShipmentTypeHHG
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounseling(appCtx, pcos, hhg, nil, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeBoatHaulAwayMoveNeedsSC creates an fully ready move with a boat haul-away shipment needing SC approval
func MakeBoatHaulAwayMoveNeedsSC(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())

	move := scenario.CreateBoatHaulAwayMoveForSC(appCtx, userUploader, moveRouter, moveInfo)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

// MakeBoatHaulAwayMoveNeedsTOOApproval creates an fully ready move with a boat haul-away shipment needing SC approval
func MakeBoatHaulAwayMoveNeedsTOOApproval(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())

	move := scenario.CreateBoatHaulAwayMoveForTOO(appCtx, userUploader, moveRouter, moveInfo)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

// MakeHHGMoveNeedsSC creates an fully ready move needing SC approval
func MakeMobileHomeMoveNeedsSC(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithMTOShipment(appCtx, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION, models.MTOShipmentTypeMobileHome, nil, locator, models.MoveStatusNeedsServiceCounseling)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

func MakeMobileHomeMoveForTOO(appCtx appcontext.AppContext) models.Move {
	hhg := models.MTOShipmentTypeHHG
	hor := models.DestinationTypeHomeOfRecord
	originDutyLocation := factory.FetchOrBuildCurrentDutyLocation(appCtx.DB())
	move := scenario.CreateMoveWithOptions(appCtx, testdatagen.Assertions{
		Order: models.Order{
			OriginDutyLocation: &originDutyLocation,
		},
		MTOShipment: models.MTOShipment{
			ShipmentType:    hhg,
			DestinationType: &hor,
		},
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: false,
		},
	})

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveNeedsServicesCounselingUSMC creates an fully ready move as USMC needing SC approval
func MakeHHGMoveNeedsServicesCounselingUSMC(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	locator := models.GenerateLocator()
	move := scenario.CreateHHGNeedsServicesCounselingUSMC3(appCtx, userUploader, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithAmendedOrders creates a move needing SC approval with amended orders
func MakeHHGMoveWithAmendedOrders(appCtx appcontext.AppContext) models.Move {
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hhg := models.MTOShipmentTypeHHG
	locator := models.GenerateLocator()
	userUploader := newUserUploader(appCtx)
	move := scenario.CreateNeedsServicesCounselingWithAmendedOrders(appCtx, userUploader, pcos, hhg, nil, locator)
	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveForSeparationNeedsSC creates an fully ready move for
// separation needing SC approval
func MakeHHGMoveForSeparationNeedsSC(appCtx appcontext.AppContext) models.Move {
	separation := internalmessages.OrdersTypeSEPARATION
	hhg := models.MTOShipmentTypeHHG
	hor := models.DestinationTypeHomeOfRecord
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounseling(appCtx, separation, hhg, &hor, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveForRetireeNeedsSC creates an fully ready move for
// separation needing SC approval
func MakeHHGMoveForRetireeNeedsSC(appCtx appcontext.AppContext) models.Move {
	retirement := internalmessages.OrdersTypeRETIREMENT
	hhg := models.MTOShipmentTypeHHG
	hos := models.DestinationTypeHomeOfSelection
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounseling(appCtx, retirement, hhg, &hos, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

func MakeMoveWithPPMShipmentReadyForFinalCloseout(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	pickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "1 First St",
				StreetAddress2: models.StringPointer("Apt 1"),
				City:           "Miami Gardens",
				State:          "FL",
				PostalCode:     "33169",
			},
		},
	}, nil)
	destinationAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				ID:             uuid.Must(uuid.NewV4()),
				StreetAddress1: "2 Second St",
				StreetAddress2: models.StringPointer("Bldg 2"),
				City:           "Key West",
				State:          "FL",
				PostalCode:     "33040",
			},
		},
	}, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          uuid.Must(uuid.NewV4()),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			PickupAddressID:             &pickupAddress.ID,
			DestinationAddressID:        &destinationAddress.ID,
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
			FinalIncentive:              models.CentPointer(50000000),
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	newmove.Orders.NewDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), newmove.Orders.NewDutyLocationID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch duty location: %w", err))
	}
	return *newmove
}

// This one is the actual function that's used for testdatagen harness(I think)
func MakeMoveWithPPMShipmentReadyForFinalCloseoutWithSIT(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	sitLocationType := models.SITLocationTypeOrigin
	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	sitDaysAllowance := 90
	pickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				PostalCode: "42444",
			},
		},
	}, nil)
	destinationAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				PostalCode: "30813",
			},
		},
	}, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusApproved,
			SITDaysAllowance:     &sitDaysAllowance,
			PickupAddressID:      &pickupAddress.ID,
			DestinationAddressID: &destinationAddress.ID,
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
			FinalIncentive:              models.CentPointer(50000000),
			SITExpected:                 models.BoolPointer(true),
			SITEstimatedEntryDate:       models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			SITEstimatedDepartureDate:   models.TimePointer(time.Date(testdatagen.GHCTestYear, time.April, 16, 0, 0, 0, 0, time.UTC)),
			SITEstimatedWeight:          models.PoundPointer(unit.Pound(1234)),
			SITEstimatedCost:            models.CentPointer(unit.Cents(12345600)),
			SITLocation:                 &sitLocationType,
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	twoMonthsAgo := threeMonthsAgo.AddDate(0, 1, 0)
	sitCost := unit.Cents(200000)
	sitItems := factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment.Shipment, &threeMonthsAgo, &twoMonthsAgo)
	sitItems = append(sitItems, factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment.Shipment, &twoMonthsAgo, nil)...)
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	for i := range sitItems {
		factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &sitCost,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			}, {
				Model:    sitItems[i],
				LinkOnly: true,
			},
		}, nil)
	}
	scenario.MakeSITExtensionsForShipment(appCtx, shipment.Shipment)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	newmove.Orders.NewDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), newmove.Orders.NewDutyLocationID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch duty location: %w", err))
	}
	return *newmove
}

func MakePPMMoveWithCloseout(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)
	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	move := scenario.CreateMoveWithCloseOut(appCtx, userUploader, moveInfo, models.AffiliationARMY)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}
	return *newmove
}

func MakePPMMoveWithCloseoutOffice(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	move := scenario.CreateMoveWithCloseoutOffice(appCtx, moveInfo, userUploader)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	var closeoutOffice models.TransportationOffice
	err = appCtx.DB().Find(&closeoutOffice, newmove.CloseoutOfficeID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch closeout office: %w", err))
	}

	newmove.CloseoutOffice = &closeoutOffice
	return *newmove
}

func MakeSubmittedMoveWithPPMShipmentForSC(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())

	move := scenario.CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter, moveInfo)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPM(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         uuid.Must(uuid.NewV4()),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	move, _ := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx appcontext.AppContext) models.Move {
	/*
	 * A service member with orders and a PPM shipment updated with an estimated weight value and estimated incentive
	 */
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &scenario.DefaultCloseoutOfficeID,
	}
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                 uuid.Must(uuid.NewV4()),
			EstimatedWeight:    models.PoundPointer(unit.Pound(4000)),
			HasProGear:         models.BoolPointer(false),
			EstimatedIncentive: models.CentPointer(unit.Cents(1000000)),
		},
	}

	move, _ := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, true, userUploader, nil, nil, assertions.PPMShipment)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMWithAboutFormComplete(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &scenario.DefaultCloseoutOfficeID,
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
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

	move, _ := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeUnsubmittedMoveWithMultipleFullPPMShipmentComplete(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &scenario.DefaultCloseoutOfficeID,
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.PPMShipmentStatusDraft,
		},
	}

	move, _ := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)

	factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMProgearWeightTicket(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
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

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMProgearWeightTicketOffice(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
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
			W2Address:                   &address,
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMProgearWeightTicketOfficeCivilian(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	order := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.Order{
				Grade: models.ServiceMemberGradeCIVILIANEMPLOYEE.Pointer(),
			},
		},
	}, nil)

	move := models.Move{
		Status:   models.MoveStatusAPPROVED,
		OrdersID: order.ID,
	}

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move:         move,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                           uuid.Must(uuid.NewV4()),
			ApprovedAt:                   &approvedAt,
			Status:                       models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:               models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:       models.StringPointer("42444"),
			ActualDestinationPostalCode:  models.StringPointer("30813"),
			HasReceivedAdvance:           models.BoolPointer(true),
			AdvanceAmountReceived:        models.CentPointer(unit.Cents(340000)),
			W2Address:                    &address,
			IsActualExpenseReimbursement: models.BoolPointer(true),
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMWeightTicketOffice(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)
	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
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
			W2Address:                   &address,
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMWeightTicketOfficeWithHHG(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)
	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
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
			W2Address:                   &address,
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMMovingExpense(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	storageStart := time.Now()
	storageEnd := storageStart.Add(7 * time.Hour * 24)
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
			ExpectedDepartureDate:       time.Date(testdatagen.GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
			SITEstimatedEntryDate:       &storageStart,
			SITEstimatedDepartureDate:   &storageEnd,
			SITExpected:                 models.BoolPointer(true),
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	twoMonthsAgo := threeMonthsAgo.AddDate(0, 1, 0)
	sitCost := unit.Cents(200000)
	sitItems := factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment.Shipment, &threeMonthsAgo, &twoMonthsAgo)
	sitItems = append(sitItems, factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment.Shipment, &twoMonthsAgo, nil)...)
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	for i := range sitItems {
		factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &sitCost,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			}, {
				Model:    sitItems[i],
				LinkOnly: true,
			},
		}, nil)
	}

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
				SITStartDate:      &storageStart,
				SITEndDate:        &storageEnd,
				SITLocation:       &sitLocation,
				WeightStored:      (*unit.Pound)(&weightStored),
			},
		},
	}, nil)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMMovingExpenseOffice(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)
	storageStart := time.Now()
	storageEnd := storageStart.Add(7 * time.Hour * 24)
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
			W2Address:                   &address,
			ExpectedDepartureDate:       time.Date(testdatagen.GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
			SITEstimatedEntryDate:       &storageStart,
			SITEstimatedDepartureDate:   &storageEnd,
			SITExpected:                 models.BoolPointer(true),
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	twoMonthsAgo := threeMonthsAgo.AddDate(0, 1, 0)
	sitCost := unit.Cents(200000)
	sitItems := factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment.Shipment, &threeMonthsAgo, &twoMonthsAgo)
	sitItems = append(sitItems, factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment.Shipment, &twoMonthsAgo, nil)...)
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	for i := range sitItems {
		factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &sitCost,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			}, {
				Model:    sitItems[i],
				LinkOnly: true,
			},
		}, nil)
	}

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
				SITStartDate:      &storageStart,
				SITEndDate:        &storageEnd,
				SITLocation:       &sitLocation,
				WeightStored:      (*unit.Pound)(&weightStored),
			},
		},
	}, nil)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMAllDocTypesOffice(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
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
			Status:                      models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

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

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

// the old serviceMemberWithOrdersAndPPMMove
func MakeDraftMoveWithPPMWithDepartureDate(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	departureDate := time.Date(2022, time.February, 01, 0, 0, 0, 0, time.UTC)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                    uuid.Must(uuid.NewV4()),
			Status:                models.PPMShipmentStatusDraft,
			EstimatedWeight:       models.PoundPointer(unit.Pound(4000)),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			ExpectedDepartureDate: departureDate,
		},
	}

	move, _ := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, nil, nil, assertions.PPMShipment)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPMShipmentAndExcessWeight(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{Gbloc: "KKFA", ProvidesCloseout: true},
		},
	}, nil)

	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            "excessweightsPPM@ppm.approved",
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        "One PPM",
		LastName:         "ExcessWeights",
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &closeoutOffice.ID,
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
			Status:                      models.PPMShipmentStatusNeedsCloseout,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			AdvanceStatus:               (*models.PPMAdvanceStatus)(models.StringPointer(string(models.PPMAdvanceStatusApproved))),
			W2Address:                   &address,
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, &assertions.Move, assertions.PPMShipment)

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
		}}, nil)
	return move
}

func MakeHHGMoveInSIT(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	twoMonthsAgo := now.AddDate(0, -2, 0)
	oneMonthAgo := now.AddDate(0, -1, 0)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsAgo, &oneMonthAgo)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &oneMonthAgo, nil)

	return move
}

// Creates an HHG move with a past Origin and Destination SIT
func HHGMoveWithPastSITs(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	fourMonthsAgo := now.AddDate(0, -4, 0)
	threeMonthsAgo := now.AddDate(0, -3, 0)
	twoMonthsAgo := now.AddDate(0, -2, 0)
	oneMonthAgo := now.AddDate(0, -1, 0)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &fourMonthsAgo, &threeMonthsAgo)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsAgo, &oneMonthAgo)

	return move
}

func MakeHHGMoveInSITNoExcessWeight(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(1350)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	twoMonthsAgo := now.AddDate(0, -2, 0)
	oneMonthAgo := now.AddDate(0, -1, 0)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsAgo, &oneMonthAgo)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &oneMonthAgo, nil)

	return move
}

func MakeHHGMoveInSITWithPendingExtension(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	twoMonthsAgo := now.AddDate(0, -2, 0)
	oneMonthAgo := now.AddDate(0, -1, 0)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsAgo, &oneMonthAgo)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &oneMonthAgo, nil)
	factory.BuildSITDurationUpdate(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
	}, nil)

	return move
}

func MakeHHGMoveInSITWithAddressChangeRequestOver50Miles(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)

	sitDaysAllowance := 90
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model: models.Entitlement{
				DependentsAuthorized: models.BoolPointer(true),
				StorageInTransit:     &sitDaysAllowance,
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

	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(now),
			},
		},
	}, nil)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  models.PoundPointer(unit.Pound(1400)),
				PrimeActualWeight:     models.PoundPointer(unit.Pound(2000)),
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	twoMonthsAgo := now.AddDate(0, -2, 0)
	oneMonthAgo := now.AddDate(0, -1, 0)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsAgo, &oneMonthAgo)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &oneMonthAgo, nil)

	newMove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newMove
}

func MakeHHGMoveInSITWithAddressChangeRequestUnder50Miles(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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

	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  models.PoundPointer(unit.Pound(1400)),
				PrimeActualWeight:     models.PoundPointer(unit.Pound(2000)),
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	twoMonthsAgo := now.AddDate(0, -2, 0)
	oneMonthAgo := now.AddDate(0, -1, 0)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsAgo, &oneMonthAgo)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &oneMonthAgo, nil)

	newMove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	return *newMove
}

func MakeHHGMoveInSITEndsToday(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")
	actualPickupDate := time.Now().AddDate(0, 0, 1)

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
				ActualPickupDate:      &actualPickupDate,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	daysAgo90 := now.AddDate(0, 0, -90)
	daysAgo45 := now.AddDate(0, 0, -45)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &daysAgo90, &daysAgo45)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &daysAgo45, nil)

	return move
}

func MakeHHGMoveInSITEndsTomorrow(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	daysAgo89 := now.AddDate(0, 0, -89)
	daysAgo44 := now.AddDate(0, 0, -44)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &daysAgo89, &daysAgo44)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &daysAgo44, nil)

	return move
}

func MakeHHGMoveInSITEndsYesterday(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	daysAgo91 := now.AddDate(0, 0, -91)
	daysAgo46 := now.AddDate(0, 0, -46)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &daysAgo91, &daysAgo46)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &daysAgo46, nil)

	return move
}

func MakeHHGMoveInSITDeparted(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	daysAgo93 := now.AddDate(0, 0, -93)
	daysAgo48 := now.AddDate(0, 0, -48)
	daysAgo5 := now.AddDate(0, 0, -5)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &daysAgo93, &daysAgo48)
	factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &daysAgo48, &daysAgo5)

	return move
}

func MakeHHGMoveInSITStartsInFuture(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	daysLater100 := now.AddDate(0, 0, 100)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &daysLater100, nil)

	return move
}

func MakeHHGMoveInSITNotApproved(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	// pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	oneMonthLater := now.AddDate(0, 1, 0)
	twoMonthsLater := now.AddDate(0, 2, 0)
	sitItems := factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &oneMonthLater, nil)
	sitItems = append(sitItems, factory.BuildDestSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsLater, nil)...)
	for i := range sitItems {
		sitItems[i].Status = models.MTOServiceItemStatusSubmitted
		err := appCtx.DB().Update(&sitItems[i])
		if err != nil {
			log.Panic(fmt.Errorf("failed to update sit service item: %w", err))
		}
	}

	return move
}

func MakeHHGMoveWithAddressChangeRequest(appCtx appcontext.AppContext) models.ShipmentAddressUpdate {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
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

	originalDeliveryAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "7 Q st",
				StreetAddress2: models.StringPointer("Apt 1"),
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30813",
			},
		},
	}, nil)

	shipmentAddressUpdate := factory.BuildShipmentAddressUpdate(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRequested,
			},
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    originalDeliveryAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	return shipmentAddressUpdate
}

func MakeHHGMoveWithAddressChangeRequestAndUnknownDeliveryAddress(appCtx appcontext.AppContext) models.ShipmentAddressUpdate {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
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

	destinationAddress := factory.BuildMinimalAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				City:       orders.OriginDutyLocation.Address.City,
				State:      orders.OriginDutyLocation.Address.State,
				PostalCode: orders.OriginDutyLocation.Address.PostalCode,
				Country:    orders.OriginDutyLocation.Address.Country,
			},
		},
	}, nil)

	shipmentAddressUpdate := factory.BuildShipmentAddressUpdate(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRequested,
			},
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    destinationAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	return shipmentAddressUpdate
}

func MakeHHGMoveWithAddressChangeRequestAndSecondDeliveryLocation(appCtx appcontext.AppContext) models.ShipmentAddressUpdate {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
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

	secondaryDeliveryAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "123 2nd Address",
			},
		},
	}, nil)

	tertiaryDeliveryAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "123 3rd Address",
			},
		},
	}, nil)

	originalDeliveryAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "7 Q st",
				StreetAddress2: models.StringPointer("Apt 1"),
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30813",
			},
		},
	}, nil)

	shipmentAddressUpdate := factory.BuildShipmentAddressUpdate(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRequested,
			},
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    tertiaryDeliveryAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.TertiaryDeliveryAddress,
		},
		{
			Model:    secondaryDeliveryAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.SecondaryDeliveryAddress,
		},
		{
			Model:    originalDeliveryAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	return shipmentAddressUpdate
}

func MakeNTSRMoveWithAddressChangeRequest(appCtx appcontext.AppContext) models.ShipmentAddressUpdate {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
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

	originalDeliveryAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "7 Q st",
				StreetAddress2: models.StringPointer("Apt 1"),
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30813",
			},
		},
	}, nil)

	now := time.Now()
	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)

	NTSRecordedWeight := unit.Pound(1400)
	serviceOrderNumber := "1234"
	shipmentAddressUpdate := factory.BuildShipmentAddressUpdate(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRequested,
			},
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				Status:                models.MTOShipmentStatusApproved,
				ShipmentType:          models.MTOShipmentTypeHHGOutOfNTSDom,
				NTSRecordedWeight:     &NTSRecordedWeight,
				ServiceOrderNumber:    &serviceOrderNumber,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
			},
		},
		{
			Model:    originalDeliveryAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	return shipmentAddressUpdate
}

func MakeMoveReadyForEDI(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	currentTime := time.Now()

	// Create Army Customer
	userInfo := newUserInfo("customer")
	userAffiliation := models.AffiliationARMY
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				Affiliation:   &userAffiliation,
				CacValidated:  true,
			},
		},
	}, nil)

	// Create LOA and TAC
	sixMonthsBefore := currentTime.AddDate(0, -6, 0)
	sixMonthsAfter := currentTime.AddDate(0, 6, 0)
	loa := factory.BuildFullLineOfAccounting(appCtx.DB(), []factory.Customization{
		{
			Model: models.LineOfAccounting{
				LoaBgnDt: &sixMonthsBefore,
				LoaEndDt: &sixMonthsAfter,
			},
		},
	}, nil)

	tac := factory.BuildTransportationAccountingCode(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				TacFnBlModCd: models.StringPointer("W"),
			},
		}, {
			Model:    loa,
			LinkOnly: true,
		},
	}, nil)

	// Create Orders
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.Order{
				TAC:       &tac.TAC,
				IssueDate: currentTime,
			},
		},
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

	// Create Move
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Status:             models.MoveStatusAPPROVED,
			},
		},
	}, nil)
	// Create Pickup Address
	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85004",
			},
		},
	}, nil)

	serviceOrderNumber := "1234"
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	ntsrShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
				ServiceOrderNumber:   &serviceOrderNumber,
			},
		},
	}, nil)

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	// Create Payment Request
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// Create Domestic linehaul service item
	dlCost := unit.Cents(80000)
	dlItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DL Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDLH,
		dlItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &dlCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Fuel surcharge service item
	fsCost := unit.Cents(10700)
	fsItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "FS Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeFSC,
		fsItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &fsCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic origin price service item
	doCost := unit.Cents(15000)
	doItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DO Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOP,
		doItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &doCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic destination price service item
	ddpCost := unit.Cents(15000)
	ddpItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DDP Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDDP,
		ddpItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &ddpCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic unpacking service item
	duCost := unit.Cents(45900)
	duItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DUPK Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "5.79",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDUPK,
		duItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &duCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// re-fetch the move so that we ensure we have exactly what is in the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

func MakeCoastGuardMoveReadyForEDI(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	currentTime := time.Now()

	// Create Coast Guard Customer
	userInfo := newUserInfo("customer")
	userAffiliation := models.AffiliationCOASTGUARD
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				Affiliation:   &userAffiliation,
				CacValidated:  true,
			},
		},
	}, nil)

	// Create LOA and TAC
	sixMonthsBefore := currentTime.AddDate(0, -6, 0)
	sixMonthsAfter := currentTime.AddDate(0, 6, 0)

	loa := factory.BuildFullLineOfAccounting(appCtx.DB(), []factory.Customization{
		{
			Model: models.LineOfAccounting{
				LoaHsGdsCd: models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeNTS),
				LoaBgnDt:   &sixMonthsBefore,
				LoaEndDt:   &sixMonthsAfter,
			},
		},
	}, nil)

	tac := factory.BuildTransportationAccountingCode(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				TacFnBlModCd: models.StringPointer("W"),
			},
		}, {
			Model:    loa,
			LinkOnly: true,
		},
	}, nil)

	// Create Orders
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.Order{
				TAC:       &tac.TAC,
				IssueDate: currentTime,
			},
		},
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

	// Create Move
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				Status:             models.MoveStatusAPPROVED,
			},
		},
	}, nil)
	// Create Pickup Address
	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85004",
			},
		},
	}, nil)

	serviceOrderNumber := "1234"
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	ntsrShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
				ServiceOrderNumber:   &serviceOrderNumber,
			},
		},
	}, nil)

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    &agentUserInfo.firstName,
				LastName:     &agentUserInfo.lastName,
				Email:        &agentUserInfo.email,
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	// Create Payment Request
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// Create Domestic linehaul service item
	dlCost := unit.Cents(80000)
	dlItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DL Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDLH,
		dlItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &dlCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Fuel surcharge service item
	fsCost := unit.Cents(10700)
	fsItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "FS Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeFSC,
		fsItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &fsCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic origin price service item
	doCost := unit.Cents(15000)
	doItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DO Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOP,
		doItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &doCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic destination price service item
	ddpCost := unit.Cents(15000)
	ddpItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DDP Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDDP,
		ddpItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &ddpCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic unpacking service item
	duCost := unit.Cents(45900)
	duItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DUPK Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "5.79",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDUPK,
		duItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &duCost,
					Status:     models.PaymentServiceItemStatusApproved,
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
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

func MakeHHGMoveInSITNoDestinationSITOutDate(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: userInfo.email,
				Active:    true,
			},
		},
	}, nil)
	customer := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: &userInfo.email,
				FirstName:     &userInfo.firstName,
				LastName:      &userInfo.lastName,
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	dependentsAuthorized := true
	sitDaysAllowance := 90
	entitlements := factory.BuildEntitlement(appCtx.DB(), []factory.Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
				StorageInTransit:     &sitDaysAllowance,
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model:    customer,
			LinkOnly: true,
		},
		{
			Model:    entitlements,
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
	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		},
	}, nil)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(1350)

	requestedPickupDate := now.AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)

	shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  &estimatedWeight,
				PrimeActualWeight:     &actualWeight,
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	agentUserInfo := newUserInfo("agent")
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{Model: models.MTOAgent{
			FirstName:    &agentUserInfo.firstName,
			LastName:     &agentUserInfo.lastName,
			Email:        &agentUserInfo.email,
			MTOAgentType: models.MTOAgentReleasing,
		},
		},
	}, nil)

	twoMonthsAgo := now.AddDate(0, -2, 0)
	oneMonthAgo := now.AddDate(0, -1, 0)
	factory.BuildOriginSITServiceItems(appCtx.DB(), move, shipment, &twoMonthsAgo, &oneMonthAgo)
	destSITItems := factory.BuildDestSITServiceItemsNoSITDepartureDate(appCtx.DB(), move, shipment, &oneMonthAgo)
	err := appCtx.DB().Update(&destSITItems)
	move.MTOServiceItems = destSITItems
	if err != nil {
		log.Panic(fmt.Errorf("failed to update sit service item: %w", err))
	}
	return move
}
