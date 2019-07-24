package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// ServeAdminFlag is the admin service flag
	ServeAdminFlag string = "serve-admin"
	// ServeSDDCFlag is the sddc service flag
	ServeSDDCFlag string = "serve-sddc"
	// ServeOrdersFlag is the orders service flag
	ServeOrdersFlag string = "serve-orders"
	// ServeDPSFlag is the DPS service flag
	ServeDPSFlag string = "serve-dps"
	// ServeAPIInternalFlag is the internal api service flag
	ServeAPIInternalFlag string = "serve-api-internal"
	// ServeAPIExternalFlag is the external api service flag
	ServeAPIExternalFlag string = "serve-api-external"
)

// InitServiceFlags initializes the service command line flags
func InitServiceFlags(flag *pflag.FlagSet) {
	flag.Bool(ServeAdminFlag, false, "Enable the Admin Service.")
	flag.Bool(ServeSDDCFlag, false, "Enable the SDDC Service.")
	flag.Bool(ServeOrdersFlag, false, "Enable the Orders Service.")
	flag.Bool(ServeDPSFlag, false, "Enable the DPS Service.")
	flag.Bool(ServeAPIInternalFlag, false, "Enable the Internal API Service.")
	flag.Bool(ServeAPIExternalFlag, false, "Enable the External API Service.")
}

// CheckServices validates these lovely service flags
func CheckServices(v *viper.Viper) error {
	adminEnabled := v.GetBool(ServeAdminFlag)
	sddcEnabled := v.GetBool(ServeSDDCFlag)
	ordersEnabled := v.GetBool(ServeOrdersFlag)
	dpsEnabled := v.GetBool(ServeDPSFlag)
	internalAPIEnabled := v.GetBool(ServeAPIInternalFlag)
	externalAPIEnabled := v.GetBool(ServeAPIExternalFlag)

	// Oops none of the flags used
	if (!adminEnabled) &&
		(!sddcEnabled) &&
		(!ordersEnabled) &&
		(!dpsEnabled) &&
		(!internalAPIEnabled) &&
		(!externalAPIEnabled) {
		return errors.New("no service was enabled")
	}

	// if Orders is enabled then the mutualTLSListener is needed too
	mutualTLSEnabled := v.GetBool(MutualTLSListenerFlag)
	fmt.Println(mutualTLSEnabled)
	if v.GetString(EnvironmentFlag) != EnvironmentDevelopment {
		if ordersEnabled && !mutualTLSEnabled ||
			!ordersEnabled && mutualTLSEnabled {
			return errors.New("for orders service to be enabled both it and the MutualTLSListener flags must be in use")
		}
	}

	return nil
}