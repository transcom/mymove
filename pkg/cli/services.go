package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// ServeAdminFlag is the admin service flag
	ServeAdminFlag string = "serve-admin"
	// ServeOrdersFlag is the orders service flag
	ServeOrdersFlag string = "serve-orders"
	// ServeAPIInternalFlag is the internal api service flag
	ServeAPIInternalFlag string = "serve-api-internal"
	// ServeGHCFlag is the ghc api service flag
	ServeGHCFlag string = "serve-api-ghc"
	// ServePrimeFlag is the prime api flag
	ServePrimeFlag string = "serve-api-prime"
	// ServeSupportFlag is the support api flag
	ServeSupportFlag string = "serve-api-support"
	// ServePrimeSimulatorFlag is the prime simulator api flag
	ServePrimeSimulatorFlag string = "serve-prime-simulator"
)

// InitServiceFlags initializes the service command line flags
func InitServiceFlags(flag *pflag.FlagSet) {
	flag.Bool(ServeAdminFlag, false, "Enable the Admin Service.")
	flag.Bool(ServeOrdersFlag, false, "Enable the Orders Service.")
	flag.Bool(ServeAPIInternalFlag, false, "Enable the Internal API Service.")
	flag.Bool(ServeGHCFlag, false, "Enable the GHC API Service.")
	flag.Bool(ServePrimeFlag, false, "Enable the Prime API Service.")
	flag.Bool(ServeSupportFlag, false, "Enable the Support Service.")
	flag.Bool(ServePrimeSimulatorFlag, false, "Enable the Prime Simulator Service.")
}

// CheckServices validates these lovely service flags
func CheckServices(v *viper.Viper) error {
	adminEnabled := v.GetBool(ServeAdminFlag)
	ordersEnabled := v.GetBool(ServeOrdersFlag)
	internalAPIEnabled := v.GetBool(ServeAPIInternalFlag)
	ghcAPIEnabled := v.GetBool(ServeGHCFlag)
	primeAPIEnabled := v.GetBool(ServePrimeFlag)
	primeSimulatorEnabled := v.GetBool(ServePrimeSimulatorFlag)

	// Oops none of the flags used
	if (!adminEnabled) &&
		(!ordersEnabled) &&
		(!internalAPIEnabled) &&
		(!ghcAPIEnabled) &&
		(!primeAPIEnabled) &&
		(!primeSimulatorEnabled) {
		return fmt.Errorf("no service was enabled")
	}

	// if Orders is enabled then the mutualTLSListener is needed too
	// if PRIME is enabled then the mutualTLSListener is needed too
	mutualTLSEnabled := v.GetBool(MutualTLSListenerFlag)
	currentEnvironment := v.GetString(EnvironmentFlag)
	devOrReviewEnvironment := currentEnvironment == EnvironmentDevelopment ||
		currentEnvironment == EnvironmentReview
	if !devOrReviewEnvironment {
		if ordersEnabled && !mutualTLSEnabled {
			return fmt.Errorf("for orders service to be enabled both %s and the %s flags must be in use", ServeOrdersFlag, MutualTLSListenerFlag)
		}
		if primeAPIEnabled && !mutualTLSEnabled {
			return fmt.Errorf("for prime service to be enabled both %s and the %s flags must be in use", ServePrimeFlag, MutualTLSListenerFlag)
		}
		if mutualTLSEnabled && !(ordersEnabled || primeAPIEnabled) {
			return fmt.Errorf("either orders or prime service must be enabled for mutualTSL to be enabled")
		}
	}

	if currentEnvironment == EnvironmentPrd && primeSimulatorEnabled {
		return fmt.Errorf("Prime Simulator cannot be enabled in production")
	}

	return nil
}
