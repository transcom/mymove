package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// Maintenance Flag
	MaintenanceFlag string = "maintenance_flag"
)

func InitMaintenanceFlags(flag *pflag.FlagSet) {
	flag.Bool(MaintenanceFlag, false, "Flag for tracking app maintenance.")
}

func CheckMaintenance(v *viper.Viper) error {
	return nil
}
