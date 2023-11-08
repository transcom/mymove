package trdm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	// TrdmGatewayRegionFlag is the TRDM Region flag
	TrdmGatewayRegionFlag string = "trdm-gateway-region"
	// GatewayURLFlag is the TRDM API Gateway URL flag
	GatewayURLFlag string = "trdm-api-gateway-url"
	// Success status code
	SuccessfulStatusCode string = "Successful"
	// Failure status code
	FailureStatusCode string = "Failure"
)

// Custom assume role provider so we can inject tests.
// See aws go sdk v2 STS assume role provider, that's what we're
// mimicking
type AssumeRoleProvider interface {
	Retrieve(ctx context.Context) (aws.Credentials, error)
}

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

func StartLastTableUpdateCron(physicalName string, logger *zap.Logger, v *viper.Viper, appCtx appcontext.AppContext, provider AssumeRoleProvider, client HTTPClient) error {

	cron := cron.New()

	// Cron tasks do not return errors, only log
	cronTask := func() {
		trdmIamRole := v.GetString(TrdmIamRoleFlag)
		region := v.GetString(TrdmGatewayRegionFlag)
		gatewayURL := v.GetString(GatewayURLFlag)

		// Initialize the request model with physicalName
		request := models.LastTableUpdateRequest{
			PhysicalName: physicalName,
		}

		// Setup response model
		lastTableUpdateResponse := models.LastTableUpdateResponse{}

		// Create gateway service
		service := NewGatewayService(client, logger, region, trdmIamRole, gatewayURL, provider)

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
		case SuccessfulStatusCode:
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
					return
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
					return
				}
			}
		case FailureStatusCode:
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

func LastTableUpdate(v *viper.Viper, appCtx appcontext.AppContext, provider AssumeRoleProvider, client HTTPClient) error {
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
	getLastTableUpdateTACErr := StartLastTableUpdateCron(transportationAccountingCode, logger, v, appCtx, provider, client)
	getLastTableUpdateLOAErr := StartLastTableUpdateCron(lineOfAccounting, logger, v, appCtx, provider, client)
	if getLastTableUpdateLOAErr != nil {
		return getLastTableUpdateLOAErr
	}
	if getLastTableUpdateTACErr != nil {
		return getLastTableUpdateTACErr
	}

	return nil
}
