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

// this is the class - note type struct and member vairbal
type ServiceMemberFactory struct {
	Model     *models.ServiceMember
	ForceUUID *uuid.UUID
}

// this is the constructor
func NewServiceMemberFactory(serviceMember models.ServiceMember, forceUUID *uuid.UUID) ServiceMemberFactory {
	return ServiceMemberFactory{&serviceMember, forceUUID}
}
func (sf ServiceMemberFactory) Create(db *pop.Connection, variants Variants) error {
	sm := sf.Model

	user := variants.User
	if isZeroUUID(variants.User.ID) {

		userFactory := NewUserFactory(models.User{}, nil)
		err := userFactory.Create(db, variants)
		if err != nil {
			return err
		}
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

// this is the class - note type struct and member variable
type UserFactory struct {
	Model     *models.User
	ForceUUID *uuid.UUID
}

// this is the constructor
func NewUserFactory(user models.User, forceUUID *uuid.UUID) UserFactory {
	return UserFactory{&user, forceUUID}
}

// this is a method
func (uf UserFactory) Create(db *pop.Connection, variants Variants) error {
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

func makeUserNew(db *pop.Connection, variants Variants) models.User {
	userFactory := NewUserFactory(models.User{}, nil)
	err := userFactory.Create(db, variants)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered creating: %v", err))
	}
	return *userFactory.Model
}
func makeServiceMemberNew(db *pop.Connection, variants Variants) models.ServiceMember {
	smFactory := NewServiceMemberFactory(models.ServiceMember{}, nil)
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
