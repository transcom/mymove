package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

type GblocAors struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	JppsoRegionID       uuid.UUID `json:"jppso_regions_id" db:"jppso_regions_id"`
	OconusRateAreaID    uuid.UUID `json:"oconus_rate_area_id" db:"oconus_rate_area_id"`
	DepartmentIndicator *string   `json:"department_indicator" db:"department_indicator"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (c GblocAors) TableName() string {
	return "gbloc_aors"
}

func FetchGblocAorsByJppsoCodeRateAreaDept(db *pop.Connection, jppsoRegionId uuid.UUID, oconusRateAreaId uuid.UUID, deptInd string) (*GblocAors, error) {
	var gblocAors GblocAors
	err := db.Q().
		InnerJoin("jppso_regions jr", "gbloc_aors.jppso_regions_id = jr.id").
		Where("gbloc_aors.oconus_rate_area_id = ?", oconusRateAreaId).
		Where("(gbloc_aors.department_indicator = ? or gbloc_aors.department_indicator is null)", deptInd).
		Where("jr.id = ?", jppsoRegionId).
		First(&gblocAors)
	if err != nil {
		return nil, err
	}
	return &gblocAors, nil
}
