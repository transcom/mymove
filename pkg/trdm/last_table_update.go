package trdm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"

	"github.com/transcom/mymove/pkg/appcontext"
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

// This is the start of the cron job. When called, it will immediately trigger the TRDM flow. After that, it will trigger again every 24 hours.
func startLastTableUpdateCron(physicalName string, logger *zap.Logger, v *viper.Viper, appCtx appcontext.AppContext, provider AssumeRoleProvider, client HTTPClient) error {

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

		// Handle the response from lastTableUpdate
		switch lastTableUpdateResponse.StatusCode {
		case SuccessfulStatusCode:
			switch physicalName {
			case LineOfAccounting:
				loas, err := FetchLOARecordsByTime(appCtx, lastTableUpdateResponse.LastUpdate)
				if err != nil {
					logger.Error("fetching loa records by time", zap.Error(err))
					return
				}
				// Check if loas are out of date
				if len(loas) > 0 {
					// Since loas were returned, we are in fact out of date
					// Trigger Get TGET data and GetTable call
					err := GetTGETData(models.GetTableRequest{
						PhysicalName:                LineOfAccounting,
						ContentUpdatedSinceDateTime: lastTableUpdateResponse.LastUpdate,
						ReturnContent:               true,
					}, *service, appCtx)
					if err != nil {
						logger.Info("successfully retrieved latest line of accounting TGET data")
					} else {
						logger.Fatal("failed to retrieve latest line of accounting TGET data")
					}
					return
				}
			case TransportationAccountingCode:
				tacs, err := FetchTACRecordsByTime(appCtx, lastTableUpdateResponse.LastUpdate)
				if err != nil {
					logger.Error("fetching tac records by time", zap.Error(err))
					return
				}
				// Check if tacs are out of date
				if len(tacs) > 0 {
					// Since tacs were returned, we are in fact out of date
					// Trigger Get TGET data and GetTable call
					err := GetTGETData(models.GetTableRequest{
						PhysicalName:                TransportationAccountingCode,
						ContentUpdatedSinceDateTime: lastTableUpdateResponse.LastUpdate,
						ReturnContent:               true,
					}, *service, appCtx)
					if err != nil {
						logger.Fatal("failed to retrieve latest transportation accounting TGET data")
					} else {
						logger.Info("successfully retrieved latest transportation accounting TGET data")
					}
					return
				}
			default:
				logger.Error("unsupported table provided")
				return
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

// Fetch Transportation Accounting Codes from DB and return the list of records if the updated_at field is < the returned LastTableUpdate updated time.
// Ex:
//
//	LastTableUpdate : 2023-08-30 15:24:13.19931
//	updated_at: 2023-08-29 15:24:13.19931
//
// Because updated_at is before LastTableUpdate the DB will return records that match this case.
//
//	returns []models.TransportationAccountingCode, error
func FetchTACRecordsByTime(appcontext appcontext.AppContext, time time.Time) ([]models.TransportationAccountingCode, error) {
	var tacCodes []models.TransportationAccountingCode
	err := appcontext.DB().Select("*").Where("updated_at < $1", time).All(&tacCodes)

	if err != nil {
		return tacCodes, errors.Wrap(err, "Fetch line items query failed")
	}

	return tacCodes, nil
}

// Fetch Line Of Accounting records from DB and return the list of records if the updated_at field is < the returned LastTableUpdate updated time.
// Ex:
//
//	LastTableUpdate : 2023-08-30 15:24:13.19931
//	updated_at: 2023-08-29 15:24:13.19931
//
// Because updated_at is before LastTableUpdate the DB will return records that match this case.
//
//	returns []models.LineOfAccounting, error
func FetchLOARecordsByTime(appcontext appcontext.AppContext, time time.Time) ([]models.LineOfAccounting, error) {
	var loa []models.LineOfAccounting
	err := appcontext.DB().Select("*").Where("updated_at < $1", time).All(&loa)
	if err != nil {
		return loa, errors.Wrap(err, "Fetch line items query failed")
	}
	return loa, nil
}
