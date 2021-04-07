package move

import (
	"strings"
	"time"

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
	// TODO: In future, add logic based on the service member's origin duty station
	// to route to send services counseling, otherwise submitted
	var err error
	if needsServiceCounseling() {
		err = move.SendToServiceCounseling()
	} else {
		submitDate := time.Now()
		err = move.Submit(submitDate)
	}
	if err != nil {
		return err
	}
	return nil
}

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
