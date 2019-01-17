package publicapi

import (
	"github.com/transcom/mymove/pkg/gen/apimessages"
)

func payloadForGBLOC(gbloc *string) *apimessages.GBLOC {
	if gbloc == nil {
		return nil
	}
	g := apimessages.GBLOC(*gbloc)
	return &g
}
