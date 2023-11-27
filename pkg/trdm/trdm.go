package trdm

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func BeginTGETFlow(v *viper.Viper, appCtx appcontext.AppContext, provider AssumeRoleProvider, client HTTPClient) error {
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)

	// These are likely to never err. Remember, errors are logged not returned in cron
	getLastTableUpdateTACErr := startLastTableUpdateCron(TransportationAccountingCode, logger, v, appCtx, provider, client)
	getLastTableUpdateLOAErr := startLastTableUpdateCron(LineOfAccounting, logger, v, appCtx, provider, client)
	if getLastTableUpdateLOAErr != nil {
		return getLastTableUpdateLOAErr
	}
	if getLastTableUpdateTACErr != nil {
		return getLastTableUpdateTACErr
	}

	return nil
}
