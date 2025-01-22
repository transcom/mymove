package factory

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type orderBuildType byte

const (
	orderBuildBasic orderBuildType = iota
	orderBuildWithoutDefaults
)

func buildOrderWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType orderBuildType) models.Order {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cOrder models.Order
	if result := findValidCustomization(customs, Order); result != nil {
		cOrder = result.Model.(models.Order)

		if result.LinkOnly {
			return cOrder
		}
	}

	// Now things get complicated
	// Find/create for the service member will create a default duty
	// location, but creating orders uses a different default duty
	// location
	//

	// Once we have created the origin duty location, we want for that
	// to be used for all subsequent factories (including
	// BuildServiceMember), so we need to remove any existing origin
	// duty location customizations and then add a LinkOnly
	// customization to the newly created duty location

	originDutyLocationCustoms := customs
	originDutyLocationChecks := []struct {
		convertFrom CustomType
		convertTo   CustomType
	}{
		{DutyLocations.OriginDutyLocation, DutyLocation},
		{TransportationOffices.OriginDutyLocation, TransportationOffice},
	}

	for _, customType := range originDutyLocationChecks {
		if result := findValidCustomization(customs, customType.convertFrom); result != nil {
			originDutyLocationCustoms =
				convertCustomizationInList(originDutyLocationCustoms,
					customType.convertFrom, customType.convertTo)
		}
	}

	originDutyLocation := BuildDutyLocation(db, originDutyLocationCustoms, nil)
	if db != nil {
		// can only do LinkOnly if we have an ID, which we won't have
		// for a stubbed duty location
		customs = replaceCustomization(customs, Customization{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &DutyLocations.OriginDutyLocation,
		})
	}

	newDutyLocationCustoms := customs
	hasNewDutyLocationCustoms := false
	newDutyLocationChecks := []struct {
		convertFrom CustomType
		convertTo   CustomType
	}{
		{DutyLocations.NewDutyLocation, DutyLocation},
		{TransportationOffices.NewDutyLocation, TransportationOffice},
	}

	for _, customType := range newDutyLocationChecks {
		if result := findValidCustomization(customs, customType.convertFrom); result != nil {
			hasNewDutyLocationCustoms = true
			newDutyLocationCustoms =
				convertCustomizationInList(newDutyLocationCustoms,
					customType.convertFrom, customType.convertTo)
		}
	}

	var newDutyLocation models.DutyLocation
	if hasNewDutyLocationCustoms {
		newDutyLocation = BuildDutyLocation(db, newDutyLocationCustoms, nil)
	} else {
		// the dev did not provide any customizations for the new duty
		// location, so use the default orders duty location trait
		newDutyLocation = FetchOrBuildOrdersDutyLocation(db)
	}

	// Find/create the user upload (and document and service member)
	//
	// This is where things get a bit hairy. BuildUserUpload needs a
	// document and calls BuildDocument (using Document customization)
	// to create one if necessary. But the customization for Orders
	// needs to distinguish between the UploadedOrders document and
	// the UploadedAmendedOrders document.
	//
	// Even hairier, the BuildDocument builds a regular service
	// member, but we need an extended service member for BuildOrder
	//
	// So ...

	// Call BuildExtendedServiceMember to create the service member

	// convert the OriginDutyLocation customs to vanilla DutyLocation
	// for BuildExtendedServiceMember
	serviceMemberCustoms := customs
	if result := findValidCustomization(customs, DutyLocations.OriginDutyLocation); result != nil {
		serviceMemberCustoms = convertCustomizationInList(customs, DutyLocations.OriginDutyLocation, DutyLocation)
	}
	serviceMember := BuildExtendedServiceMember(db, serviceMemberCustoms, traits)

	if db != nil {
		// Now we need  a LinkOnly customization for the created
		// ServiceMember
		// can only do LinkOnly if we have an ID, which we won't have
		// for a stubbed service member
		customs = replaceCustomization(customs, Customization{
			Model:    serviceMember,
			LinkOnly: true,
		})
	}

	// Find the customizations for UploadedOrders and build the
	// uploadedOrders
	var uploadedOrders models.Document
	uploadedOrdersCustoms := customs
	needsUploadedOrdersUserUpload := true
	if result := findValidCustomization(customs, Documents.UploadedOrders); result != nil {
		// the dev provided UploadedOrders customizations
		// If this is a LinkOnly UploadedOrders, we do not need to
		// build any user uploads
		if result.LinkOnly {
			needsUploadedOrdersUserUpload = false
		}
		//, so use them
		uploadedOrdersCustoms = convertCustomizationInList(customs,
			Documents.UploadedOrders, Document)
	}
	uploadedOrders = BuildDocument(db, uploadedOrdersCustoms, nil)

	// Now we have the document properly customized, but if we call
	// BuildUserUpload with the provided customizations, it won't know
	// we have already created the document. So now we prepend a
	// LinkOnly Document customization
	if db != nil {
		// can only do LinkOnly if we have an ID, which we won't have
		// for a stubbed document
		customs = replaceCustomization(customs, Customization{
			Model:    uploadedOrders,
			LinkOnly: true,
			Type:     &Documents.UploadedOrders,
		})
	}

	if needsUploadedOrdersUserUpload {
		// Now call BuildUserUpload with our re-jiggered customs (only
		// available if db != nil)
		uploadCustoms := customs
		if result := findValidCustomization(customs, Documents.UploadedOrders); result != nil {

			uploadCustoms = convertCustomizationInList(customs, Documents.UploadedOrders, Document)
		}
		userUpload := BuildUserUpload(db, uploadCustoms, traits)
		// make sure we append the upload to the uploadedOrders document
		uploadedOrders.UserUploads = append(uploadedOrders.UserUploads, userUpload)
	}

	entitlement := BuildEntitlement(db, customs, traits)

	// by default, the amended orders document is not provided. It is
	// only created if the dev provides the
	// Documents.UploadedAmendedOrders customization
	uploadedAmendedOrdersCustoms := customs
	var amendedOrdersDocument *models.Document
	hasAmendedOrdersCustoms := findValidCustomization(customs, Documents.UploadedAmendedOrders)
	// only basic builds include amended orders
	if buildType == orderBuildBasic && hasAmendedOrdersCustoms != nil {
		uploadedAmendedOrdersCustoms =
			convertCustomizationInList(uploadedAmendedOrdersCustoms,
				Documents.UploadedAmendedOrders, Document)
		doc := BuildDocument(db, uploadedAmendedOrdersCustoms, nil)
		amendedOrdersDocument = &doc
	}

	defaultOrdersNumber := "ORDER3"
	defaultTACNumber := "F8E1"
	defaultDepartmentIndicator := "AIR_AND_SPACE_FORCE"
	defaultGrade := models.ServiceMemberGradeE1
	defaultHasDependents := false
	defaultSpouseHasProGear := false
	defaultOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	defaultOrdersTypeDetail := internalmessages.OrdersTypeDetail("HHG_PERMITTED")
	defaultDestinationGbloc := "AGFM"
	destinationGbloc := &defaultDestinationGbloc
	testYear := 2018
	defaultIssueDate := time.Date(testYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	defaultReportByDate := time.Date(testYear, time.August, 1, 0, 0, 0, 0, time.UTC)
	defaultStatus := models.OrderStatusDRAFT
	defaultOriginDutyLocationGbloc := "KKFA"
	originDutyLocationGbloc := &defaultOriginDutyLocationGbloc
	defaultSupplyAndServicesCostEstimate := models.SupplyAndServicesCostEstimate
	defaultMethodOfPayment := models.MethodOfPayment
	defaultNAICS := models.NAICS
	contractor := FetchOrBuildDefaultContractor(db, customs, traits)
	defaultPackingAndShippingInstructions := models.InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + models.InstructionsAfterContractNumber

	var ordersNumber *string
	var tac *string
	var departmentsIndicator *string
	var ordersTypeDetail *internalmessages.OrdersTypeDetail

	// the basic build time adds additional defaults
	if buildType == orderBuildBasic {
		ordersNumber = &defaultOrdersNumber
		tac = &defaultTACNumber
		departmentsIndicator = &defaultDepartmentIndicator
		ordersTypeDetail = &defaultOrdersTypeDetail
	}

	if db != nil {
		// make sure the origin duty location address is loaded as it
		// may not be
		if originDutyLocation.Address.PostalCode == "" {
			err := db.EagerPreload("Address", "Address.Country").Find(&originDutyLocation, originDutyLocation.ID)
			if err != nil {
				log.Panicf("Error loading duty location by id %s: %s\n", originDutyLocation.ID.String(), err)
			}
		}
		originPostalCodeToGBLOC := FetchOrBuildPostalCodeToGBLOC(db, originDutyLocation.Address.PostalCode, "KKFA")
		originDutyLocationGbloc = &originPostalCodeToGBLOC.GBLOC

		if newDutyLocation.Address.PostalCode == "" {
			err := db.EagerPreload("Address", "Address.Country").Find(&newDutyLocation, newDutyLocation.ID)
			if err != nil {
				log.Panicf("Error loading duty location by id %s: %s\n", newDutyLocation.ID.String(), err)
			}
			destinationPostalCodeToGBLOC := FetchOrBuildPostalCodeToGBLOC(db, newDutyLocation.Address.PostalCode, "AGFM")
			destinationGbloc = &destinationPostalCodeToGBLOC.GBLOC
		}
	}

	order := models.Order{
		ServiceMember:                  serviceMember,
		ServiceMemberID:                serviceMember.ID,
		NewDutyLocation:                newDutyLocation,
		NewDutyLocationID:              newDutyLocation.ID,
		DestinationGBLOC:               destinationGbloc,
		UploadedOrders:                 uploadedOrders,
		UploadedOrdersID:               uploadedOrders.ID,
		IssueDate:                      defaultIssueDate,
		ReportByDate:                   defaultReportByDate,
		OrdersType:                     defaultOrdersType,
		OrdersNumber:                   ordersNumber,
		HasDependents:                  defaultHasDependents,
		SpouseHasProGear:               defaultSpouseHasProGear,
		Status:                         defaultStatus,
		TAC:                            tac,
		DepartmentIndicator:            departmentsIndicator,
		Grade:                          &defaultGrade,
		Entitlement:                    &entitlement,
		EntitlementID:                  &entitlement.ID,
		OriginDutyLocation:             &originDutyLocation,
		OriginDutyLocationID:           &originDutyLocation.ID,
		OrdersTypeDetail:               ordersTypeDetail,
		OriginDutyLocationGBLOC:        originDutyLocationGbloc,
		SupplyAndServicesCostEstimate:  defaultSupplyAndServicesCostEstimate,
		MethodOfPayment:                defaultMethodOfPayment,
		NAICS:                          defaultNAICS,
		PackingAndShippingInstructions: defaultPackingAndShippingInstructions,
	}

	if amendedOrdersDocument != nil {
		order.UploadedAmendedOrders = amendedOrdersDocument
		order.UploadedAmendedOrdersID = &amendedOrdersDocument.ID
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&order, cOrder)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &order)
	}

	return order
}

// BuildOrder creates an Order.
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildOrder(db *pop.Connection, customs []Customization, traits []Trait) models.Order {
	return buildOrderWithBuildType(db, customs, traits, orderBuildBasic)
}

// BuildOrderWithoutDefaults creates an Order that only includes fields that the
// server member would have supplied prior to uploading their documents
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildOrderWithoutDefaults(db *pop.Connection, customs []Customization, traits []Trait) models.Order {
	return buildOrderWithBuildType(db, customs, traits, orderBuildWithoutDefaults)
}

// ------------------------
//      TRAITS
// ------------------------

// GetTraitHasDependents returns a customization to enable dependents on an order
func GetTraitHasDependents() []Customization {
	return []Customization{
		{
			Model: models.Order{
				HasDependents: true,
			},
		},
	}
}
