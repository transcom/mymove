package factory

import (
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
