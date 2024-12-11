package entitlements

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type weightAllotmentFetcher struct {
}

// NewWeightAllotmentFetcher returns a new weight allotment fetcher
func NewWeightAllotmentFetcher() services.WeightAllotmentFetcher {
	return &weightAllotmentFetcher{}
}

func (waf *weightAllotmentFetcher) GetWeightAllotment(appCtx appcontext.AppContext, grade string) (*models.HHGAllowance, error) {
	var hhgAllowance models.HHGAllowance
	err := appCtx.DB().
		RawQuery(`
          SELECT hhg_allowances.*
          FROM hhg_allowances
          INNER JOIN pay_grades ON hhg_allowances.pay_grade_id = pay_grades.id
          WHERE pay_grades.grade = $1
          LIMIT 1
        `, grade).
		First(&hhgAllowance)
	if err != nil {
		return nil, apperror.NewQueryError("HHGAllowance", err, fmt.Sprintf("Error retrieving HHG allowance for grade: %s", grade))
	}

	return &hhgAllowance, nil
}
