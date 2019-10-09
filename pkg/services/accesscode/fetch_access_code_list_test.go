package accesscode

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/pagination"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func (suite *AccessCodeServiceSuite) TestFetchAccessCodeListNoFilterNoAssociation() {
	ppmMove := models.SelectedMoveTypePPM
	code1 := "CODE12"
	accessCode1 := models.AccessCode{
		Code:     code1,
		MoveType: &ppmMove,
	}
	suite.MustSave(&accessCode1)
	hhgMove := models.SelectedMoveTypeHHG
	code2 := "12CODE"
	accessCode2 := models.AccessCode{
		Code:     code2,
		MoveType: &hhgMove,
	}
	suite.MustSave(&accessCode2)
	var queryFilters []services.QueryFilter
	var associations []services.QueryAssociation
	newAssociations := query.NewQueryAssociations(associations)
	queryBuilder := query.NewQueryBuilder(suite.DB())
	lf := NewAccessCodeListFetcher(queryBuilder)

	acs, err := lf.FetchAccessCodeList(queryFilters, newAssociations, defaultPagination())

	suite.NoError(err)
	suite.Len(acs, 2)
	var codes []string
	for _, e := range acs {
		codes = append(codes, e.Code)
	}
	suite.Contains(codes, code1)
	suite.Contains(codes, code2)
}

func (suite *AccessCodeServiceSuite) TestFetchAccessCodeListWithFilter() {
	ppmMove := models.SelectedMoveTypePPM
	code1 := "CODE12"
	accessCode1 := models.AccessCode{
		Code:     code1,
		MoveType: &ppmMove,
	}
	suite.MustSave(&accessCode1)
	hhgMove := models.SelectedMoveTypeHHG
	code2 := "12CODE"
	accessCode2 := models.AccessCode{
		Code:     code2,
		MoveType: &hhgMove,
	}
	suite.MustSave(&accessCode2)
	var queryFilters []services.QueryFilter
	queryFilters = append(queryFilters, query.NewQueryFilter("move_type", "=", "PPM"))
	var associations []services.QueryAssociation
	newAssociations := query.NewQueryAssociations(associations)
	queryBuilder := query.NewQueryBuilder(suite.DB())
	lf := NewAccessCodeListFetcher(queryBuilder)

	acs, err := lf.FetchAccessCodeList(queryFilters, newAssociations, defaultPagination())

	suite.NoError(err)
	suite.Len(acs, 1)
	suite.Equal(code1, acs[0].Code)
}

func (suite *AccessCodeServiceSuite) TestFetchAccessCodeListWithAssociation() {
	m := testdatagen.MakeDefaultMove(suite.DB())
	id, _ := uuid.NewV4()
	sm := m.Orders.ServiceMember
	ac := models.AccessCode{
		ID:              id,
		ServiceMemberID: &sm.ID,
		ServiceMember:   sm,
		Code:            "ABCXYZ",
		MoveType:        m.SelectedMoveType,
	}
	assertions := testdatagen.Assertions{
		AccessCode: ac,
	}
	testdatagen.MakeAccessCode(suite.DB(), assertions)
	var queryFilters []services.QueryFilter
	associations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember.Orders.Moves"),
	}
	newAssociations := query.NewQueryAssociations(associations)
	queryBuilder := query.NewQueryBuilder(suite.DB())
	lf := NewAccessCodeListFetcher(queryBuilder)

	acs, err := lf.FetchAccessCodeList(queryFilters, newAssociations, defaultPagination())

	suite.NoError(err)
	suite.Len(acs, 1)
	suite.Equal(*sm.Edipi, *acs[0].ServiceMember.Edipi)
}
