package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PayGrade represents a customer's pay grade (Including civilian)
type PayGrade struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Grade            string    `json:"grade" db:"grade"`
	GradeDescription *string   `json:"grade_description" db:"grade_description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" method
func (pg PayGrade) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Name: "Grade", Field: pg.Grade},
	), nil
}

// PayGrades is a slice of PayGrade
type PayGrades []PayGrade

func GetPayGradesForAffiliation(db *pop.Connection, affiliation string) (PayGrades, error) {
	var payGrades PayGrades

	err := db.Q().All(&payGrades)
	if err != nil {
		return nil, err
	}

	// Define excluded grades per affiliation
	excludedGrades := map[string][]string{
		string(AffiliationAIRFORCE): {string(ServiceMemberGradeMIDSHIPMAN)},
		string(AffiliationARMY):     {string(ServiceMemberGradeMIDSHIPMAN)}, //, string(ServiceMemberGradeAVIATIONCADET)},
		// string(AffiliationNAVY):       {string(ServiceMemberGradeAVIATIONCADET)},
		// string(AffiliationCOASTGUARD): {string(ServiceMemberGradeAVIATIONCADET)},
		string(AffiliationMARINES):    {string(ServiceMemberGradeMIDSHIPMAN)}, //, string(ServiceMemberGradeAVIATIONCADET)},
		string(AffiliationSPACEFORCE): {string(ServiceMemberGradeMIDSHIPMAN)}, //, string(ServiceMemberGradeAVIATIONCADET)},
	}

	if grades, ok := excludedGrades[affiliation]; ok {
		payGrades.RemoveByGrades(grades)
	}
	//  else {
	// 	// Default exclusions for unknown affiliations
	// 	payGrades.RemoveByGrades([]string{
	// 		string(ServiceMemberGradeACADEMYCADET),
	// 		string(ServiceMemberGradeMIDSHIPMAN),
	// 		string(ServiceMemberGradeAVIATIONCADET),
	// 	})
	// }

	return payGrades, nil
}

// RemoveByGrades removes all PayGrades with a Grade in the given list.
func (pgs *PayGrades) RemoveByGrades(gradesToRemove []string) {
	gradeMap := make(map[string]struct{}, len(gradesToRemove))
	for _, g := range gradesToRemove {
		gradeMap[g] = struct{}{}
	}

	newList := make(PayGrades, 0, len(*pgs))
	for _, pg := range *pgs {
		if _, found := gradeMap[pg.Grade]; !found {
			newList = append(newList, pg)
		}
	}
	*pgs = newList
}
