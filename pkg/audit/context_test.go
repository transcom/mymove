package audit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
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

	suite.Run("WithAuditUser returns context with audit user", func() {
		var emptyUser models.User
		user := factory.BuildDefaultUser(suite.DB())
		context := context.Background()

		// If nothing is added to the context an empty user should be returned
		returnedEmptyUser := RetrieveAuditUserFromContext(context)
		suite.Equal(emptyUser, returnedEmptyUser)

		// If an audit user is added to the context that audit user should be returned
		context = WithAuditUser(context, user)
		returnedUser := RetrieveAuditUserFromContext(context)
		suite.Equal(user, returnedUser)
	})
}
