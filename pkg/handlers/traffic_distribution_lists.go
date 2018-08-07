package handlers

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
)

/*
 * --------------------------------------------
 * The code below is for the INTERNAL REST API.
 * --------------------------------------------
 */

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

func publicPayloadForTrafficDistributionListModel(tdl *models.TrafficDistributionList) *apimessages.TrafficDistributionList {
	if tdl == nil {
		return nil
	}
	return &apimessages.TrafficDistributionList{
		SourceRateArea:    swag.String(tdl.SourceRateArea),
		DestinationRegion: swag.String(tdl.DestinationRegion),
		CodeOfService:     swag.String(tdl.CodeOfService),
	}
}
