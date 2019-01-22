package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestCreatePPMHandler() {
	user1 := testdatagen.MakeDefaultServiceMember(suite.DB())
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHGPPM

	move, verrs, locErr := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.False(verrs.HasAny(), "failed to create new move")
	suite.Nil(locErr)

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.AuthenticateRequest(request, orders.ServiceMember)

	newPPMPayload := internalmessages.CreatePersonallyProcuredMovePayload{
		WeightEstimate:   swag.Int64(12),
		PickupPostalCode: swag.String("00112"),
		DaysInStorage:    swag.Int64(3),
	}

	newPPMParams := ppmop.CreatePersonallyProcuredMoveParams{
		MoveID:                              strfmt.UUID(move.ID.String()),
		CreatePersonallyProcuredMovePayload: &newPPMPayload,
		HTTPRequest:                         request,
	}

	handler := CreatePersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(newPPMParams)
	// assert we got back the 201 response
	createdResponse := response.(*ppmop.CreatePersonallyProcuredMoveCreated)
	createdIssuePayload := createdResponse.Payload
	suite.NotNil(createdIssuePayload.ID)

	// Next try the wrong user
	request = suite.AuthenticateRequest(request, user1)
	newPPMParams.HTTPRequest = request

	badUserResponse := handler.Handle(newPPMParams)
	suite.CheckResponseForbidden(badUserResponse)

	// Now try a bad move
	newPPMParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newPPMParams)
	suite.CheckResponseNotFound(badMoveResponse)

}

func (suite *HandlerSuite) TestSubmitPPMHandler() {
	t := suite.T()

	// create a ppm
	move1 := testdatagen.MakeDefaultMove(suite.DB())
	ppm := models.PersonallyProcuredMove{
		MoveID:         move1.ID,
		Move:           move1,
		WeightEstimate: swag.Int64(1),
		Status:         models.PPMStatusDRAFT,
	}

	verrs, err := suite.DB().ValidateAndCreate(&ppm)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	req := httptest.NewRequest("POST", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move1.Orders.ServiceMember)

	submitPPMParams := ppmop.SubmitPersonallyProcuredMoveParams{
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
		HTTPRequest:              req,
	}

	// submit the PPM
	handler := SubmitPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(submitPPMParams)
	okResponse := response.(*ppmop.SubmitPersonallyProcuredMoveOK)
	submitPPMPayload := okResponse.Payload

	suite.Require().Equal(submitPPMPayload.Status, internalmessages.PPMStatusSUBMITTED, "PPM should have been submitted")
}

func (suite *HandlerSuite) TestIndexPPMHandler() {

	t := suite.T()

	// Given: moves and associated PPMs
	move1 := testdatagen.MakeDefaultMove(suite.DB())
	move2 := testdatagen.MakeDefaultMove(suite.DB())

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move1.ID,
		Move:           move1,
		WeightEstimate: swag.Int64(1),
		Status:         models.PPMStatusDRAFT,
	}
	ppm2 := models.PersonallyProcuredMove{
		MoveID:         move1.ID,
		Move:           move1,
		WeightEstimate: swag.Int64(2),
		Status:         models.PPMStatusDRAFT,
	}
	otherPPM := models.PersonallyProcuredMove{
		MoveID:         move2.ID,
		Move:           move2,
		WeightEstimate: swag.Int64(4),
		Status:         models.PPMStatusDRAFT,
	}

	verrs, err := suite.DB().ValidateAndCreate(&ppm1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	verrs, err = suite.DB().ValidateAndCreate(&ppm2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	verrs, err = suite.DB().ValidateAndCreate(&otherPPM)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move1.Orders.ServiceMember)

	indexPPMParams := ppmop.IndexPersonallyProcuredMovesParams{
		MoveID:      strfmt.UUID(move1.ID.String()),
		HTTPRequest: req,
	}

	handler := IndexPersonallyProcuredMovesHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(indexPPMParams)

	// assert we got back the 201 response
	okResponse := response.(*ppmop.IndexPersonallyProcuredMovesOK)
	indexPPMPayload := okResponse.Payload

	for _, ppm := range indexPPMPayload {
		if *ppm.ID == *handlers.FmtUUID(otherPPM.ID) {
			t.Error("We should only have got back ppms associated with this move")
		}
	}
	if len(indexPPMPayload) != 2 {
		t.Error("We should have gotten back two good ones. ")
	}

}

func (suite *HandlerSuite) TestPatchPPMHandler() {
	scenario.RunRateEngineScenario1(suite.DB())
	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")

	initialWeight := swag.Int64(4100)
	newWeight := swag.Int64(4105)

	// Date picked essentialy at random, but needs to be within TestYear
	newMoveDate := time.Date(testdatagen.TestYear, time.November, 10, 23, 0, 0, 0, time.UTC)
	initialMoveDate := newMoveDate.Add(-2 * 24 * time.Hour)

	hasAdditionalPostalCode := swag.Bool(true)
	newHasAdditionalPostalCode := swag.Bool(false)
	additionalPickupPostalCode := swag.String("90210")

	hasSit := swag.Bool(true)
	newHasSit := swag.Bool(false)
	daysInStorage := swag.Int64(3)
	newPickupPostalCode := swag.String("32168")
	newDestinationPostalCode := swag.String("29401")

	move := testdatagen.MakeDefaultMove(suite.DB())

	newAdvanceWorksheet := models.Document{
		ServiceMember:   move.Orders.ServiceMember,
		ServiceMemberID: move.Orders.ServiceMemberID,
	}
	suite.MustSave(&newAdvanceWorksheet)

	ppm1 := models.PersonallyProcuredMove{
		MoveID:                     move.ID,
		Move:                       move,
		Size:                       &initialSize,
		WeightEstimate:             initialWeight,
		PlannedMoveDate:            &initialMoveDate,
		HasAdditionalPostalCode:    hasAdditionalPostalCode,
		AdditionalPickupPostalCode: additionalPickupPostalCode,
		HasSit:                     hasSit,
		DaysInStorage:              daysInStorage,
		Status:                     models.PPMStatusDRAFT,
		AdvanceWorksheet:           newAdvanceWorksheet,
		AdvanceWorksheetID:         &newAdvanceWorksheet.ID,
	}
	suite.MustSave(&ppm1)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:                    &newSize,
		WeightEstimate:          newWeight,
		PlannedMoveDate:         handlers.FmtDatePtr(&newMoveDate),
		HasAdditionalPostalCode: newHasAdditionalPostalCode,
		PickupPostalCode:        newPickupPostalCode,
		DestinationPostalCode:   newDestinationPostalCode,
		HasSit:                  newHasSit,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest:                        req,
		MoveID:                             strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	handler.SetPlanner(route.NewTestingPlanner(900))
	response := handler.Handle(patchPPMParams)

	// assert we got back the 201 response
	okResponse := response.(*ppmop.PatchPersonallyProcuredMoveOK)
	patchPPMPayload := okResponse.Payload

	suite.Equal(*patchPPMPayload.Size, newSize, "Size should have been updated.")
	suite.Equal(patchPPMPayload.WeightEstimate, newWeight, "Weight should have been updated.")

	suite.Equal(patchPPMPayload.PickupPostalCode, newPickupPostalCode, "PickupPostalCode should have been updated.")
	suite.Equal(patchPPMPayload.DestinationPostalCode, newDestinationPostalCode, "DestinationPostalCode should have been updated.")
	suite.Nil(patchPPMPayload.AdditionalPickupPostalCode, "AdditionalPickupPostalCode should have been updated to nil.")
	suite.Equal(*(*time.Time)(patchPPMPayload.PlannedMoveDate), newMoveDate, "MoveDate should have been updated.")
	suite.Nil(patchPPMPayload.DaysInStorage, "AdditionalPostalCode should have been updated to nil.")
	suite.Equal(*patchPPMPayload.Mileage, int64(900), "Mileage should have been set to 900")
}

func (suite *HandlerSuite) TestPatchPPMHandlerSetWeightLater() {
	t := suite.T()
	scenario.RunRateEngineScenario1(suite.DB())

	weight := swag.Int64(4100)

	// Date picked essentialy at random, but needs to be within TestYear
	moveDate := time.Date(testdatagen.TestYear, time.November, 10, 23, 0, 0, 0, time.UTC)

	pickupPostalCode := swag.String("32168")
	destinationPostalCode := swag.String("29401")

	move := testdatagen.MakeDefaultMove(suite.DB())

	ppm1 := models.PersonallyProcuredMove{
		MoveID:                move.ID,
		Move:                  move,
		PlannedMoveDate:       &moveDate,
		PickupPostalCode:      pickupPostalCode,
		DestinationPostalCode: destinationPostalCode,
		Status:                models.PPMStatusDRAFT,
	}
	suite.MustSave(&ppm1)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	payload := &internalmessages.PatchPersonallyProcuredMovePayload{
		WeightEstimate: weight,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest:                        req,
		MoveID:                             strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: payload,
	}

	handler := PatchPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	handler.SetPlanner(route.NewTestingPlanner(900))
	response := handler.Handle(patchPPMParams)

	// assert we got back the 201 response
	okResponse := response.(*ppmop.PatchPersonallyProcuredMoveOK)
	patchPPMPayload := okResponse.Payload

	if patchPPMPayload.WeightEstimate != weight {
		t.Error("Weight should have been updated.")
	}

	suite.Assertions.Equal(int64(900), *patchPPMPayload.Mileage)
	suite.Assertions.Equal(int64(242246), *patchPPMPayload.IncentiveEstimateMin)
	suite.Assertions.Equal(int64(267746), *patchPPMPayload.IncentiveEstimateMax)
	suite.Assertions.Nil(patchPPMPayload.EstimatedStorageReimbursement)
	suite.Assertions.Equal(int64(0), *patchPPMPayload.PlannedSitMax)
	suite.Assertions.Equal(int64(97785), *patchPPMPayload.SitMax)

	// Now check that SIT values update when days in storage is set
	hasSit := swag.Bool(true)
	daysInStorage := swag.Int64(3)
	*payload = internalmessages.PatchPersonallyProcuredMovePayload{
		HasSit:        hasSit,
		DaysInStorage: daysInStorage,
	}

	response = handler.Handle(patchPPMParams)
	// assert we got back the 201 response
	okResponse = response.(*ppmop.PatchPersonallyProcuredMoveOK)
	patchPPMPayload = okResponse.Payload

	suite.Assertions.Equal("$32.60", *patchPPMPayload.EstimatedStorageReimbursement)
	suite.Assertions.Equal(int64(3260), *patchPPMPayload.PlannedSitMax)
}

func (suite *HandlerSuite) TestPatchPPMHandlerWrongUser() {
	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")
	initialWeight := swag.Int64(1)
	newWeight := swag.Int64(5)

	// Date picked essentialy at random, but needs to be within TestYear
	newMoveDate := time.Date(testdatagen.TestYear, time.November, 10, 23, 0, 0, 0, time.UTC)
	initialMoveDate := newMoveDate.Add(-2 * 24 * time.Hour)

	user2 := testdatagen.MakeDefaultServiceMember(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())

	ppm1 := models.PersonallyProcuredMove{
		MoveID:          move.ID,
		Move:            move,
		Size:            &initialSize,
		WeightEstimate:  initialWeight,
		PlannedMoveDate: &initialMoveDate,
		Status:          models.PPMStatusDRAFT,
	}
	suite.MustSave(&ppm1)

	req := httptest.NewRequest("PATCH", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, user2)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:            &newSize,
		WeightEstimate:  newWeight,
		PlannedMoveDate: handlers.FmtDatePtr(&newMoveDate),
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest:                        req,
		MoveID:                             strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(patchPPMParams)

	suite.CheckResponseForbidden(response)
}

// TODO: no response is returned when the moveid doesn't match. How did this ever work?
func (suite *HandlerSuite) TestPatchPPMHandlerWrongMoveID() {
	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")
	initialWeight := swag.Int64(1)
	newWeight := swag.Int64(5)

	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders1 := testdatagen.MakeDefaultOrder(suite.DB())

	selectedMoveType := models.SelectedMoveTypeHHGPPM

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err, "Failed to save move")
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	move2, verrs, err := orders1.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err, "Failed to save move")
	suite.False(verrs.HasAny(), "failed to validate move")
	move2.Orders = orders1

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move2.ID,
		Move:           *move2,
		Size:           &initialSize,
		WeightEstimate: initialWeight,
		Status:         models.PPMStatusDRAFT,
	}
	suite.MustSave(&ppm1)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, orders.ServiceMember)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:           &newSize,
		WeightEstimate: newWeight,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest:                        req,
		MoveID:                             strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(patchPPMParams)
	suite.CheckResponseForbidden(response)

}

func (suite *HandlerSuite) TestPatchPPMHandlerNoMove() {
	t := suite.T()

	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")
	initialWeight := swag.Int64(1)
	newWeight := swag.Int64(5)

	move := testdatagen.MakeDefaultMove(suite.DB())

	badMoveID := uuid.Must(uuid.NewV4())

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		Move:           move,
		Size:           &initialSize,
		WeightEstimate: initialWeight,
		Status:         models.PPMStatusDRAFT,
	}
	suite.MustSave(&ppm1)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:           &newSize,
		WeightEstimate: newWeight,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest:                        req,
		MoveID:                             strfmt.UUID(badMoveID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(patchPPMParams)

	// assert we got back the badrequest response
	_, ok := response.(*ppmop.PatchPersonallyProcuredMoveBadRequest)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
}

func (suite *HandlerSuite) TestPatchPPMHandlerAdvance() {
	t := suite.T()

	initialSize := internalmessages.TShirtSize("S")
	initialWeight := swag.Int64(1)

	move := testdatagen.MakeDefaultMove(suite.DB())

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		Move:           move,
		Size:           &initialSize,
		WeightEstimate: initialWeight,
		Status:         models.PPMStatusDRAFT,
	}
	suite.MustSave(&ppm1)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	// First, create an advance
	truth := true
	var initialAmount int64
	initialAmount = 1000
	initialMethod := internalmessages.MethodOfReceiptMILPAY
	initialAdvance := internalmessages.Reimbursement{
		RequestedAmount: &initialAmount,
		MethodOfReceipt: &initialMethod,
	}

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		HasRequestedAdvance: &truth,
		Advance:             &initialAdvance,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest:                        req,
		MoveID:                             strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(patchPPMParams)

	created, ok := response.(*ppmop.PatchPersonallyProcuredMoveOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	suite.Require().Equal(internalmessages.ReimbursementStatusDRAFT, *created.Payload.Advance.Status, "expected Draft")
	suite.Require().Equal(initialAmount, *created.Payload.Advance.RequestedAmount, "expected amount to shine through.")

	// Then, update the advance
	var newAmount int64
	newAmount = 9999999
	badStatus := internalmessages.ReimbursementStatusREQUESTED
	payload.Advance.RequestedAmount = &newAmount
	payload.Advance.Status = &badStatus

	response = handler.Handle(patchPPMParams)

	// assert we got back the created response
	updated, ok := response.(*ppmop.PatchPersonallyProcuredMoveOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	suite.Require().Equal(internalmessages.ReimbursementStatusDRAFT, *updated.Payload.Advance.Status, "expected Draft still")
	suite.Require().Equal(newAmount, *updated.Payload.Advance.RequestedAmount, "expected amount to be updated")

}

func (suite *HandlerSuite) TestPatchPPMHandlerEdgeCases() {
	t := suite.T()

	initialSize := internalmessages.TShirtSize("S")
	initialWeight := swag.Int64(1)

	move := testdatagen.MakeDefaultMove(suite.DB())

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		Move:           move,
		Size:           &initialSize,
		WeightEstimate: initialWeight,
		Status:         models.PPMStatusDRAFT,
	}
	suite.MustSave(&ppm1)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	// First, try and set has_requested_advance without passing in an advance
	truth := true
	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		HasRequestedAdvance: &truth,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest:                        req,
		MoveID:                             strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(patchPPMParams)

	suite.CheckResponseBadRequest(response)

	// Then, try and create an advance without setting has requested advance
	var initialAmount int64
	initialAmount = 1000
	initialMethod := internalmessages.MethodOfReceiptMILPAY
	initialAdvance := internalmessages.Reimbursement{
		RequestedAmount: &initialAmount,
		MethodOfReceipt: &initialMethod,
	}
	payload = internalmessages.PatchPersonallyProcuredMovePayload{
		Advance: &initialAdvance,
	}

	response = handler.Handle(patchPPMParams)

	created, ok := response.(*ppmop.PatchPersonallyProcuredMoveOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	suite.Require().Equal(internalmessages.ReimbursementStatusDRAFT, *created.Payload.Advance.Status, "expected Draft")
	suite.Require().Equal(initialAmount, *created.Payload.Advance.RequestedAmount, "expected amount to shine through.")
}

func (suite *HandlerSuite) TestRequestPPMPayment() {
	t := suite.T()

	initialSize := internalmessages.TShirtSize("S")
	initialWeight := swag.Int64(1)

	move := testdatagen.MakeDefaultMove(suite.DB())

	err := move.Submit()
	if err != nil {
		t.Fatal("Should transition.")
	}
	err = move.Approve()
	if err != nil {
		t.Fatal("Should transition.")
	}
	err = move.Complete()
	if err != nil {
		t.Fatal("Should transition.")
	}

	suite.MustSave(&move)

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		Move:           move,
		Size:           &initialSize,
		WeightEstimate: initialWeight,
		Status:         models.PPMStatusDRAFT,
	}
	err = ppm1.Submit()
	if err != nil {
		t.Fatal("Should transition.")
	}
	err = ppm1.Approve()
	if err != nil {
		t.Fatal("Should transition.")
	}

	suite.MustSave(&ppm1)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	requestPaymentParams := ppmop.RequestPPMPaymentParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm1.ID.String()),
	}

	handler := RequestPPMPaymentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(requestPaymentParams)

	created, ok := response.(*ppmop.RequestPPMPaymentOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	suite.Require().Equal(internalmessages.PPMStatusPAYMENTREQUESTED, created.Payload.Status, "expected payment requested")

}

func (suite *HandlerSuite) TestRequestPPMExpenseSummaryHandler() {
	t := suite.T()
	// When: There is a move, ppm, move document and 2 expense docs
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	sm := ppm.Move.Orders.ServiceMember

	assertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   "OK",
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}

	testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions)

	req := httptest.NewRequest("GET", "/fake/path", nil)
	req = suite.AuthenticateRequest(req, sm)

	requestExpenseSumParams := ppmop.RequestPPMExpenseSummaryParams{
		HTTPRequest:              req,
		PersonallyProcuredMoveID: strfmt.UUID(ppm.ID.String()),
	}

	handler := RequestPPMExpenseSummaryHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(requestExpenseSumParams)

	expenseSummary, ok := response.(*ppmop.RequestPPMExpenseSummaryOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
	// Then: expect the following values to be equal
	suite.Assertions.Equal(internalmessages.MovingExpenseTypeCONTRACTEDEXPENSE, expenseSummary.Payload.Categories[0].Category)
	suite.Assertions.Equal(int64(5178), expenseSummary.Payload.Categories[0].PaymentMethods.GTCC)
	suite.Assertions.Equal(int64(5178), expenseSummary.Payload.Categories[0].Total)
	suite.Assertions.Equal(int64(5178), expenseSummary.Payload.GrandTotal.PaymentMethodTotals.GTCC)
	suite.Assertions.Equal(int64(5178), expenseSummary.Payload.GrandTotal.Total)
}
