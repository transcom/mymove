package user

import "github.com/transcom/mymove/pkg/models"

type testOfficeUserQueryBuilder struct {
	fakeFetchOne func(model interface{}) error
}

func (t *testOfficeUserQueryBuilder) FetchOne(model interface{}, field string, value interface{}) error {
	m := t.fakeFetchOne(model)
	return m
}

func (suite *UserServiceSuite) TestFetchOfficeUser() {
	testEmail := "test@example.com"
	fakeFetchOne := func(model interface{}) error {
		*model = models.OfficeUser{Email: testEmail}
		return nil
	}

	builder := &testOfficeUserQueryBuilder{
		fakeFetchOne: fakeFetchOne,
	}
	fetcher := NewOfficeUserFetcher(builder)

	officeUser, err := fetcher.FetchOfficeUser("id", "1")

	suite.NoError(err)
	suite.Equal(testEmail, officeUser.Email)
}
