package serviceitem

import (
	"testing"

	"github.com/gobuffalo/validate"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type testServiceItemQueryBuilder struct {
	fakeFetchOne  func(model interface{}) error
	fakeCreateOne func(models interface{}) (*validate.Errors, error)
}

func (t *testServiceItemQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (t *testServiceItemQueryBuilder) CreateOne(model interface{}) (*validate.Errors, error) {
	return nil, nil
}

type ServiceItemServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestUserSuite(t *testing.T) {

	hs := &ServiceItemServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}
