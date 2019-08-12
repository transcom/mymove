package utilities

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/db/dbfmt"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type UtilitiesSuite struct {
	testingsuite.PopTestSuite
}

func (suite *UtilitiesSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestUtilitiesSuite(t *testing.T) {
	hs := &UtilitiesSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, hs)
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithoutDeletedAtWithoutAssociations() {
	//model without deleted_at with no associations
	user := testdatagen.MakeDefaultUser(suite.DB())

	dbfmt.Println(user)

	err := SoftDestroy(suite.DB(), user)

	suite.NoError(err)
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithoutDeletedAtWithAssociations() {
	// model without deleted_at with associations
	// service member

}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithDeletedAtWithAssociations() {
	// model with deleted_at with associations
	// uploads

}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithDeletedAtWithoutAssociations() {
	//model with deleted_at with no associations
	expenseDocumentModel := testdatagen.MakeMovingExpenseDocument(suite.DB(), testdatagen.Assertions{
		MovingExpenseDocument: models.MovingExpenseDocument{
			MovingExpenseType:    models.MovingExpenseTypeCONTRACTEDEXPENSE,
			PaymentMethod:        "GTCC",
			RequestedAmountCents: unit.Cents(10000),
		},
	})

	suite.MustSave(&expenseDocumentModel)

	err := SoftDestroy(suite.DB(), expenseDocumentModel)
	suite.NoError(err)
}
