package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// EIAKeyFlag is the EIA Key Flag
	EIAKeyFlag string = "eia-key"
	// EIAURLFlag is the EIA URL Flag
	EIAURLFlag string = "eia-url"
)

// InitEIAFlags initializes EIA command line flags
func InitEIAFlags(flag *pflag.FlagSet) {
	flag.String(EIAURLFlag, "https://api.eia.gov/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W", "URL for Energy Information Administration (EIA) Open Data API")
	flag.String(EIAKeyFlag, "", "Key for Energy Information Administration (EIA) Open Data API")
}

// CheckEIA validates EIA command line flags
func CheckEIA(v *viper.Viper) error {
	eiaURL := v.GetString(EIAURLFlag)
	if eiaURL != "https://api.eia.gov/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W" {
		return fmt.Errorf("invalid EIA Open Data URL %s, expecting https://api.eia.gov/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W", eiaURL)
	}

	eiaKey := v.GetString(EIAKeyFlag)
	if len(eiaKey) != 32 {
		return fmt.Errorf("expected EIA Open Data API key to be 32 characters long; key is %d chars", len(eiaKey))
	}
	return nil
}
