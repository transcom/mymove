package ppmcloseout

import (
	"github.com/transcom/mymove/pkg/services"
)

// ppmDocumentFetcher is the concrete implementation of the services.PPMDocumentFetcher interface
type ppmCloseoutFetcher struct{}

// NewPPMDocumentFetcher creates a new struct
func NewPPMCloseoutFetcher() services.PPMCloseout {
	return &ppmCloseoutFetcher{}
}

func (p *ppmCloseoutFetcher) GetPPMCloseout() {

}
