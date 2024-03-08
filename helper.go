package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/AlecAivazis/survey/v2"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v3"
)

const (
	Prod        = "Prod"
	QA          = "QA"
	Dev         = "Dev"
	Prod3       = "Prod3"
	SelfManaged = "SelfManaged"
)

const (
	MigratorService = "Migrator"
	NextGenService  = "NextGen"
	TemplateService = "Template"
	PipelineService = "Pipeline"
)

var urlMap = map[string]map[string]string{
	Prod: {
		PipelineService: "https://app.harness.io/gateway/pipeline",
		TemplateService: "https://app.harness.io/gateway/template",
		MigratorService: "https://app.harness.io/gateway/ng-migration/api/ng-migration",
		NextGenService:  "https://app.harness.io/gateway/ng",
	},
	QA: {
		PipelineService: "https://qa.harness.io/gateway/pipeline",
		TemplateService: "https://qa.harness.io/gateway/template",
		MigratorService: "https://qa.harness.io/gateway/ng-migration/api/ng-migration",
		NextGenService:  "https://qa.harness.io/gateway/ng",
	},
	Dev: {
		PipelineService: "https://localhost:8181/pipeline",
		TemplateService: "https://localhost:8181/template",
		MigratorService: "https://localhost:9080/api/ng-migration",
		NextGenService:  "https://localhost:8181/ng",
	},
	Prod3: {
		PipelineService: "https://app3.harness.io/gateway/pipeline",
		TemplateService: "https://app3.harness.io/gateway/template",
		MigratorService: "https://app3.harness.io/gateway/ng-migration/api/ng-migration",
		NextGenService:  "https://app3.harness.io/gateway/ng",
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

	return fmt.Sprintf("%s/%s?%s", GetBaseUrl(environment, service), endpoint, params)
}

func GetUrl(environment string, service string, path string, accountId string) string {
	return fmt.Sprintf("%s/%s?accountIdentifier=%s", GetBaseUrl(environment, service), path, accountId)
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
	for k := range dict {
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

func ToLowerCase(input string) string {
	return strings.ToLower(input)
}

func ToSnakeCase(input string) string {
	return strcase.ToSnake(input)
}

func TrimQuotes(input string) string {
	if strings.HasPrefix(input, "\"") && strings.HasSuffix(input, "\"") {
		return input[1 : len(input)-1]
	}

	return input
}

func GenerateHarnessUIFormatIdentifier(name string) string {
	pattern := regexp.MustCompile(`<\+([^>]+)>`)
	matches := pattern.FindAllStringSubmatch(name, -1)
	preservedSubstrings := make(map[string]string)
	for _, match := range matches {
		placeholder := match[1]
		preservedKey := fmt.Sprintf("<+%s>", placeholder)
		preservedSubstrings[preservedKey] = match[0]
		name = strings.ReplaceAll(name, preservedKey, "PLACEHOLDER")
	}

	name = removeAccents(name)
	name = stripStartingChars(name)
	name = stripSpecialChars(name)
	name = strings.ReplaceAll(name, " ", "_")

	for _, preservedValue := range preservedSubstrings {
		name = strings.ReplaceAll(name, "PLACEHOLDER", preservedValue)
	}

	return name
}

func removeAccents(input string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, input)
	if err != nil {
		return ""
	}
	return output
}

func stripStartingChars(s string) string {
	return regexp.MustCompile("^[0-9-$]*").ReplaceAllString(s, "")
}

func stripSpecialChars(s string) string {
	return regexp.MustCompile("[^0-9a-zA-Z_$ ]").ReplaceAllString(s, "")
}

func ToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if s == "" {
		return s
	}
	pattern := regexp.MustCompile(`<\+([^>]+)>`)
	matches := pattern.FindAllStringSubmatch(s, -1)
	preservedSubstrings := make(map[string]string)
	for _, match := range matches {
		placeholder := match[1]
		preservedKey := fmt.Sprintf("<+%s>", placeholder)
		preservedSubstrings[preservedKey] = match[0]
		s = strings.ReplaceAll(s, preservedKey, "PLACEHOLDER")
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

	treatedString := n.String()
	for _, preservedValue := range preservedSubstrings {
		treatedString = strings.ReplaceAll(treatedString, "PLACEHOLDER", preservedValue)
	}

	return treatedString
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
	url := GetUrlWithQueryParams(migrationReq.Environment, MigratorService, entity, map[string]string{
		AccountIdentifier: migrationReq.Account,
		"appId":           migrationReq.AppId,
	})
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

func GetBaseUrl(environment string, service string) string {
	if environment == "Prod1" || environment == "Prod2" {
		environment = "Prod"
	}
	if environment != SelfManaged {
		url := urlMap[environment][service]
		if len(url) == 0 {
			log.Fatalf("invalid environment value - %s", environment)
		}
		return url
	}
	var url string
	switch service {
	case PipelineService:
		url = migrationReq.BaseUrl + "/pipeline"
	case TemplateService:
		url = migrationReq.BaseUrl + "/template"
	case NextGenService:
		url = migrationReq.BaseUrl + "/ng"
	case MigratorService:
		url = migrationReq.BaseUrl + "/ng-migration/api/ng-migration"
	default:
		panic("Unknown service! Please contact Harness support")
	}
	log.Debugf("BaseUrl for SelfManaged - %s", url)
	return url
}

func GetEntityIds(entity string, idsString string, namesString string) ([]string, error) {
	ids := Split(idsString, ",")
	if len(ids) > 0 {
		return ids, nil
	}
	names := Split(namesString, ",")
	if len(names) == 0 {
		return nil, nil
	}
	nameToIdMap, err := GetEntityNameIdMap(entity)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, item := range names {
		itemId, ok := nameToIdMap[item]
		if !ok {
			continue
		}
		result = append(result, itemId)
	}
	return result, nil
}

func GetEntityNameIdMap(entity string) (map[string]string, error) {
	items, err := listEntities(entity)
	if err != nil {
		return nil, err
	}

	var nameToIdMap = make(map[string]string)
	for _, item := range items {
		nameToIdMap[item.Name] = item.Id
	}
	return nameToIdMap, err
}

func MigrateEntities(promptConfirm bool, scopes []string, pluralValue string, entityType EntityType) (err error) {
	promptConfirm = PromptOrgAndProject(scopes) || promptConfirm
	logMigrationDetails()
	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	importType := ImportType("ALL")
	var ids []string
	if !migrationReq.All {
		importType = "SPECIFIC"
		ids, err = GetEntityIds(pluralValue, migrationReq.Identifiers, migrationReq.Names)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to get ids of the %s", pluralValue))
		}
		if len(ids) == 0 {
			log.Fatal(fmt.Sprintf("No %s found with given names/ids", pluralValue))
		}
	}
	log.Info(fmt.Sprintf("Importing the %s....", pluralValue))
	scope := AccountScope
	if len(migrationReq.AppId) > 0 {
		scope = AppScope
	}
	CreateEntities(getReqBody(entityType, Filter{
		AppId: migrationReq.AppId,
		Type:  importType,
		Ids:   ids,
		Scope: scope,
	}))
	log.Info(fmt.Sprintf("Imported the %s.", pluralValue))

	return nil
}

func LoadYamlFromFile(filePath string) map[string]string {
	filePath = strings.TrimSpace(filePath)
	if len(filePath) == 0 {
		return nil
	}
	yFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]string)
	err = yaml.Unmarshal(yFile, &data)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully loaded %d custom expressions from the file", len(data))
	return data
}

func LoadOverridesFromFile(filePath string) map[string]EntityOverrideInput {
	filePath = strings.TrimSpace(filePath)
	if len(filePath) == 0 {
		return nil
	}
	yFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var data OverrideFileData
	err = yaml.Unmarshal(yFile, &data)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully loaded %d overrides & %d settings from the file", len(data.Overrides), len(data.Settings))

	if len(data.Overrides) == 0 {
		return nil
	}

	var overrides = make(map[string]EntityOverrideInput)
	var nameToIdMap = make(map[string]map[string]string)
	for i, override := range data.Overrides {
		assertNotBlank(override.Type, fmt.Sprintf("Type cannot be blank in overrides for index - %d", i))
		assertNotAllBlank(fmt.Sprintf("Name, Identifier & Scope are blank in overrides for index - %d", i), override.Name, override.Identifier, override.Scope)
		assertAllowedValues(override.Type, []string{UserGroups, Template, Connector, Secret, Service, Environment, Workflow, Pipeline, SecretManager}, fmt.Sprintf("Only a few types of entities support overrides for index %d", i))
		if len(strings.TrimSpace(override.ID)) > 0 {
			overrides[fmt.Sprintf("CgEntityId(id=%s, type=%s)", override.ID, override.Type)] = EntityOverrideInput{
				Name:       override.Name,
				Identifier: override.Identifier,
				Scope:      override.Scope,
			}
		} else {
			assertNotBlank(override.FirstGenName, fmt.Sprintf("Both firstGen name & ID fields cannot be blank in overrides for index %d", i))
			if len(nameToIdMap[override.Type]) == 0 {
				nameToIdMap[override.Type], err = GetEntityNameIdMap(GetEndpointFromType(override.Type))
				if err != nil {
					log.Fatal(fmt.Sprintf("Failed to fetch ids from names for - %s", override.Type), err)
				}
			}
			id, ok := nameToIdMap[override.Type][override.FirstGenName]
			if !ok {
				log.Fatal(fmt.Sprintf("Failed to fetch id for name %s of type - %s", override.FirstGenName, override.Type))
			}
			overrides[fmt.Sprintf("CgEntityId(id=%s, type=%s)", id, override.Type)] = EntityOverrideInput{
				Name:       override.Name,
				Identifier: override.Identifier,
				Scope:      override.Scope,
			}
		}
	}

	return overrides
}

type ReplaceSection struct {
	Old string `yaml:"old"`
	New string `yaml:"new"`
}

func LoadCustomeStringsFromFile(filePath string) map[string]string {
	filePath = strings.TrimSpace(filePath)
	if len(filePath) == 0 {
		return nil
	}
	// Read the entire file
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Info(fmt.Sprintf("error reading file: %v", err))
	}

	// Unmarshal YAML content into a slice of maps
	var replaceSections []map[string]string
	err = yaml.Unmarshal(content, &replaceSections)
	if err != nil {
		log.Info(fmt.Sprintf("error unmarshalling YAML: %v", err))
	}

	// Create a map to store the accumulated sections
	mergedSections := make(map[string]string)

	// Process each section
	for _, replaceSection := range replaceSections {
		oldValue, oldExists := replaceSection["old"]
		newValue, newExists := replaceSection["new"]

		if oldExists && newExists {
			delete(replaceSection, "old")
			delete(replaceSection, "new")
			mergedSections[oldValue] = newValue
		} else {
			log.Info("'old' or 'new' keys not found in a section")
		}
	}
	return mergedSections
}

func GetEndpointFromType(entityType string) string {
	if entityType == UserGroups {
		return "usergroups"
	}
	return strings.ToLower(entityType + "s")
}

func LoadSettingsFromFile(filePath string) []Setting {
	filePath = strings.TrimSpace(filePath)
	if len(filePath) == 0 {
		return nil
	}
	yFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	var data OverrideFileData
	err = yaml.Unmarshal(yFile, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data.Settings
}

func assertNotBlank(value string, message string) {
	if len(strings.TrimSpace(value)) == 0 {
		log.Fatal(message)
	}
}

func assertAllowedValues(value string, allowed []string, message string) {
	if !slices.Contains(allowed, value) {
		log.Fatal(message)
	}
}

func assertNotAllBlank(message string, values ...*string) {
	for _, value := range values {
		if value != nil && len(strings.TrimSpace(*value)) > 0 {
			return
		}
	}
	log.Fatal(message)
}
