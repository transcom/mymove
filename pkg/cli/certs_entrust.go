package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// EntrustL1KCertFlag is the Entrust L1K Cert Flag
	EntrustL1KCertFlag string = "entrust-l1k-cert"
	// EntrustG2CertFlag is the Entrust G2 Cert Flag
	EntrustG2CertFlag string = "entrust-g2-cert"
)

// InitEntrustCertFlags initializes the Certificate Flags
func InitEntrustCertFlags(flag *pflag.FlagSet) {
	flag.String(EntrustL1KCertFlag, "", "The Entrust L1K certificate.")
	flag.String(EntrustG2CertFlag, "", "The Entrust G2 certificate.")
}

// CheckEntrustCert validates Cert command line flags
func CheckEntrustCert(v *viper.Viper) error {

	dbEnv := v.GetString(DbEnvFlag)
	isDevOrTest := dbEnv == "development" || dbEnv == "test"
	if devlocalCAPath := v.GetString(DevlocalCAFlag); isDevOrTest && devlocalCAPath == "" {
		return errors.Errorf("No devlocal CA path defined")
	}

	entrustG2CertString := v.GetString(EntrustG2CertFlag)
	if len(entrustG2CertString) == 0 {
		return errors.Errorf("%s is missing", EntrustG2CertFlag)
	}

	entrustL1KCertString := v.GetString(EntrustL1KCertFlag)
	if len(entrustL1KCertString) == 0 {
		return errors.Errorf("%s is missing", EntrustL1KCertFlag)
	}

	return nil
}
