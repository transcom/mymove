package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// DimensionType determines what type of dimension for a service item
type DimensionType string

const (
	// DimensionTypeItem describes dimensions for an ITEM type
	DimensionTypeItem DimensionType = "ITEM"
	// DimensionTypeCrate  describes dimensions for a CRATE type
	DimensionTypeCrate DimensionType = "CRATE"
)

// MTOServiceItemDimension is an object representing dimensions for a service item.
type MTOServiceItemDimension struct {
	ID               uuid.UUID             `db:"id"`
	MTOServiceItem   MTOServiceItem        `belongs_to:"mto_service_items" fk_id:"mto_service_item_id"`
	MTOServiceItemID uuid.UUID             `db:"mto_service_item_id"`
	Type             DimensionType         `db:"type"`
	Length           unit.ThousandthInches `db:"length_thousandth_inches"`
	Height           unit.ThousandthInches `db:"height_thousandth_inches"`
	Width            unit.ThousandthInches `db:"width_thousandth_inches"`
	CreatedAt        time.Time             `db:"created_at"`
	UpdatedAt        time.Time             `db:"updated_at"`
}

// MTOServiceItemDimensions is a slice containing MTOServiceItemDimension.
type MTOServiceItemDimensions []MTOServiceItemDimension

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOServiceItemDimension) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MTOServiceItemID, Name: "MTOServiceItemID"})
	vs = append(vs, &validators.StringInclusion{Field: string(m.Type), Name: "Type", List: []string{
		string(DimensionTypeItem),
		string(DimensionTypeCrate),
	}})
	vs = append(vs, &validators.IntIsGreaterThan{Field: int(m.Length), Name: "Length", Compared: -1})
	vs = append(vs, &validators.IntIsGreaterThan{Field: int(m.Width), Name: "Width", Compared: -1})
	vs = append(vs, &validators.IntIsGreaterThan{Field: int(m.Height), Name: "Height", Compared: -1})

	return validate.Validate(vs...), nil
}

// TableName overrides the table name used by Pop.
func (m MTOServiceItemDimension) TableName() string {
	return "mto_service_item_dimensions"
}

// Volume calculates Length x Height x Width
func (m *MTOServiceItemDimension) Volume() unit.CubicThousandthInch {
	volume := m.Length * m.Width * m.Height
	return unit.CubicThousandthInch(volume)
}
