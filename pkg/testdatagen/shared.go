package testdatagen

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"

	"github.com/transcom/mymove/pkg/random"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/imdario/mergo"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/uploader"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// Assertions defines assertions about what the data contains
type Assertions struct {
	AccessCode                               models.AccessCode
	Address                                  models.Address
	AdminUser                                models.AdminUser
	BackupContact                            models.BackupContact
	Contractor                               models.Contractor
	DestinationAddress                       models.Address
	DestinationDutyStation                   models.DutyStation
	DistanceCalculation                      models.DistanceCalculation
	Document                                 models.Document
	DutyStation                              models.DutyStation
	ElectronicOrder                          models.ElectronicOrder
	ElectronicOrdersRevision                 models.ElectronicOrdersRevision
	Entitlement                              models.Entitlement
	FuelEIADieselPrice                       models.FuelEIADieselPrice
	File                                     afero.File
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
	Order                                    models.Order
	Organization                             models.Organization
	OriginDutyStation                        models.DutyStation
	PaymentRequest                           models.PaymentRequest
	PaymentRequestToInterchangeControlNumber models.PaymentRequestToInterchangeControlNumber
	PaymentServiceItem                       models.PaymentServiceItem
	PaymentServiceItemParam                  models.PaymentServiceItemParam
	PaymentServiceItemParams                 models.PaymentServiceItemParams
	PersonallyProcuredMove                   models.PersonallyProcuredMove
	PickupAddress                            models.Address
	PrimeUpload                              models.PrimeUpload
	PrimeUploader                            *uploader.PrimeUploader
	ProofOfServiceDoc                        models.ProofOfServiceDoc
	ReContract                               models.ReContract
	ReContractYear                           models.ReContractYear
	ReDomesticServiceArea                    models.ReDomesticServiceArea
	ReDomesticLinehaulPrice                  models.ReDomesticLinehaulPrice
	Reimbursement                            models.Reimbursement
	ReRateArea                               models.ReRateArea
	ReService                                models.ReService
	ReZip3                                   models.ReZip3
	SecondaryPickupAddress                   models.Address
	SecondaryDeliveryAddress                 models.Address
	ServiceItemParamKey                      models.ServiceItemParamKey
	ServiceParam                             models.ServiceParam
	SignedCertification                      models.SignedCertification
	ServiceMember                            models.ServiceMember
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
	WebhookNotification                      models.WebhookNotification
	WebhookSubscription                      models.WebhookSubscription
	Zip3Distance                             models.Zip3Distance
}

func stringPointer(s string) *string {
	return &s
}

func poundPointer(p unit.Pound) *unit.Pound {
	return &p
}

func mustCreate(db *pop.Connection, model interface{}, stub bool) {
	if stub {
		return
	}

	verrs, err := db.ValidateAndCreate(model)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("Validation errors encountered saving %#v: %v", model, verrs))
	}
}

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("Validation errors encountered saving %#v: %v", model, verrs))
	}
}

func noErr(err error) {
	if err != nil {
		log.Panic(fmt.Errorf("Error encountered: %v", err))
	}
}

// isZeroUUID determines whether a UUID is its zero value
func isZeroUUID(testID uuid.UUID) bool {
	return testID == uuid.Nil
}

// mergeModels merges src into dst, if non-zero values are present
// dst should be a pointer the struct you are merging into
func mergeModels(dst, src interface{}) {
	noErr(
		mergo.Merge(dst, src, mergo.WithOverride, mergo.WithTransformers(customTransformer{})),
	)
}

// MergeModels exposes the private function mergeModels
func MergeModels(dst, src interface{}) {
	mergeModels(dst, src)
}

// Source chars for random string
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// Returns a random alphanumeric string of specified length
func makeRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		randInt, err := random.GetRandomInt(len(letterBytes))
		if err != nil {
			log.Panicf("failed to create random string %v", err)
			return ""
		}
		b[i] = letterBytes[randInt]

	}
	return string(b)
}

// Fixture opens a file from the testdata dir
func Fixture(name string) afero.File {
	fixtureDir := "testdata"
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(fmt.Errorf("failed to get current directory: %s", err))
	}

	fixturePath := path.Join(cwd, "pkg/testdatagen", fixtureDir, name)
	// #nosec This will only be using test data
	file, err := os.Open(filepath.Clean(fixturePath))
	if err != nil {
		log.Panic(fmt.Errorf("Error opening local file: %v", err))
	}

	return file
}

// customTransformer handles testing for zero values in structs that mergo can't normally deal with
type customTransformer struct {
}

// Transformer checks if src is not a zero value, then overwrites dst
func (t customTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	// UUID comparison
	if typ == reflect.TypeOf(uuid.UUID{}) || typ == reflect.TypeOf(&uuid.UUID{}) {
		return func(dst, src reflect.Value) error {
			// We need to cast the actual value to validate
			var srcIsValid bool
			if src.Kind() == reflect.Ptr {
				srcID := src.Interface().(*uuid.UUID)
				srcIsValid = !src.IsNil() && !isZeroUUID(*srcID)
			} else {
				srcID := src.Interface().(uuid.UUID)
				srcIsValid = !isZeroUUID(srcID)
			}
			if dst.CanSet() && srcIsValid {
				dst.Set(src)
			}
			return nil
		}
	}
	// time.Time comparison
	if typ == reflect.TypeOf(time.Time{}) || typ == reflect.TypeOf(&time.Time{}) {
		return func(dst, src reflect.Value) error {
			srcIsValid := false
			// Either it's a non-nil pointer or a non-pointer
			if src.Kind() != reflect.Ptr || !src.IsNil() {
				isZeroMethod := src.MethodByName("IsZero")
				srcIsValid = !isZeroMethod.Call([]reflect.Value{})[0].Bool()
			}
			if dst.CanSet() && srcIsValid {
				dst.Set(src)
			}
			return nil
		}
	}
	return nil
}
