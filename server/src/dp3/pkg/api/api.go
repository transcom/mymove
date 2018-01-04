package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	version1Mux.HandleFunc(pat.Get("/swagger.yml"), swaggerYAMLHandler)
	apiMux.Handle(pat.New("/v1/*"), version1Mux)

	return apiMux
}

// Incoming body for POST /issues
type incomingIssue struct {
	Body string `json:"body"`
}

func swaggerYAMLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WEOIFNWEOIFNWIOEFNOWINF")
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
	http.ServeFile(w, r, "./swagger.yml")
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
