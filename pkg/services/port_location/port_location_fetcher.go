package portlocation

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type portLocationFetcher struct {
}

// NewPortLocationFetcher returns a new port location fetcher
func NewPortLocationFetcher() services.PortLocationFetcher {
	return &portLocationFetcher{}
}

func (p *portLocationFetcher) FetchPortLocationByPortCode(appCtx appcontext.AppContext, portCode string) (*models.PortLocation, error) {
	portLocation := models.PortLocation{}
	err := appCtx.DB().Eager("Port").Where("is_active = TRUE").InnerJoin("ports p", "port_id = p.id").Where("p.port_code = $1", portCode).First(&portLocation)
	if err != nil {
		return nil, apperror.NewQueryError("PortLocation", err, "")
	}
	return &portLocation, err
}
