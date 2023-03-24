package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
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

func MkDir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func EndsWith(str string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

func Set(str []string) []string {
	var result []string
	dict := make(map[string]bool)
	for _, val := range str {
		dict[val] = true
	}
	for k, _ := range dict {
		result = append(result, k)
	}
	return result
}

func WriteToFile(absFilePath string, content []byte) error {
	f, err := os.Create(absFilePath)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write(content)
	return err
}

func ReadFile(absFilePath string) (string, error) {
	d, err := os.ReadFile(absFilePath)
	if err != nil {
		return "", err
	}
	return string(d), err
}

func ToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if s == "" {
		return s
	}
	n := strings.Builder{}
	n.Grow(len(s))
	capNext := false
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		vIsNum := v >= '0' && v <= '9'
		if vIsNum && i == 0 {
			n.WriteByte('_')
		}
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

func Split(str string, sep string) (result []string) {
	if len(strings.TrimSpace(str)) == 0 {
		return
	}
	result = strings.Split(str, sep)
	for i, s := range result {
		result[i] = strings.TrimSpace(s)
	}
	return
}

func listEntities(entity string) (data []BaseEntityDetail, err error) {
	url := GetUrl(migrationReq.Environment, MIGRATOR, entity, migrationReq.Account)
	resp, err := Get(url, migrationReq.Auth)
	if err != nil {
		return
	}

	byteData, err := json.Marshal(resp.Resource)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return
	}
	return
}
