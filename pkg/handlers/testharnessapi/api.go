package testharnessapi

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen/testharness"
)

type InternalServerError struct {
	// The error string
	//
	// Required: true
	Error string `json:"error"`
}

type BaseTestHarnessHandler struct {
	handlers.HandlerConfig
}

func NewDefaultBuilder(handlerConfig handlers.HandlerConfig) http.Handler {
	return handlerConfig.AuditableAppContextFromRequestBasicHandler(
		func(appCtx appcontext.AppContext, w http.ResponseWriter, r *http.Request) error {
			action := chi.URLParam(r, "action")

			response, err := testharness.Dispatch(appCtx, action)
			if err != nil {
				appCtx.Logger().Error("Testharness error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				response = InternalServerError{
					Error: err.Error(),
				}
			}

			// if the accept header starts with text/html, assume this
			// is a human using a browser and return something vaguely
			// human readable
			if strings.HasPrefix(r.Header.Get("Accept"), "text/html") {
				w.Header().Set("content-type", "text/html")
				t := template.Must(template.New("users").Parse(`
				  <html>
				  <head>
					<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css" integrity="sha384-rbsA2VBKQhggwzxH7pPCaAqO46MgnOM80zW1RWuH61DGLwZJEdK2Kadq2F9CUG65" crossorigin="anonymous">
				  </head>
				  <body class="py-4">
					<div class="container">
					  <div class="row mb-3">
						<pre>{{.}}</pre>
					  </div>
					</div> <!-- container -->
				  </body>
				  </html>
				`))
				b, err := json.MarshalIndent(response, "", "\t")
				if err != nil {
					return err
				}
				return t.Execute(w, string(b))

			}

			w.Header().Set("content-type", "application/json")
			return json.NewEncoder(w).Encode(response)
		})
}

func NewBuilderList(handlerConfig handlers.HandlerConfig) http.Handler {
	return handlerConfig.AuditableAppContextFromRequestBasicHandler(
		func(appCtx appcontext.AppContext, w http.ResponseWriter, r *http.Request) error {
			actions := testharness.Actions()
			t := template.Must(template.New("actions").Parse(`
	  <html>
	  <head>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css" integrity="sha384-rbsA2VBKQhggwzxH7pPCaAqO46MgnOM80zW1RWuH61DGLwZJEdK2Kadq2F9CUG65" crossorigin="anonymous">
	  </head>
	  <body class="py-4">
		<div class="container">
			<h1>Welcome to the Testharness data scenarios <span class="badge bg-primary text-bg-primary">New</span></h1>
		  <div class="alert alert-info">
			<p>These scenarios were introduced in <a href="https://transcom.github.io/mymove-docs/docs/adrs/use-playwright/">ADR 0076</a> and addressed further in <a href="https://transcom.github.io/mymove-docs/docs/adrs/deprecating-devseed-scenarios">ADR 0083</a>. These scenarios can be used to create Moves in the MilMove system when working locally or with Ephemeral deployments. <mark>Click any of the buttons below to create a Move data scenario.</mark></p>
		  </div>
		  <div class="row mb-3">
			<div class="col-md-8">
			{{range .}}
			<form method="post" action="/testharness/build/{{.}}">
				<button type="submit" value="{{.}}">{{.}}</button>
			</form>
			{{end}}
			</div>
		  </div>
		</div> <!-- container -->
	  </body>
	  </html>
	`))
			w.Header().Add("Content-type", "text/html")
			return t.Execute(w, actions)

		})
}
