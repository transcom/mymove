package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/apperror"
)

const (
	BasePeriodYear1 string = "Base Period Year 1"
	BasePeriodYear2 string = "Base Period Year 2"
	BasePeriodYear3 string = "Base Period Year 3"
	OptionPeriod1   string = "Option Period 1"
	OptionPeriod2   string = "Option Period 2"
	OptionPeriod3   string = "Option Period 3"
	AwardTerm1      string = "Award Term 1"
	AwardTerm2      string = "Award Term 2"
	AwardTerm       string = "Award Term"
	OptionPeriod    string = "Option Period"
	BasePeriodYear  string = "Base Period Year"
)

type ExpectedEscalationPriceContractsCount struct {
	ExpectedAmountOfContractYearsForCalculation     int
	ExpectedAmountOfBasePeriodYearsForCalculation   int
	ExpectedAmountOfOptionPeriodYearsForCalculation int
	ExpectedAmountOfAwardTermsForCalculation        int
}

// ReContractYear represents a single "year" of a contract
type ReContractYear struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	ContractID           uuid.UUID `json:"contract_id" db:"contract_id"`
	Name                 string    `json:"name" db:"name"`
	StartDate            time.Time `json:"start_date" db:"start_date"`
	EndDate              time.Time `json:"end_date" db:"end_date"`
	Escalation           float64   `json:"escalation" db:"escalation"`
	EscalationCompounded float64   `json:"escalation_compounded" db:"escalation_compounded"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`

	// Associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
}

// TableName overrides the table name used by Pop.
func (r ReContractYear) TableName() string {
	return "re_contract_years"
}

type ReContractYears []ReContractYear

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReContractYear) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
		&validators.TimeIsPresent{Field: r.StartDate, Name: "StartDate"},
		&validators.TimeIsPresent{Field: r.EndDate, Name: "EndDate"},
		&validators.TimeAfterTime{FirstTime: r.EndDate, FirstName: "EndDate", SecondTime: r.StartDate, SecondName: "StartDate"},
		&Float64IsPresent{Field: r.Escalation, Name: "Escalation"},
		&Float64IsGreaterThan{Field: r.Escalation, Name: "Escalation", Compared: 0},
		&Float64IsPresent{Field: r.EscalationCompounded, Name: "EscalationCompounded"},
		&Float64IsGreaterThan{Field: r.EscalationCompounded, Name: "EscalationCompounded", Compared: 0},
	), nil
}

func GetExpectedEscalationPriceContractsCount(contractYearName string, hasOptionYear3 bool) (ExpectedEscalationPriceContractsCount, error) {
	switch contractYearName {
	case BasePeriodYear1:
		return ExpectedEscalationPriceContractsCount{
			ExpectedAmountOfContractYearsForCalculation:     1,
			ExpectedAmountOfBasePeriodYearsForCalculation:   1,
			ExpectedAmountOfOptionPeriodYearsForCalculation: 0,
			ExpectedAmountOfAwardTermsForCalculation:        0,
		}, nil
	case BasePeriodYear2:
		return ExpectedEscalationPriceContractsCount{
			ExpectedAmountOfContractYearsForCalculation:     2,
			ExpectedAmountOfBasePeriodYearsForCalculation:   2,
			ExpectedAmountOfOptionPeriodYearsForCalculation: 0,
			ExpectedAmountOfAwardTermsForCalculation:        0,
		}, nil
	case BasePeriodYear3:
		return ExpectedEscalationPriceContractsCount{
			ExpectedAmountOfContractYearsForCalculation:     3,
			ExpectedAmountOfBasePeriodYearsForCalculation:   3,
			ExpectedAmountOfOptionPeriodYearsForCalculation: 0,
			ExpectedAmountOfAwardTermsForCalculation:        0,
		}, nil
	case OptionPeriod1:
		return ExpectedEscalationPriceContractsCount{
			ExpectedAmountOfContractYearsForCalculation:     4,
			ExpectedAmountOfBasePeriodYearsForCalculation:   3,
			ExpectedAmountOfOptionPeriodYearsForCalculation: 1,
			ExpectedAmountOfAwardTermsForCalculation:        0,
		}, nil
	case OptionPeriod2:
		return ExpectedEscalationPriceContractsCount{
			ExpectedAmountOfContractYearsForCalculation:     5,
			ExpectedAmountOfBasePeriodYearsForCalculation:   3,
			ExpectedAmountOfOptionPeriodYearsForCalculation: 2,
			ExpectedAmountOfAwardTermsForCalculation:        0,
		}, nil
	case OptionPeriod3:
		return ExpectedEscalationPriceContractsCount{
			ExpectedAmountOfContractYearsForCalculation:     6,
			ExpectedAmountOfBasePeriodYearsForCalculation:   3,
			ExpectedAmountOfOptionPeriodYearsForCalculation: 3,
			ExpectedAmountOfAwardTermsForCalculation:        0,
		}, nil
	case AwardTerm1:
		if hasOptionYear3 {
			return ExpectedEscalationPriceContractsCount{
				ExpectedAmountOfContractYearsForCalculation:     7,
				ExpectedAmountOfBasePeriodYearsForCalculation:   3,
				ExpectedAmountOfOptionPeriodYearsForCalculation: 3,
				ExpectedAmountOfAwardTermsForCalculation:        1,
			}, nil
		} else {
			return ExpectedEscalationPriceContractsCount{
				ExpectedAmountOfContractYearsForCalculation:     6,
				ExpectedAmountOfBasePeriodYearsForCalculation:   3,
				ExpectedAmountOfOptionPeriodYearsForCalculation: 2,
				ExpectedAmountOfAwardTermsForCalculation:        1,
			}, nil
		}
	case AwardTerm2:
		if hasOptionYear3 {
			return ExpectedEscalationPriceContractsCount{
				ExpectedAmountOfContractYearsForCalculation:     8,
				ExpectedAmountOfBasePeriodYearsForCalculation:   3,
				ExpectedAmountOfOptionPeriodYearsForCalculation: 3,
				ExpectedAmountOfAwardTermsForCalculation:        2,
			}, nil
		} else {
			return ExpectedEscalationPriceContractsCount{
				ExpectedAmountOfContractYearsForCalculation:     7,
				ExpectedAmountOfBasePeriodYearsForCalculation:   3,
				ExpectedAmountOfOptionPeriodYearsForCalculation: 2,
				ExpectedAmountOfAwardTermsForCalculation:        2,
			}, nil
		}
	default:
		err := apperror.NewInternalServerError(fmt.Sprintf("Unexpected contract year name %s.", contractYearName))
		return ExpectedEscalationPriceContractsCount{}, err
	}
}
