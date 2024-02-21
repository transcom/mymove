package factory

import (
	"time"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *FactorySuite) TestBuildSignedCertification() {
	suite.Run("Successful creation of default notification", func() {
		// Under test:      BuildSignedCertification
		// Mocked:          None
		// Set up:          Create an SignedCertification with no customizations or traits
		// Expected outcome:SignedCertification should be created with default values

		signedCertification := BuildSignedCertification(suite.DB(), nil, nil)

		defaultCertificationType := models.SignedCertificationTypeSHIPMENT
		defaultCertificationText := "LEGAL TEXT"
		defaultSignature := "SIGNATURE"
		defaultDate := testdatagen.NextValidMoveDate
		// VALIDATE RESULTS
		suite.False(signedCertification.MoveID.IsNil())
		suite.False(signedCertification.SubmittingUserID.IsNil())
		suite.NotNil(signedCertification.CertificationType)
		suite.Equal(defaultCertificationType, *signedCertification.CertificationType)
		suite.Equal(defaultCertificationText, signedCertification.CertificationText)
		suite.Equal(defaultSignature, signedCertification.Signature)
		suite.Equal(defaultDate, signedCertification.Date)
	})

	suite.Run("Successful creation of a signedCertification with move customization", func() {
		// Under test:      BuildSignedCertification
		// Set up:          Create an SignedCertification with customized
		// attributes and Move
		// Expected outcome:Notofication should be created with custom
		// attributes

		move := BuildMove(suite.DB(), nil, nil)
		customCertificationType := models.SignedCertificationTypePPMPAYMENT
		customSignedCertification := models.SignedCertification{
			CertificationType: &customCertificationType,
			CertificationText: "CUSTOM TEXT",
			Signature:         "CUSTOM SIGNATURE",
			Date:              time.Now().Add(time.Hour * 48),
		}
		signedCertification := BuildSignedCertification(suite.DB(), []Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: customSignedCertification,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(move.ID, signedCertification.MoveID)
		suite.Equal(move.Orders.ServiceMember.UserID, signedCertification.SubmittingUserID)
		suite.NotNil(signedCertification.CertificationType)
		suite.Equal(*customSignedCertification.CertificationType,
			*signedCertification.CertificationType)
		suite.Equal(customSignedCertification.CertificationText,
			signedCertification.CertificationText)
		suite.Equal(customSignedCertification.Signature,
			signedCertification.Signature)
		suite.Equal(customSignedCertification.Date,
			signedCertification.Date)
	})

	suite.Run("Successful creation of a signedCertification with service member customization", func() {
		// Under test:      BuildSignedCertification
		// Set up:          Create an SignedCertification with customized
		// attributes and ServiceMember
		// Expected outcome:Notofication should be created with custom
		// attributes

		serviceMember := BuildExtendedServiceMember(suite.DB(), nil, nil)
		customCertificationType := models.SignedCertificationTypePPMPAYMENT
		customSignedCertification := models.SignedCertification{
			CertificationType: &customCertificationType,
			CertificationText: "CUSTOM TEXT",
			Signature:         "CUSTOM SIGNATURE",
			Date:              time.Now().Add(time.Hour * 48),
		}
		signedCertification := BuildSignedCertification(suite.DB(), []Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
			{
				Model: customSignedCertification,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(serviceMember.UserID, signedCertification.SubmittingUserID)
		move, err := models.FetchMove(suite.DB(), &auth.Session{}, signedCertification.MoveID)
		suite.NoError(err)
		suite.Equal(serviceMember.ID, move.Orders.ServiceMemberID)
		suite.NotNil(signedCertification.CertificationType)
		suite.Equal(*customSignedCertification.CertificationType,
			*signedCertification.CertificationType)
		suite.Equal(customSignedCertification.CertificationText,
			signedCertification.CertificationText)
		suite.Equal(customSignedCertification.Signature,
			signedCertification.Signature)
		suite.Equal(customSignedCertification.Date,
			signedCertification.Date)
	})

	suite.Run("Successful creation of stubbed signedCertification", func() {
		// Under test:      BuildSignedCertification
		// Set up:          Create a stubbed signedCertification, but don't pass in a db
		// Expected outcome:SignedCertification should be created with
		// stubbed service member, no signedCertification should be created in database
		precount, err := suite.DB().Count(&models.SignedCertification{})
		suite.NoError(err)

		signedCertification := BuildSignedCertification(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(signedCertification.MoveID.IsNil())
		suite.True(signedCertification.SubmittingUserID.IsNil())
		suite.Equal("SIGNATURE", signedCertification.Signature)

		// Count how many signedCertification are in the DB, no new
		// signedCertifications should have been created
		count, err := suite.DB().Count(&models.SignedCertification{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}
