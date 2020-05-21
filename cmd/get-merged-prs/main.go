package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

type Pipelines struct {
	Items []Pipeline `json:"items"`
}

type Pipeline struct {
	ID  string            `json:"id"`
	VCS map[string]string `json:"vcs"`
}

func (p *Pipelines) Fetch() {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "https://circleci.com/api/v2/project/github/transcom/mymove/pipeline", nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Circle-Token", os.Getenv("CIRCLE_TOKEN"))
	q := req.URL.Query()
	q.Add("branch", "master")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bodyBytes, &p)
}

func (p *Pipeline) Workflow() *Workflow {
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://circleci.com/api/v2/pipeline/%s/workflow", p.ID), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Circle-Token", os.Getenv("CIRCLE_TOKEN"))
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	workflow := &Workflow{PipelineID: p.ID}
	json.Unmarshal(bodyBytes, workflow)
	return workflow
}

type Workflows []Workflow

type Workflow struct {
	PipelineID string
	Items      []map[string]string `json:"items"`
}

func (w *Workflow) Attributes() map[string]string {
	return w.Items[0]
}

func (w *Workflow) Status() string {
	return w.Attributes()["status"]
}

func (w *Workflow) Success() bool {
	return w.Status() == "success"
}

func (w *Workflow) OnHold() bool {
	return w.Status() == "on_hold"
}

func main() {
	var latestDeployedCommit string
	var latestOnHoldCommit string

	var pipelines Pipelines
	pipelines.Fetch()
	items := pipelines.Items
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
	commit2, _, err := ghClient.Repositories.GetCommit(ctx, "transcom", "mymove", latestOnHoldCommit)

	if err != nil {
		panic(err)
	}

	before := commit2.Commit.Author.Date.Add(1 * time.Minute).Format(time.RFC3339)
	after := commit.Commit.Author.Date.Add(1 * time.Minute).Format(time.RFC3339)

	issues, _, err := ghClient.Search.Issues(ctx, fmt.Sprintf("repo:transcom/mymove is:pr is:closed base:master merged:%s..%s sort:updated-asc", after, before), nil)
	if err != nil {
		fmt.Println(err)
	}

	for _, issue := range issues.Issues {
		fmt.Printf("%s\n", *issue.Title)
		fmt.Printf("%s\n", *issue.URL)
	}
}
