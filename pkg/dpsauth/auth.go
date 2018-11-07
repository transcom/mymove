package dpsauth

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// SetCookiePath is the path for this resource
const SetCookiePath = "/dps_auth/set_cookie"

// SetCookieHandler handles setting the DPS auth cookie and redirecting to DPS
type SetCookieHandler struct {
	logger *zap.Logger
}

// NewSetCookieHandler creates a new SetCookieHandler
func NewSetCookieHandler(logger *zap.Logger) SetCookieHandler {
	return SetCookieHandler{logger: logger}
}

func (h SetCookieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***In ServeHTTP***")

	loginGovID := r.URL.Query().Get("login_gov_id")
	cookieName := r.URL.Query().Get("cookie_name")
	cookie, err := LoginGovIDToCookie(loginGovID)
	if err != nil {
		h.logger.Error("Converting user ID to cookie value", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	cookie.Name = cookieName
	//cookie.Domain = ".sddc.army.mil"
	cookie.Path = "/"
	fmt.Println(cookie.String())
	w.Header().Set("Set-Cookie", cookie.String())

	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Redirect(w, r, "https://github.com", http.StatusSeeOther)
}
