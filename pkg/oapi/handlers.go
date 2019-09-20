package oapi

import "github.com/labstack/echo/v4"

type Server struct {
}

func (s Server) GetMoveOrderIndex(ctx echo.Context) error {
	panic("implement me")
}

func (s Server) CreateMoveOrder(ctx echo.Context) error {
	panic("implement me")
}

func (s Server) DeleteMoveOrder(ctx echo.Context, id string) error {
	panic("implement me")
}

func (s Server) GetMoveOrder(ctx echo.Context, id string) error {
	panic("implement me")
}

func (s Server) UpdateMoveOrder(ctx echo.Context, id string) error {
	panic("implement me")
}

func (s Server) UpdateMoveOrderStatus(ctx echo.Context, id string) error {
	panic("implement me")
}

func (s Server) DeleteLineItem(ctx echo.Context, moveOrderID string, lineItemID string) error {
	panic("implement me")
}

func (s Server) GetLineItemIndex(ctx echo.Context, moveOrderID string, lineItemID string) error {
	panic("implement me")
}

func (s Server) UpdateLineItem(ctx echo.Context, moveOrderID string, lineItemID string) error {
	panic("implement me")
}

func (s Server) UpdateLineItemStatus(ctx echo.Context, moveOrderID string, lineItemID string) error {
	panic("implement me")
}

func (s Server) GetPaymentRequests(ctx echo.Context) error {
	panic("implement me")
}

func (s Server) CreatePaymentRequest(ctx echo.Context) error {
	panic("implement me")
}

func (s Server) FetchPaymentRequest(ctx echo.Context, paymentRequestID string) error {
	panic("implement me")
}

func (s Server) UpdatePaymentRequest(ctx echo.Context, paymentRequestID string) error {
	panic("implement me")
}

func (s Server) UpdatePaymentRequestStatus(ctx echo.Context, paymentRequestID string) error {
	panic("implement me")
}

func NewServer() *Server {
	return &Server{}
}