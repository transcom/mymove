package factory

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/random"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// CONSTANTS

// DefaultContractCode is the default contract code for testing
const DefaultContractCode = "TRUSS_TEST"

// Customization type is the building block for passing in customizations and traits
type Customization struct {
	Model interface{} // The model that the factory will build
	Type  *CustomType // Custom type, usually same as model,
	// except for models that appear in multiple fields
	LinkOnly       bool        // Tells factory to just link in this model, do not create
	ExtendedParams interface{} // Some models need extra data for creation, pointer to extended params
}

// CustomType is a string that represents what kind of customization it is
type CustomType string

// Model customizations will each have a "type" here
// This does not have to match the model type but generally will
// You can have CustomType like ResidentialAddress to define specifically
// where this address will get created and nested
var Address CustomType = "Address"
var AdditionalDocuments CustomType = "AdditionalDocuments"
var AdminUser CustomType = "AdminUser"
var AuditHistory CustomType = "AuditHistory"
var BackupContact CustomType = "BackupContact"
var BoatShipment CustomType = "BoatShipment"
var City CustomType = "City"
var ClientCert CustomType = "ClientCert"
var Contractor CustomType = "Contractor"
var Country CustomType = "Country"
var CustomerSupportRemark CustomType = "CustomerSupportRemark"
var Document CustomType = "Document"
var DutyLocation CustomType = "DutyLocation"
var Entitlement CustomType = "Entitlement"
var UBAllowance CustomType = "UBAllowances"
var EvaluationReport CustomType = "EvaluationReport"
var LineOfAccounting CustomType = "LineOfAccounting"
var MobileHome CustomType = "MobileHome"
var Move CustomType = "Move"
var MovingExpense CustomType = "MovingExpense"
var MTOAgent CustomType = "MTOAgent"
var MTOServiceItem CustomType = "MTOServiceItem"
var MTOServiceItemDimension CustomType = "MTOServiceItemDimension"
var MTOShipment CustomType = "MTOShipment"
var Notification CustomType = "Notification"
var OfficePhoneLine CustomType = "OfficePhoneLine"
var OfficeUser CustomType = "OfficeUser"
var Order CustomType = "Order"
var Organization CustomType = "Organization"
var PPMShipment CustomType = "PPMShipment"
var PaymentRequest CustomType = "PaymentRequest"
var PaymentServiceItem CustomType = "PaymentServiceItem"
var PaymentServiceItemParam CustomType = "PaymentServiceItemParam"
var PaymentRequestToInterchangeControlNumber CustomType = "PaymentRequestToInterchangeControlNumber"
var Port CustomType = "Port"
var PortLocation CustomType = "PortLocation"
var PostalCodeToGBLOC CustomType = "PostalCodeToGBLOC"
var PrimeUpload CustomType = "PrimeUpload"
var ProgearWeightTicket CustomType = "ProgearWeightTicket"
var ProofOfServiceDoc CustomType = "ProofOfServiceDoc"
var ReService CustomType = "ReService"
var Role CustomType = "Role"
var ServiceItemParamKey CustomType = "ServiceItemParamKey"
var ServiceParam CustomType = "ServiceParam"
var ServiceMember CustomType = "ServiceMember"
var ServiceRequestDocument CustomType = "ServiceRequestDocument"
var ServiceRequestDocumentUpload CustomType = "ServiceRequestDocumentUpload"
var ShipmentAddressUpdate CustomType = "ShipmentAddressUpdate"
var SignedCertification CustomType = "SignedCertification"
var SITDurationUpdate CustomType = "SITDurationUpdate"
var State CustomType = "State"
var StorageFacility CustomType = "StorageFacility"
var TransportationAccountingCode CustomType = "TransportationAccountingCode"
var TransportationOffice CustomType = "TransportationOffice"
var TransportationOfficeAssignment CustomType = "TransportationOfficeAssignment"
var Upload CustomType = "Upload"
var UserUpload CustomType = "UserUpload"
var User CustomType = "User"
var UsersRoles CustomType = "UsersRoles"
var WebhookNotification CustomType = "WebhookNotification"
var WeightTicket CustomType = "WeightTicket"
var UsPostRegionCity CustomType = "UsPostRegionCity"
var UsersPrivileges CustomType = "UsersPrivileges"
var Privilege CustomType = "Privilege"

// defaultTypesMap allows us to assign CustomTypes for most default types
var defaultTypesMap = map[string]CustomType{
	"models.Address":                                  Address,
	"models.AdminUser":                                AdminUser,
	"factory.TestDataAuditHistory":                    AuditHistory,
	"models.BackupContact":                            BackupContact,
	"models.BoatShipment":                             BoatShipment,
	"models.City":                                     City,
	"models.ClientCert":                               ClientCert,
	"models.Contractor":                               Contractor,
	"models.Country":                                  Country,
	"models.CustomerSupportRemark":                    CustomerSupportRemark,
	"models.Document":                                 Document,
	"models.DutyLocation":                             DutyLocation,
	"models.Entitlement":                              Entitlement,
	"models.UBAllowances":                             UBAllowance,
	"models.EvaluationReport":                         EvaluationReport,
	"models.LineOfAccounting":                         LineOfAccounting,
	"models.MobileHome":                               MobileHome,
	"models.Move":                                     Move,
	"models.MovingExpense":                            MovingExpense,
	"models.MTOAgent":                                 MTOAgent,
	"models.MTOServiceItem":                           MTOServiceItem,
	"models.MTOServiceItemDimension":                  MTOServiceItemDimension,
	"models.MTOShipment":                              MTOShipment,
	"models.Notification":                             Notification,
	"models.OfficePhoneLine":                          OfficePhoneLine,
	"models.OfficeUser":                               OfficeUser,
	"models.Order":                                    Order,
	"models.Organization":                             Organization,
	"models.PaymentRequest":                           PaymentRequest,
	"models.PaymentServiceItem":                       PaymentServiceItem,
	"models.PaymentServiceItemParam":                  PaymentServiceItemParam,
	"models.PaymentRequestToInterchangeControlNumber": PaymentRequestToInterchangeControlNumber,
	"models.PPMShipment":                              PPMShipment,
	"models.Port":                                     Port,
	"models.PortLocation":                             PortLocation,
	"models.PostalCodeToGBLOC":                        PostalCodeToGBLOC,
	"models.PrimeUpload":                              PrimeUpload,
	"models.ProgearWeightTicket":                      ProgearWeightTicket,
	"models.ProofOfServiceDoc":                        ProofOfServiceDoc,
	"models.ReService":                                ReService,
	"models.ServiceItemParamKey":                      ServiceItemParamKey,
	"models.ServiceMember":                            ServiceMember,
	"models.ServiceRequestDocument":                   ServiceRequestDocument,
	"models.ServiceRequestDocumentUpload":             ServiceRequestDocumentUpload,
	"models.ServiceParam":                             ServiceParam,
	"models.SignedCertification":                      SignedCertification,
	"models.ShipmentAddressUpdate":                    ShipmentAddressUpdate,
	"models.SITDurationUpdate":                        SITDurationUpdate,
	"models.State":                                    State,
	"models.StorageFacility":                          StorageFacility,
	"models.TransportationAccountingCode":             TransportationAccountingCode,
	"models.UsPostRegionCity":                         UsPostRegionCity,
	"models.TransportationOffice":                     TransportationOffice,
	"models.TransportationOfficeAssignment":           TransportationOfficeAssignment,
	"models.Upload":                                   Upload,
	"models.UserUpload":                               UserUpload,
	"models.User":                                     User,
	"models.UsersRoles":                               UsersRoles,
	"models.WebhookNotification":                      WebhookNotification,
	"models.WeightTicket":                             WeightTicket,
	"roles.Role":                                      Role,
	"models.UsersPrivileges":                          UsersPrivileges,
	"models.Privilege":                                Privilege,
}

// Instead of nesting structs, we create specific CustomTypes here to give devs
// a code-completion friendly way to select the right type

// addressGroup is a grouping of all address related fields
type addressGroup struct {
	PickupAddress                 CustomType
	DeliveryAddress               CustomType
	SecondaryPickupAddress        CustomType
	SecondaryDeliveryAddress      CustomType
	TertiaryPickupAddress         CustomType
	TertiaryDeliveryAddress       CustomType
	ResidentialAddress            CustomType
	BackupMailingAddress          CustomType
	DutyLocationAddress           CustomType
	DutyLocationTOAddress         CustomType
	SITOriginHHGOriginalAddress   CustomType
	SITOriginHHGActualAddress     CustomType
	SITDestinationFinalAddress    CustomType
	SITDestinationOriginalAddress CustomType
	W2Address                     CustomType
	OriginalAddress               CustomType
	NewAddress                    CustomType
}

// Addresses is the struct to access the various fields externally
var Addresses = addressGroup{
	PickupAddress:                 "PickupAddress",
	DeliveryAddress:               "DeliveryAddress",
	SecondaryPickupAddress:        "SecondaryPickupAddress",
	SecondaryDeliveryAddress:      "SecondaryDeliveryAddress",
	TertiaryPickupAddress:         "TertiaryPickupAddress",
	TertiaryDeliveryAddress:       "TertiaryDeliveryAddress",
	ResidentialAddress:            "ResidentialAddress",
	BackupMailingAddress:          "BackupMailingAddress",
	DutyLocationAddress:           "DutyLocationAddress",
	DutyLocationTOAddress:         "DutyLocationTOAddress",
	SITOriginHHGOriginalAddress:   "SITOriginHHGOriginalAddress",
	SITOriginHHGActualAddress:     "SITOriginHHGActualAddress",
	SITDestinationFinalAddress:    "SITDestinationFinalAddress",
	SITDestinationOriginalAddress: "SITDestinationOriginalAddress",
	W2Address:                     "W2Address",
	OriginalAddress:               "OriginalAddress",
	NewAddress:                    "NewAddress",
}

// dimensionGroup is a grouping of all the Dimension related fields
type dimensionGroup struct {
	CrateDimension CustomType
	ItemDimension  CustomType
}

// Dimensions is the struct to access the fields externally
var Dimensions = dimensionGroup{
	// MTOServiceItems may include:
	CrateDimension: "CrateDimension",
	ItemDimension:  "ItemDimension",
}

type documentGroup struct {
	UploadedOrders        CustomType
	UploadedAmendedOrders CustomType
}

var Documents = documentGroup{
	// Orders may include:
	UploadedOrders:        "UploadedOrders",
	UploadedAmendedOrders: "UploadedAmendedOrders",
}

// dutyLocationsGroup is a grouping of all the duty location related fields
type dutyLocationsGroup struct {
	OriginDutyLocation CustomType
	NewDutyLocation    CustomType
}

// DutyLocations is the struct to access the fields externally
var DutyLocations = dutyLocationsGroup{
	// Orders may include:
	OriginDutyLocation: "OriginDutyLocation",
	NewDutyLocation:    "NewDutyLocation",
}

// transportationOfficeGroup is a grouping of all the transportation office related fields
type transportationOfficeGroup struct {
	OriginDutyLocation CustomType
	NewDutyLocation    CustomType
	CloseoutOffice     CustomType
	CounselingOffice   CustomType
}

// TransportationOffices is the struct to access the fields externally
var TransportationOffices = transportationOfficeGroup{
	OriginDutyLocation: "OriginDutyLocationTransportationOffice",
	NewDutyLocation:    "NewDutyLocationTransportationOffice",
	CloseoutOffice:     "CloseoutOffice",
	CounselingOffice:   "CounselingOffice",
}

type officeUserGroup struct {
	SCAssignedUser  CustomType
	TIOAssignedUser CustomType
	TOOAssignedUser CustomType
}

var OfficeUsers = officeUserGroup{
	SCAssignedUser:  "SCAssignedUser",
	TIOAssignedUser: "TIOAssignedUser",
	TOOAssignedUser: "TOOAssignedUser",
}

// type officeUserGroup struct {
// 	TIOAssignedUser CustomType
// }

// var OfficeUsers = officeUserGroup{
// 	TIOAssignedUser: "TIOAssignedUser",
// }

// uploadGroup is a grouping of all the upload related fields
type uploadGroup struct {
	UploadTypePrime CustomType
	UploadTypeUser  CustomType
}

// Uploads is the struct to access the fields externally
var Uploads = uploadGroup{
	UploadTypePrime: "UploadTypePrime",
	UploadTypeUser:  "UploadTypeUser",
}

// portLocationGroup is a grouping of all the port related fields
type portLocationGroup struct {
	PortOfDebarkation CustomType
	PortOfEmbarkation CustomType
}

// PortLocations is the struct to access the fields externally
var PortLocations = portLocationGroup{
	PortOfDebarkation: "PODLocation",
	PortOfEmbarkation: "POELocation",
}

// Below are errors returned by various functions

var ErrNestedModel = errors.New("NESTED_MODEL_ERROR")

// Trait is a function that returns a set of customizations
// Every Trait should start with GetTrait for discoverability
type Trait func() []Customization

// assignType uses the model name to assign the CustomType
// if it's already assigned, do not reassign
func assignType(custom *Customization) error {
	if custom.Type != nil {
		return nil
	}
	// Get the model and check that it's a struct
	model := custom.Model
	mv := reflect.ValueOf(model)
	if mv.Kind() != reflect.Struct {
		return fmt.Errorf("Customization.Model field had type %v - should contain a struct", mv.Kind())
	}
	// Get the model type and find the default type
	typestring, ok := defaultTypesMap[mv.Type().String()]
	if ok {
		custom.Type = &typestring
	} else {
		return fmt.Errorf("Customization.Model field had type %v which is not supported in defaultTypesMap", mv.Type().String())
	}
	return nil
}

// setDefaultTypes assigns types to all customizations in the list provided
func setDefaultTypes(clist []Customization) {
	for idx := 0; idx < len(clist); idx++ {
		if err := assignType(&clist[idx]); err != nil {
			log.Panic(err.Error())
		}
	}
}

// linkOnlyHasID ensures LinkOnly customizations have an ID
func linkOnlyHasID(clist []Customization) error {
	for idx := 0; idx < len(clist); idx++ {
		if clist[idx].LinkOnly && !hasID(clist[idx].Model) {
			return fmt.Errorf("Customization was LinkOnly but the Model had no ID. LinkOnly models must have ID")
		}
	}
	return nil
}

// setupCustomizations prepares the customizations customs for the factory
// by applying and merging the traits.
// customs is a slice that will be modified by setupCustomizations.
//
//   - Ensures linkOnly customizations have ID
//   - Merges customizations and traits
//   - Assigns default types to all default customizations
//   - Ensure there's only one customization per type
func setupCustomizations(customs []Customization, traits []Trait) []Customization {

	// Ensure LinkOnly customizations all have ID
	err := linkOnlyHasID(customs)
	if err != nil {
		log.Panic(err)
	}

	// Merge customizations with traits (also sets default types)
	customs = mergeCustomization(customs, traits)

	// Ensure unique customizations
	err = isUnique(customs)
	if err != nil {
		log.Panic(err)
	}

	return customs

}

// isUnique ensures there's only one customization per type.
// Requirement: All customizations should already have a type assigned.
func isUnique(customs []Customization) error {
	// validate that there are no repeat CustomTypes
	m := make(map[CustomType]int)
	for i, custom := range customs {
		if custom.Type == nil {
			log.Panic("All customizations should have type.")
		}
		// if custom type already exists
		idx, exists := m[*custom.Type]
		if exists {
			return fmt.Errorf("found more than one instance of %s Customization at index %d and %d",
				*custom.Type, idx, i)
		}
		// Add to hashmap
		m[*custom.Type] = i
	}
	return nil

}

// findCustomWithIdx is a helper function to find a customization of a specific type and its index
// Returns:
//   - index of the found customization
//   - pointer to the customization
func findCustomWithIdx(customs []Customization, customType CustomType) (int, *Customization) {
	for i := 0; i < len(customs); i++ {
		if customs[i].Type != nil && *customs[i].Type == customType {
			return i, &customs[i]
		}
	}
	return -1, nil
}

// toStructPtr takes an interface wrapping a struct and
// returns an interface wrapping a pointer to the struct
// For e.g. interface{}(models.User) → interface{}(*models.User)
func toStructPtr(obj interface{}) interface{} {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp.Interface()
}

// toInterfacePtr takes an interface wrapping a pointer to a struct and
// returns an interface wrapping the struct
// For e.g. interface{}(*models.User) → interface{}(models.User)
func toInterfacePtr(obj interface{}) interface{} {
	rv := reflect.ValueOf(obj).Elem()
	return rv.Interface()
}

// mergeInterfaces transforms the interfaces to match what mergeModels expects
func mergeInterfaces(model1 interface{}, model2 interface{}) interface{} {
	// This wrapper function is needed because merge models expects the first
	// model to be an interface containing a pointer to struct.
	// This function converts
	//    an interface containing the struct
	//    → an interface containing a pointer to the struct
	modelPtr := toStructPtr(model1)
	testdatagen.MergeModels(modelPtr, model2)
	model := toInterfacePtr(modelPtr)
	return model
}

func hasID(model interface{}) bool {
	mv := reflect.ValueOf(model)

	// mv should be a model of type struct
	if mv.Kind() != reflect.Struct || !strings.HasPrefix(mv.Type().String(), "models.") {
		log.Panic("Expecting interface containing a model")
	}

	return !mv.FieldByName("ID").IsZero()
}

// mergeCustomization takes the original set of customizations
// and merges with the traits.
// Required : All customizations MUST have types
//
// The order of application is
//   - Earlier traits override later traits in the trait list
//   - Customizations override the traits
//
// So if you have [trait1, trait2] customization
// and all three contain the same object:
//   - trait 1 will override trait 2 (so start with the highest priority)
//   - customization will override trait 2
func mergeCustomization(customs []Customization, traits []Trait) []Customization {
	// Make sure all customs have a proper type
	setDefaultTypes(customs)

	// Iterate list of traits, each could return a list of customizations
	for _, trait := range traits {

		// Get customizations and set default types
		traitCustomizations := trait()
		setDefaultTypes(traitCustomizations)

		// for each trait custom, merge or replace the one in user supplied customizations
		for _, traitCustom := range traitCustomizations {
			j, callerCustom := findCustomWithIdx(customs, *traitCustom.Type)
			if callerCustom != nil {
				// If a customization is marked as LinkOnly, it means we use that precreated object
				// Therefore we can't merge a trait with it, as we don't update fields on pre-created
				// objects. So we only merge if LinkOnly is false.
				if !callerCustom.LinkOnly {
					result := mergeInterfaces(traitCustom.Model, callerCustom.Model)
					callerCustom.Model = result
					customs[j] = *callerCustom
				}
			} else {
				customs = append(customs, traitCustom)
			}
		}
	}
	return customs
}

// Helper function for when you need to elevate the type of customization from say
// ResidentialAddress to Address before you call BuildAddress
// This is a little finicky because we want to be careful not to harm the existing list
// TBD should we validate again here?
func convertCustomizationInList(customs []Customization, from CustomType, to CustomType) []Customization {
	if _, custom := findCustomWithIdx(customs, to); custom != nil {
		log.Panic(fmt.Errorf("a customization of type %s already exists", to))
	}
	if idx, custom := findCustomWithIdx(customs, from); custom != nil {
		// Create a slice in new memory
		var newCustoms []Customization
		// Populate with copies of objects
		newCustoms = append(newCustoms, customs...)
		// Update the type
		newCustoms[idx].Type = &to
		return newCustoms
	}
	log.Panic(fmt.Errorf("no customization of type %s found", from))
	return nil
}

func findValidCustomization(customs []Customization, customType CustomType) *Customization {
	_, custom := findCustomWithIdx(customs, customType)
	if custom == nil {
		return nil
	}

	// Else check that the customization is valid
	if err := checkNestedModels(*custom); err != nil {
		if errors.Is(err, ErrNestedModel) {
			if !custom.LinkOnly {
				log.Panic(fmt.Errorf("errors encountered in customization for %s: %w", *custom.Type, err))
			}
		} else {
			log.Panic(fmt.Errorf("errors encountered in customization for %s: %w", *custom.Type, err))
		}
	}
	return custom
}

func replaceCustomization(customs []Customization, newCustom Customization) []Customization {
	err := linkOnlyHasID([]Customization{newCustom})
	if err != nil {
		log.Panic(err)
	}

	// assign the type as all customizations should have a type
	if err := assignType(&newCustom); err != nil {
		log.Panic(err.Error())
	}
	// See if an existing customization exists with the type
	ndx, _ := findCustomWithIdx(customs, *newCustom.Type)
	if ndx >= 0 {
		// Found a customization for the provided model and we need to
		// replace it
		customs[ndx] = newCustom
	} else {
		// Did not find an existing customization, append it
		customs = append(customs, newCustom)
	}

	return customs
}

// Caller should have already setup Customizations using setupCustomizations
func removeCustomization(customs []Customization, customType CustomType) []Customization {
	// See if an existing customization exists with the type
	ndx, _ := findCustomWithIdx(customs, customType)
	if ndx >= 0 {
		// Found a customization for the provided model and we need to remove it
		// Order shouldn't matter because the setupCustomizations should have already merged customizations and traits
		// and ensured that there's only one customization per type
		// Replace the customization we want to remove with the last customization in the slice
		customs[ndx] = customs[len(customs)-1]
		customs = customs[:len(customs)-1]
	}

	return customs
}

// checkNestedModels ensures we have no nested models.
// - If a field is a struct that is a model, it should be empty
// - If a field is a pointer to struct, and that struct is a model it should be nil
func checkNestedModels(c interface{}) error {

	// c IS THE CUSTOMIZATION, SHOULD NOT BE A POINTER
	// c SHOULD NOT BE A POINTER
	if reflect.ValueOf(c).Kind() == reflect.Pointer {
		return fmt.Errorf("function expects a struct, received a pointer")
	}

	// mv IS THE MODEL VALUE, SHOULD NOT BE EMPTY
	mv := reflect.ValueOf(c).FieldByName("Model") // get the interface
	if mv.IsNil() {
		return fmt.Errorf("customization must contain a model")
	}
	mv = mv.Elem() // get the model from the interface

	// mv SHOULD BE A STRUCT
	if mv.Kind() == reflect.Struct {
		numberOfFields := mv.NumField()
		mt := mv.Type() // get the model type

		// CHECK ALL FIELDS IN THE STRUCT
		for i := 0; i < numberOfFields; i++ {
			fieldName := mt.Field(i).Name
			field := mv.Field(i)

			// There are a couple conditions we want to check for
			// - If a field is a struct that is a model, it should be empty
			// - If a field is a pointer to struct, and that struct is a model it should be nil

			// IF A FIELD IS A MODELS STRUCT - SHOULD BE EMPTY
			ft := field.Type()
			if field.Kind() == reflect.Struct && strings.HasPrefix(ft.String(), "models.") {
				if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
					return fmt.Errorf("%s cannot be populated, no nested models allowed: %w", fieldName, ErrNestedModel)
				}
			}

			// IF A FIELD IS A POINTER TO A MODELS STRUCT - SHOULD ALSO BE EMPTY
			if field.Kind() == reflect.Pointer {
				nf := field.Elem()
				if !field.IsNil() && nf.Kind() == reflect.Struct && strings.HasPrefix(nf.Type().String(), "models.") {
					return fmt.Errorf("%s cannot be populated, no nested models allowed: %w", fieldName, ErrNestedModel)
				}
			}
		}
	}
	return nil
}

func mustCreate(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndCreate(model)
	if err != nil {
		log.Panic(fmt.Errorf("errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("validation errors encountered saving %#v: %v", model, verrs))
	}
}

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		log.Panic(fmt.Errorf("errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("validation errors encountered saving %#v: %v", model, verrs))
	}
}

// RandomEdipi creates a random Edipi for a service member
func RandomEdipi() string {
	low := 1000000000
	high := 9999999999
	randInt, err := random.GetRandomIntAddend(low, high)
	if err != nil {
		log.Panicf("failure to generate random Edipi %v", err)
	}
	return strconv.Itoa(low + int(randInt))
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

func setBoolPtr(customBoolPtr *bool, defaultBool bool) *bool {
	result := &defaultBool
	if customBoolPtr != nil {
		result = customBoolPtr
	}
	return result
}

// FixtureOpen opens a file from the testdata dir
func FixtureOpen(name string) afero.File {
	fixtureDir := "testdata"
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(fmt.Errorf("failed to get current directory: %s", err))
	}

	// if this is called from inside another package, we want to drop everything
	// including and after 'pkg'
	// 'a/b/c/pkg/something/' → 'a/b/c/'
	cwd = strings.Split(cwd, "pkg")[0]

	fixturePath := path.Join(cwd, "pkg/testdatagen", fixtureDir, name)
	file, err := os.Open(filepath.Clean(fixturePath))
	if err != nil {
		log.Panic(fmt.Errorf("error opening local file: %v", err))
	}

	return file
}
