package models

import "github.com/tiaguinho/gosoap"

type SoapCaller interface {
	Call(m string, p gosoap.Params) (res *gosoap.Response, err error)
}
