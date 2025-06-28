package entitlements

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type weightAllotmentFetcher struct {
}

// NewWeightAllotmentFetcher returns a new weight allotment fetcher
func NewWeightAllotmentFetcher() services.WeightAllotmentFetcher {
	return &weightAllotmentFetcher{}
}

func (waf *weightAllotmentFetcher) GetWeightAllotment(appCtx appcontext.AppContext, grade string, ordersType internalmessages.OrdersType) (models.WeightAllotment, error) {
	// Check order allotment first
	if ordersType == internalmessages.OrdersTypeSTUDENTTRAVEL { // currently only applies to student travel order that limits overall authorized weight
		entitlement, err := waf.GetWeightAllotmentByOrdersType(appCtx, ordersType)
		if err != nil {
			return models.WeightAllotment{}, err
		}
		return entitlement, nil
	}

	// Continue if the orders type is not student travel
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
		return models.WeightAllotment{}, apperror.NewQueryError("HHGAllowance", err, fmt.Sprintf("Error retrieving HHG allowance for grade: %s", grade))
	}

	maxGunSafeWeightAllowance, err := models.GetMaxGunSafeAllowance(appCtx)
	if err != nil {
		return models.WeightAllotment{}, err
	}

	// Convert HHGAllowance to WeightAllotment
	weightAllotment := models.WeightAllotment{
		TotalWeightSelf:               hhgAllowance.TotalWeightSelf,
		TotalWeightSelfPlusDependents: hhgAllowance.TotalWeightSelfPlusDependents,
		ProGearWeight:                 hhgAllowance.ProGearWeight,
		ProGearWeightSpouse:           hhgAllowance.ProGearWeightSpouse,
		GunSafeWeight:                 maxGunSafeWeightAllowance,
	}

	return weightAllotment, nil
}

var ordersTypeToAllotmentAppParamName = map[internalmessages.OrdersType]string{
	internalmessages.OrdersTypeSTUDENTTRAVEL: "studentTravelHhgAllowance",
}

// Helper func to enforce strict unmarshal of application param values into a given  interface
func strictUnmarshal(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	// Fail on unknown fields
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func (waf *weightAllotmentFetcher) GetWeightAllotmentByOrdersType(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType) (models.WeightAllotment, error) {
	if paramName, ok := ordersTypeToAllotmentAppParamName[ordersType]; ok {
		// We currently store orders allotment overrides as an application parameter
		// as it is a current one-off use case introduced by E-06189
		var jsonData json.RawMessage
		err := appCtx.DB().RawQuery(`
		SELECT parameter_json
		FROM application_parameters
		WHERE parameter_name = $1
		`, paramName).First(&jsonData)

		if err != nil {
			return models.WeightAllotment{}, fmt.Errorf("failed to fetch weight allotment for orders type %s: %w", ordersType, err)
		}

		// Convert the JSON data to the WeightAllotment struct
		var weightAllotment models.WeightAllotment
		err = strictUnmarshal(jsonData, &weightAllotment)
		if err != nil {
			return models.WeightAllotment{}, fmt.Errorf("failed to parse weight allotment JSON for orders type %s: %w", ordersType, err)
		}

		return weightAllotment, nil
	}
	return models.WeightAllotment{}, fmt.Errorf("no entitlement found for orders type %s", ordersType)
}

func (waf *weightAllotmentFetcher) GetAllWeightAllotments(appCtx appcontext.AppContext) (map[internalmessages.OrderPayGrade]models.WeightAllotment, error) {
	var hhgAllowances models.HHGAllowances
	err := appCtx.DB().
		Eager("PayGrade").
		All(&hhgAllowances)
	if err != nil {
		return nil, apperror.NewQueryError("HHGAllowances", err, "Error retrieving all HHG allowances")
	}

	maxGunSafeWeightAllowance, err := models.GetMaxGunSafeAllowance(appCtx)
	if err != nil {
		return nil, err
	}

	weightAllotments := make(map[internalmessages.OrderPayGrade]models.WeightAllotment)

	for _, hhgAllowance := range hhgAllowances {
		// Convert HHGAllowance to WeightAllotment
		weightAllotment := models.WeightAllotment{
			TotalWeightSelf:               hhgAllowance.TotalWeightSelf,
			TotalWeightSelfPlusDependents: hhgAllowance.TotalWeightSelfPlusDependents,
			ProGearWeight:                 hhgAllowance.ProGearWeight,
			ProGearWeightSpouse:           hhgAllowance.ProGearWeightSpouse,
			GunSafeWeight:                 maxGunSafeWeightAllowance,
		}

		grade := internalmessages.OrderPayGrade(hhgAllowance.PayGrade.Grade)
		weightAllotments[grade] = weightAllotment
	}
	return weightAllotments, nil
}
