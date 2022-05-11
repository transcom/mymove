package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestMTOServiceItemDimension() {
	suite.T().Run("test valid MTOServiceItemDimension", func(t *testing.T) {
		mtoServiceItemDimensionID := uuid.Must(uuid.NewV4())

		validMTOServiceItemDimension := models.MTOServiceItemDimension{
			MTOServiceItemID: mtoServiceItemDimensionID,
			Type:             models.DimensionTypeCrate,
			Length:           0,
			Height:           0,
			Width:            0,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOServiceItemDimension, expErrors)
	})

	suite.T().Run("test invalid MTOServiceItemDimension", func(t *testing.T) {
		validMTOServiceItemDimension := models.MTOServiceItemDimension{
			MTOServiceItemID: uuid.Nil,
			Type:             "NOT VALID",
			Length:           -1,
			Height:           -1,
			Width:            -1,
		}
		expErrors := map[string][]string{
			"mtoservice_item_id": {"MTOServiceItemID can not be blank."},
			"type":               {"Type is not in the list [ITEM, CRATE]."},
			"height":             {"-1 is not greater than -1."},
			"length":             {"-1 is not greater than -1."},
			"width":              {"-1 is not greater than -1."},
		}
		suite.verifyValidationErrors(&validMTOServiceItemDimension, expErrors)
	})

	suite.T().Run("correct volume is calculated by Volume function", func(t *testing.T) {
		validMTOServiceItemDimension := models.MTOServiceItemDimension{
			Length: 6000,
			Height: 10000,
			Width:  15000,
		}
		dimensionsPointer := &validMTOServiceItemDimension
		suite.Equal(unit.CubicThousandthInch(900000000000), dimensionsPointer.Volume())
	})
}
