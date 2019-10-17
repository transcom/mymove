package upload

import (
	"errors"
	"reflect"
	"testing"

	"github.com/transcom/mymove/pkg/services/pagination"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testUploadQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
}

func (t *testUploadQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) error {
	m := t.fakeFetchMany(model)
	return m
}

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func (suite *UploadsServiceSuite) TestFetchUploads() {
	suite.T().Run("if uploads are fetched, it should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
			listOfUploads := reflect.ValueOf(model).Elem()
			listOfUploads.Set(reflect.Append(listOfUploads, reflect.ValueOf(models.Upload{ID: id})))
			return nil
		}

		builder := &testUploadQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewUploadFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
		queryAssociations := []services.QueryAssociation{}
		associations := query.NewQueryAssociations(queryAssociations)

		uploadRecords, err := fetcher.FetchUploads(filters, associations, defaultPagination())

		suite.NoError(err)
		suite.Equal(id, uploadRecords[0].ID)
	})

	suite.T().Run("if there is an error, we get it with zero uploads", func(t *testing.T) {
		fakeFetchWithAssociations := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testUploadQueryBuilder{
			fakeFetchMany: fakeFetchWithAssociations,
		}
		fetcher := NewUploadFetcher(builder)
		queryAssociations := []services.QueryAssociation{}
		associations := query.NewQueryAssociations(queryAssociations)

		uploads, err := fetcher.FetchUploads([]services.QueryFilter{}, associations, defaultPagination())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.Uploads{}, uploads)
	})
}