package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"os"
	"strconv"
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

func PostReq(reqUrl string, auth string, body interface{}) ([]byte, error) {
	postBody, _ := json.Marshal(body)
	requestBody := bytes.NewBuffer(postBody)
	log.WithFields(log.Fields{
		"body": string(postBody),
	}).Debug("The request body")
	req, err := http.NewRequest("POST", reqUrl, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("received non 200 response code. The response code was " + strconv.Itoa(resp.StatusCode))
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"body": string(respBody),
	}).Debug("The response body")
	return respBody, nil
}

func GetUrl(environment string, service string, path string, accountId string) string {
	return fmt.Sprintf("%s/api/ng-migration/%s?accountIdentifier=%s", urlMap[environment][service], path, accountId)
}

func MakeAPICall(url string, auth string, body interface{}) ([]byte, error) {
	resp, err := PostReq(url, auth, body)
	if err != nil {
		log.Fatalln("There was error. Exiting...", err)
	}
	return resp, err
}

func CreateEntity(url string, auth string, body RequestBody) {
	resp, err := PostReq(url, auth, body)
	if err != nil {
		log.Fatalln("There was error while migrating. Exiting...", err)
	}

	respBody := MigrationResponseBody{}
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		log.Fatalln("There was error while parsing the response from server. Exiting...")
	}
	if len(respBody.Resource.Errors) == 0 {
		return
	}
	log.Info("Here are the errors while migrating - ")
	for i := range respBody.Resource.Errors {
		e := respBody.Resource.Errors[i]
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
