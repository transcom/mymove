package models_test

import (
	"strings"

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
	suite.True(strings.Contains(err.Error(), "FETCH_FORBIDDEN"))
}
