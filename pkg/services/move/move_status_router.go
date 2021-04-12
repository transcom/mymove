package move

import (
	"strings"

	"github.com/gobuffalo/pop/v5"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveStatusRouter struct {
	db *pop.Connection
}

// NewMoveStatusRouter creates a new moveStatusRouter service
func NewMoveStatusRouter(db *pop.Connection) services.MoveStatusRouter {
	return &moveStatusRouter{db}
}

//FetchOrder retrieves a Move if it is visible for a given locator
func (f moveStatusRouter) RouteMove(move *models.Move) error {
	var err error

	if needsServiceCounseling() {
		err = move.SendToServiceCounseling()
	} else {
		err = move.Submit()
	}
	if err != nil {
		return err
	}
	return nil
}

// TODO: Replace the code in this function to determine whether or not the move
// needs service counseling based on the service member's origin duty station.
// Then remove all code related to the service counseling feature flag here and
// in pkg/cli/featureflag.go, and remove any references to
// `FEATURE_FLAG_SERVICE_COUNSELING` from the entire project.
// You'll need to update the test setup in TestSubmitMoveForServiceCounselingHandler
// so that the move's origin duty station will trigger service counseling.
func needsServiceCounseling() bool {
	logger := zap.NewNop()
	flag := pflag.CommandLine
	v := viper.New()
	err := v.BindPFlags(flag)
	if err != nil {
		logger.Fatal("could not bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return v.GetBool(cli.FeatureFlagServiceCounseling)
}
