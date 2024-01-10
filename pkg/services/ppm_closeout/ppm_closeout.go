package ppmcloseout

import (
	"github.com/transcom/mymove/pkg/services"
)

// ppmDocumentFetcher is the concrete implementation of the services.PPMDocumentFetcher interface
type ppmCloseout struct{}

// NewPPMDocumentFetcher creates a new struct
func NewPPMCloseout() services.PPMCloseout {
	return &ppmCloseout{}
}

func (p *ppmCloseout) GetPPMCloseout() {

}
