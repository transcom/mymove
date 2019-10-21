package services

import (
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/unit"
)

// DomesticServicePricingData contains the input data for pricing domestic linehaul
type DomesticServicePricingData struct {
	MoveDate    time.Time
	Distance    unit.Miles
	Weight      unit.Pound
	ServiceArea string
}

// MarshalLogObject allows DomesticServicePricingData to be logged by zap
func (d DomesticServicePricingData) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddTime("MoveDate", d.MoveDate)
	encoder.AddInt("Distance", d.Distance.Int())
	encoder.AddInt("Weight", d.Weight.Int())
	encoder.AddString("ServiceArea", d.ServiceArea)
	return nil
}

// DomesticLinehaulPricer prices domestic linehaul for a GHC move
//go:generate mockery -name DomesticLinehaulPricer
type DomesticLinehaulPricer interface {
	PriceDomesticLinehaul(DomesticServicePricingData) (unit.Cents, error)
}
