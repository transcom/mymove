// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to end an asynchronous connection pertaining to file formatting
// RA: Given the functions causing the lint errors are used to end a running asynchronous connection and
// RA: the relevant code is used for generating patch notes for internal consumption, it does not present a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

// Pipelines represents multiple CircleCI pipelines
type Pipelines struct {
	Items         []Pipeline `json:"items"`
	NextPageToken string     `json:"next_page_token"`
}

// Pipeline represents a single CircleCI pipeline
type Pipeline struct {
	ID  string            `json:"id"`
	VCS map[string]string `json:"vcs"`
}

// Fetch gets all pipelines for the mymove project
func (p *Pipelines) Fetch(pageToken string) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "https://circleci.com/api/v2/project/github/transcom/mymove/pipeline", nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Circle-Token", os.Getenv("CIRCLE_TOKEN"))
	q := req.URL.Query()
	// include page token if it is there
	// cannot include both branch and token as the page token encodes the branch
	if pageToken != "" {
		q.Add("page-token", pageToken)
	} else {
		q.Add("branch", "master")
	}
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &p)
	if err != nil {
		panic(err)
	}
}

// Workflow gets the workflow information associated with a CircleCI Pipeline
func (p *Pipeline) Workflow() *Workflow {
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://circleci.com/api/v2/pipeline/%s/workflow", p.ID), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Circle-Token", os.Getenv("CIRCLE_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	workflow := &Workflow{PipelineID: p.ID}
	err = json.Unmarshal(bodyBytes, workflow)
	if err != nil {
		panic(err)
	}
	return workflow
}

// Workflows represents a collection of CircleCI workflows
type Workflows []Workflow

// Workflow represents a CircleCI workflow
type Workflow struct {
	PipelineID string
	Items      []map[string]string `json:"items"`
}

// Attributes gets the attributes of a workflow
func (w *Workflow) Attributes() map[string]string {
	return w.Items[0]
}

// Status gets the status of a workflow
func (w *Workflow) Status() string {
	return w.Attributes()["status"]
}

// Success checks if the status is "success"
func (w *Workflow) Success() bool {
	return w.Status() == "success"
}

// OnHold checks if the status is "on_hold"
func (w *Workflow) OnHold() bool {
	return w.Status() == "on_hold"
}

func main() {
	var latestDeployedCommit string
	var latestOnHoldCommit string

	var pipelines Pipelines
	pipelines.Fetch("")
	items := pipelines.Items
	// get the next page because last success might be a while ago, like on a Monday
	var pipelinesPage2 Pipelines
	pipelinesPage2.Fetch(pipelines.NextPageToken)
	items = append(items, pipelinesPage2.Items...)

	var workflows Workflows
	for _, item := range items {
		pipelineWorkflow := *item.Workflow()
		workflows = append(workflows, pipelineWorkflow)
	}

	var latestOnHoldWorkflow Workflow
	for _, workflow := range workflows {
		if workflow.OnHold() {
			latestOnHoldWorkflow = workflow
			break
		}
	}

	var latestOnHoldPipeline Pipeline
	for _, pipeline := range items {
		if pipeline.ID == latestOnHoldWorkflow.PipelineID {
			latestOnHoldPipeline = pipeline
			break
		}
	}

	latestOnHoldCommit = latestOnHoldPipeline.VCS["revision"]

	var latestSuccessfulWorkflow Workflow
	for _, workflow := range workflows {
		if workflow.Success() {
			latestSuccessfulWorkflow = workflow
			break
		}
	}

	var latestSuccessfulPipeline Pipeline
	for _, pipeline := range items {
		if pipeline.ID == latestSuccessfulWorkflow.PipelineID {
			latestSuccessfulPipeline = pipeline
			break
		}
	}

	latestDeployedCommit = latestSuccessfulPipeline.VCS["revision"]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)

	commit, _, err := ghClient.Repositories.GetCommit(ctx, "transcom", "mymove", latestDeployedCommit)

	if err != nil {
		panic(err)
	}

	commit2, _, err := ghClient.Repositories.GetCommit(ctx, "transcom", "mymove", latestOnHoldCommit)

	if err != nil {
		panic(err)
	}

	before := commit2.Commit.Author.Date.Add(1 * time.Minute).Format(time.RFC3339)
	after := commit.Commit.Author.Date.Add(1 * time.Minute).Format(time.RFC3339)

	issues, _, err := ghClient.Search.Issues(ctx, fmt.Sprintf("repo:transcom/mymove is:pr is:closed base:master merged:%s..%s sort:updated-asc", after, before), nil)
	if err != nil {
		var githubError *github.ErrorResponse

		if errors.As(err, &githubError) {
			fmt.Println("Nothing to deploy!")
		} else {
			fmt.Println(err)
		}
	}

	for _, issue := range issues.Issues {
		fmt.Printf("%s\n", *issue.Title)
		fmt.Printf("%s\n", *issue.HTMLURL)
	}
}
