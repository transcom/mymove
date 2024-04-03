package ppmshipment

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PPMShipmentSuite struct {
	*testingsuite.PopTestSuite
	filesToClose []afero.File
}

func TestPPMShipmentServiceSuite(t *testing.T) {
	ts := &PPMShipmentSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

// setUpMockPPMShipmentUpdater sets up the input mock PPMShipmentUpdater to return the given return values once. It is
// meant as a helper to not have to remember the exact mocking syntax each time, and to cut a tiny bit of boilerplate.
func setUpMockPPMShipmentUpdater(
	mockPPMShipmentUpdater *mocks.PPMShipmentUpdater,
	appCtx appcontext.AppContext,
	ppmShipment *models.PPMShipment,
	returnValue ...interface{},
) {
	mockPPMShipmentUpdater.
		On(
			"UpdatePPMShipmentWithDefaultCheck",
			appCtx,
			ppmShipment,
			ppmShipment.Shipment.ID,
		).
		Return(returnValue...).
		Once()
}

// setUpMockPPMShipmentFetcher sets up the input mock PPMShipmentFetcher to return the given return values once. It is
// meant as a helper to not have to remember the exact mocking syntax each time, and to cut a tiny bit of boilerplate.
func setUpMockPPMShipmentFetcher(
	mockPPMShipmentFetcher *mocks.PPMShipmentFetcher,
	appCtx appcontext.AppContext,
	ppmShipmentID uuid.UUID,
	eagerPreloadAssociations []string,
	postloadAssociations []string,
	returnValue ...interface{},
) {
	mockPPMShipmentFetcher.
		On(
			"GetPPMShipment",
			appCtx,
			ppmShipmentID,
			eagerPreloadAssociations,
			postloadAssociations,
		).
		Return(returnValue...).
		Once()
}
