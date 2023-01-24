package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func GetAccountSummary(*cli.Context) error {
	_ = PromptEnvDetails()
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
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Processing"
	s.Start()
	for {
		time.Sleep(time.Second)
		url := GetUrlWithQueryParams(migrationReq.Environment, MIGRATOR, "discover/summary/async-result", map[string]string{
			"accountIdentifier": migrationReq.Account,
			"requestId":         reqId,
		})
		resp, err := Get(url, migrationReq.Auth)
		if err != nil {
			s.Stop()
			log.Fatal("Failed to fetch account summary")
		}
		if resp.Resource.Status == "ERROR" {
			s.Stop()
			log.Fatal("Failed to fetch account summary")
		}
		if resp.Resource.Status == "DONE" {
			s.Stop()
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
			renderMultipleSummaries(k, v.Count, []SubSummary{{
				Title: "Deployment Types",
				Data:  v.DeploymentTypeSummary,
			}, {
				Title: "Cloud Providers",
				Data:  v.CloudProviderTypeSummary,
			}})
		case Workflow:
			renderMultipleSummaries(k, v.Count, []SubSummary{{
				Title: "Workflow Types",
				Data:  v.TypeSummary,
			}, {
				Title: "Steps",
				Data:  v.StepTypeSummary,
			}})
		case Service:
			renderMultipleSummaries(k, v.Count, []SubSummary{{
				Title: "Service Types",
				Data:  v.DeploymentTypeSummary,
			}, {
				Title: "Artifact Types",
				Data:  v.ArtifactTypeSummary,
			}})
		case ApplicationManifest:
			renderMultipleSummaries(k, v.Count, []SubSummary{{
				Title: "Types",
				Data:  v.KindSummary,
			}, {
				Title: "Store",
				Data:  v.StoreSummary,
			}})
		case Connector:
			v.Count = v.Count - v.TypeSummary["STRING"]
			delete(v.TypeSummary, "STRING")
			renderSummaryWithCount(k, v.Count, v.TypeSummary)
		case Account:
		case Application:
		case Pipeline:
		case Secret:
		case Environment:
		case WorkflowExecution:
		default:
			renderSummaryWithCount(k, v.Count, v.TypeSummary)
		}
	}
}

func renderMultipleSummaries(title string, count int64, summaries []SubSummary) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	header := fmt.Sprintf("%s (%d)", title, count)
	t.AppendHeader(table.Row{header, header}, rowConfigAutoMerge)
	for _, summary := range summaries {
		if len(summary.Data) > 0 {
			var rows []table.Row
			for k, v := range summary.Data {
				rows = append(rows, table.Row{k, v})
			}
			t.SetOutputMirror(os.Stdout)
			t.AppendRow(table.Row{summary.Title, summary.Title}, rowConfigAutoMerge)
			t.AppendSeparator()
			t.AppendRows(rows)
			t.AppendSeparator()
			t.SetStyle(table.StyleLight)
			t.SetColumnConfigs([]table.ColumnConfig{
				{Number: 1, AlignHeader: text.AlignCenter},
			})
		}
	}
	t.Render()
}

func renderSummaryWithCount(title string, count int64, data map[string]int64) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	if len(data) > 0 {
		t := table.NewWriter()
		header := fmt.Sprintf("%s (%d)", title, count)
		t.AppendHeader(table.Row{header, header}, rowConfigAutoMerge)
		var rows []table.Row
		for k, v := range data {
			rows = append(rows, table.Row{k, v})
		}
		t.SetOutputMirror(os.Stdout)
		t.AppendRows(rows)
		t.AppendSeparator()
		t.SetStyle(table.StyleLight)
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AlignHeader: text.AlignCenter},
		})
		t.Render()
	}
}

func renderTable(title string, data map[string]int64) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	if len(data) > 0 {
		var rows []table.Row
		for k, v := range data {
			rows = append(rows, table.Row{k, v})
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{title, title}, rowConfigAutoMerge)
		t.AppendRows(rows)
		t.AppendSeparator()
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AlignHeader: text.AlignCenter},
		})
		t.SetStyle(table.StyleLight)
		t.Render()
	}
}

type SubSummary struct {
	Title string
	Data  map[string]int64
}
