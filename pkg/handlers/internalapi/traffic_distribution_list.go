package internalapi

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForTrafficDistributionListModel(tdl *models.TrafficDistributionList) *internalmessages.TrafficDistributionList {
	if tdl == nil {
		return nil
	}
	return &internalmessages.TrafficDistributionList{
		SourceRateArea:    swag.String(tdl.SourceRateArea),
		DestinationRegion: swag.String(tdl.DestinationRegion),
		CodeOfService:     swag.String(tdl.CodeOfService),
	}
}
