package models_test

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_SignedCertificationValidations() {
	signedCert := &SignedCertification{}

	expErrors := map[string][]string{
		"certification_text": {"CertificationText can not be blank."},
		"signature":          {"Signature can not be blank."},
	}

	suite.verifyValidationErrors(signedCert, expErrors)
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

	certificationType := SignedCertificationTypePPMPAYMENT
	signedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	sc, err := FetchSignedCertificationsPPMPayment(suite.DB(), session, move.ID)
	suite.Nil(err)
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

	certificationType := SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	signedCertificationType := SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: SignedCertification{
			MoveID:                   otherPpm.Move.ID,
			SubmittingUserID:         otherSm.UserID,
			PersonallyProcuredMoveID: &otherPpm.ID,
			CertificationType:        &signedCertificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	_, err := FetchSignedCertificationsPPMPayment(suite.DB(), session, otherPpm.MoveID)
	suite.Equal(errors.Cause(err), ErrFetchForbidden)
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

	ppmPayment := SignedCertificationTypePPMPAYMENT
	ppmPaymentsignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &ppmPayment,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	ppmCert := SignedCertificationTypePPM
	ppmSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &ppmCert,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	hhgCert := SignedCertificationTypeHHG
	hhgSignedCertification := testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &hhgCert,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	scs, err := FetchSignedCertifications(suite.DB(), session, move.ID)
	var ids []uuid.UUID
	for _, sc := range scs {
		ids = append(ids, sc.ID)
	}

	suite.Len(scs, 3)
	suite.Nil(err)
	suite.ElementsMatch(ids, []uuid.UUID{hhgSignedCertification.ID, ppmSignedCertification.ID, ppmPaymentsignedCertification.ID})
}
