package paperwork

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmComputerParams struct {
	Weight                                unit.Pound
	OriginPickupZip5                      string
	OriginDutyStationZip5                 string
	DestinationZip5                       string
	DistanceMilesFromOriginPickupZip      int
	DistanceMilesFromOriginDutyStationZip int
	Date                                  time.Time
	DaysInSIT                             int
}

type mockPPMComputer struct {
	costComputation   rateengine.CostComputation
	err               error
	ppmComputerParams []ppmComputerParams
}

func (mppmc *mockPPMComputer) ComputeLowestCostPPMMove(weight unit.Pound, originPickupZip5 string, originDutyStationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyStationZip int, date time.Time, daysInSit int) (cost rateengine.CostComputation, err error) {
	mppmc.ppmComputerParams = append(mppmc.ppmComputerParams, ppmComputerParams{
		Weight:                                weight,
		OriginPickupZip5:                      originPickupZip5,
		OriginDutyStationZip5:                 originDutyStationZip5,
		DestinationZip5:                       destinationZip5,
		DistanceMilesFromOriginPickupZip:      distanceMilesFromOriginPickupZip,
		DistanceMilesFromOriginDutyStationZip: distanceMilesFromOriginDutyStationZip,
		Date:                                  date,
		DaysInSIT:                             daysInSit,
	})
	return mppmc.costComputation, mppmc.err
}

func (mppmc *mockPPMComputer) CalledWith() []ppmComputerParams {
	return mppmc.ppmComputerParams
}

func (suite *PaperworkSuite) TestComputeObligationsParams() {
	ppmComputer := NewSSWPPMComputer(&mockPPMComputer{})
	pickupPostalCode := "85369"
	destinationPostalCode := "31905"
	ppm := models.PersonallyProcuredMove{
		PickupPostalCode:      &pickupPostalCode,
		DestinationPostalCode: &destinationPostalCode,
	}
	noPPM := models.ShipmentSummaryFormData{PersonallyProcuredMoves: models.PersonallyProcuredMoves{}}
	missingZip := models.ShipmentSummaryFormData{PersonallyProcuredMoves: models.PersonallyProcuredMoves{{}}}
	missingActualMoveDate := models.ShipmentSummaryFormData{PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm}}

	_, err1 := ppmComputer.ComputeObligations(noPPM, route.NewTestingPlanner(10))
	_, err2 := ppmComputer.ComputeObligations(missingZip, route.NewTestingPlanner(10))
	_, err3 := ppmComputer.ComputeObligations(missingActualMoveDate, route.NewTestingPlanner(10))

	suite.NotNil(err1)
	suite.NotNil(err2)
	suite.NotNil(err3)
}

func (suite *PaperworkSuite) TestTestComputeObligations() {
	miles := 100
	totalWeightEntitlement := unit.Pound(1000)
	ppmRemainingEntitlement := unit.Pound(2000)
	planner := route.NewTestingPlanner(miles)
	origMoveDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	actualDate := time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC)
	pickupPostalCode := "85369"
	destinationPostalCode := "31905"
	cents := unit.Cents(1000)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate:      &origMoveDate,
			ActualMoveDate:        &actualDate,
			PickupPostalCode:      &pickupPostalCode,
			DestinationPostalCode: &destinationPostalCode,
			TotalSITCost:          &cents,
		},
	})

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "31905",
	}
	suite.MustSave(&address)

	stationName := "New Duty Station"
	station := models.DutyStation{
		Name:        stationName,
		Affiliation: internalmessages.AffiliationAIRFORCE,
		AddressID:   address.ID,
		Address:     address,
	}
	suite.MustSave(&station)

	orderID := uuid.Must(uuid.NewV4())
	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:               orderID,
			NewDutyStationID: station.ID,
			NewDutyStation:   station,
		},
	})

	currentDutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	params := models.ShipmentSummaryFormData{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm},
		WeightAllotment:         models.SSWMaxWeightEntitlement{TotalWeight: totalWeightEntitlement},
		PPMRemainingEntitlement: ppmRemainingEntitlement,
		CurrentDutyStation:      currentDutyStation,
		Order:                   order,
	}
	suite.Run("TestComputeObligations", func() {
		mockComputer := mockPPMComputer{
			costComputation: rateengine.CostComputation{GCC: 100, SITMax: 20000},
		}
		ppmComputer := NewSSWPPMComputer(&mockComputer)
		expectMaxObligationParams := ppmComputerParams{
			Weight:                                totalWeightEntitlement,
			OriginPickupZip5:                      pickupPostalCode,
			OriginDutyStationZip5:                 currentDutyStation.Address.PostalCode,
			DestinationZip5:                       destinationPostalCode,
			DistanceMilesFromOriginPickupZip:      miles,
			DistanceMilesFromOriginDutyStationZip: miles,
			Date:                                  actualDate,
			DaysInSIT:                             0,
		}
		expectActualObligationParams := ppmComputerParams{
			Weight:                                ppmRemainingEntitlement,
			OriginPickupZip5:                      pickupPostalCode,
			OriginDutyStationZip5:                 currentDutyStation.Address.PostalCode,
			DestinationZip5:                       destinationPostalCode,
			DistanceMilesFromOriginPickupZip:      miles,
			DistanceMilesFromOriginDutyStationZip: miles,
			Date:                                  actualDate,
			DaysInSIT:                             0,
		}
		cost, err := ppmComputer.ComputeObligations(params, planner)

		suite.NoError(err)
		calledWith := mockComputer.CalledWith()
		suite.Equal(*ppm.TotalSITCost, cost.ActualObligation.SIT)
		suite.Equal(expectActualObligationParams, calledWith[0])
		suite.Equal(expectMaxObligationParams, calledWith[1])
	})

	suite.Run("TestComputeObligations when actual PPM SIT exceeds MaxSIT", func() {
		mockComputer := mockPPMComputer{
			costComputation: rateengine.CostComputation{SITMax: unit.Cents(500)},
		}
		ppmComputer := NewSSWPPMComputer(&mockComputer)
		obligations, err := ppmComputer.ComputeObligations(params, planner)

		suite.NoError(err)
		suite.Equal(unit.Cents(500), obligations.ActualObligation.SIT)
	})

	suite.Run("TestComputeObligations when there is no actual PPM SIT", func() {
		ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				OriginalMoveDate:      &origMoveDate,
				ActualMoveDate:        &actualDate,
				PickupPostalCode:      &pickupPostalCode,
				DestinationPostalCode: &destinationPostalCode,
			},
		})
		currentDutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
		shipmentSummaryFormParams := models.ShipmentSummaryFormData{
			PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm},
			WeightAllotment:         models.SSWMaxWeightEntitlement{TotalWeight: totalWeightEntitlement},
			CurrentDutyStation:      currentDutyStation,
			Order:                   order,
		}
		mockComputer := mockPPMComputer{
			costComputation: rateengine.CostComputation{SITMax: unit.Cents(500)},
		}
		ppmComputer := NewSSWPPMComputer(&mockComputer)
		obligations, err := ppmComputer.ComputeObligations(shipmentSummaryFormParams, planner)

		suite.NoError(err)
		suite.Equal(unit.Cents(0), obligations.ActualObligation.SIT)
	})

	suite.Run("TestCalcError", func() {
		mockComputer := mockPPMComputer{err: errors.New("ERROR")}
		ppmComputer := SSWPPMComputer{&mockComputer}
		_, err := ppmComputer.ComputeObligations(params, planner)

		suite.NotNil(err)
	})
}
