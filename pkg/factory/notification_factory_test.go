package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildNotification() {
	suite.Run("Successful creation of default notification", func() {
		// Under test:      BuildNotification
		// Mocked:          None
		// Set up:          Create an Notification with no customizations or traits
		// Expected outcome:Notification should be created with default values

		notification := BuildNotification(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.False(notification.ServiceMemberID.IsNil())
		suite.False(notification.ServiceMember.ID.IsNil())
		suite.Equal("", notification.SESMessageID)
		suite.Equal(models.NotificationTypes(""), notification.NotificationType)
	})

	suite.Run("Successful creation of a notification with customization", func() {
		// Under test:      BuildNotification
		// Set up:          Create an Notification with a customized
		// attributes and ServiceMember
		// Expected outcome:Notofication should be created with custom
		// attributes

		serviceMember := BuildServiceMember(suite.DB(), nil, nil)
		customNotification := models.Notification{
			SESMessageID:     "123",
			NotificationType: models.MovePaymentReminderEmail,
		}
		notification := BuildNotification(suite.DB(), []Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
			{
				Model: customNotification,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(serviceMember.ID, notification.ServiceMemberID)
		suite.Equal(serviceMember.ID, notification.ServiceMember.ID)
		suite.Equal(customNotification.SESMessageID, notification.SESMessageID)
		suite.Equal(customNotification.NotificationType, notification.NotificationType)
	})

	suite.Run("Successful creation of stubbed notification", func() {
		// Under test:      BuildNotification
		// Set up:          Create a stubbed notification, but don't pass in a db
		// Expected outcome:Notification should be created with
		// stubbed service member, no notification should be created in database
		precount, err := suite.DB().Count(&models.Notification{})
		suite.NoError(err)

		notification := BuildNotification(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(notification.ServiceMemberID.IsNil())
		suite.True(notification.ServiceMember.ID.IsNil())
		suite.NotNil(notification.ServiceMember.Edipi)
		suite.Equal("", notification.SESMessageID)
		suite.Equal(models.NotificationTypes(""), notification.NotificationType)

		// Count how many notification are in the DB, no new
		// notifications should have been created
		count, err := suite.DB().Count(&models.Notification{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}
