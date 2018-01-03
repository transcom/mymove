package api

import (
	"encoding/json"
	"github.com/markbates/pop"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"
	"net/http"
)

// pkg global variable for db connection
var dbConnection *pop.Connection

// Create db connection
func Init(dbInitialConnection *pop.Connection) {
	dbConnection = dbInitialConnection
}

// Mux creates the API router and returns it for inclusion in the app router
func Mux() *goji.Mux {

	apiMux := goji.SubMux()

	version1Mux := goji.SubMux()
	version1Mux.HandleFunc(pat.Post("/issues"), submitIssueHandler)
	apiMux.Handle(pat.New("/v1/*"), version1Mux)

	return apiMux
}

// Incoming body for POST /issues
type incomingIssue struct {
	Body string `json:"body"`
}

// Response to POST /issues
type newIssueResponse struct {
	ID int64 `json:"id"`
}

func submitIssueHandler(w http.ResponseWriter, r *http.Request) {
	var newIssue incomingIssue

	if err := json.NewDecoder(r.Body).Decode(&newIssue); err != nil {
		zap.L().Error("Json decode", zap.Error(err))
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
	} else {
		// fmt.Println(newIssue)
		resp := newIssueResponse{1}
		responseJSON, err := json.Marshal(resp)

		if err != nil {
			zap.L().Error("Encode message", zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		} else {
			w.Write(responseJSON)
		}
	}
}
