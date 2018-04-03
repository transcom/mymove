package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestSSNEncryption() {
	t := suite.T()

	user, err := testdatagen.MakeUser(suite.db)
	if err != nil {
		t.Error("creating a user should have worked")
	}
	fakeSSN := "123-12-1234"

	mySSN, err := BuildSocialSecurityNumber(user.ID, fakeSSN)
	if err != nil {
		t.Error("don't expect an error here")
	}
	if mySSN.EncryptedHash == fakeSSN {
		t.Error("The encrypted hash should *not* be the same as the SSN")
	}

	shouldMatch := mySSN.Matches(fakeSSN)

	if !shouldMatch {
		t.Error("the source SSN should match the hash")
	}

	shouldNotMatch := mySSN.Matches("321-21-4321")
	if shouldNotMatch {
		t.Error("A different SSN should not match the hash")
	}

	suite.mustSave(&mySSN)
}

func (suite *ModelSuite) TestSSNFormat() {
	t := suite.T()

	user, err := testdatagen.MakeUser(suite.db)
	if err != nil {
		t.Error("creating a user should have worked")
	}

	sneakySSNs := []string{
		"123-121234",
		"123121234",
		"123 12 1234",
		"123-1  21234",
		"123.12.1234",
		"123-1 2-1234",
	}

	for _, sneakySSN := range sneakySSNs {
		_, err = BuildSocialSecurityNumber(user.ID, sneakySSN)
		if err != ErrSSNBadFormat {
			t.Error("Expected the bad formatter error.")
		}
	}

}

func (suite *ModelSuite) TestSSNSalt() {
	t := suite.T()

	user, err := testdatagen.MakeUser(suite.db)
	if err != nil {
		t.Error("creating a user should have worked")
	}
	fakeSSN := "123-12-1234"

	mySSN, err := BuildSocialSecurityNumber(user.ID, fakeSSN)
	if err != nil {
		t.Error("don't expect an error here")
	}
	if mySSN.EncryptedHash == fakeSSN {
		t.Error("The encrypted hash should *not* be the same as the SSN")
	}

	shouldMatch := mySSN.Matches(fakeSSN)

	if !shouldMatch {
		t.Error("the source SSN should match the hash")
	}

	secondSSN, err := BuildSocialSecurityNumber(user.ID, fakeSSN)
	if err != nil {
		t.Error("dont' expect an error here")
	}

	if secondSSN.EncryptedHash == mySSN.EncryptedHash {
		t.Error("These hashes should be salted, every one should be different.")
	}

	shouldMatch = secondSSN.Matches(fakeSSN)
	if !shouldMatch {
		t.Error("Even though the hash is different, it should still match our source SSN.")
	}
}

func (suite *ModelSuite) TestRawSSNNotAllowed() {
	t := suite.T()

	user, err := testdatagen.MakeUser(suite.db)
	if err != nil {
		t.Error("creating a user should have worked")
	}

	sneakySSNs := []string{
		"123-12-1234",
		"123-121234",
		"123121234",
		"123 12 1234",
		"123-1  21234",
		"123.12.1234",
		"123_12-1234",
	}

	for _, sneakySSN := range sneakySSNs {
		mySSN := SocialSecurityNumber{
			UserID:        user.ID,
			EncryptedHash: sneakySSN,
		}

		verrs, err := suite.db.ValidateAndCreate(&mySSN)
		if !verrs.HasAny() {
			t.Error("It should not be possible to save an SSN to the db.")
		}
		if err != nil {
			t.Error("It shouldn't error here though, it's a validation issue")
		}
	}

}
