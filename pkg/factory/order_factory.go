package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildOrder creates a Order.
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildOrder(db *pop.Connection, customs []Customization, traits []Trait) models.Order {
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
	if result := findValidCustomization(customs, DutyLocations.OriginDutyLocation); result != nil {
		originDutyLocationCustoms =
			convertCustomizationInList(originDutyLocationCustoms,
				DutyLocations.OriginDutyLocation, DutyLocation)
	}
	var originDutyLocation models.DutyLocation
	if db != nil {
		originDutyLocation = BuildDutyLocation(db, originDutyLocationCustoms, nil)
		customs = replaceCustomization(customs, Customization{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &DutyLocations.OriginDutyLocation,
		})
	}

	var newDutyLocation models.DutyLocation
	if result := findValidCustomization(customs, DutyLocations.NewDutyLocation); result != nil {
		// the dev provided customizations for the new duty location,
		// so use them
		newDutyLocationCustoms :=
			convertCustomizationInList(customs,
				DutyLocations.NewDutyLocation, DutyLocation)
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
		customs = replaceCustomization(customs, Customization{
			Model:    serviceMember,
			LinkOnly: true,
		})
	}

	// Find the customizations for UploadedOrders and build the
	// uploadedOrders
	var uploadedOrders models.Document
	uploadedOrdersCustoms := customs
	if result := findValidCustomization(customs, Documents.UploadedOrders); result != nil {
		// the dev provided UploadedOrders customizations, so use them
		uploadedOrdersCustoms = convertCustomizationInList(customs,
			Documents.UploadedOrders, Document)
	}
	uploadedOrders = BuildDocument(db, uploadedOrdersCustoms, nil)

	// Now we have the document properly customized, but if we call
	// BuildUserUpload with the provided customizations, it won't know
	// we have already created the document. So now we prepend a
	// LinkOnly Document customization
	if db != nil {
		customs = replaceCustomization(customs, Customization{
			Model:    uploadedOrders,
			LinkOnly: true,
			Type:     &Documents.UploadedOrders,
		})
	}

	// Now call BuildUserUpload with our re-jiggered customs (only
	// available if db != nil)
	uploadCustoms := customs
	if result := findValidCustomization(customs, Documents.UploadedOrders); result != nil {

		uploadCustoms = convertCustomizationInList(customs, Documents.UploadedOrders, Document)
	}
	userUpload := BuildUserUpload(db, uploadCustoms, traits)
	// make sure we append the upload to the uploadedOrders document
	uploadedOrders.UserUploads = append(uploadedOrders.UserUploads, userUpload)

	entitlement := BuildEntitlement(db, customs, traits)

	// by default, the amended orders document is not provided. It is
	// only created if the dev provides the
	// Documents.UploadedAmendedOrders customization
	uploadedAmendedOrdersCustoms := customs
	var amendedOrdersDocument *models.Document
	if result := findValidCustomization(customs, Documents.UploadedAmendedOrders); result != nil {
		uploadedAmendedOrdersCustoms =
			convertCustomizationInList(uploadedAmendedOrdersCustoms,
				Documents.UploadedAmendedOrders, Document)
		doc := BuildDocument(db, uploadedAmendedOrdersCustoms, nil)
		amendedOrdersDocument = &doc
	}

	defaultOrdersNumber := "ORDER3"
	defaultTACNumber := "F8E1"
	defaultDepartmentIndicator := "AIR_FORCE"
	defaultGrade := "E_1"
	defaultHasDependents := false
	defaultSpouseHasProGear := false
	defaultOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	defaultOrdersTypeDetail := internalmessages.OrdersTypeDetail("HHG_PERMITTED")
	testYear := 2018
	defaultIssueDate := time.Date(testYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	defaultReportByDate := time.Date(testYear, time.August, 1, 0, 0, 0, 0, time.UTC)
	defaultStatus := models.OrderStatusDRAFT

	order := models.Order{
		ServiceMember:        serviceMember,
		ServiceMemberID:      serviceMember.ID,
		NewDutyLocation:      newDutyLocation,
		NewDutyLocationID:    newDutyLocation.ID,
		UploadedOrders:       uploadedOrders,
		UploadedOrdersID:     uploadedOrders.ID,
		IssueDate:            defaultIssueDate,
		ReportByDate:         defaultReportByDate,
		OrdersType:           defaultOrdersType,
		OrdersNumber:         &defaultOrdersNumber,
		HasDependents:        defaultHasDependents,
		SpouseHasProGear:     defaultSpouseHasProGear,
		Status:               defaultStatus,
		TAC:                  &defaultTACNumber,
		DepartmentIndicator:  &defaultDepartmentIndicator,
		Grade:                &defaultGrade,
		Entitlement:          &entitlement,
		EntitlementID:        &entitlement.ID,
		OriginDutyLocation:   &originDutyLocation,
		OriginDutyLocationID: &originDutyLocation.ID,
		OrdersTypeDetail:     &defaultOrdersTypeDetail,
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
