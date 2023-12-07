package trdm

import (
	"context"
	"database/sql"
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
				loaDataOutOfDate, ourLastUpdate, caseErr := TGETLOADataOutOfDate(appCtx, lastTableUpdateResponse.LastUpdate)
				if caseErr != nil {
					logger.Error("fetching loa records by time", zap.Error(err))
					return
				}
				if ourLastUpdate == nil {
					logger.Fatal("our last update appears to be nil and no errors were returned")
					return
				}
				// Check if loas are out of date
				if loaDataOutOfDate {
					// Trigger Get TGET data and GetTable call
					caseErr := GetTGETData(models.GetTableRequest{
						PhysicalName:                LineOfAccounting,
						ContentUpdatedSinceDateTime: *ourLastUpdate,
						ReturnContent:               true,
					}, lastTableUpdateResponse.LastUpdate, *service, appCtx, logger)
					if caseErr != nil {
						logger.Fatal("failed to retrieve latest line of accounting TGET data", zap.Error(err))
					} else {
						logger.Info("successfully retrieved latest line of accounting TGET data")
					}
					return
				}
			case TransportationAccountingCode:
				tacDataOutOfDate, ourLastUpdate, caseErr := TGETTACDataOutOfDate(appCtx, lastTableUpdateResponse.LastUpdate)
				if caseErr != nil {
					logger.Error("fetching tac records by time", zap.Error(err))
					return
				}
				if ourLastUpdate == nil {
					logger.Fatal("our last update appears to be nil and no errors were returned")
					return
				}
				// Check if tacs are out of date
				if tacDataOutOfDate {
					// Trigger Get TGET data and GetTable call
					caseErr := GetTGETData(models.GetTableRequest{
						PhysicalName:                TransportationAccountingCode,
						ContentUpdatedSinceDateTime: *ourLastUpdate,
						ReturnContent:               true,
					}, lastTableUpdateResponse.LastUpdate, *service, appCtx, logger)
					if caseErr != nil {
						logger.Fatal("failed to retrieve latest transportation accounting TGET data", zap.Error(err))
					} else {
						logger.Info("successfully retrieved latest transportation accounting TGET data")
					}
					return
				}
			default:
				logger.Error("unsupported table provided", zap.String("responseBody", string(body)))
				return
			}
		case FailureStatusCode:
			logger.Error("trdm api gateway request failed, please inspect the trdm gateway logs", zap.Error(err))
			return
		default:
			logger.Error("unexpected api gateway request failure response, please inspect the trdm gateway logs", zap.Error(err))
			return
		}
	}

	// Schedule the task to run every day at midnight
	res, err := cron.AddFunc("0 0 * * *", cronTask)
	if err != nil {
		return fmt.Errorf("error adding cron task: %s, %v", err.Error(), res)
	}

	cron.Start()
	return nil
}

// Fetching all TAC records by time to check if our TAC data is out of date is incredibly inefficient. So instead we will check if a record
// exists.
// Ex:
//
//	LastTableUpdate : 2023-08-30 15:24:13.19931
//	updated_at: 2023-08-29 15:24:13.19931
//
// Because updated_at is before LastTableUpdate the DB will return true because this means out TGET data is out of date
//
//	returns bool, latestUpdateTime, error
func TGETTACDataOutOfDate(appcontext appcontext.AppContext, timeToCheck time.Time) (bool, *time.Time, error) {
	var tac models.TransportationAccountingCode

	// Get the most recent TAC record
	err := appcontext.DB().
		Order("updated_at DESC").
		First(&tac)

	if err != nil {
		// Check if it just so happens that our DB is empty (Such as if we're in a brand new environment)
		if err == sql.ErrNoRows {
			tenYearsAgo := time.Now().AddDate(-10, 0, 0)
			// Return TGET data out of date and gather 10 years of TGET data
			return true, &tenYearsAgo, nil
		}
		// Else, our error is not because the DB is empty. It's because the query failed.
		return false, nil, errors.Wrap(err, "db query for TAC failed")
	}

	// Compare the latest update time with the provided time. If our latest TAC entry has
	// an updated at value less than the provided time, that means our TGET data is out of date.
	// Otherwise, if the update at value is equal to or newer than the provided time then our data is
	// up to date.
	isOutOfDate := tac.UpdatedAt.Before(timeToCheck)

	return isOutOfDate, &tac.UpdatedAt, nil
}

// Fetching all LOA records by time to check if our LOA data is out of date is incredibly inefficient. So instead we will check if a record
// exists.
// Ex:
//
//	LastTableUpdate : 2023-08-30 15:24:13.19931
//	updated_at: 2023-08-29 15:24:13.19931
//
// Because updated_at is before LastTableUpdate the DB will return true because this means out TGET data is out of date
//
//	returns bool, latestUpdatedAt, error
func TGETLOADataOutOfDate(appcontext appcontext.AppContext, timeToCheck time.Time) (bool, *time.Time, error) {
	var loa models.LineOfAccounting

	// Get the most recent LOA record
	err := appcontext.DB().
		Order("updated_at DESC").
		First(&loa)

	if err != nil {
		// Check if it just so happens that our DB is empty (Such as if we're in a brand new environment)
		if err == sql.ErrNoRows {
			tenYearsAgo := time.Now().AddDate(-10, 0, 0)
			// Return TGET data out of date and gather 10 years of TGET data
			return true, &tenYearsAgo, nil
		}
		// Else, our error is not because the DB is empty. It's because the query failed.
		return false, nil, errors.Wrap(err, "db query for LOA failed")
	}

	// Compare the latest update time with the provided time. If our latest LOA entry has
	// an updated at value less than the provided time, that means our TGET data is out of date.
	// Otherwise, if the update at value is equal to or newer than the provided time then our data is
	// up to date.
	isOutOfDate := loa.UpdatedAt.Before(timeToCheck)

	return isOutOfDate, &loa.UpdatedAt, nil
}

// This type should be used as an array
type MissingWeek struct {
	// Matching dates should not occur because that would mean our data is up to date and the function using this type should never be called
	StartOfWeek time.Time // Example: AUG 01 OR AUG 05 OR AUG 01 (Read these top down)
	EndOfWeek   time.Time // Example: AUG 07 OR AUG 07 OR NOV 01 (Read these top down)
}

// This func won't inherently check from our DB but it will receive the latest update from a table within that DB to compare to TRDM
// It will allow us to split a request that would be 1 request of 1 month+ of data into 4 requests at 1 week each
func FetchWeeksOfMissingTime(ourLastUpdate time.Time, trdmLastUpdate time.Time) ([]MissingWeek, error) {
	if trdmLastUpdate.Before(ourLastUpdate) {
		return nil, errors.New("the provided parameters are out of order")
	}
	// Matching dates should not occur because that would mean our data is up to date and the function using this type should never be called
	var missingWeeks []MissingWeek // Individual start and end dates of each week to be used in the TRDM filter so we can grab 1 week at a time

	// Create a startOfWeek for each loop iteration.
	// If that start of week is not after the last update in TRDM then there are still weeks or days in between ourLastUpdate and trdmLastUpdate
	// Then set the new startOfWeek to the end of that week and run the loop again.
	// This loop will run every time until the startOfWeek finally matches trdmLastUpdate, meaning we have found all missing weeks
	// This loop can be modified to be startOfWeek.Before(trdmLastUpdate) to specifically find only the weeks, but we want individual days too
	for startOfWeek := ourLastUpdate; !startOfWeek.After(trdmLastUpdate); startOfWeek = startOfWeek.AddDate(0, 0, 7) {
		endOfWeek := startOfWeek.AddDate(0, 0, 6) // Add 6 days because it's already a day. 1 + 6 = 7, and 7 days are in a week
		// Set endOfWeek to the last second of the last day so it is truly the end of the week
		endOfWeek = endOfWeek.Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond)

		if endOfWeek.After(trdmLastUpdate) {
			// If it is the last week, set to the last second of trdmLastUpdate instead of trying to pull data that doesn't exist yet
			endOfWeek = trdmLastUpdate.Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond)
		}

		missingWeeks = append(missingWeeks, MissingWeek{StartOfWeek: startOfWeek, EndOfWeek: endOfWeek}) // Not always full weeks, sometimes split in half
	}

	return missingWeeks, nil
}
