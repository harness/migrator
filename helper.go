package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"os"
)

const (
	Prod  string = "Prod"
           	QA           = "QA"
	Dev          = "Dev"
	Prod3        = "Prod3"
)

const (
	MIGRATOR string = "Migrator"
	NG              = "NextGen"
)

var urlMap = map[string]map[string]string{
	Prod: {
		MIGRATOR: "https://app.harness.io/gateway/ng-migration",
		NG:       "https://app.harness.io/gateway/ng",
	},
	QA: {
		MIGRATOR: "https://qa.harness.io/gateway/ng-migration",
		NG:       "https://qa.harness.io/gateway/ng",
	},
	Dev: {
		MIGRATOR: "https://localhost:9080",
		NG:       "https://localhost:8181/ng",
	},
	Prod3: {
		MIGRATOR: "https://app3.harness.io/gateway/ng-migration",
		NG:       "https://app3.harness.io/gateway/ng",
	},
}

func TextInput(question string) string {
	var text = ""
	prompt := &survey.Input{
		Message: question,
	}
	err := survey.AskOne(prompt, &text, survey.WithValidator(survey.Required))
	if err != nil {
		log.Error(err.Error())
		os.Exit(0)
	}
	return text
}

func SelectInput(question string, options []string, defaultValue interface{}) string {
	var text = ""
	prompt := &survey.Select{
		Message: question,
		Options: options,
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &text, survey.WithValidator(survey.Required))
	if err != nil {
		log.Error(err.Error())
		os.Exit(0)
	}
	return text
}

func ConfirmInput(question string) bool {
	confirm := false
	prompt := &survey.Confirm{
		Message: question,
	}
	_ = survey.AskOne(prompt, &confirm)
	return confirm
}

func GetUrlWithQueryParams(environment string, service string, endpoint string, queryParams map[string]string) string {
	params := ""
	for k, v := range queryParams {
		params = params + k + "=" + v + "&"
	}

	return fmt.Sprintf("%s/api/ng-migration/%s?%s", urlMap[environment][service], endpoint, params)
}

func GetUrl(environment string, service string, path string, accountId string) string {
	return fmt.Sprintf("%s/api/ng-migration/%s?accountIdentifier=%s", urlMap[environment][service], path, accountId)
}

func CreateEntity(url string, auth string, body RequestBody) {
	resp, err := Post(url, auth, body)
	if err != nil {
		log.Fatalln("There was error while migrating. Exiting...", err)
	}

	if len(resp.Resource.Stats) > 0 {
		var rows []table.Row
		for k, v := range resp.Resource.Stats {
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

	if len(resp.Resource.Errors) == 0 {
		return
	}
	log.Info("Here are the errors while migrating - ")
	for i := range resp.Resource.Errors {
		e := resp.Resource.Errors[i]
		if len(e.Entity.Id) > 0 {
			log.WithFields(log.Fields{
				"type":  e.Entity.Type,
				"appId": e.Entity.AppId,
				"id":    e.Entity.Id,
				"name":  e.Entity.Name,
			}).Error(e.Message)
		} else {
			log.Error(e.Message)
		}
	}
}

func getOrDefault(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func ContainsAny[E comparable](source []E, values []E) bool {
	for i := range values {
		v := values[i]
		if slices.Contains(source, v) {
			return true
		}
	}
	return false
}
