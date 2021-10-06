package auth

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

const (
	// OfficeTestHost
	OfficeTestHost string = "office.example.com"
	// MilTestHost
	MilTestHost string = "mil.example.com"
	// OrdersTestHost
	OrdersTestHost string = "orders.example.com"
	// DpsTestHost
	DpsTestHost string = "dps.example.com"
	// SddcTestHost
	SddcTestHost string = "sddc.example.com"
	// AdminTestHost
	AdminTestHost string = "admin.example.com"
	// PrimeTestHost
	PrimeTestHost string = "prime.example.com"
)

// ApplicationTestServername is a collection of the test servernames
func ApplicationTestServername() ApplicationServername {
	appnames := ApplicationServername{
		MilServername:    MilTestHost,
		OfficeServername: OfficeTestHost,
		AdminServername:  AdminTestHost,
		OrdersServername: OrdersTestHost,
		DpsServername:    DpsTestHost,
		SddcServername:   SddcTestHost,
		PrimeServername:  PrimeTestHost,
	}
	return appnames
}

type authSuite struct {
	testingsuite.BaseTestSuite
	logger *zap.Logger
}

func TestAuthSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &authSuite{logger: logger}
	suite.Run(t, hs)
}
