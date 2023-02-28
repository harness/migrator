package main

import (
	"encoding/json"
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
	resource, err := getResource(resp.Resource)
	if err != nil || len(resource.RequestId) == 0 {
		log.Fatal("Failed to fetch account summary")
		return err
	}
	reqId := resource.RequestId
	log.Infof("The request id is - %s", reqId)
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
		resource, err := getResource(resp.Resource)
		if err != nil {
			s.Stop()
			log.Fatal("Failed to fetch account summary")
		}
		if resource.Status == "ERROR" {
			s.Stop()
			log.Fatal("Failed to fetch account summary")
		}
		if resource.Status == "DONE" {
			s.Stop()
			summary, err := getSummary(resource)
			if err != nil {
				s.Stop()
				log.Fatal("Failed to fetch account summary")
			}
			renderSummary(summary.Summary)
			break
		}
	}
	return nil
}

func getResource(data interface{}) (resource Resource, err error) {
	byteData, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteData, &resource)
	if err != nil {
		return
	}
	return
}

func getSummary(resource Resource) (summary SummaryResponse, err error) {
	byteData, err := json.Marshal(resource.ResponsePayload)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteData, &summary)
	if err != nil {
		return
	}
	return
}

func renderSummary(summary map[string]EntitySummary) {
	if len(summary["ACCOUNT"].Name) > 0 {
		fmt.Printf("Account Name - %s\n\n", summary["ACCOUNT"].Name)
	}
	renderTable("Summary", map[string]interface{}{
		"Applications": summary[Application].Count,
		"Pipelines":    summary[Pipeline].Count,
		"Workflows":    summary[Workflow].Count,
		"Environments": summary[Environment].Count,
		"Services":     summary[Service].Count,
		"Connectors":   summary[Connector].Count,
		"Secrets":      summary[Secret].Count,
	})
	for k, v := range summary {
		switch k {
		case Infrastructure:
			renderMultipleSummaries(k, 2, v.Count, []SubSummary{{
				Title: "Deployment Types",
				Data:  mapToArray(v.DeploymentTypeSummary),
			}, {
				Title: "Cloud Providers",
				Data:  mapToArray(v.CloudProviderTypeSummary),
			}})
		case Workflow:
			renderMultipleSummaries(k, 3, v.Count, []SubSummary{{
				Title: "Workflow Types",
				Data:  mapToArrayWithDefaultStatus(v.TypeSummary, "SUPPORTED"),
			}, {
				Title: "Steps",
				Data:  toArray(v.StepsSummary),
			}})
		case Service:
			renderMultipleSummaries(k, 3, v.Count, []SubSummary{{
				Title: "Service Types",
				Data:  toArray(v.DeploymentsSummary),
			}, {
				Title: "Artifact Types",
				Data:  toArray(v.ArtifactsSummary),
			}})
		case ApplicationManifest:
			renderMultipleSummaries(k, 2, v.Count, []SubSummary{{
				Title: "Types",
				Data:  mapToArray(v.KindSummary),
			}, {
				Title: "Store",
				Data:  mapToArray(v.StoreSummary),
			}})
		case Connector:
			v.Count = v.Count - v.TypesSummary["STRING"].Count
			delete(v.TypeSummary, "STRING")
			delete(v.TypesSummary, "STRING")
			renderSummaryWithCountAndStatus(k, v.Count, v.TypesSummary)
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

func renderMultipleSummaries(title string, cols int, count int64, summaries []SubSummary) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	header := fmt.Sprintf("%s (%d)", title, count)
	t.AppendHeader(makeArray(cols, header), rowConfigAutoMerge)
	for _, summary := range summaries {
		if len(summary.Data) > 0 {
			var rows []table.Row
			for _, v := range summary.Data {
				rows = append(rows, v)
			}
			t.SetOutputMirror(os.Stdout)
			t.AppendRow(makeArray(cols, summary.Title), rowConfigAutoMerge)
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

func renderSummaryWithCountAndStatus(title string, count int64, data map[string]SummaryDetails) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	if len(data) > 0 {
		t := table.NewWriter()
		header := fmt.Sprintf("%s (%d)", title, count)
		t.AppendHeader(table.Row{header, header, header}, rowConfigAutoMerge)
		var rows []table.Row
		for k, v := range data {
			rows = append(rows, table.Row{k, v.Count, v.Status})
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

func renderTable(title string, data map[string]interface{}) {
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

func mapToArray(dict map[string]int64) [][]interface{} {
	var data [][]interface{}
	for k, v := range dict {
		data = append(data, []interface{}{k, v})
	}
	return data
}

func mapToArrayWithDefaultStatus(dict map[string]int64, status string) [][]interface{} {
	var data [][]interface{}
	for k, v := range dict {
		data = append(data, []interface{}{k, v, status})
	}
	return data
}

type SubSummary struct {
	Title string
	Data  [][]interface{}
}

func makeArray(size int, defaultValue interface{}) []interface{} {
	var data []interface{}
	for i := 0; i < size; i++ {
		data = append(data, defaultValue)
	}
	return data
}

func toArray(summaries map[string]SummaryDetails) (data [][]interface{}) {
	for k, v := range summaries {
		data = append(data, []interface{}{k, v.Count, v.Status})
	}
	return
}
