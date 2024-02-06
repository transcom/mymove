package testdatagen

import (
	"log"
	"strconv"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/random"
)

// RandomEdipi creates a random Edipi for a service member
func RandomEdipi() string {
	low := 1000000000
	high := 9999999999
	randInt, err := random.GetRandomIntAddend(low, high)
	if err != nil {
		log.Panicf("Failure to generate random Edipi %v", err)
	}
	return strconv.Itoa(low + int(randInt))
}

// makeServiceMember creates a single ServiceMember
// If not provided, it will also create an associated
// - User
// - ResidentialAddress
//
// Deprecated: use factory.BuildServiceMember
func makeServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	aServiceMember := assertions.ServiceMember
	user := aServiceMember.User
	agency := aServiceMember.Affiliation
	email := "leo_spaceman_sm@example.com"
	currentAddressID := aServiceMember.ResidentialAddressID
	currentAddress := aServiceMember.ResidentialAddress

	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.ServiceMember.UserID) {
		if assertions.User.OktaEmail == "" {
			assertions.User.OktaEmail = email
		}
		user = MakeDefaultUser(db)
	}
	if assertions.User.OktaEmail != "" {
		email = assertions.User.OktaEmail
	}

	if agency == nil {
		army := models.AffiliationARMY
		agency = &army
	}

	if currentAddressID == nil || isZeroUUID(*currentAddressID) {
		newAddress := MakeDefaultAddress(db)
		currentAddressID = &newAddress.ID
		currentAddress = &newAddress
	}

	randomEdipi := RandomEdipi()

	serviceMember := models.ServiceMember{
		UserID:               user.ID,
		User:                 user,
		Edipi:                models.StringPointer(randomEdipi),
		Affiliation:          agency,
		FirstName:            models.StringPointer("Leo"),
		LastName:             models.StringPointer("Spacemen"),
		Telephone:            models.StringPointer("212-123-4567"),
		PersonalEmail:        &email,
		ResidentialAddressID: currentAddressID,
		ResidentialAddress:   currentAddress,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceMember, assertions.ServiceMember)

	mustCreate(db, &serviceMember, assertions.Stub)

	return serviceMember
}

// makeExtendedServiceMember creates a single ServiceMember
// If not provided it will also create an associated
//   - User,
//   - ResidentialAddress
//   - BackupMailingAddress
//   - DutyLocation
//   - BackupContact
//
// Deprecated: use factory.BuildExtendedServiceMember
func makeExtendedServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	affiliation := assertions.ServiceMember.Affiliation
	if affiliation == nil {
		army := models.AffiliationARMY
		affiliation = &army
	}
	residentialAddress := MakeDefaultAddress(db)
	backupMailingAddress := MakeAddress2(db, assertions)

	dutyLocation := assertions.OriginDutyLocation
	if isZeroUUID(dutyLocation.ID) {
		dutyLocation = fetchOrMakeDefaultCurrentDutyLocation(db)
	}

	gbloc, err := models.FetchGBLOCForPostalCode(db, dutyLocation.Address.PostalCode)

	// Duty location must have a GBLOC associated to the postal code
	// Check for an existing GBLOC and make one if it doesn't exist
	if gbloc.GBLOC == "" || err != nil {
		makePostalCodeToGBLOC(db, dutyLocation.Address.PostalCode, "KKFA")
	}

	// Combine extended SM defaults with assertions
	smDefaults := models.ServiceMember{
		Edipi:                  models.StringPointer(RandomEdipi()),
		Affiliation:            affiliation,
		ResidentialAddressID:   &residentialAddress.ID,
		BackupMailingAddressID: &backupMailingAddress.ID,
		EmailIsPreferred:       models.BoolPointer(true),
		Telephone:              models.StringPointer("555-555-5555"),
	}

	mergeModels(&smDefaults, assertions.ServiceMember)

	assertions.ServiceMember = smDefaults

	serviceMember := makeServiceMember(db, assertions)

	contactAssertions := Assertions{
		BackupContact: models.BackupContact{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
	}
	backupContact := makeBackupContact(db, contactAssertions)
	serviceMember.BackupContacts = append(serviceMember.BackupContacts, backupContact)
	if !assertions.Stub {
		MustSave(db, &serviceMember)
	}

	return serviceMember
}
