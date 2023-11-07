package trdm

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
)

const (
	// TrdmIamRoleFlag is the TRDM IAM Role flag
	TrdmIamRoleFlag string = "trdm-iam-role"
	// TrdmRegionFlag is the TRDM Region flag
	TrdmRegionFlag string = "trdm-region"
	// GatewayURLFlag is the TRDM API Gateway URL flag
	GatewayURLFlag string = "trdm-api-gateway-url"
	// Success status code
	successfulStatusCode string = "Successful"
	// Failure status code
	failureStatusCode string = "Failure"
)

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

func StartLastTableUpdateCron(physicalName string, logger *zap.Logger, v *viper.Viper, tlsConfig *tls.Config, appCtx appcontext.AppContext) error {

	cron := cron.New()

	// Cron tasks do not return errors, only log
	cronTask := func() {
		roleFlag := v.GetString(TrdmIamRoleFlag)
		regionFlag := v.GetString(TrdmRegionFlag)
		gatewayURL := v.GetString(GatewayURLFlag)

		// Obtain creds for signing
		creds, err := retrieveCredentials(regionFlag, roleFlag, logger)
		if err != nil {
			logger.Error("retrieving aws creds", zap.Error(err))
			return
		}

		// Initialize the request model with physicalName
		request := models.LastTableUpdateRequest{
			PhysicalName: physicalName, // assuming physicalName is available in this scope
		}

		// Setup response model
		lastTableUpdateResponse := models.LastTableUpdateResponse{}

		// Setup client
		tr := &http.Transport{TLSClientConfig: tlsConfig}
		httpClient := &http.Client{Transport: tr, Timeout: time.Duration(30) * time.Second}

		// Create gateway service
		service := NewGatewayService(httpClient, logger, regionFlag, roleFlag, gatewayURL, &creds)

		// Fire off to retrieve the latest table update, compare that to our own internal latest update records,
		// and then call getTable if there is new data found. The getTable call will happen inside of this chain
		httpResp, err := service.gatewayLastTableUpdate(request)
		if err != nil {
			logger.Error("gateway last table update", zap.Error(err))
			return
		}
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			logger.Error("could not read lastTableUpdate response body", zap.Error(err))
			return
		}
		defer httpResp.Body.Close()
		err = json.Unmarshal(body, &lastTableUpdateResponse)
		if err != nil {
			logger.Error("could not unmarshal body into lastTableUpdateResponse", zap.Error(err))
			return
		}

		switch lastTableUpdateResponse.StatusCode {
		case successfulStatusCode:
			switch physicalName {
			case lineOfAccounting:
				loas, err := FetchLOARecordsByTime(appCtx, lastTableUpdateResponse.LastUpdate)
				if err != nil {
					logger.Error("fetching loa records by time", zap.Error(err))
					return
				}
				// Check if loas are out of date
				if len(loas) > 0 {
					// Since loas were returned, we are in fact out of date
					// TODO: GetTable
				}
			}
		case transportationAccountingCode:
			tacs, err := FetchTACRecordsByTime(appCtx, lastTableUpdateResponse.LastUpdate)
			if err != nil {
				logger.Error("fetching tac records by time", zap.Error(err))
				return
			}
			// Check if tacs are out of date
			if len(tacs) > 0 {
				// Since tacs were returned, we are in fact out of date
				// TODO: GetTable
			}
		case failureStatusCode:
			logger.Error("trdm api gateway request failed, please inspect the trdm gateway logs")
			return
		default:
			logger.Error("unexpected api gateway request failure response, please inspect the trdm gateway logs")
			return
		}
	}

	// Run the task immediately
	cronTask()

	// Schedule the task to run every 24 hours
	res, err := cron.AddFunc("@every 24h00m00s", cronTask)
	if err != nil {
		return fmt.Errorf("error adding cron task: %s, %v", err.Error(), res)
	}

	cron.Start()
	return nil
}

func LastTableUpdate(v *viper.Viper, tlsConfig *tls.Config, appCtx appcontext.AppContext) error {
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)

	// TODO: Turn back on when implementing getTable
	// DB connection
	// dbConnection, err := cli.InitDatabase(v, logger)
	// if err != nil {
	// 	return err
	// }

	// These are likely to never err. Remember, errors are logged not returned in cron
	getLastTableUpdateTACErr := StartLastTableUpdateCron(transportationAccountingCode, logger, v, tlsConfig, appCtx)
	getLastTableUpdateLOAErr := StartLastTableUpdateCron(lineOfAccounting, logger, v, tlsConfig, appCtx)
	if getLastTableUpdateLOAErr != nil {
		return getLastTableUpdateLOAErr
	}
	if getLastTableUpdateTACErr != nil {
		return getLastTableUpdateTACErr
	}

	return nil
}
