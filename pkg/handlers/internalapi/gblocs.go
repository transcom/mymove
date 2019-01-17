package internalapi

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func payloadForGBLOC(gbloc *string) *internalmessages.GBLOC {
	if gbloc == nil {
		return nil
	}
	g := internalmessages.GBLOC(*gbloc)
	return &g
}
