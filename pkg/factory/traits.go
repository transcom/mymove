package factory

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
)

// GetTraitNavy is a sample GetTraitFunc
func GetTraitNavy() []Customization {
	navy := models.AffiliationNAVY
	return []Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &navy,
			},
		},
	}
}

// GetTraitArmy is a sample GetTraitFunc
func GetTraitArmy() []Customization {
	army := models.AffiliationARMY
	return []Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &army,
			},
		},
	}
}

// GetTraitAddress2 is a sample GetTraitFunc
func GetTraitAddress2() []Customization {
	return []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 Any Avenue",
				StreetAddress2: swag.String("P.O. Box 9876"),
				StreetAddress3: swag.String("c/o Some Person"),
				City:           "Fairfield",
				State:          "CA",
				PostalCode:     "94535",
				Country:        swag.String("US"),
			},
		},
	}
}

// GetTraitAddress3 is a sample GetTraitFunc
func GetTraitAddress3() []Customization {
	return []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: swag.String("P.O. Box 1234"),
				StreetAddress3: swag.String("c/o Another Person"),
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50309",
				Country:        swag.String("US"),
			},
		},
	}
}

// GetTraitAddress4 is a sample GetTraitFunc
func GetTraitAddress4() []Customization {
	return []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 Over There Avenue",
				StreetAddress2: swag.String("P.O. Box 1234"),
				StreetAddress3: swag.String("c/o Another Person"),
				City:           "Houston",
				State:          "TX",
				PostalCode:     "77083",
				Country:        swag.String("US"),
			},
		},
	}
}
