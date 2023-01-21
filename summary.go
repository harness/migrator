package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func GetAccountSummary(*cli.Context) error {
	_ = PromptEnvDetails()
	logMigrationDetails()

	url := GetUrl(migrationReq.Environment, MIGRATOR, "discover/summary/async", migrationReq.Account)
	return handleSummary(url)
}

func GetAppSummary(*cli.Context) error {
	_ = PromptEnvDetails()
	if len(migrationReq.AppId) == 0 {
		migrationReq.AppId = TextInput("Please provide the application ID - ")
	}
	url := GetUrlWithQueryParams(migrationReq.Environment, MIGRATOR, "discover/summary/async", map[string]string{
		"accountIdentifier": migrationReq.Account,
		"appId":             migrationReq.AppId,
	})
	return handleSummary(url)
}

func handleSummary(url string) error {
	resp, err := Get(url, migrationReq.Auth)
	if err != nil {
		log.Fatal("Failed to fetch account summary")
	}
	if len(resp.Resource.RequestId) == 0 {
		log.Fatal("Failed to fetch account summary")
	}
	reqId := resp.Resource.RequestId
	for {
		time.Sleep(time.Second)
		url := GetUrlWithQueryParams(migrationReq.Environment, MIGRATOR, "discover/summary/async-result", map[string]string{
			"accountIdentifier": migrationReq.Account,
			"requestId":         reqId,
		})
		resp, err := Get(url, migrationReq.Auth)
		if err != nil {
			log.Fatal("Failed to fetch account summary")
		}
		if resp.Resource.Status == "ERROR" {
			log.Fatal("Failed to fetch account summary")
		}
		if resp.Resource.Status == "PROCESSING" {
			print("Please wait we are processing the request...")
		}
		if resp.Resource.Status == "DONE" {
			renderSummary(resp.Resource.ResponsePayload.Summary)
			break
		}
	}
	return nil
}

func renderSummary(summary map[string]EntitySummary) {
	renderTable("Summary", map[string]int64{
		"Pipelines":    summary[Pipeline].Count,
		"Applications": summary[Application].Count,
		"Secrets":      summary[Secret].Count,
		"Environments": summary[Environment].Count,
	})
	for k, v := range summary {
		switch k {
		case Infrastructure:
			renderTableWithCount(k, v.Count, v.DeploymentTypeSummary)
			renderTableWithCount(k, v.Count, v.CloudProviderTypeSummary)
		case Workflow:
			renderTableWithCount(k, v.Count, v.TypeSummary)
			renderTableWithCount(k, v.Count, v.StepTypeSummary)
		case Service:
			renderTableWithCount(k, v.Count, v.DeploymentTypeSummary)
			renderTableWithCount(k, v.Count, v.ArtifactTypeSummary)
		case ApplicationManifest:
			renderTableWithCount(k, v.Count, v.KindSummary)
			renderTableWithCount(k, v.Count, v.StoreSummary)
		case Connector:
			delete(v.TypeSummary, "STRING")
			renderTableWithCount(k, v.Count, v.TypeSummary)
		case Account:
		case Application:
		case Pipeline:
		case Secret:
		case Environment:
		case WorkflowExecution:
		default:
			renderTableWithCount(k, v.Count, v.TypeSummary)
		}
	}
}

func renderTableWithCount(title string, count int64, data map[string]int64) {
	if len(data) > 0 {
		var rows []table.Row
		for k, v := range data {
			rows = append(rows, table.Row{k, v})
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{fmt.Sprintf("%s (%d)", title, count), ""})
		t.AppendRows(rows)
		t.AppendSeparator()
		t.SetStyle(table.StyleLight)
		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})
		t.Render()
	}
}

func renderTable(title string, data map[string]int64) {
	if len(data) > 0 {
		var rows []table.Row
		for k, v := range data {
			rows = append(rows, table.Row{k, v})
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{title, ""})
		t.AppendRows(rows)
		t.AppendSeparator()
		t.SetStyle(table.StyleLight)
		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})
		t.Render()
	}
}
