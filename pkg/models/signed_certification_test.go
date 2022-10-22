package models_test

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
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
				SubmittingUserID:         uuid.Must(uuid.NewV4()),
				MoveID:                   uuid.Must(uuid.NewV4()),
				PersonallyProcuredMoveID: &uuid.Nil,
				PpmID:                    &uuid.Nil,
				CertificationType:        &blankCertType,
				CertificationText:        "Lorem ipsum dolor sit amet...",
				Signature:                "Best Customer",
				Date:                     testdatagen.NextValidMoveDate,
			},
			expectedErrs: map[string][]string{
				"personally_procured_move_id": {"PersonallyProcuredMoveID can not be blank."},
				"ppm_id":                      {"PpmID can not be blank."},
				"certification_type":          {fmt.Sprintf("CertificationType is not in the list [%s].", validCertTypes)},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.signedCertification, tc.expectedErrs)
		})
	}
}

func (suite *ModelSuite) TestFetchSignedCertificationsPPMPayment() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember

	session := &auth.Session{
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
		ApplicationName: auth.MilApp,
	}

	certificationType := models.SignedCertificationTypePPMPAYMENT
	signedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	sc, err := models.FetchSignedCertificationsPPMPayment(suite.DB(), session, move.ID)
	suite.NoError(err)
	suite.Equal(signedCertification.ID, sc.ID)
}

func (suite *ModelSuite) TestFetchSignedCertificationsPPMPaymentAuth() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember
	otherPpm := testdatagen.MakeDefaultPPM(suite.DB())
	otherSm := otherPpm.Move.Orders.ServiceMember

	session := &auth.Session{
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
		ApplicationName: auth.MilApp,
	}

	certificationType := models.SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	signedCertificationType := models.SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   otherPpm.Move.ID,
			SubmittingUserID:         otherSm.UserID,
			PersonallyProcuredMoveID: &otherPpm.ID,
			CertificationType:        &signedCertificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	_, err := models.FetchSignedCertificationsPPMPayment(suite.DB(), session, otherPpm.MoveID)
	suite.Equal(errors.Cause(err), models.ErrFetchForbidden)
}

func (suite *ModelSuite) TestFetchSignedCertifications() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember

	session := &auth.Session{
		UserID:          sm.UserID,
		ServiceMemberID: sm.ID,
		ApplicationName: auth.MilApp,
	}

	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	ppmPaymentsignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &ppmPayment,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	ppmCert := models.SignedCertificationTypeSHIPMENT
	ppmSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &ppmCert,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	hhgCert := models.SignedCertificationTypeSHIPMENT
	hhgSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &hhgCert,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	scs, err := models.FetchSignedCertifications(suite.DB(), session, move.ID)
	var ids []uuid.UUID
	for _, sc := range scs {
		ids = append(ids, sc.ID)
	}

	suite.Len(scs, 3)
	suite.NoError(err)
	suite.ElementsMatch(ids, []uuid.UUID{hhgSignedCertification.ID, ppmSignedCertification.ID, ppmPaymentsignedCertification.ID})
}
