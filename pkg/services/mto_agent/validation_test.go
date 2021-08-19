package mtoagent

import (
	"context"
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *MTOAgentServiceSuite) TestMergeAgent() {
	firstName := "Jason"
	lastName := "Ash"
	email := "jason.ash@example.com"
	phone := "202-555-9301"
	oldAgent := models.MTOAgent{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     &email,
		Phone:     &phone,
	}

	newFirstName := "First"
	newEmail := "email@email.email"
	newPhone := ""

	successAgent := oldAgent
	successAgent.FirstName = &newFirstName
	successAgent.Email = &newEmail
	successAgent.Phone = &newPhone

	newAgent := mergeAgent(successAgent, &oldAgent)

	suite.Equal(*newAgent.FirstName, *successAgent.FirstName)
	suite.Equal(*newAgent.Email, *successAgent.Email)
	suite.Nil(newAgent.Phone)

	// Checking that the old agent instances weren't changed:
	suite.NotEqual(*newAgent.FirstName, *oldAgent.FirstName)
	suite.NotNil(oldAgent.Phone)
}

func (suite *MTOAgentServiceSuite) TestValidateMTOAgent() {
	na := models.MTOAgent{ID: uuid.Must(uuid.NewV4())}
	oa := models.MTOAgent{ID: uuid.Must(uuid.NewV4())}
	sh := models.MTOShipment{ID: uuid.Must(uuid.NewV4())}

	// these checks just ensure the parameters are being passed as expected
	checkNew := mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		suite.Equal(newAgent.ID, na.ID)
		return nil
	})
	checkOld := mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		suite.Equal(oldAgent.ID, oa.ID)
		return nil
	})
	checkShip := mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		suite.Equal(shipment.ID, sh.ID)
		return nil
	})

	checkEmpty := mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		verrs := validate.NewErrors()
		return verrs
	})
	checkVerr := mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		verrs := validate.NewErrors()
		verrs.Add("forceVERR", "forced")
		return verrs
	})
	checkErr := mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		return fmt.Errorf("forced error, not of type *validate.Errors")
	})
	checkSkip := mtoAgentValidatorFunc(func(ctx context.Context, newAgent models.MTOAgent, oldAgent *models.MTOAgent, shipment *models.MTOShipment) error {
		suite.Fail("should not have been called after a non-verr short-circuit")
		return nil
	})

	testCases := map[string]struct {
		checks []mtoAgentValidator
		verf   func(error)
	}{
		"happy path": {
			[]mtoAgentValidator{
				checkNew,
				checkOld,
				checkShip,
				checkEmpty,
			},
			func(err error) {
				suite.NoError(err)
			},
		},
		"short circuit": {
			[]mtoAgentValidator{
				checkVerr,
				checkEmpty,
				checkErr,
				checkSkip,
			},
			func(err error) {
				suite.Error(err)
				switch verr := err.(type) {
				case *validate.Errors:
					suite.Fail("did not expect a *validate.Errors", "%v", verr)
				}
			},
		},
		"only verrs": {
			[]mtoAgentValidator{
				checkVerr,
				checkEmpty,
				checkVerr,
				checkEmpty,
			},
			func(err error) {
				suite.Error(err)
				switch e := err.(type) {
				case services.InvalidInputError:
					suite.True(e.ValidationErrors.HasAny())
					suite.Contains(e.ValidationErrors.Keys(), "forceVERR")
				default:
					suite.IsType(services.InvalidInputError{}, err)
				}
			},
		},
	}

	for name, tc := range testCases {
		suite.Run(name, func() {
			tc.verf(validateMTOAgent(context.Background(), na, &oa, &sh, tc.checks...))
		})
	}
}
