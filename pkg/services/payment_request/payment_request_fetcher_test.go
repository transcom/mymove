package paymentrequest

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testPaymentRequestQueryBuilder struct {
	fakeFetchOne func(model interface{}) error
}

func (t *testPaymentRequestQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequest() {
	suite.T().Run("If a payment request is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testPaymentRequestQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		fetcher := NewPaymentRequestFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		adminUser, err := fetcher.FetchPaymentRequest(filters)

		suite.NoError(err)
		suite.Equal(id, adminUser.ID)
	})

	suite.T().Run("if there is an error, we get it with zero payment request", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testPaymentRequestQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewPaymentRequestFetcher(builder)

		paymentRequest, err := fetcher.FetchPaymentRequest([]services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.PaymentRequest{}, paymentRequest)
	})
}
