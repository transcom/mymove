package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildAuditHistory() {

	suite.Run("Successful creation of default audit history", func() {
		// Under test:      BuildAuditHistory
		// Mocked:          None
		// Set up:          Create an audit history with no customizations or traits
		// Expected outcome:audit history should be created with default values

		// SETUP
		// Create a default audit history to compare values

		fakeContext := models.JSONMap{}
		fakeOldData := models.JSONMap{"Locator": "asdf"}
		fakeChangedData := models.JSONMap{"Locator": "fdsa"}

		defaulthistory := TestDataAuditHistory{
			SchemaName:    "public",
			TableNameDB:   "moves",
			RelID:         16592,
			Context:       &fakeContext,
			ContextID:     models.StringPointer("1234567"),
			TransactionID: models.Int64Pointer(1234),
			ClientQuery:   models.StringPointer("select 1;"),
			Action:        "UPDATE",
			EventName:     models.StringPointer("setFinancialReviewFlag"),
			OldData:       &fakeOldData,
			ChangedData:   &fakeChangedData,
		}

		// FUNCTION UNDER TEST
		history := BuildAuditHistory(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaulthistory.SchemaName, history.SchemaName)
		suite.Equal(defaulthistory.TableNameDB, history.TableNameDB)
		suite.Equal(defaulthistory.RelID, history.RelID)
		suite.Equal(defaulthistory.Context, history.Context)
		suite.Equal(defaulthistory.ContextID, history.ContextID)
		suite.Equal(defaulthistory.TransactionID, history.TransactionID)
		suite.Equal(defaulthistory.ClientQuery, history.ClientQuery)
		suite.Equal(defaulthistory.Action, history.Action)
		suite.Equal(defaulthistory.EventName, history.EventName)
		suite.Equal(defaulthistory.OldData, history.OldData)
		suite.Equal(defaulthistory.ChangedData, history.ChangedData)
		suite.Nil(history.ObjectID)
		suite.Nil(history.SessionUserID)
	})

	suite.Run("Successful creation of customized audit history", func() {
		// Under test:      BuildAuditHistory
		// Mocked:          None
		// Set up:          Create audit history with customization
		// Expected outcome:audit history should be created with customized values

		// SETUP
		oldData := models.JSONMap{"Locator": "qwer"}
		changedData := models.JSONMap{"Locator": "rewq"}

		customHistory := TestDataAuditHistory{
			SchemaName:    "public1",
			TableNameDB:   "test_table",
			RelID:         55555,
			ObjectID:      models.UUIDPointer(uuid.Must(uuid.NewV4())),
			SessionUserID: models.UUIDPointer(uuid.Must(uuid.NewV4())),
			ContextID:     models.StringPointer("1234567"),
			TransactionID: models.Int64Pointer(1234),
			ClientQuery:   models.StringPointer("select 2;"),
			Action:        "INSERT",
			EventName:     models.StringPointer("testEvent"),
			OldData:       &oldData,
			ChangedData:   &changedData,
		}

		// FUNCTION UNDER TEST
		history := BuildAuditHistory(suite.DB(), []Customization{
			{Model: customHistory},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customHistory.SchemaName, history.SchemaName)
		suite.Equal(customHistory.TableNameDB, history.TableNameDB)
		suite.Equal(customHistory.RelID, history.RelID)
		suite.Equal(*customHistory.ContextID, *history.ContextID)
		suite.Equal(*customHistory.TransactionID, *history.TransactionID)
		suite.Equal(*customHistory.ClientQuery, *history.ClientQuery)
		suite.Equal(customHistory.Action, history.Action)
		suite.Equal(*customHistory.EventName, *history.EventName)
		suite.Equal(customHistory.OldData, history.OldData)
		suite.Equal(customHistory.ChangedData, history.ChangedData)
		suite.Equal(*customHistory.ObjectID, *history.ObjectID)
		suite.Equal(*customHistory.SessionUserID, *history.SessionUserID)
	})

	suite.Run("Successful creation of stubbed audit history", func() {
		// Under test:      BuildAuditHistory
		// Set up:          Create a stubbed audit history
		// Expected outcome:No new audit history should be created

		// Check num historyDurationUpdates
		precount, err := suite.DB().Count(&TestDataAuditHistory{})
		suite.NoError(err)

		history := BuildAuditHistory(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(history.ID.IsNil())

		// Count how many notification are in the DB, no new
		// audit history should have been created
		count, err := suite.DB().Count(&TestDataAuditHistory{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}
