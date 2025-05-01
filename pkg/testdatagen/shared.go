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
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/random"
	"github.com/transcom/mymove/pkg/uploader"
)

// Assertions defines assertions about what the data contains
type Assertions struct {
	Address                                  models.Address
	AdminUser                                models.AdminUser
	BackupContact                            models.BackupContact
	ClientCert                               models.ClientCert
	Contractor                               models.Contractor
	CustomerSupportRemark                    models.CustomerSupportRemark
	DestinationAddress                       models.Address
	DistanceCalculation                      models.DistanceCalculation
	Document                                 models.Document
	DutyLocation                             models.DutyLocation
	EdiError                                 models.EdiError
	ElectronicOrder                          models.ElectronicOrder
	ElectronicOrdersRevision                 models.ElectronicOrdersRevision
	Entitlement                              models.Entitlement
	EvaluationReport                         models.EvaluationReport
	File                                     afero.File
	FuelEIADieselPrice                       models.FuelEIADieselPrice
	GHCDieselFuelPrice                       models.GHCDieselFuelPrice
	Invoice                                  models.Invoice
	LineOfAccounting                         models.LineOfAccounting
	Move                                     models.Move
	MovingExpense                            models.MovingExpense
	MTOAgent                                 models.MTOAgent
	MTOServiceItem                           models.MTOServiceItem
	MTOServiceItemCustomerContact            models.MTOServiceItemCustomerContact
	MTOServiceItemDimension                  models.MTOServiceItemDimension
	MTOServiceItemDimensionCrate             models.MTOServiceItemDimension
	MTOShipment                              models.MTOShipment
	Notification                             models.Notification
	OfficeUser                               models.OfficeUser
	Order                                    models.Order
	Organization                             models.Organization
	OriginDutyLocation                       models.DutyLocation
	PaymentRequest                           models.PaymentRequest
	PaymentRequestToInterchangeControlNumber models.PaymentRequestToInterchangeControlNumber
	PaymentServiceItem                       models.PaymentServiceItem
	PaymentServiceItemParam                  models.PaymentServiceItemParam
	PaymentServiceItemParams                 models.PaymentServiceItemParams
	PickupAddress                            models.Address
	PPMShipment                              models.PPMShipment
	PrimeUpload                              models.PrimeUpload
	PrimeUploader                            *uploader.PrimeUploader
	ProgearWeightTicket                      models.ProgearWeightTicket
	ProofOfServiceDoc                        models.ProofOfServiceDoc
	ReContract                               models.ReContract
	ReContractYear                           models.ReContractYear
	ReDomesticAccessorialPrice               models.ReDomesticAccessorialPrice
	ReDomesticLinehaulPrice                  models.ReDomesticLinehaulPrice
	ReDomesticOtherPrice                     models.ReDomesticOtherPrice
	ReDomesticServiceArea                    models.ReDomesticServiceArea
	ReDomesticServiceAreaPrice               models.ReDomesticServiceAreaPrice
	ReTaskOrderFee                           models.ReTaskOrderFee
	Reimbursement                            models.Reimbursement
	Report                                   models.EvaluationReport
	ReportViolation                          models.ReportViolation
	ReRateArea                               models.ReRateArea
	ReService                                models.ReService
	Reweigh                                  models.Reweigh
	ReZip3                                   models.ReZip3
	ReZip5RateArea                           models.ReZip5RateArea
	Role                                     roles.Role
	SecondaryDeliveryAddress                 models.Address
	SecondaryPickupAddress                   models.Address
	TertiaryDeliveryAddress                  models.Address
	TertiaryPickupAddress                    models.Address
	ServiceItemParamKey                      models.ServiceItemParamKey
	ServiceMember                            models.ServiceMember
	ServiceParam                             models.ServiceParam
	SignedCertification                      models.SignedCertification
	SITDurationUpdate                        models.SITDurationUpdate
	StorageFacility                          models.StorageFacility
	Stub                                     bool
	TransportationAccountingCode             models.TransportationAccountingCode
	TransportationOffice                     models.TransportationOffice
	Upload                                   models.Upload
	Uploader                                 *uploader.Uploader
	UploadUseZeroBytes                       bool
	User                                     models.User
	UsersRoles                               models.UsersRoles
	UserUpload                               models.UserUpload
	UserUploader                             *uploader.UserUploader
	Violation                                models.PWSViolation
	WebhookNotification                      models.WebhookNotification
	WebhookSubscription                      models.WebhookSubscription
	WeightTicket                             models.WeightTicket
	Zip3Distance                             models.Zip3Distance
}

func mustCreate(db *pop.Connection, model interface{}, stub bool) {
	if stub {
		return
	}

	verrs, err := db.ValidateAndCreate(model)
	if err != nil {
		log.Panic(fmt.Errorf("errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("validation errors encountered saving %#v: %v", model, verrs))
	}
}

func Save(db *pop.Connection, model interface{}) error {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		return errors.Wrap(err, "errors encountered saving model")
	}
	if verrs.HasAny() {
		return errors.Errorf("validation errors encountered saving model: %v", verrs)
	}
	return nil
}

func MustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		log.Panic(fmt.Errorf("errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("validation errors encountered saving %#v: %v", model, verrs))
	}
}

func noErr(err error) {
	if err != nil {
		log.Panic(fmt.Errorf("error encountered: %v", err))
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
func MakeRandomString(n int) string {
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

const numberBytes = "0123456789"

// MakeRandomNumberString makes a random numeric string of specified length
func MakeRandomNumberString(n int) string {
	b := make([]byte, n)
	for i := range b {
		randInt, err := random.GetRandomInt(len(numberBytes))
		if err != nil {
			log.Panicf("failed to create random string %v", err)
			return ""
		}
		b[i] = numberBytes[randInt]

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
		log.Panic(fmt.Errorf("error opening local file: %v", err))
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

// EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation checks for ServiceMember in assertions, or creates one if
// none exists. Several of the document functions need a service member, but they don't always share assertions, look
// at the same assertion, or create the service members in the same ways. We'll check now to see if we already have one
// created, and if not, create one that we can place in the assertions for all the rest.
func EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation(db *pop.Connection, assertions Assertions) (Assertions, error) {
	if !assertions.Stub && assertions.ServiceMember.CreatedAt.IsZero() || assertions.ServiceMember.ID.IsNil() {
		var err error
		serviceMember, err := makeExtendedServiceMember(db, assertions)
		if err != nil {
			return Assertions{}, err
		}

		assertions.ServiceMember = serviceMember
		assertions.Order.ServiceMemberID = serviceMember.ID
		assertions.Order.ServiceMember = serviceMember
		assertions.Document.ServiceMemberID = serviceMember.ID
		assertions.Document.ServiceMember = serviceMember
	} else {
		assertions.Order.ServiceMemberID = assertions.ServiceMember.ID
		assertions.Order.ServiceMember = assertions.ServiceMember
		assertions.Document.ServiceMemberID = assertions.ServiceMember.ID
		assertions.Document.ServiceMember = assertions.ServiceMember
	}

	return assertions, nil
}

// GetOrCreateDocument checks if a document exists. If it does, it returns it, otherwise, it creates it
func GetOrCreateDocument(db *pop.Connection, document models.Document, assertions Assertions) (models.Document, error) {
	if assertions.Stub && document.CreatedAt.IsZero() || document.ID.IsNil() {
		// Ensure our doc is associated with the expected ServiceMember
		document.ServiceMemberID = assertions.ServiceMember.ID
		document.ServiceMember = assertions.ServiceMember
		// Set generic Document to have the specific assertions that were passed in
		assertions.Document = document

		return makeDocument(db, assertions)
	}

	return document, nil
}

// getOrCreateUpload checks if an upload exists. If it does, it returns it, otherwise, it creates it.
func getOrCreateUpload(db *pop.Connection, upload models.UserUpload, assertions Assertions) (models.UserUpload, error) {
	if assertions.Stub && upload.CreatedAt.IsZero() || upload.ID.IsNil() {
		// Set generic UserUpload to have the specific assertions that were passed in
		assertions.UserUpload = upload

		return makeUserUpload(db, assertions)
	}

	return upload, nil
}

// GetOrCreateDocumentWithUploads checks if a document exists. If it doesn't, it creates it. Then checks if the document
// has any uploads. If not, creates an upload associated with the document. Returns the document at the end. This
// function expects to get a specific document assertion since we're dealing with multiple documents in this overall
// file.
//
// Usage example:
//
//	emptyDocument := GetOrCreateDocumentWithUploads(db, assertions.WeightTicket.EmptyDocument, assertions)
func GetOrCreateDocumentWithUploads(db *pop.Connection, document models.Document, assertions Assertions) (models.Document, error) {
	// hang on to UserUploads, if any, for later
	userUploads := document.UserUploads

	// Ensure our doc is associated with the expected ServiceMember
	document.ServiceMemberID = assertions.ServiceMember.ID
	document.ServiceMember = assertions.ServiceMember

	var err error
	doc, err := GetOrCreateDocument(db, document, assertions)
	if err != nil {
		return models.Document{}, err
	}

	// Clear out doc.UserUploads because we'll be looping over the assertions that were passed in and potentially
	// creating data from those. It's easier to start with a clean slate than to track which ones were already created
	// vs which ones are newly created.
	doc.UserUploads = nil

	// Try getting or creating any uploads that were passed in via specific assertions
	for _, userUpload := range userUploads {
		// In case these weren't already set, set them so that they point at the correct document.
		userUpload.DocumentID = &doc.ID
		userUpload.Document = doc

		var err error
		upload, err := getOrCreateUpload(db, userUpload, assertions)
		if err != nil {
			return models.Document{}, err
		}

		doc.UserUploads = append(doc.UserUploads, upload)
	}

	// If at the end we still don't have an upload, we'll just create the default one.
	if len(doc.UserUploads) == 0 {
		// This will be overriding the assertions locally only because we have a copy rather than a pointer
		assertions.UserUpload.DocumentID = &doc.ID
		assertions.UserUpload.Document = doc

		var err error
		upload, err := makeUserUpload(db, assertions)
		if err != nil {
			return models.Document{}, err
		}

		doc.UserUploads = append(doc.UserUploads, upload)
	}

	return doc, nil
}
