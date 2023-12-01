package trdm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/loa"
	"github.com/transcom/mymove/pkg/parser/tac"
)

func GetTGETData(getTableRequest models.GetTableRequest, lastTableUpdateResponseLastUpdate time.Time, service GatewayService, appCtx appcontext.AppContext, logger *zap.Logger) error {
	// Setup response model
	getTableResponse := models.GetTableResponse{}

	// See how many weeks of data we're looking to gather before firing off to getTable
	missingWeeks, err := FetchWeeksOfMissingTime(getTableRequest.ContentUpdatedSinceDateTime, lastTableUpdateResponseLastUpdate)
	if err != nil {
		logger.Error("failed to identify missing weeks of data", zap.Error(err))
		return err
	}

	switch {
	case len(missingWeeks) > 2:
		// Since we're requesting more than 2 weeks of data, we need to split it up to not overload the response bodies.
		for _, week := range missingWeeks {
			weekRequest := getTableRequest
			// Add the first and second datetime filters based on the missing week
			// These are not provided inside of getTableRequest.
			weekRequest.ContentUpdatedSinceDateTime = week.StartOfWeek
			weekRequest.ContentUpdatedOnOrBeforeDateTime = &week.EndOfWeek

			// Fire off this weeks request
			resp, err := service.gatewayGetTable(weekRequest)
			if err != nil {
				logger.Fatal("failed to call gatewayGetTable for week", zap.Error(err), zap.Time("startOfWeek", week.StartOfWeek), zap.Time("endOfWeek", week.EndOfWeek))
				return err
			}
			// Read it
			if resp.Body == nil {
				return errors.New("received empty body response from API gateway")
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("failed to read response body", zap.Error(err))
				return err
			}
			defer resp.Body.Close()

			// Parse it into getTableResponse model
			err = json.Unmarshal(body, &getTableResponse)
			if err != nil {
				logger.Error("failed to unmarshal response body into getTableResponse type", zap.Error(err))
				return err
			}

			// Parse the attachment, this will also store it in the DB if all goes well
			err = parseGetTableResponse(appCtx, getTableResponse.Attachment, getTableRequest.PhysicalName)
			if err != nil {
				logger.Error("failed to parseGetTableResponse and store it into the database", zap.Error(err))
				return err
			}
			logger.Info("retrieving trdm TGET data successful for week", zap.String("request physicalName", weekRequest.PhysicalName), zap.Time("startOfWeek", week.StartOfWeek), zap.Time("endOfWeek", week.EndOfWeek))
		}
	default:
		// Forward model to getTable to gather TGET data
		resp, err := service.gatewayGetTable(getTableRequest)
		if err != nil {
			logger.Error("failed to call gatewayGetTable", zap.Error(err))
			return err
		}
		// Read it
		if resp.Body == nil {
			return errors.New("received empty body response from API gateway")
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("failed to read response body", zap.Error(err))
			return err
		}
		defer resp.Body.Close()

		// Parse it into getTableResponse model
		err = json.Unmarshal(body, &getTableResponse)
		if err != nil {
			logger.Error("failed to unmarshal response body into getTableResponse type", zap.Error(err))
			return err
		}

		// Parse the attachment, this will also store it in the DB if all goes well
		err = parseGetTableResponse(appCtx, getTableResponse.Attachment, getTableRequest.PhysicalName)
		if err != nil {
			logger.Error("failed to parseGetTableResponse and store it into the database", zap.Error(err))
			return err
		}
	}
	logger.Info("retrieving trdm TGET data successful", zap.String("request physicalName", getTableRequest.PhysicalName))
	return nil
}

// Parses pipedelimited file attachment from GetTable webservice and saves records to database
//
//	returns error
func parseGetTableResponse(appcontext appcontext.AppContext, attachment []byte, physicalName string) error {
	reader := bytes.NewReader(attachment)
	switch physicalName {
	case LineOfAccounting:
		loaCodes, err := loa.Parse(reader)
		if err != nil {
			return err
		}
		err = createLoaCodes(appcontext, loaCodes)
		if err != nil {
			return err
		}
	case TransportationAccountingCode:
		tacCodes, err := tac.Parse(reader)
		// Consolidate duplicates
		consolidatedTacs := tac.ConsolidateDuplicateTACsDesiredFromTRDM(tacCodes)
		if err != nil {
			return err
		}
		if err = createTacCodes(appcontext, consolidatedTacs); err != nil {
			return err
		}
	default:
		return errors.New("provided physical name is not valid for TGET data")
	}
	return nil
}

// Saves new TAC Code slice to DB
func createTacCodes(appcontext appcontext.AppContext, tacCodes []models.TransportationAccountingCode) error {
	saveErr := appcontext.DB().Create(tacCodes)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

// Saves new LOA Code slice to DB
func createLoaCodes(appcontext appcontext.AppContext, loa []models.LineOfAccounting) error {
	saveErr := appcontext.DB().Create(loa)
	if saveErr != nil {
		return saveErr
	}
	return nil
}
