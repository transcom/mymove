package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// makeOrder creates a single Order and associated data.
//
// Deprecated: use factory.BuildOrder
func makeOrder(db *pop.Connection, assertions Assertions) (models.Order, error) {
	// Create new relational data if not provided
	sm := assertions.Order.ServiceMember
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.Order.ServiceMemberID) {
		var err error
		sm, err = makeExtendedServiceMember(db, assertions)
		if err != nil {
			return models.Order{}, nil
		}
	}

	dutyLocation := assertions.Order.NewDutyLocation
	// Note above
	if isZeroUUID(assertions.Order.NewDutyLocationID) {
		var err error
		dutyLocation, err = fetchOrMakeDefaultNewOrdersDutyLocation(db)
		if err != nil {
			return models.Order{}, err
		}
	}

	document := assertions.Order.UploadedOrders
	// Note above
	if isZeroUUID(assertions.Order.UploadedOrdersID) {
		fullDocumentAssertions := Assertions{
			Document: models.Document{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
			},
		}

		mergeModels(&fullDocumentAssertions, assertions)

		var err error
		document, err = makeDocument(db, fullDocumentAssertions)
		if err != nil {
			return models.Order{}, err
		}

		fullUserUploadAssertions := Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &document.ID,
				Document:   document,
				UploaderID: sm.UserID,
			},
			UserUploader: assertions.UserUploader,
		}

		mergeModels(&fullUserUploadAssertions, assertions)

		u, err := makeUserUpload(db, fullUserUploadAssertions)
		if err != nil {
			return models.Order{}, err
		}

		document.UserUploads = append(document.UserUploads, u)
	}

	defaultOrderNumber := "ORDER3"
	ordersNumber := assertions.Order.OrdersNumber
	if ordersNumber == nil {
		ordersNumber = &defaultOrderNumber
	}

	defaultTACNumber := "F8E1"
	TAC := assertions.Order.TAC
	if TAC == nil {
		TAC = &defaultTACNumber
	}

	defaultDepartmentIndicator := "AIR_AND_SPACE_FORCE"
	departmentIndicator := assertions.Order.DepartmentIndicator
	if departmentIndicator == nil {
		departmentIndicator = &defaultDepartmentIndicator
	}
	hasDependents := assertions.Order.HasDependents || false
	spouseHasProGear := assertions.Order.SpouseHasProGear || false
	grade := models.ServiceMemberGradeE1

	entitlement := assertions.Entitlement
	if isZeroUUID(entitlement.ID) {
		assertions.Order.Grade = &grade
		entitlement = makeEntitlement(db, assertions)
	}

	originDutyLocation := assertions.OriginDutyLocation
	if isZeroUUID(originDutyLocation.ID) {
		var err error
		originDutyLocation, err = makeDutyLocation(db, assertions)
		if err != nil {
			return models.Order{}, err
		}
	}

	gbloc, err := models.FetchGBLOCForPostalCode(db, originDutyLocation.Address.PostalCode)
	if gbloc.GBLOC == "" || err != nil {
		gbloc = makePostalCodeToGBLOC(db, originDutyLocation.Address.PostalCode, "KKFA")
	}

	orderTypeDetail := assertions.Order.OrdersTypeDetail
	hhgPermittedString := internalmessages.OrdersTypeDetail("HHG_PERMITTED")
	if orderTypeDetail == nil || *orderTypeDetail == "" {
		orderTypeDetail = &hhgPermittedString
	}

	// Added as a stopgap solution to populate these new fields
	// This testdatagen function is still being utilized through the MakeMTOServiceItemCustomerContact func
	defaultSupplyAndServicesCostEstimate := models.SupplyAndServicesCostEstimate
	defaultMethodOfPayment := models.MethodOfPayment
	defaultNAICS := models.NAICS
	contractor := fetchOrMakeContractor(db, Assertions{})
	defaultPackingAndShippingInstructions := models.InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + models.InstructionsAfterContractNumber

	order := models.Order{
		ServiceMember:                  sm,
		ServiceMemberID:                sm.ID,
		NewDutyLocation:                dutyLocation,
		NewDutyLocationID:              dutyLocation.ID,
		UploadedOrders:                 document,
		UploadedOrdersID:               document.ID,
		UploadedAmendedOrders:          assertions.Order.UploadedAmendedOrders,
		UploadedAmendedOrdersID:        assertions.Order.UploadedAmendedOrdersID,
		IssueDate:                      time.Date(TestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
		ReportByDate:                   time.Date(TestYear, time.August, 1, 0, 0, 0, 0, time.UTC),
		OrdersType:                     internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:                   ordersNumber,
		HasDependents:                  hasDependents,
		SpouseHasProGear:               spouseHasProGear,
		Status:                         models.OrderStatusDRAFT,
		TAC:                            TAC,
		DepartmentIndicator:            departmentIndicator,
		Grade:                          &grade,
		Entitlement:                    &entitlement,
		EntitlementID:                  &entitlement.ID,
		OriginDutyLocation:             &originDutyLocation,
		OriginDutyLocationID:           &originDutyLocation.ID,
		OrdersTypeDetail:               orderTypeDetail,
		OriginDutyLocationGBLOC:        &gbloc.GBLOC,
		SupplyAndServicesCostEstimate:  defaultSupplyAndServicesCostEstimate,
		MethodOfPayment:                defaultMethodOfPayment,
		NAICS:                          defaultNAICS,
		PackingAndShippingInstructions: defaultPackingAndShippingInstructions,
	}

	// Overwrite values with those from assertions
	mergeModels(&order, assertions.Order)

	mustCreate(db, &order, assertions.Stub)

	return order, nil
}
