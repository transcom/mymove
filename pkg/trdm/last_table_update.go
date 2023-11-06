package trdm

import (
	"crypto/tls"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
)

// const successfulStatusCode = "Successful"

// Date/time value is used in conjunction with the contentUpdatedSinceDateTime column in the getTable method.
type GetLastTableUpdater interface {
	GetLastTableUpdate(appCtx appcontext.AppContext, physicalName string) error
}

// FetchAllTACRecords queries and fetches all transportation_accounting_codes
func FetchAllTACRecords(appcontext appcontext.AppContext) ([]models.TransportationAccountingCode, error) {
	var tacCodes []models.TransportationAccountingCode
	query := `SELECT * FROM transportation_accounting_codes`
	err := appcontext.DB().RawQuery(query).All(&tacCodes)
	if err != nil {
		return tacCodes, errors.Wrap(err, "Fetch line items query failed")
	}

	return tacCodes, nil
}

func StartLastTableUpdateCron() error {
	cron := cron.New()

	cronTask := func() {
		// TODO:
	}

	res, err := cron.AddFunc("@every 24h00m00s", cronTask)
	if err != nil {
		return fmt.Errorf("error adding cron task: %s, %v", err.Error(), res)
	}
	cron.Start()
	return nil
}

func LastTableUpdate(v *viper.Viper, _ *tls.Config) error {
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)

	// TODO: Turn back on when replacing with rest
	// DB connection
	// dbConnection, err := cli.InitDatabase(v, logger)
	// if err != nil {
	// 	return err
	// }

	// appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	// tr := &http.Transport{TLSClientConfig: tlsConfig}
	// httpClient := &http.Client{Transport: tr, Timeout: time.Duration(30) * time.Second}

	// TODO: Replace with api gateway call

	// TODO: Replace with REST
	// getLastTableUpdateTACErr := NewTRDMGetLastTableUpdate(transportationAccountingCode, tacBodyID, certificate, rsaKey, soapClient).GetLastTableUpdate(appCtx, transportationAccountingCode)
	// getLastTableUpdateLOAErr := NewTRDMGetLastTableUpdate(lineOfAccounting, loaBodyID, certificate, rsaKey, soapClient).GetLastTableUpdate(appCtx, lineOfAccounting)
	// if getLastTableUpdateLOAErr != nil {
	// 	return getLastTableUpdateLOAErr
	// }
	// if getLastTableUpdateTACErr != nil {
	// 	return getLastTableUpdateTACErr
	// }

	// TODO: Replace with REST
	// cronErrTAC := StartLastTableUpdateCron(appCtx, certificate, publicPem, rsaKey, transportationAccountingCode, soapClient)
	// cronErrLOA := StartLastTableUpdateCron(appCtx, certificate, publicPem, rsaKey, lineOfAccounting, soapClient)

	// if cronErrLOA != nil {
	// 	return cronErrLOA
	// }
	// if cronErrTAC != nil {
	// 	return cronErrTAC
	// }
	return nil
}
