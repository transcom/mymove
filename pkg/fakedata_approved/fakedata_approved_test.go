package fakedata

import (
	"fmt"
)

type fakeDataTestCase struct {
	firstName      string
	lastName       string
	address        string
	phone          string
	email          string
	expected       bool
	expectedStrict bool
}

var fakeDataTestCases = []fakeDataTestCase{
	// STRICT MATCH IS TRUE
	{
		// 0
		firstName:      "Jason",
		lastName:       "Ash",
		address:        "448 Washington Blvd NE",
		phone:          "999-999-9999",
		email:          "test@email.com",
		expected:       true,
		expectedStrict: true,
	},
	{
		// 1
		firstName:      "Gregory",
		lastName:       "Van der Heide",
		address:        "6622 Airport Way S #1430",
		phone:          "123-555-9999",
		email:          "test@example.com",
		expected:       true,
		expectedStrict: true,
	},
	{
		// 2
		firstName:      "Christopher",
		lastName:       "Swinglehurst-Walters",
		address:        "4124 Apache Dr, Apt 18C",
		phone:          "456-555-9359",
		email:          "test@truss.works",
		expected:       true,
		expectedStrict: true,
	},
	{
		// 3
		firstName:      "Jayden",
		lastName:       "Jackson Jr.",
		address:        "441 SW Río de la Plata Drive",
		phone:          "456-555-9359",
		email:          "test@email.com",
		expected:       true,
		expectedStrict: true,
	},
	// STRICT MATCH IS FALSE
	{
		// 4
		firstName:      "Christopher",
		lastName:       "Swinglehurst Walters",   //"Swinglehurst-Walters"
		address:        "4124 Apache Dr Apt 18C", //"4124 Apache Dr, Apt 18C",
		phone:          "456-555-9359",
		email:          "test@email.com",
		expected:       true,
		expectedStrict: false,
	},
	{
		// 5
		firstName:      "Jayden",
		lastName:       "Jackson Jr",              // "Jackson Jr."
		address:        "6622 Airport Way S 1430", //"6622 Airport Way S #1430"
		phone:          "456-555-9359",
		email:          "test@email.com",
		expected:       true,
		expectedStrict: false,
	},
	{
		// 6
		firstName:      "Barbara",
		lastName:       "St Juste",                                 //"St. Juste"
		address:        "1292 Orchard Terrace, Building C Unit 10", //"1292 Orchard Terrace, Building C, Unit 10",
		phone:          "456-555-9359",
		email:          "test@email.com",
		expected:       true,
		expectedStrict: false,
	},
	// FAKE DATA IS NOT VALID
	{
		// 7
		firstName:      "Paul", //"Jason"
		lastName:       "Ash",
		address:        "448 Washington NE", //"448 Washington Blvd NE",
		phone:          "999-199-9999",
		email:          "test@google.com",
		expected:       false,
		expectedStrict: false,
	},
	{
		// 8
		firstName:      "Gregxry",                //"Gregory"
		lastName:       "Vaz der Zeide",          //"Van der Heide"
		address:        "99 Airport Way S #1430", //"6622 Airport Way S #1430"
		phone:          "123-55-9999",
		email:          "test@ex.com",
		expected:       false,
		expectedStrict: false,
	},
}

func (suite *FakeDataSuite) TestFakeDataTestCases() {
	var result bool
	var err error
	for i, testCase := range fakeDataTestCases {
		result, err = IsValidFakeDataFullName(testCase.firstName, testCase.lastName)
		if !suite.Equal(testCase.expected, result) {
			suite.Fail(fmt.Sprintf("Failure on IsValidFakeDataFullName test case %d (0 indexed)\n", i))
		}
		suite.NoError(err)
		result, err = IsValidFakeDataFullNameStrict(testCase.firstName, testCase.lastName)
		if !suite.Equal(testCase.expectedStrict, result) {
			suite.Fail(fmt.Sprintf("Failure on IsValidFakeDataFullNameStrict test case %d (0 indexed)\n", i))
		}
		suite.NoError(err)
		result, err = IsValidFakeDataName(fmt.Sprintf("%s %s", testCase.firstName, testCase.lastName))
		if !suite.Equal(testCase.expected, result) {
			suite.Fail(fmt.Sprintf("Failure on IsValidFakeDataName test case %d (0 indexed)\n", i))
		}
		suite.NoError(err)
		result, err = IsValidFakeDataAddress(testCase.address)
		if !suite.Equal(testCase.expected, result) {
			suite.Fail(fmt.Sprintf("Failure on IsValidFakeDataAddress test case %d (0 indexed)\n", i))
		}
		suite.NoError(err)
		result, err = IsValidFakeDataAddressStrict(testCase.address)
		if !suite.Equal(testCase.expectedStrict, result) {
			suite.Fail(fmt.Sprintf("Failure on IsValidFakeDataAddressStrict test case %d (0 indexed)\n", i))
		}
		suite.NoError(err)
		result, err = IsValidFakeDataPhone(testCase.phone)
		if !suite.Equal(testCase.expected, result) {
			suite.Fail(fmt.Sprintf("Failure on IsValidFakeDataPhone test case %d (0 indexed)\n", i))
		}
		suite.NoError(err)
		result, err = IsValidFakeDataEmail(testCase.email)
		if !suite.Equal(testCase.expected, result) {
			suite.Fail(fmt.Sprintf("Failure on IsValidFakeDataEmail test case %d (0 indexed)\n", i))
		}
		suite.NoError(err)
	}
}
