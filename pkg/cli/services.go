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
	// ServeGHCFlag is the ghc api service flag
	ServeGHCFlag string = "serve-api-ghc"
	// ServePrimeFlag is the prime api flag
	ServePrimeFlag string = "serve-api-prime"
	// ServeSupportFlag is the prime api flag
	ServeSupportFlag string = "serve-api-support"
)

// InitServiceFlags initializes the service command line flags
func InitServiceFlags(flag *pflag.FlagSet) {
	flag.Bool(ServeAdminFlag, false, "Enable the Admin Service.")
	flag.Bool(ServeSDDCFlag, false, "Enable the SDDC Service.")
	flag.Bool(ServeOrdersFlag, false, "Enable the Orders Service.")
	flag.Bool(ServeDPSFlag, false, "Enable the DPS Service.")
	flag.Bool(ServeAPIInternalFlag, false, "Enable the Internal API Service.")
	flag.Bool(ServeGHCFlag, false, "Enable the GHC API Service.")
	flag.Bool(ServePrimeFlag, false, "Enable the Prime API Service.")
	flag.Bool(ServeSupportFlag, false, "Enable the Support Service.")
}

// CheckServices validates these lovely service flags
func CheckServices(v *viper.Viper) error {
	adminEnabled := v.GetBool(ServeAdminFlag)
	sddcEnabled := v.GetBool(ServeSDDCFlag)
	ordersEnabled := v.GetBool(ServeOrdersFlag)
	dpsEnabled := v.GetBool(ServeDPSFlag)
	internalAPIEnabled := v.GetBool(ServeAPIInternalFlag)
	ghcAPIEnabled := v.GetBool(ServeGHCFlag)
	primeAPIEnabled := v.GetBool(ServePrimeFlag)

	// Oops none of the flags used
	if (!adminEnabled) &&
		(!sddcEnabled) &&
		(!ordersEnabled) &&
		(!dpsEnabled) &&
		(!internalAPIEnabled) &&
		(!ghcAPIEnabled) &&
		(!primeAPIEnabled) {
		return errors.New("no service was enabled")
	}

	// if DPS is enabled then the mutualTLSListener is needed too
	// if Orders is enabled then the mutualTLSListener is needed too
	// if PRIME is enabled then the mutualTLSListener is needed too
	mutualTLSEnabled := v.GetBool(MutualTLSListenerFlag)
	currentEnvironment := v.GetString(EnvironmentFlag)
	devOrReviewEnvironment := currentEnvironment == EnvironmentDevelopment ||
		currentEnvironment == EnvironmentReview
	if !devOrReviewEnvironment {
		if dpsEnabled && !mutualTLSEnabled {
			return errors.New(fmt.Sprintf("for dps service to be enabled both %s and the %s flags must be in use", ServeDPSFlag, MutualTLSListenerFlag))
		}
		if ordersEnabled && !mutualTLSEnabled {
			return errors.New(fmt.Sprintf("for orders service to be enabled both %s and the %s flags must be in use", ServeOrdersFlag, MutualTLSListenerFlag))
		}
		if primeAPIEnabled && !mutualTLSEnabled {
			return errors.New(fmt.Sprintf("for prime service to be enabled both %s and the %s flags must be in use", ServePrimeFlag, MutualTLSListenerFlag))
		}
		if mutualTLSEnabled && !(dpsEnabled || ordersEnabled || primeAPIEnabled) {
			return errors.New("either dps, orders or prime service must be enabled for mutualTSL to be enabled")
		}
	}

	return nil
}
