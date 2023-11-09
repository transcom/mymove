package trdm

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/loa"
	"github.com/transcom/mymove/pkg/parser/tac"
)

func getTGETData(getTableRequest models.GetTableRequest, service GatewayService, appCtx appcontext.AppContext) error {
	// Setup response model
	getTableResponse := models.GetTableResponse{}

	// Forward model to getTable to gather TGET data
	resp, err := service.gatewayGetTable(getTableRequest)
	if err != nil {
		return err
	}
	// Read it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse it into getTableResponse model
	err = json.Unmarshal(body, &getTableResponse)
	if err != nil {
		return err
	}

	// Parse the attachment, this will also store it in the DB if all goes well
	err = parseGetTableResponse(appCtx, getTableResponse.Attachment, getTableRequest.PhysicalName)
	if err != nil {
		return err
	}

	return nil
}

// Parses pipedelimited file attachment from GetTable webservice and saves records to database
//
//	returns error
func parseGetTableResponse(appcontext appcontext.AppContext, attachment []byte, physicalName string) error {
	reader := bytes.NewReader(attachment)
	switch physicalName {
	case lineOfAccounting:
		loaCodes, err := loa.Parse(reader)
		if err != nil {
			return err
		}
		err = saveLoaCodes(appcontext, loaCodes)
		if err != nil {
			return err
		}
	case transportationAccountingCode:
		tacCodes, err := tac.Parse(reader)
		consolidatedTacs := tac.ConsolidateDuplicateTACsDesiredFromTRDM(tacCodes)
		if err != nil {
			return err
		}
		if err = saveTacCodes(appcontext, consolidatedTacs); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

// Saves TAC Code slice to DB and updates records
func saveTacCodes(appcontext appcontext.AppContext, tacCodes []models.TransportationAccountingCode) error {
	saveErr := appcontext.DB().Update(tacCodes)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

// Saves LOA Code slice to DB and updates records
func saveLoaCodes(appcontext appcontext.AppContext, loa []models.LineOfAccounting) error {
	saveErr := appcontext.DB().Update(loa)
	if saveErr != nil {
		return saveErr
	}
	return nil
}
