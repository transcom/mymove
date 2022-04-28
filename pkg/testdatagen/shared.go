package testdatagen

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/transcom/mymove/pkg/random"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/imdario/mergo"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/uploader"

	"github.com/transcom/mymove/pkg/models"
)

// Assertions defines assertions about what the data contains
type Assertions struct {
	Address                                  models.Address
	AdminUser                                models.AdminUser
	TestDataAuditHistory                     TestDataAuditHistory
	BackupContact                            models.BackupContact
	Contractor                               models.Contractor
	DestinationAddress                       models.Address
	DistanceCalculation                      models.DistanceCalculation
	Document                                 models.Document
	DutyLocation                             models.DutyLocation
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
	MTOServiceItemDimensionCrate             models.MTOServiceItemDimension
	MTOServiceItemCustomerContact            models.MTOServiceItemCustomerContact
	MTOShipment                              models.MTOShipment
	Notification                             models.Notification
	WeightTicketSetDocument                  models.WeightTicketSetDocument
	OfficeUser                               models.OfficeUser
	Order                                    models.Order
	Organization                             models.Organization
	OriginDutyLocation                       models.DutyLocation
	PaymentRequest                           models.PaymentRequest
	PaymentRequestToInterchangeControlNumber models.PaymentRequestToInterchangeControlNumber
	PaymentServiceItem                       models.PaymentServiceItem
	PaymentServiceItemParam                  models.PaymentServiceItemParam
	PaymentServiceItemParams                 models.PaymentServiceItemParams
	PersonallyProcuredMove                   models.PersonallyProcuredMove
	PickupAddress                            models.Address
	PPMShipment                              models.PPMShipment
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
	Reweigh                                  models.Reweigh
	ReZip3                                   models.ReZip3
	Role                                     roles.Role
	SecondaryPickupAddress                   models.Address
	SecondaryDeliveryAddress                 models.Address
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
	Zip3Distance                             models.Zip3Distance
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

func Save(db *pop.Connection, model interface{}) error {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		return errors.Wrap(err, "Errors encountered saving model")
	}
	if verrs.HasAny() {
		return errors.Errorf("Validation errors encountered saving model: %v", verrs)
	}
	return nil
}

func MustSave(db *pop.Connection, model interface{}) {
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

// ConvertUUIDStringToUUID takes a uuid string and converts it to uuid.UUID; panic if uuid string isn't a valid UUID.
func ConvertUUIDStringToUUID(uuidString string) uuid.UUID {
	return uuid.Must(uuid.FromString(uuidString))
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

	// if this is called from inside another package remove so we're left with the parent dir
	cwd = strings.Split(cwd, "pkg")[0]

	fixturePath := path.Join(cwd, "pkg/testdatagen", fixtureDir, name)
	file, err := os.Open(filepath.Clean(fixturePath))
	if err != nil {
		log.Panic(fmt.Errorf("Error opening local file: %v", err))
	}

	return file
}

// FixtureRuntimeFile allows us to include a fixture like a PDF in the test
func FixtureRuntimeFile(name string) *runtime.File {
	fixtureDir := "testdatagen/testdata"
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	fixturePath := path.Join(cwd, "..", "..", fixtureDir, name)

	file, err := os.Open(filepath.Clean(fixturePath))
	if err != nil {
		log.Panic("Error opening fixture file", zap.Error(err))
	}

	info, err := file.Stat()
	if err != nil {
		log.Panic("Error accessing fixture stats", zap.Error(err))
	}

	header := multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	returnFile := &runtime.File{
		Header: &header,
		Data:   file,
	}
	return returnFile
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

// CurrentDateWithoutTime returns a pointer to a time.Time, stripped of any time info (so date only).
func CurrentDateWithoutTime() *time.Time {
	currentTime := time.Now()
	year, month, day := currentTime.Date()

	currentDate := time.Date(year, month, day, 0, 0, 0, 0, currentTime.Location())

	return &currentDate
}
