package trdm

import "github.com/tiaguinho/gosoap"

//go:generate mockery --name SoapCaller --outpkg trdmmocks --output ./trdmmocks
type SoapCaller interface {
	Call(m string, p gosoap.SoapParams) (res *gosoap.Response, err error)
}
