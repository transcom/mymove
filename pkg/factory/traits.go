package factory

import "github.com/transcom/mymove/pkg/models"

// GetTraitActiveUser returns a customization to enable active on a user
func GetTraitActiveUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
			},
			Type: &User,
		},
	}
}

// GetTraitNavy is a sample GetTraitFunc
func GetTraitNavy() []Customization {
	navy := models.AffiliationNAVY
	return []Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &navy,
			},
			Type: &ServiceMember,
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
			Type: &ServiceMember,
		},
	}
}
