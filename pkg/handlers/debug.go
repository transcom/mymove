package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
)

func NewDebugSessionsHandler(handlerConfig HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := handlerConfig.AppContextFromRequest(r)
		m := make(map[string]map[string]auth.Session)

		msm := handlerConfig.GetMilSessionManager()
		msessions := make(map[string]auth.Session)
		err := msm.Iterate(r.Context(), func(ctx context.Context) error {
			obj := msm.Get(ctx, "session")
			session, ok := obj.(auth.Session)
			if ok {
				token := msm.Token(ctx)
				msessions[token] = session
			}
			return nil
		})
		if err != nil {
			appCtx.Logger().Error("Error iteration mil sessions", zap.Error(err))
		}
		m[handlerConfig.AppNames().MilServername] = msessions

		osm := handlerConfig.GetOfficeSessionManager()
		osessions := make(map[string]auth.Session)
		err = osm.Iterate(r.Context(), func(ctx context.Context) error {
			obj := osm.Get(ctx, "session")
			session, ok := obj.(auth.Session)
			if ok {
				token := osm.Token(ctx)
				osessions[token] = session
			}
			return nil
		})
		if err != nil {
			appCtx.Logger().Error("Error iteration office sessions", zap.Error(err))
		}

		m[handlerConfig.AppNames().OfficeServername] = osessions

		asm := handlerConfig.GetAdminSessionManager()
		asessions := make(map[string]auth.Session)
		err = asm.Iterate(r.Context(), func(ctx context.Context) error {
			obj := asm.Get(ctx, "session")
			session, ok := obj.(auth.Session)
			if ok {
				token := asm.Token(ctx)
				asessions[token] = session
			}
			return nil
		})
		if err != nil {
			appCtx.Logger().Error("Error iteration admin sessions", zap.Error(err))
		}
		m[handlerConfig.AppNames().AdminServername] = asessions

		w.Header().Add("Content-type", "application/json")

		err = json.NewEncoder(w).Encode(m)
		if err != nil {
			appCtx.Logger().Error("redis write error", zap.Error(err))
			http.Error(w, "redis write error", http.StatusInternalServerError)
			return
		}
	}
}
