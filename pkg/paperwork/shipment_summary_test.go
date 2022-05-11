package paperwork

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmComputerParams struct {
	Weight                                 unit.Pound
	OriginPickupZip5                       string
	OriginDutyLocationZip5                 string
	DestinationZip5                        string
	DistanceMilesFromOriginPickupZip       int
	DistanceMilesFromOriginDutyLocationZip int
	Date                                   time.Time
	DaysInSIT                              int
}

type mockPPMComputer struct {
	costDetails       rateengine.CostDetails
	err               error
	ppmComputerParams []ppmComputerParams
}

func (mppmc *mockPPMComputer) ComputePPMMoveCosts(appCtx appcontext.AppContext, weight unit.Pound, originPickupZip5 string, originDutyLocationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyLocationZip int, date time.Time, daysInSit int) (cost rateengine.CostDetails, err error) {
	mppmc.ppmComputerParams = append(mppmc.ppmComputerParams, ppmComputerParams{
		Weight:                                 weight,
		OriginPickupZip5:                       originPickupZip5,
		OriginDutyLocationZip5:                 originDutyLocationZip5,
		DestinationZip5:                        destinationZip5,
		DistanceMilesFromOriginPickupZip:       distanceMilesFromOriginPickupZip,
		DistanceMilesFromOriginDutyLocationZip: distanceMilesFromOriginDutyLocationZip,
		Date:                                   date,
		DaysInSIT:                              daysInSit,
	})
	return mppmc.costDetails, mppmc.err
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

	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(10, nil)
	_, err1 := ppmComputer.ComputeObligations(suite.AppContextForTest(), noPPM, planner)
	_, err2 := ppmComputer.ComputeObligations(suite.AppContextForTest(), missingZip, planner)
	_, err3 := ppmComputer.ComputeObligations(suite.AppContextForTest(), missingActualMoveDate, planner)

	suite.NotNil(err1)
	suite.Equal("missing ppm", err1.Error())

	suite.NotNil(err2)
	suite.Equal("missing required address parameter", err2.Error())

	suite.NotNil(err3)
	suite.Equal("missing required original move date parameter", err3.Error())
}

func (suite *PaperworkSuite) TestComputeObligations() {
	miles := 100
	totalWeightEntitlement := unit.Pound(1000)
	ppmRemainingEntitlement := unit.Pound(2000)
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(miles, nil)
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

	locationName := "New Duty Location"
	location := models.DutyLocation{
		Name:      locationName,
		AddressID: address.ID,
		Address:   address,
	}
	suite.MustSave(&location)

	orderID := uuid.Must(uuid.NewV4())
	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:                orderID,
			NewDutyLocationID: location.ID,
			NewDutyLocation:   location,
		},
	})

	currentDutyLocation := testdatagen.FetchOrMakeDefaultCurrentDutyLocation(suite.DB())
	params := models.ShipmentSummaryFormData{
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm},
		WeightAllotment:         models.SSWMaxWeightEntitlement{TotalWeight: totalWeightEntitlement},
		PPMRemainingEntitlement: ppmRemainingEntitlement,
		CurrentDutyLocation:     currentDutyLocation,
		Order:                   order,
	}
	suite.Run("TestComputeObligations", func() {
		var costDetails = make(rateengine.CostDetails)
		costDetails["pickupLocation"] = &rateengine.CostDetail{
			Cost:      rateengine.CostComputation{GCC: 100, SITMax: 20000},
			IsWinning: true,
		}
		costDetails["originDutyLocation"] = &rateengine.CostDetail{
			Cost:      rateengine.CostComputation{GCC: 200, SITMax: 30000},
			IsWinning: true,
		}

		mockComputer := mockPPMComputer{
			costDetails: costDetails,
		}
		ppmComputer := NewSSWPPMComputer(&mockComputer)
		expectMaxObligationParams := ppmComputerParams{
			Weight:                                 totalWeightEntitlement,
			OriginPickupZip5:                       pickupPostalCode,
			OriginDutyLocationZip5:                 currentDutyLocation.Address.PostalCode,
			DestinationZip5:                        destinationPostalCode,
			DistanceMilesFromOriginPickupZip:       miles,
			DistanceMilesFromOriginDutyLocationZip: miles,
			Date:                                   origMoveDate,
			DaysInSIT:                              0,
		}
		expectActualObligationParams := ppmComputerParams{
			Weight:                                 ppmRemainingEntitlement,
			OriginPickupZip5:                       pickupPostalCode,
			OriginDutyLocationZip5:                 currentDutyLocation.Address.PostalCode,
			DestinationZip5:                        destinationPostalCode,
			DistanceMilesFromOriginPickupZip:       miles,
			DistanceMilesFromOriginDutyLocationZip: miles,
			Date:                                   origMoveDate,
			DaysInSIT:                              0,
		}
		cost, err := ppmComputer.ComputeObligations(suite.AppContextForTest(), params, planner)

		suite.NoError(err)
		calledWith := mockComputer.CalledWith()
		suite.Equal(*ppm.TotalSITCost, cost.ActualObligation.SIT)
		suite.Equal(expectActualObligationParams, calledWith[0])
		suite.Equal(expectMaxObligationParams, calledWith[1])
	})

	suite.Run("TestComputeObligations when actual PPM SIT exceeds MaxSIT", func() {
		var costDetails = make(rateengine.CostDetails)
		costDetails["pickupLocation"] = &rateengine.CostDetail{
			Cost:      rateengine.CostComputation{SITMax: unit.Cents(500)},
			IsWinning: true,
		}
		costDetails["originDutyLocation"] = &rateengine.CostDetail{
			Cost:      rateengine.CostComputation{SITMax: unit.Cents(600)},
			IsWinning: false,
		}
		mockComputer := mockPPMComputer{
			costDetails: costDetails,
		}
		ppmComputer := NewSSWPPMComputer(&mockComputer)
		obligations, err := ppmComputer.ComputeObligations(suite.AppContextForTest(), params, planner)

		suite.NoError(err)
		suite.Equal(unit.Cents(500), obligations.ActualObligation.SIT)
	})

	suite.Run("TestComputeObligations when there is no actual PPM SIT", func() {
		var costDetails = make(rateengine.CostDetails)
		costDetails["pickupLocation"] = &rateengine.CostDetail{
			Cost:      rateengine.CostComputation{SITMax: unit.Cents(500)},
			IsWinning: true,
		}
		costDetails["originDutyLocation"] = &rateengine.CostDetail{
			Cost:      rateengine.CostComputation{SITMax: unit.Cents(600)},
			IsWinning: false,
		}
		mockComputer := mockPPMComputer{
			costDetails: costDetails,
		}

		ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				OriginalMoveDate:      &origMoveDate,
				ActualMoveDate:        &actualDate,
				PickupPostalCode:      &pickupPostalCode,
				DestinationPostalCode: &destinationPostalCode,
			},
		})
		currentDutyLocation := testdatagen.FetchOrMakeDefaultCurrentDutyLocation(suite.DB())
		shipmentSummaryFormParams := models.ShipmentSummaryFormData{
			PersonallyProcuredMoves: models.PersonallyProcuredMoves{ppm},
			WeightAllotment:         models.SSWMaxWeightEntitlement{TotalWeight: totalWeightEntitlement},
			CurrentDutyLocation:     currentDutyLocation,
			Order:                   order,
		}
		ppmComputer := NewSSWPPMComputer(&mockComputer)
		obligations, err := ppmComputer.ComputeObligations(suite.AppContextForTest(), shipmentSummaryFormParams, planner)

		suite.NoError(err)
		suite.Equal(unit.Cents(0), obligations.ActualObligation.SIT)
	})

	suite.Run("TestCalcError", func() {
		mockComputer := mockPPMComputer{err: errors.New("ERROR")}
		ppmComputer := SSWPPMComputer{&mockComputer}
		_, err := ppmComputer.ComputeObligations(suite.AppContextForTest(), params, planner)

		suite.NotNil(err)
	})
}
