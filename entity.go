package main

import (
	"encoding/json"
	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var skipLogs = []string{
	"already exists in the parent folder",
	"Duplicate identifier, please try again with a new identifier",
	"already exists in the account",
	"already exists in this scope",
	"] already exists  appId=",
}

func CreateEntities(body RequestBody) {
	reqId, err := QueueCreateEntity(body)
	if err != nil {
		return
	}
	PollForCompletion(reqId)
}

func QueueCreateEntity(body RequestBody) (reqId string, err error) {
	url := GetUrl(migrationReq.Environment, MigratorService, "save/async", migrationReq.Account)
	resp, err := Post(url, migrationReq.Auth, body)
	if err != nil {
		log.Fatal("Failed to create the entities", err)
	}
	resource, err := getResource(resp.Resource)
	if err != nil || len(resource.RequestId) == 0 {
		log.Fatal("Failed to create the entities", err)
		return
	}
	reqId = resource.RequestId
	log.Infof("The request id is - %s", reqId)
	return
}

func PollForCompletion(reqId string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Processing"
	s.Start()
	for {
		time.Sleep(time.Second * 10)
		url := GetUrlWithQueryParams(migrationReq.Environment, MigratorService, "save/async-result", map[string]string{
			AccountIdentifier: migrationReq.Account,
			"requestId":       reqId,
		})
		resp, err := Get(url, migrationReq.Auth)
		if err != nil {
			s.Stop()
			log.Fatal("Failed to create the entities", err)
		}
		resource, err := getResource(resp.Resource)
		if err != nil {
			s.Stop()
			log.Fatal("Failed to create the entities", err)
		}
		if resource.Status == "ERROR" {
			s.Stop()
			log.Fatal("Failed to create the entities")
		}
		if resource.Status == "DONE" {
			s.Stop()
			saveSummary, err := getSaveSummary(resource)
			if err != nil {
				log.Fatal("Failed to create the entities", err)
			}
			renderSaveSummary(saveSummary)
			break
		}
	}
}

func getSaveSummary(resource Resource) (summary SaveSummary, err error) {
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

func renderSaveSummary(saveSummary SaveSummary) {
	if len(saveSummary.Stats) > 0 {
		var rows []table.Row
		for k, v := range saveSummary.Stats {
			rows = append(rows, table.Row{k, v.SuccessfullyMigrated, v.AlreadyMigrated})
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"", "Successfully Migrated", "Already Migrated"})
		t.AppendRows(rows)
		t.AppendSeparator()
		t.SetStyle(table.StyleLight)
		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})
		t.Render()
	}

	if len(saveSummary.SkipDetails) > 0 {
		log.Info("Here are the details of entities that got skipped while migrating - ")
		for i := range saveSummary.SkipDetails {
			w := saveSummary.SkipDetails[i]
			logWithDetails(log.WarnLevel, w.Entity, w.Reason)
		}
	}

	if len(saveSummary.Errors) > 0 {
		log.Info("Here are the errors while migrating - ")
		for i := range saveSummary.Errors {
			e := saveSummary.Errors[i]
			level := log.ErrorLevel
			// log as debug if the error is in skipLogs
			for _, v := range skipLogs {
				if strings.Contains(e.Message, v) {
					level = log.DebugLevel
					break
				}
			}
			logWithDetails(level, e.Entity, e.Message)
		}
	}

	if len(saveSummary.SkippedExpressionsList) > 0 {
		var rows []table.Row
		for _, v := range saveSummary.SkippedExpressionsList {
			rows = append(rows, table.Row{v.EntityType, v.Identifier, v.OrgIdentifier, v.ProjectIdentifier, strings.Join(v.Expressions, "\n")})
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Type", "Identifier", "Org", "Project", "Expressions"})
		t.AppendRows(rows)
		t.AppendSeparator()
		t.SetStyle(table.StyleLight)
		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})
		t.Render()
	}
}

func logWithDetails(level log.Level, entity CurrentGenEntity, message string) {
	if len(entity.Id) > 0 {
		log.WithFields(log.Fields{
			"type":  entity.Type,
			"appId": entity.AppId,
			"id":    entity.Id,
			"name":  entity.Name,
		}).Log(level, message)
	} else {
		log.WithFields(log.Fields{}).Log(level, message)
	}
}
