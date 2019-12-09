package models_test

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGenerateReferenceID() {
	r, err := models.GenerateReferenceID(suite.DB())
	suite.NotNil(r)
	referenceID := *r
	suite.NoError(err)
	firstNum, _ := strconv.Atoi(strings.Split(referenceID, "-")[0])
	secondNum, _ := strconv.Atoi(strings.Split(referenceID, "-")[1])
	suite.Equal(reflect.TypeOf(referenceID).String(), "string")
	suite.Equal(firstNum >= 0 && firstNum <= 9999, true)
	suite.Equal(secondNum >= 0 && secondNum <= 9999, true)
	suite.Equal(string(referenceID[4]), "-")
}

func (suite *ModelSuite) TestGenerateContractorID() {
	contractor := testdatagen.MakePrimeContractor(suite.DB(), testdatagen.Assertions{})
	mto := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:           uuid.FromStringOrNil("5d4b25bb-eb04-4c03-9a81-ee0398cb779e"),
			ContractorID: &contractor.ID,
			Contractor:   &contractor,
		},
	})
	testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{
		ServiceItem: models.ServiceItem{MoveTaskOrder: mto}},
	)
	testdatagen.MakeEntitlement(suite.DB(), testdatagen.Assertions{
		GHCEntitlement: models.GHCEntitlement{MoveTaskOrder: &mto}},
	)

	testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:           uuid.FromStringOrNil("1c030e51-b5be-40a2-80bf-97a330891307"),
			Status:       models.MoveTaskOrderStatusDraft,
			ContractorID: &contractor.ID,
			Contractor:   &contractor,
		},
	})

	mtoID := "1c030e51-b5be-40a2-80bf-97a330891308"

	testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:     uuid.FromStringOrNil(mtoID),
			Status: models.MoveTaskOrderStatusDraft,
		},
	})

	order := models.MoveTaskOrder{}

	err := suite.DB().Find(&order, mtoID)
	log.Println("Order Contractor", order.Contractor)
	suite.NoError(err)

	suite.Equal(&contractor.ID, mto.ContractorID)
}