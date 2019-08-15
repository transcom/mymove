package utilities

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
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

func (suite *UtilitiesSuite) TestSoftDestroy_NotModel() {
	accessCodeFetcher := &mocks.AccessCodeFetcher{}

	err := SoftDestroy(suite.DB(), &accessCodeFetcher)

	suite.Error(err)
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithoutDeletedAtWithoutAssociations() {
	//model without deleted_at with no associations
	user := testdatagen.MakeDefaultUser(suite.DB())

	err := SoftDestroy(suite.DB(), &user)

	suite.Error(err)
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

	err := SoftDestroy(suite.DB(), &expenseDocumentModel)
	suite.NoError(err)
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithoutDeletedAtWithAssociations() {
	// model without deleted_at with associations
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	suite.MustSave(&serviceMember)

	err := SoftDestroy(suite.DB(), &serviceMember)
	suite.Error(err)
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithDeletedAtWithAssociations() {
	// model with deleted_at with associations
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusPAYMENTREQUESTED,
		},
	})
	move := ppm.Move
	moveDoc := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
				Status:                   models.MoveDocumentStatusOK,
			},
		})
	suite.MustSave(&moveDoc)

	emptyWeight := unit.Pound(1000)
	fullWeight := unit.Pound(2500)
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDoc.ID,
		MoveDocument:             moveDoc,
		EmptyWeight:              &emptyWeight,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight,
		FullWeightTicketMissing:  false,
		VehicleNickname:          "My Car",
		VehicleOptions:           "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	suite.MustSave(&weightTicketSetDocument)
	err := SoftDestroy(suite.DB(), &moveDoc)
	suite.NoError(err)

}
