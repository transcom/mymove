package featureflag

const (
	// Checks if the mobile home FF is enabled
	DomesticMobileHome string = "mobile_home"

	// Toggles service items on/off completely for mobile home shipments
	DomesticMobileHomeDOPEnabled       string = "domestic_mobile_home_origin_price_enabled"
	DomesticMobileHomeDDPEnabled       string = "domestic_mobile_home_destination_price_enabled"
	DomesticMobileHomePackingEnabled   string = "domestic_mobile_home_packing_enabled"
	DomesticMobileHomeUnpackingEnabled string = "domestic_mobile_home_unpacking_enabled"

	// Toggles whether or not the DMHF is applied to these service items for Mobile Home shipments (if they are not toggled off by the above flags)
	DomesticMobileHomeDOPFactor       string = "domestic_mobile_home_factor_origin_price"
	DomesticMobileHomeDDPFactor       string = "domestic_mobile_home_factor_destination_price"
	DomesticMobileHomePackingFactor   string = "domestic_mobile_home_factor_packing"
	DomesticMobileHomeUnpackingFactor string = "domestic_mobile_home_factor_unpacking"
)
