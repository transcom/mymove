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

	shouldMatch := mySSN.MatchesRawSSN(fakeSSN)

	if !shouldMatch {
		t.Error("the source SSN should match the hash")
	}

	shouldNotMatch := mySSN.MatchesRawSSN("321-21-4321")
	if shouldNotMatch {
		t.Error("A different SSN should not match the hash")
	}

	suite.mustSave(&mySSN)
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

	shouldMatch := mySSN.MatchesRawSSN(fakeSSN)

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

	shouldMatch = secondSSN.MatchesRawSSN(fakeSSN)
	if !shouldMatch {
		t.Error("Even though the hash is different, it should still match our source SSN.")
	}
}
