package audit

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type AuditSuite struct {
	*testingsuite.PopTestSuite
}

func TestAuditSuite(t *testing.T) {

	hs := &AuditSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *AuditSuite) TestContextFunctions() {
	suite.Run("WithEventName returns context with event name", func() {
		eventName := "eventName"
		context := context.Background()

		// If nothing is added to the context an empty string should be returned
		emptyStringEventName := RetrieveEventNameFromContext(context)
		suite.Equal("", emptyStringEventName)

		// If an event name is added to the context that event name should be returned
		context = WithEventName(context, eventName)
		returnedEventName := RetrieveEventNameFromContext(context)
		suite.Equal(eventName, returnedEventName)
	})

	suite.Run("WithAuditUserID returns context with audit user", func() {
		var nilUUID uuid.UUID
		user := factory.BuildDefaultUser(suite.DB())
		context := context.Background()

		// If nothing is added to the context an empty user should be returned
		returnedNilUserID := RetrieveAuditUserIDFromContext(context)
		suite.Equal(nilUUID, returnedNilUserID)

		// If an audit user is added to the context that audit user should be returned
		context = WithAuditUserID(context, user.ID)
		returnedUserID := RetrieveAuditUserIDFromContext(context)
		suite.Equal(user.ID, returnedUserID)
	})
}
