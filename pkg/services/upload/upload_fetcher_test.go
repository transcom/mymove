package upload

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testUploadQueryBuilder struct {
	fakeFetchWithAssociations func(model interface{}) error
}

func (t *testUploadQueryBuilder) FetchWithAssociations(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations) error {
	m := t.fakeFetchWithAssociations(model)
	return m
}

func (suite *UploadsServiceSuite) TestFetchUpload() {
	suite.T().Run("if the upload is fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchWithAssociations := func(model interface{}) error {
			listOfUploads := reflect.ValueOf(model).Elem()
			listOfUploads.Set(reflect.Append(listOfUploads, reflect.ValueOf(models.Upload{ID: id})))
			return nil
		}

		builder := &testUploadQueryBuilder{
			fakeFetchWithAssociations: fakeFetchWithAssociations,
		}

		fetcher := NewUploadFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
		queryAssociations := []services.QueryAssociation{}
		associations := query.NewQueryAssociations(queryAssociations)

		upload, err := fetcher.FetchUploads(filters, associations)

		suite.NoError(err)
		suite.Equal(id, upload[0].ID)
	})

	suite.T().Run("if there is an error, we get it with zero uploads", func(t *testing.T) {
		fakeFetchWithAssociations := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testUploadQueryBuilder{
			fakeFetchWithAssociations: fakeFetchWithAssociations,
		}
		fetcher := NewUploadFetcher(builder)
		queryAssociations := []services.QueryAssociation{}
		associations := query.NewQueryAssociations(queryAssociations)

		uploads, err := fetcher.FetchUploads([]services.QueryFilter{}, associations)

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.Uploads{}, uploads)
	})
}