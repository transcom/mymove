package testdatagen

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/uploader"
)

func UserFactory(db *pop.Connection, variants Variants) models.User {
	loginGovUUID := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "first.last@login.gov.test",
		Active:        false,
	}

	// Overwrite values with those from assertions
	mergeModels(&user, variants.User)

	mustCreate(db, &user, variants.Stub)

	return user
}

func ServiceMemberFactory(db *pop.Connection, variants Variants) models.ServiceMember {
	// aServiceMember := variants.ServiceMember
	// currentAddressID := aServiceMember.ResidentialAddressID
	// currentAddress := aServiceMember.ResidentialAddress
	// Check that no nested values in sm

	// ID is required because it must be populated for Eager saving to work.
	user := variants.User
	if isZeroUUID(variants.User.ID) {
		user = UserFactory(db, variants)
	}

	// currentAddress := variants.Addresses.ServiceMember__CurrentAddress
	// if isZeroUUID(currentAddress.ID) {}
	// 	currentAddress := AddressFactory(db, variants)
	// }

	army := models.AffiliationARMY
	randomEdipi := RandomEdipi()
	rank := models.ServiceMemberRankE1
	email := "leo_spaceman_sm@example.com"

	serviceMember := models.ServiceMember{
		UserID:        user.ID,
		User:          user,
		Edipi:         swag.String(randomEdipi),
		Affiliation:   &army,
		FirstName:     swag.String("Leo"),
		LastName:      swag.String("Spacemen"),
		Telephone:     swag.String("212-123-4567"),
		PersonalEmail: &email,
		Rank:          &rank,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceMember, variants.ServiceMember)

	mustCreate(db, &serviceMember, variants.Stub)

	return serviceMember
}

// this is the class - note type struct and member vairbal
type ServiceMemberFac struct {
	Model     *models.ServiceMember
	ForceUUID *uuid.UUID
}

// this is the constructor
func NewServiceMemberFac(serviceMember models.ServiceMember, forceUUID *uuid.UUID) ServiceMemberFac {
	return ServiceMemberFac{&serviceMember, forceUUID}
}
func (sf ServiceMemberFac) Create(db *pop.Connection, variants Variants) error {
	sm := sf.Model

	user := variants.User
	if isZeroUUID(variants.User.ID) {

		userFactory := NewUserFac(models.User{}, nil)
		userFactory.Create(db, variants)
		user = *userFactory.Model
	}

	army := models.AffiliationARMY
	randomEdipi := RandomEdipi()
	rank := models.ServiceMemberRankE1
	email := "leo_spaceman_sm@example.com"

	sm.UserID = user.ID
	sm.User = user
	sm.Edipi = swag.String(randomEdipi)
	sm.Affiliation = &army
	sm.FirstName = swag.String("Leo")
	sm.LastName = swag.String("Spacemen")
	sm.Telephone = swag.String("212-123-4567")
	sm.PersonalEmail = &email
	sm.Rank = &rank

	// Overwrite values with those from assertions
	mergeModels(sm, variants.ServiceMember)

	mustCreate(db, sm, variants.Stub)

	return nil
}

// this is the class - note type struct and member vairbal
type UserFac struct {
	Model     *models.User
	ForceUUID *uuid.UUID
}

// this is the constructor
func NewUserFac(user models.User, forceUUID *uuid.UUID) UserFac {
	return UserFac{&user, forceUUID}
}

// this is a method
func (uf UserFac) Create(db *pop.Connection, variants Variants) error {
	user := uf.Model

	loginGovUUID := uuid.Must(uuid.NewV4())
	user.LoginGovUUID = &loginGovUUID
	user.LoginGovEmail = "first.last@login.gov.test"
	user.Active = false

	// Overwrite values with those from assertions
	mergeModels(user, variants.User)

	mustCreate(db, user, variants.Stub)

	return nil
}

func CheckNestedStruct(edges interface{}) {
	rv := reflect.ValueOf(edges)
	// uidField := rv.FieldByName("id")
	// fmt.Println(uidField.Type())
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldKind := field.Kind()
		// fieldType := field.Type()
		if fieldKind == reflect.Struct {

			CheckNestedStruct(field)
		}
	}
}

func makeUserX(db *pop.Connection, variants Variants) models.User {
	userFactory := NewUserFac(models.User{}, nil)
	err := userFactory.Create(db, variants)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered creating: %v", err))
	}
	return *userFactory.Model
}
func makeSMX(db *pop.Connection, variants Variants) models.ServiceMember {
	smFactory := NewServiceMemberFac(models.ServiceMember{}, nil)
	err := smFactory.Create(db, variants)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered creating: %v", err))
	}
	return *smFactory.Model
}

type Assertor interface {
	Merge(model interface{}) interface{}
}
type Variants struct {
	Address                                  models.Address
	AdminUser                                models.AdminUser
	TestDataAuditHistory                     TestDataAuditHistory
	BackupContact                            models.BackupContact
	ClientCert                               models.ClientCert
	Contractor                               models.Contractor
	DistanceCalculation                      models.DistanceCalculation
	Document                                 models.Document
	DutyLocation                             models.DutyLocation
	ElectronicOrder                          models.ElectronicOrder
	ElectronicOrdersRevision                 models.ElectronicOrdersRevision
	Entitlement                              models.Entitlement
	EvaluationReport                         models.EvaluationReport
	FuelEIADieselPrice                       models.FuelEIADieselPrice
	File                                     afero.File
	GHCDieselFuelPrice                       models.GHCDieselFuelPrice
	Invoice                                  models.Invoice
	Move                                     models.Move
	MoveDocument                             models.MoveDocument
	MovingExpenseDocument                    models.MovingExpenseDocument
	MTOAgent                                 models.MTOAgent
	MTOServiceItem                           models.MTOServiceItem
	MTOServiceItemDimension                  models.MTOServiceItemDimension
	MTOServiceItemCustomerContact            models.MTOServiceItemCustomerContact
	MTOShipment                              models.MTOShipment
	Notification                             models.Notification
	WeightTicketSetDocument                  models.WeightTicketSetDocument
	OfficeUser                               models.OfficeUser
	CustomerSupportRemark                    models.CustomerSupportRemark
	Order                                    models.Order
	Organization                             models.Organization
	OriginDutyLocation                       models.DutyLocation
	PaymentRequest                           models.PaymentRequest
	PaymentRequestToInterchangeControlNumber models.PaymentRequestToInterchangeControlNumber
	MTOServiceItemDimensionCrate             models.MTOServiceItemDimension
	PaymentServiceItem                       models.PaymentServiceItem
	DestinationAddress                       models.Address
	PaymentServiceItemParam                  models.PaymentServiceItemParam
	PaymentServiceItemParams                 models.PaymentServiceItemParams
	PersonallyProcuredMove                   models.PersonallyProcuredMove
	PPMShipment                              models.PPMShipment
	PrimeUpload                              models.PrimeUpload
	PrimeUploader                            *uploader.PrimeUploader
	ProofOfServiceDoc                        models.ProofOfServiceDoc
	ReContract                               models.ReContract
	ReContractYear                           models.ReContractYear
	ReDomesticLinehaulPrice                  models.ReDomesticLinehaulPrice
	ReDomesticOtherPrice                     models.ReDomesticOtherPrice
	ReDomesticServiceArea                    models.ReDomesticServiceArea
	ReDomesticServiceAreaPrice               models.ReDomesticServiceAreaPrice
	Reimbursement                            models.Reimbursement
	ReRateArea                               models.ReRateArea
	ReService                                models.ReService
	Reweigh                                  models.Reweigh
	ReZip3                                   models.ReZip3
	Role                                     roles.Role
	ServiceItemParamKey                      models.ServiceItemParamKey
	ServiceParam                             models.ServiceParam
	SignedCertification                      models.SignedCertification
	SITExtension                             models.SITExtension
	ServiceMember                            models.ServiceMember
	StorageFacility                          models.StorageFacility
	Stub                                     bool
	Tariff400ngServiceArea                   models.Tariff400ngServiceArea
	Tariff400ngItem                          models.Tariff400ngItem
	Tariff400ngItemRate                      models.Tariff400ngItemRate
	Tariff400ngZip3                          models.Tariff400ngZip3
	TrafficDistributionList                  models.TrafficDistributionList
	TransportationAccountingCode             models.TransportationAccountingCode
	TransportationOffice                     models.TransportationOffice
	TransportationServiceProvider            models.TransportationServiceProvider
	TransportationServiceProviderPerformance models.TransportationServiceProviderPerformance
	Upload                                   models.Upload
	UploadUseZeroBytes                       bool
	Uploader                                 *uploader.Uploader
	UserUpload                               models.UserUpload
	UserUploader                             *uploader.UserUploader
	User                                     models.User
	UsersRoles                               models.UsersRoles
	WebhookNotification                      models.WebhookNotification
	WebhookSubscription                      models.WebhookSubscription
	WeightTicket                             models.WeightTicket
	Zip3Distance                             models.Zip3Distance
	ServiceMemberCurrentAddress              models.Address
}
