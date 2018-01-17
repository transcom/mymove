package api

import (
	"encoding/json"
	"net/http"

	"github.com/markbates/pop"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"dp3/pkg/models"
)

// pkg global variable for db connection
var dbConnection *pop.Connection

// Init the API package with its database connection
func Init(dbInitialConnection *pop.Connection) {
	dbConnection = dbInitialConnection
}

// Mux creates the API router and returns it for inclusion in the app router
func Mux() *goji.Mux {

	apiMux := goji.SubMux()

	version1Mux := goji.SubMux()
	version1Mux.HandleFunc(pat.Post("/issues"), submitIssueHandler)
	version1Mux.HandleFunc(pat.Get("/issues"), indexIssueHandler)
	apiMux.Handle(pat.New("/v1/*"), version1Mux)

	return apiMux
}

// Incoming body for POST /issues
type incomingIssue struct {
	Body string `json:"body"`
}

func submitIssueHandler(w http.ResponseWriter, r *http.Request) {
	var incomingIssue incomingIssue

	if err := json.NewDecoder(r.Body).Decode(&incomingIssue); err != nil {
		zap.L().Error("Json decode", zap.Error(err))
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
	} else {
		// Create the issue in the database
		newIssue := models.Issue{
			Body: incomingIssue.Body,
		}
		if err := dbConnection.Create(&newIssue); err != nil {
			zap.L().Error("DB Insertion", zap.Error(err))
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
		} else {
			responseJSON, err := json.Marshal(newIssue)
			if err != nil {
				zap.L().Error("Encode Response", zap.Error(err))
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusCreated)
				w.Write(responseJSON)
			}
		}
	}
}

func indexIssueHandler(w http.ResponseWriter, r *http.Request) {
	// Query all issues in the db
	issues := []models.Issue{}
	if err := dbConnection.All(&issues); err != nil {
		zap.L().Error("DB Query", zap.Error(err))
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
	}

	responseJSON, err := json.Marshal(issues)
	if err != nil {
		zap.L().Error("Encode issues", zap.Error(err))
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	} else {
		w.Write(responseJSON)
	}

}
