package mtoagent

import (
	"context"
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *MTOAgentServiceSuite) TestValidationRules() {
	suite.Run("checkShipmentID", func() {
		suite.Run("success", func() {
			newA := models.MTOAgent{MTOShipmentID: uuid.Must(uuid.NewV4())}
			testCases := map[string]struct {
				newA models.MTOAgent
				oldA *models.MTOAgent
			}{
				"create": {
					newA: newA,
					oldA: nil,
				},
				"update": {
					newA: newA,
					oldA: &models.MTOAgent{MTOShipmentID: newA.MTOShipmentID},
				},
			}
			for name, tc := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(context.Background(), tc.newA, tc.oldA, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newA models.MTOAgent
				oldA *models.MTOAgent
			}{
				"create": {
					newA: models.MTOAgent{},
					oldA: nil,
				},
				"update": {
					newA: models.MTOAgent{MTOShipmentID: id},
					oldA: &models.MTOAgent{},
				},
			}
			for name, tc := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(context.Background(), tc.newA, tc.oldA, nil)
					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(verr.HasAny())
						suite.Contains(verr.Keys(), "mtoShipmentID")
					default:
						suite.Failf("expected *validate.Errors", "%t - %v", err, err)
					}
				})
			}
		})
	})

	suite.Run("checkAgentID", func() {
		suite.Run("success", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newA models.MTOAgent
				oldA *models.MTOAgent
			}{
				"create": {
					newA: models.MTOAgent{},
					oldA: nil,
				},
				"update": {
					newA: models.MTOAgent{ID: id},
					oldA: &models.MTOAgent{ID: id},
				},
			}
			for name, tc := range testCases {
				suite.Run(name, func() {
					err := checkAgentID().Validate(context.Background(), tc.newA, tc.oldA, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newA models.MTOAgent
				oldA *models.MTOAgent
				verr bool
			}{
				"create": {
					newA: models.MTOAgent{ID: id},
					oldA: nil,
					verr: true,
				},
				"update": {
					newA: models.MTOAgent{ID: id},
					oldA: &models.MTOAgent{ID: uuid.Must(uuid.NewV4())},
					verr: false,
				},
			}
			for name, tc := range testCases {
				suite.Run(name, func() {
					err := checkAgentID().Validate(context.Background(), tc.newA, tc.oldA, nil)
					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(tc.verr, "expected something other than a *validate.Errors type")
						suite.Contains(verr.Keys(), "ID")
					default:
						suite.False(tc.verr, "expected a *validate.Errors: %t - naid %s", err, tc.newA.ID)
					}
				})
			}

		})
	})

	suite.Run("checkContactInfo", func() {
		firstName := "Jason"
		lastName := "Ash"
		email := "jason.ash@example.com"
		phone := "202-555-9301"
		oldAgent := &models.MTOAgent{
			FirstName: &firstName,
			LastName:  &lastName,
			Email:     &email,
			Phone:     &phone,
		}

		suite.Run("success", func() {
			firstName := "Carol"
			lastName := ""
			email := ""
			phone := "234-555-4567"

			agent := models.MTOAgent{
				FirstName: &firstName,
				LastName:  &lastName,
				Email:     &email,
				Phone:     &phone,
			}

			err := checkContactInfo().Validate(context.Background(), agent, oldAgent, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.NoVerrs(verr)
			default:
				suite.Failf("expected *validate.Errs", "%v", err)
			}
		})

		// Test unsuccessful check for contact info
		suite.Run("failure", func() {
			firstName := ""
			email := ""
			phone := ""

			agent := models.MTOAgent{
				FirstName: &firstName,
				Email:     &email,
				Phone:     &phone,
			}

			err := checkContactInfo().Validate(context.Background(), agent, oldAgent, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Contains(verr.Keys(), "firstName")
				suite.Contains(verr.Keys(), "contactInfo")
			default:
				suite.Failf("expected *validate.Errs", "%v", err)
			}
		})
	})

	suite.Run("checkAgentType", func() {
		shipment := models.MTOShipment{
			ID: uuid.Must(uuid.NewV4()),
		}
		oldAgent := models.MTOAgent{
			ID:            uuid.Must(uuid.NewV4()),
			MTOShipmentID: shipment.ID,
			MTOShipment:   shipment,
			MTOAgentType:  models.MTOAgentReceiving,
		}
		shipment.MTOAgents = models.MTOAgents{
			oldAgent,
		}

		suite.Run("success", func() {
			lastName := "Baker"
			testCases := map[string]struct {
				newA models.MTOAgent
				oldA *models.MTOAgent
				ship *models.MTOShipment
			}{
				"unchanged MTOAgentType": {
					newA: models.MTOAgent{LastName: &lastName},
					ship: &shipment,
				},
				"valid MTOAgentType change": {
					newA: models.MTOAgent{
						ID:           oldAgent.ID,
						MTOAgentType: models.MTOAgentReleasing, // oldAgent is RECEIVING, so we're switching types
					},
					oldA: &oldAgent,
					ship: &shipment,
				},
			}
			for name, tc := range testCases {
				suite.Run(name, func() {
					err := checkAgentType().Validate(context.Background(), tc.newA, tc.oldA, tc.ship)
					suite.NoError(err, "Unexpected error from checkAgentType: %v", err)
				})
			}
		})

		suite.Run("failure", func() {
			maxed := models.MTOShipment{
				ID: uuid.Must(uuid.NewV4()),
			}
			maxed.MTOAgents = models.MTOAgents{
				models.MTOAgent{
					MTOShipmentID: maxed.ID,
					MTOShipment:   maxed,
					MTOAgentType:  models.MTOAgentReceiving,
				},
				models.MTOAgent{
					MTOShipmentID: maxed.ID,
					MTOShipment:   maxed,
					MTOAgentType:  models.MTOAgentReleasing,
				},
			}

			testCases := map[string]struct {
				newA models.MTOAgent
				ship *models.MTOShipment
				verf func(error)
			}{
				"agent type collision": {
					newA: models.MTOAgent{
						// No ID because we're simulating a create
						MTOAgentType: oldAgent.MTOAgentType, // oldAgent is RECEIVING, so this is the same type
					},
					ship: &shipment,
					verf: func(err error) {
						suite.Error(err, "Unexpectedly no error from checkAgentType with duplicated MTOAgentType")
						suite.IsType(services.ConflictError{}, err)
						suite.Contains(err.Error(), models.MTOAgentReceiving)
					},
				},
				"incorrect usage": {
					ship: nil,
					verf: func(err error) {
						suite.IsType(services.ImplementationError{}, err)
					},
				},
				"maxed out number of agents": {
					newA: models.MTOAgent{
						// No ID because we're simulating a create
						MTOAgentType: models.MTOAgentReceiving, // value doesn't matter, but we need one to validate
					},
					ship: &maxed,
					verf: func(err error) {
						suite.Error(err, "Unexpectedly no error from checkAgentType with max number of agents")
						suite.IsType(services.ConflictError{}, err)
						suite.Contains(err.Error(), "This shipment already has 2 agents - no more can be added")
					},
				},
			}

			for name, tc := range testCases {
				suite.Run(name, func() {
					err := checkAgentType().Validate(context.Background(), tc.newA, nil, tc.ship)
					tc.verf(err)
				})
			}
		})
	})

	suite.Run("checkPrimeAvailability", func() {
		shipment := models.MTOShipment{
			ID: uuid.Must(uuid.NewV4()),
		}
		checkHappy := primeFunc(func(uuid.UUID) (bool, error) {
			return true, nil
		})
		checkError := primeFunc(func(uuid.UUID) (bool, error) {
			return false, fmt.Errorf("forced")
		})
		testCases := map[string]struct {
			check services.MoveTaskOrderChecker
			ship  *models.MTOShipment
			err   bool
		}{
			"happy": {
				check: checkHappy,
				ship:  &shipment,
				err:   false,
			},
			"error": {
				check: checkError,
				ship:  &shipment,
				err:   true,
			},
			"misused": {
				check: checkHappy,
				ship:  nil,
				err:   true,
			},
		}

		for name, tc := range testCases {
			suite.Run(name, func() {
				err := checkPrimeAvailability(tc.check).Validate(context.Background(), models.MTOAgent{}, nil, tc.ship)
				if err == nil {
					if tc.err {
						suite.Fail("expected error")
					}
					return
				}
				suite.IsType(services.NotFoundError{}, err)
			})
		}

	})
}

type primeFunc func(uuid.UUID) (bool, error)

func (fn primeFunc) MTOAvailableToPrime(id uuid.UUID) (bool, error) {
	return fn(id)
}
