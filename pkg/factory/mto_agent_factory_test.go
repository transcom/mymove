package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildMTOAgent() {
	suite.Run("Successful creation of default agent", func() {
		// Under test:      BuildMTOAgent
		// Set up:          Create an agent with no customizations or traits
		// Expected outcome: agent should be created with default values

		// SETUP
		defaultAgent := models.MTOAgent{
			FirstName:    models.StringPointer("Jason"),
			LastName:     models.StringPointer("Ash"),
			Email:        models.StringPointer("jason.ash@example.com"),
			Phone:        models.StringPointer("202-555-9301"),
			MTOAgentType: models.MTOAgentReleasing,
		}

		agent := BuildMTOAgent(suite.DB(), nil, nil)

		suite.Equal(defaultAgent.FirstName, agent.FirstName)
		suite.Equal(defaultAgent.LastName, agent.LastName)
		suite.Equal(defaultAgent.Email, agent.Email)
		suite.Equal(defaultAgent.Phone, agent.Phone)
		suite.Equal(defaultAgent.MTOAgentType, agent.MTOAgentType)
		suite.False(agent.MTOShipmentID.IsNil())
		suite.NotNil(agent.MTOShipment)
		suite.False(agent.MTOShipment.ID.IsNil())
		suite.Equal(models.MTOAgentReleasing, agent.MTOAgentType)
	})

	suite.Run("Successful creation of custom agent", func() {
		// Under test:      BuildMTOAgent
		// Set up:          Create an agent and pass custom fields
		// Expected outcome: agent should be created with custom values

		// SETUP
		customAgent := models.MTOAgent{
			FirstName:    models.StringPointer("Riley"),
			LastName:     models.StringPointer("Baker"),
			Email:        models.StringPointer("rbaker@example.com"),
			Phone:        models.StringPointer("555-555-5555"),
			MTOAgentType: models.MTOAgentReceiving,
		}
		customMove := models.Move{
			Locator: "AAAA",
		}
		customShipment := models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		}

		// CALL FUNCTION UNDER TEST
		agent := BuildMTOAgent(suite.DB(), []Customization{
			{Model: customAgent},
			{Model: customMove},
			{Model: customShipment},
		}, nil)

		suite.Equal(customAgent.FirstName, agent.FirstName)
		suite.Equal(customAgent.LastName, agent.LastName)
		suite.Equal(customAgent.Email, agent.Email)
		suite.Equal(customAgent.Phone, agent.Phone)
		suite.Equal(customAgent.MTOAgentType, agent.MTOAgentType)
		suite.Equal(customMove.Locator, agent.MTOShipment.MoveTaskOrder.Locator)
		suite.False(agent.MTOShipmentID.IsNil())
		suite.NotNil(agent.MTOShipment)
		suite.False(agent.MTOShipment.ID.IsNil())
	})

	suite.Run("Successful creation of stubbed agent", func() {
		// Under test:      BuildMTOAgent
		// Set up:          Create a stubbed agent
		// Expected outcome:No new agent should be created in the db
		precount, err := suite.DB().Count(&models.MTOAgent{})
		suite.NoError(err)

		agent := BuildMTOAgent(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(agent.MTOShipmentID.IsNil())
		suite.NotNil(agent.MTOShipment)
		suite.True(agent.MTOShipment.ID.IsNil())

		// Count how many notification are in the DB, no new
		// notifications should have been created
		count, err := suite.DB().Count(&models.MTOAgent{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful return of linkOnly agent", func() {
		// Under test:       BuildMTOAgent
		// Set up:           Pass in a linkOnly agent
		// Expected outcome: No new agent should be created.

		// Check num agent records
		precount, err := suite.DB().Count(&models.MTOAgent{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		agent := BuildMTOAgent(suite.DB(), []Customization{
			{
				Model: models.MTOAgent{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.MTOAgent{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, agent.ID)
	})
}
