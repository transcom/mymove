package models_test

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestSignedCertificationValidations() {
	blankCertType := models.SignedCertificationType("")
	validCertTypes := strings.Join(models.AllowedSignedCertificationTypes, ", ")
	testCases := map[string]struct {
		signedCertification models.SignedCertification
		expectedErrs        map[string][]string
	}{
		"Can Validate Successfully": {
			signedCertification: models.SignedCertification{
				SubmittingUserID:  uuid.Must(uuid.NewV4()),
				MoveID:            uuid.Must(uuid.NewV4()),
				CertificationText: "Lorem ipsum dolor sit amet...",
				Signature:         "Best Customer",
				Date:              testdatagen.NextValidMoveDate,
			},
			expectedErrs: nil,
		},
		"Catches Missing Required Fields": {
			signedCertification: models.SignedCertification{},
			expectedErrs: map[string][]string{
				"submitting_user_id": {"SubmittingUserID can not be blank."},
				"move_id":            {"MoveID can not be blank."},
				"certification_text": {"CertificationText can not be blank."},
				"signature":          {"Signature can not be blank."},
				"date":               {"Date can not be blank."},
			},
		},
		"Validates Optional Fields": {
			signedCertification: models.SignedCertification{
				SubmittingUserID:  uuid.Must(uuid.NewV4()),
				MoveID:            uuid.Must(uuid.NewV4()),
				PpmID:             &uuid.Nil,
				CertificationType: &blankCertType,
				CertificationText: "Lorem ipsum dolor sit amet...",
				Signature:         "Best Customer",
				Date:              testdatagen.NextValidMoveDate,
			},
			expectedErrs: map[string][]string{
				"ppm_id":             {"PpmID can not be blank."},
				"certification_type": {fmt.Sprintf("CertificationType is not in the list [%s].", validCertTypes)},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.signedCertification, tc.expectedErrs, nil)
		})
	}
}

func (suite *ModelSuite) TestFetchSignedCertificationsPPMPayment() {

	move := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)

	sm := move.Orders.ServiceMember

	session := &auth.Session{
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
		ApplicationName: auth.MilApp,
	}

	certificationType := models.SignedCertificationTypePPMPAYMENT
	signedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &certificationType,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)

	sc, err := models.FetchSignedCertificationsPPMPayment(suite.DB(), session, move.ID)
	suite.NoError(err)
	suite.Equal(signedCertification.ID, sc.ID)
}

func (suite *ModelSuite) TestFetchSignedCertificationsPPMPaymentAuth() {
	ppmMove1 := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)
	sm := ppmMove1.Orders.ServiceMember
	ppmMove2 := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)

	session := &auth.Session{
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
		ApplicationName: auth.MilApp,
	}

	certificationType := models.SignedCertificationTypePPMPAYMENT
	factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    ppmMove1,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &certificationType,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)

	signedCertificationType := models.SignedCertificationTypePPMPAYMENT
	factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    ppmMove2,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &signedCertificationType,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)

	_, err := models.FetchSignedCertificationsPPMPayment(suite.DB(), session, ppmMove2.ID)
	suite.Equal(errors.Cause(err), models.ErrFetchForbidden)
}

func (suite *ModelSuite) TestFetchSignedCertifications() {
	move := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)
	sm := move.Orders.ServiceMember

	session := &auth.Session{
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
		ApplicationName: auth.MilApp,
	}

	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	ppmPaymentsignedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &ppmPayment,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)
	ppmCert := models.SignedCertificationTypeSHIPMENT
	ppmSignedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &ppmCert,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)
	hhgCert := models.SignedCertificationTypeSHIPMENT
	hhgSignedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &hhgCert,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)

	scs, err := models.FetchSignedCertifications(suite.DB(), session, move.ID)
	var ids []uuid.UUID
	for _, sc := range scs {
		ids = append(ids, sc.ID)
	}

	suite.Len(scs, 3)
	suite.NoError(err)
	suite.ElementsMatch(ids, []uuid.UUID{hhgSignedCertification.ID, ppmSignedCertification.ID, ppmPaymentsignedCertification.ID})
}

func (suite *ModelSuite) TestFetchSignedCertificationsByType() {
	move := factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)
	sm := move.Orders.ServiceMember

	session := &auth.Session{
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
		ApplicationName: auth.MilApp,
	}

	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	ppmPaymentsignedCertification := factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &ppmPayment,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
				PpmID:             models.UUIDPointer(move.MTOShipments[0].PPMShipment.ID),
			},
		},
	}, nil)

	scs, err := models.FetchSignedCertificationPPMByType(suite.DB(), session, move.ID, move.MTOShipments[0].PPMShipment.ID, models.SignedCertificationTypePPMPAYMENT)
	var ids []uuid.UUID
	for _, sc := range scs {
		ids = append(ids, sc.ID)
	}

	suite.Len(scs, 1)
	suite.NoError(err)
	suite.ElementsMatch(ids, []uuid.UUID{ppmPaymentsignedCertification.ID})
}
