package internalapi

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	calendarop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"net/http/httptest"
	"time"
)

func (suite *HandlerSuite) TestShowUnavailableMoveDatesHandler() {
	req := httptest.NewRequest("GET", "/calendar/unavailable_move_dates", nil)

	params := calendarop.ShowUnavailableMoveDatesParams{
		HTTPRequest: req,
		StartDate:   strfmt.Date(time.Date(2018, 9, 26, 0, 0, 0, 0, time.UTC)),
	}

	unavailableDates := []strfmt.Date{
		strfmt.Date(time.Date(2018, 9, 26, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 29, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 9, 30, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 1, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 2, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 6, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 7, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 13, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 14, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 20, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 21, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 27, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 28, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 3, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 4, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 10, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 11, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 17, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 18, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 24, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 11, 25, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 1, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 2, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 8, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 22, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 12, 23, 0, 0, 0, 0, time.UTC)),
	}

	showHandler := ShowUnavailableMoveDatesHandler{handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())}
	response := showHandler.Handle(params)

	suite.IsType(&calendarop.ShowUnavailableMoveDatesOK{}, response)
	okResponse := response.(*calendarop.ShowUnavailableMoveDatesOK)

	suite.Equal(unavailableDates, okResponse.Payload)
}

func (suite *HandlerSuite) TestShowMoveDatesSummaryHandler() {
	dutyStationAddress := testdatagen.MakeAddress(suite.TestDB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "Fort Gordon",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30813",
			Country:        swag.String("United States"),
		},
	})

	dutyStation := testdatagen.MakeDutyStation(suite.TestDB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:        "Fort Sam Houston",
			Affiliation: internalmessages.AffiliationARMY,
			AddressID:   dutyStationAddress.ID,
			Address:     dutyStationAddress,
		},
	})

	rank := internalmessages.ServiceMemberRankE4
	serviceMember := testdatagen.MakeServiceMember(suite.TestDB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Rank:          &rank,
			DutyStationID: &dutyStation.ID,
			DutyStation:   dutyStation,
		},
	})

	newDutyStationAddress := testdatagen.MakeAddress(suite.TestDB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "n/a",
			City:           "San Antonio",
			State:          "TX",
			PostalCode:     "78234",
			Country:        swag.String("United States"),
		},
	})

	newDutyStation := testdatagen.MakeDutyStation(suite.TestDB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:        "Fort Gordon",
			Affiliation: internalmessages.AffiliationARMY,
			AddressID:   newDutyStationAddress.ID,
			Address:     newDutyStationAddress,
		},
	})

	move := testdatagen.MakeMove(suite.TestDB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID:  serviceMember.ID,
			ServiceMember:    serviceMember,
			ReportByDate:     time.Date(2018, 10, 20, 0, 0, 0, 0, time.UTC),
			NewDutyStationID: newDutyStation.ID,
			NewDutyStation:   newDutyStation,
		},
	})

	path := fmt.Sprintf("/calendar/%s/move_dates", move.ID.String())
	req := httptest.NewRequest("GET", path, nil)

	params := calendarop.ShowMoveDatesSummaryParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
		MoveDate:    strfmt.Date(time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC)),
	}

	context := handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1125))

	showHandler := ShowMoveDatesSummaryHandler{context}
	response := showHandler.Handle(params)

	suite.IsType(&calendarop.ShowMoveDatesSummaryOK{}, response)
	okResponse := response.(*calendarop.ShowMoveDatesSummaryOK)

	moveDates := internalmessages.MoveDatesSummaryPayload{
		Pack: []strfmt.Date{
			strfmt.Date(time.Date(2018, 9, 27, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 1, 0, 0, 0, 0, time.UTC)),
		},
		Pickup: []strfmt.Date{
			strfmt.Date(time.Date(2018, 10, 2, 0, 0, 0, 0, time.UTC)),
		},
		Transit: []strfmt.Date{
			strfmt.Date(time.Date(2018, 10, 2, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 3, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 4, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 5, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 8, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 9, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 11, 0, 0, 0, 0, time.UTC)),
			strfmt.Date(time.Date(2018, 10, 12, 0, 0, 0, 0, time.UTC)),
		},
		Delivery: []strfmt.Date{
			strfmt.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
		},
		Report: []strfmt.Date{
			strfmt.Date(move.Orders.ReportByDate),
		},
	}

	suite.Equal(moveDates, *okResponse.Payload)
}
