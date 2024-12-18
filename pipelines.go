package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const spinnaker string = "spinnaker"
const legacy string = "harness-legacy"
const authBasic string = "basic"
const authx509 string = "x509"

func migratePipelines(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if migrationReq.Platform == spinnaker {
		return migrateSpinnakerPipelines()
	}

	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app containing the pipeline -")
	}

	if len(migrationReq.WorkflowScope) == 0 {
		promptConfirm = true
		migrationReq.WorkflowScope = SelectInput("Scope for workflow to be migrated as templates:", scopes, Project)
	}

	if len(migrationReq.PipelineIds) == 0 && !migrationReq.All {
		allPipelinesConfirm := ConfirmInput("No pipelines provided. This defaults to migrating all pipelines within the application. Do you want to proceed?")
		if !allPipelinesConfirm {
			promptConfirm = true
			migrationReq.PipelineIds = TextInput("Provide the pipelines that you wish to import as template as comma separated values(e.g. pipeline1,pipeline2)")
		}
	}

	promptConfirm = PromptOrgAndProject([]string{Project}) || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with pipeline migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	// Migrating the pipelines
	log.Info("Importing the pipelines....")
	var pipelineIds []string
	if len(migrationReq.PipelineIds) > 0 {
		pipelineIds = Split(migrationReq.PipelineIds, ",")
	}
	CreateEntities(getReqBody(Pipeline, Filter{
		PipelineIds: pipelineIds,
		AppId:       migrationReq.AppId,
	}))
	log.Info("Imported the pipelines.")

	return nil
}

func migrateSpinnakerPipelines() error {
	authMethod := authBasic
	if len(migrationReq.Cert) > 0 {
		authMethod = authx509
	}
	log.Info("Importing the application....")
	if len(migrationReq.PipelineJson) == 0 && len(migrationReq.SpinnakerHost) == 0 {
		migrationReq.SpinnakerHost = TextInput("Please provide spinnaker host : ")
	}
	if len(migrationReq.PipelineJson) == 0 && len(migrationReq.SpinnakerAppName) == 0 {
		migrationReq.SpinnakerAppName = TextInput("Please provide the Spinnaker application name : ")
	}
	if !migrationReq.All {
		if len(migrationReq.PipelineJson) == 0 {
			migrationReq.PipelineName = TextInput("Please provide the Spinnaker pipeline name : ")
		}
	}

	logSpinnakerMigrationDetails(authMethod)
	confirm := ConfirmInput("Do you want to proceed with pipeline migration?")
	if !confirm {
		log.Fatal("Aborting...")
	}

	var jsonBody []byte
	var pipelines []map[string]interface{}
	var err error

	if len(migrationReq.PipelineName) > 0 {
		jsonBody, err = getSinglePipeline(authMethod, migrationReq.PipelineName)
	} else {
		if len(migrationReq.PipelineJson) > 0 {
			// Read from file
			jsonBody, err = os.ReadFile(migrationReq.PipelineJson)
			if err != nil {
				return fmt.Errorf("failed to read pipeline JSON file: %v", err)
			}
		} else {
			jsonBody, err = getAllPipelines(authMethod, migrationReq.SpinnakerAppName)
		}
	}

	if err != nil {
		return err
	}

	pipelines, err = normalizeJsonArray(jsonBody)
	if err != nil {
		return err
	}

	if len(migrationReq.PipelineJson) == 0 {
		pipelines, err = fetchDependentPipelines(pipelines, err, authMethod)
	}
	if err != nil {
		return err
	}
	dryRun := migrationReq.DryRun
	plan := migrationReq.Plan
	payload := map[string]interface{}{
		"pipelines": pipelines, // Expecting pipelines as []map[string]interface{}
		"dryRun":    dryRun,    // dryRun as a bool
		"planOnly":  plan,
	}

	_, err = createSpinnakerPipelines(payload, dryRun)
	return err
}

func checkUnsupportedStages(payload map[string]interface{}, supportedStages []string) {
	unsupportedStagesMap := make(map[string][]string)

	// Convert supported stages to lowercase for case-insensitive comparison
	supportedMap := make(map[string]bool)
	for _, stage := range supportedStages {
		supportedMap[strings.ToLower(stage)] = true
	}

	// Regular expression to remove non-alphanumeric characters
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	// Iterate over each pipeline
	for _, pipeline := range payload["pipelines"].([]map[string]interface{}) {
		unsupportedStages := []string{}

		// Get stages from the pipeline
		if stages, ok := pipeline["stages"].([]interface{}); ok {
			for _, stage := range stages {
				// Ensure stage is a map and retrieve the type
				if stageMap, ok := stage.(map[string]interface{}); ok {
					if stageType, ok := stageMap["type"].(string); ok {
						// Remove special characters and convert to lowercase
						cleanStageType := strings.ToLower(re.ReplaceAllString(stageType, ""))

						// Check if the cleaned stage type is unsupported
						if !supportedMap[cleanStageType] {
							unsupportedStages = append(unsupportedStages, stageType) // Add original type for clarity in logs
						}
					}
				}
			}
		}

		// Add to results if there are unsupported stages
		if len(unsupportedStages) > 0 {
			if name, ok := pipeline["name"].(string); ok {
				unsupportedStagesMap[name] = unsupportedStages
			}
		}
	}

	if len(unsupportedStagesMap) > 0 {
		for pipelineID, stages := range unsupportedStagesMap {
			// Construct a single log message for each pipeline
			message := fmt.Sprintf("\n Pipeline with id: %s\n has unsupported Stages:\n", pipelineID)
			for _, stage := range stages {
				message += fmt.Sprintf("  - %s\n", stage)
			}
			log.Warn(message)
		}
	} else {
		log.Info("All stages in all pipelines are supported.")
	}

}

func fetchDependentPipelines(pipelines []map[string]interface{}, err error, authMethod string) ([]map[string]interface{}, error) {
	var pipelinesToRemove []int
	for _, pipeline := range pipelines {
		stages, ok := pipeline["stages"].([]interface{})
		if !ok {
			fmt.Println("Error: Unable to assert 'stages' to the correct type.")
			continue
		}

		for _, stage := range stages {
			s := stage.(map[string]interface{})
			stageType, okType := s["type"].(string)
			pipelineId, okId := s["pipeline"].(string)

			if okType && stageType == "pipeline" && okId {
				pipelinesToRemove, err = addDependentPipelineRecursive(pipelines, s, pipelineId, pipelinesToRemove, err, authMethod)
				log.Info(fmt.Printf("Updated stage with pipeline ID %s\n", pipelineId))
			}
		}
	}
	uniqueElementsToRemove := uniqueSliceElements(pipelinesToRemove)
	pipelines = deleteElements(pipelines, uniqueElementsToRemove)
	return pipelines, err
}

func uniqueSliceElements[T comparable](inputSlice []T) []T {
	uniqueSlice := make([]T, 0, len(inputSlice))
	seen := make(map[T]bool, len(inputSlice))
	for _, element := range inputSlice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}
	return uniqueSlice
}

func addDependentPipelineRecursive(pipelines []map[string]interface{}, pipelineStage map[string]interface{}, pipelineId string, pipelinesToRemove []int, err error, authMethod string) ([]int, error) {

	i, p, err := findPipelineIndexById(pipelines, pipelineId)

	if true {

		if p == nil {
			var appName = pipelineStage["application"].(string)
			p, err = findPipelineById(authMethod, appName, pipelineId)
		} else {
			pipelinesToRemove = append(pipelinesToRemove, i)
			//pipelines = deleteElements(pipelines, pipelinesToRemove)
		}
		stages, ok := p["stages"].([]interface{})
		if !ok {
			fmt.Println("error: unable to assert 'stages' to the correct type.")
			return nil, errors.New("unable to assert 'stages' to the correct type")
		}
		for _, stage := range stages {
			s := stage.(map[string]interface{})
			stageType, okType := s["type"].(string)
			pId, okId := s["pipeline"].(string)

			if okType && stageType == "pipeline" && okId {
				pipelineStage["dependentPipeline"] = p
				pipelinesToRemove, err := addDependentPipelineRecursive(pipelines, s, pId, pipelinesToRemove, err, authMethod)
				if err != nil {
					return nil, err
				}
				return pipelinesToRemove, nil
			} else {
				pipelineStage["dependentPipeline"] = p
				return pipelinesToRemove, nil
			}
		}

	}
	return nil, err

}
func deleteElements(slice []map[string]interface{}, indices []int) []map[string]interface{} {
	// Sort indices in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))

	for _, index := range indices {
		slice = append(slice[:index], slice[index+1:]...)
	}

	return slice
}

// / normalizeJsonArray returns an array of 1 element if body is an object, otherwise returns the existing array
func normalizeJsonArray(body []byte) ([]map[string]interface{}, error) {
	var temp interface{}
	err := json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}

	var normalizedData []map[string]interface{}

	switch v := temp.(type) {
	case map[string]interface{}:
		// If the data is a single object, wrap it in a slice
		normalizedData = append(normalizedData, v)
	case []interface{}:
		// If the data is an array, convert each element to a map[string]interface{} and append to the slice
		for _, item := range v {
			if mapItem, ok := item.(map[string]interface{}); ok {
				normalizedData = append(normalizedData, mapItem)
			} else {
				return nil, fmt.Errorf("array element is not a JSON object")
			}
		}
	default:
		return nil, fmt.Errorf("unexpected data type")
	}
	return normalizedData, nil
}

func BulkRemovePipelines(*cli.Context) error {
	promptConfirm := PromptEnvDetails()
	promptConfirm = PromptOrgAndProject([]string{Project}) || promptConfirm
	names := Split(migrationReq.Names, ",")
	identifiers := Split(migrationReq.Identifiers, ",")

	if migrationReq.All {
		identifiers = []string{}
		pipelines := getPipelines(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier)
		for _, pipeline := range pipelines {
			identifiers = append(identifiers, pipeline.Identifier)
		}
	}

	if len(names) == 0 && len(identifiers) == 0 {
		log.Fatal("No names or identifiers for the pipelines provided. Aborting")
	}
	if len(names) > 0 && len(identifiers) > 0 {
		log.Fatal("Both names and identifiers for the pipelines provided. Aborting")
	}

	n := len(identifiers)
	if len(names) > 0 {
		n = len(names)
	}
	if promptConfirm {
		confirm := ConfirmInput("Are you sure you want to proceed with deletion of " + strconv.Itoa(n) + " pipelines?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	if len(names) > 0 {
		pipelines := getPipelines(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier)
		for _, name := range names {
			id := findPipelineIdByName(pipelines, name)
			if len(id) > 0 {
				identifiers = append(identifiers, id)
			}
		}
		log.Debugf("Valid identifiers for the given names are - %s", identifiers)
	}

	for _, identifier := range identifiers {
		deletePipeline(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier, identifier)
	}
	log.Info("Finished operation for all given pipelines")
	return nil
}

func deletePipeline(orgId string, projectId string, pipelineId string) {
	queryParams := map[string]string{
		ProjectIdentifier: projectId,
		OrgIdentifier:     orgId,
		AccountIdentifier: migrationReq.Account,
	}
	url := GetUrlWithQueryParams(migrationReq.Environment, PipelineService, fmt.Sprintf("api/pipelines/%s", pipelineId), queryParams)

	log.Infof("Deleting the pipeline with identifier %s", pipelineId)

	_, err := Delete(url, migrationReq.Auth, nil)

	if err == nil {
		log.Infof("Successfully deleted the pipeline - %s", pipelineId)
	} else {
		log.Errorf("Failed to delete the pipeline - %s", pipelineId)
	}
}

func getPipelines(orgId string, projectId string) []PipelineDetails {
	queryParams := map[string]string{
		ProjectIdentifier: projectId,
		OrgIdentifier:     orgId,
		AccountIdentifier: migrationReq.Account,
		"size":            "1000",
	}
	url := GetUrlWithQueryParams(migrationReq.Environment, PipelineService, "api/pipelines/list", queryParams)
	resp, err := Post(url, migrationReq.Auth, FilterRequestBody{FilterType: "PipelineSetup"})
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to fetch pipelines", err)
	}
	byteData, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatal("Failed to fetch pipelines", err)
	}
	var pipelineListBody PipelineListBody
	err = json.Unmarshal(byteData, &pipelineListBody)
	if err != nil {
		log.Fatal("Failed to fetch pipelines", err)
	}
	return pipelineListBody.Pipelines
}

func findPipelineIdByName(pipelines []PipelineDetails, name string) string {
	for _, o := range pipelines {
		if o.Name == name {
			return o.Identifier
		}
	}
	return ""
}

func getAllPipelines(authMethod string, appName string) ([]byte, error) {
	return GetWithAuth(migrationReq.SpinnakerHost, "applications/"+appName+"/pipelineConfigs", authMethod, migrationReq.Auth64, migrationReq.Cert, migrationReq.Key, migrationReq.AllowInsecureReq)
}

// this is because there's no endpoint in gate to fetch pipeline config based on a pipeline id
func findPipelineById(authMethod string, appName string, pipelineId string) (map[string]interface{}, error) {

	var jsonBody []byte
	var err error
	var pipelines []map[string]interface{}
	jsonBody, err = getAllPipelines(authMethod, appName)
	if err != nil {
		return nil, err
	}
	pipelines, err = normalizeJsonArray(jsonBody)

	if err != nil {
		return nil, err
	}

	for _, p := range pipelines {
		if id, ok := p["id"].(string); ok && id == pipelineId {
			return p, nil
		}
	}
	return nil, errors.New("spinnaker Pipeline not found by id")
}

func findPipelineIndexById(pipelines []map[string]interface{}, pipelineId string) (int, map[string]interface{}, error) {

	for i, p := range pipelines {
		if id, ok := p["id"].(string); ok && id == pipelineId {
			return i, p, nil
		}
	}
	return -1, nil, errors.New("spinnaker Pipeline not found by id")
}

func getSinglePipeline(authMethod string, name string) ([]byte, error) {
	return GetWithAuth(migrationReq.SpinnakerHost, "applications/"+migrationReq.SpinnakerAppName+"/pipelineConfigs/"+name, authMethod, migrationReq.Auth64, migrationReq.Cert, migrationReq.Key, migrationReq.AllowInsecureReq)
}

func getSupportedStages() ([]string, error) {
	url := GetUrl(migrationReq.Environment, MigratorService, "spinnaker/pipelines/stages", migrationReq.Account)
	resp, err := Get(url, migrationReq.Auth)
	if err != nil {
		log.Warnf("failed to fetch supported stages: %v", err)
		return nil, err
	}
	byteData, err := json.Marshal(resp.Resource)
	if err != nil {
		log.Warnf("failed to parse supported stages: %v", err)
	}

	var stages []string
	if err := json.Unmarshal(byteData, &stages); err != nil {
		log.Warnf("failed to parse supported stages: %v", err)
		return nil, err
	}

	return stages, nil
}

func createSpinnakerPipelines(pipelines map[string]interface{}, dryRun bool) (reqId string, err error) {
	queryParams := map[string]string{
		ProjectIdentifier: migrationReq.ProjectIdentifier,
		OrgIdentifier:     migrationReq.OrgIdentifier,
		AccountIdentifier: migrationReq.Account,
	}
	if !dryRun && migrationReq.Environment != Dev {
		err = CheckProjectExistsAndCreate()
		if err != nil {
			return "", err
		}
	}

	j, err := json.MarshalIndent(pipelines, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal pipelines JSON: %v", err)
	}
	str := string(j)
	log.Info(str)
	url := GetUrlWithQueryParams(migrationReq.Environment, MigratorService, "spinnaker/pipelines", queryParams)
	resp, err := Post(url, migrationReq.Auth, pipelines)
	if err != nil {
		return "", fmt.Errorf("failed to post pipelines: %v", err)
	}
	resource, err := getResource(resp.Resource)
	if err != nil {
		return "", fmt.Errorf("failed to get resource: %v", err)
	}
	hasErrors := resource.Errors != nil && len(resource.Errors) > 0
	if hasErrors {
		printAllErrors(pipelines, resource.Errors)
	}
	if len(resource.RequestId) != 0 {
		reqId = resource.RequestId
		log.Infof("The request id is - %s", reqId)
	}
	if resource.SkipDetails != nil && len(resource.SkipDetails) > 0 {
		jsonData, err := json.MarshalIndent(resource.SkipDetails, "", "    ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal resource skip details JSON: %v", err)
		}
		jsonString := string(jsonData)
		log.Warnf("Entity not migrated : %s", jsonString)
	}
	// Pretty print successfullyMigratedDetails
	if len(resource.SuccessfullyMigratedDetails) > 0 {
		printCreatedEntities(resource.SuccessfullyMigratedDetails)
		if migrationReq.Plan {
			saveProjectFiles(resource.SuccessfullyMigratedDetails, migrationReq.ProjectIdentifier)
		}
	} else {
		return "", fmt.Errorf("spinnaker migration failed")
	}
	if !dryRun && migrationReq.Environment != Dev {
		reconcilePipeline(resp, queryParams)
	}

	if !dryRun {
		log.Info("Spinnaker migration completed")
	} else {
		log.Info("Note: This was a dry run of the spinnaker migration")
	}

	return reqId, nil
}

func saveProjectFiles(details []SuccessfullyMigratedDetail, projectIdentifier string) {
	// Create the project folder if it doesn't exist
	projectFolder := filepath.Join(".", projectIdentifier)
	if _, err := os.Stat(projectFolder); os.IsNotExist(err) {
		err := os.MkdirAll(projectFolder, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating directory %s: %v", projectFolder, err)
		}
		log.Printf("Created directory: %s", projectFolder)
	}

	for _, detail := range details {
		if len(detail.AdditionalInfo) > 0 {
			// Generate the filename with the full path under the project folder
			filename := filepath.Join(projectFolder, fmt.Sprintf("%s_%s.yaml", detail.NgEntityDetail.EntityType, detail.NgEntityDetail.Identifier))

			// Decode the Base64 string
			decodedInfo, err := base64.StdEncoding.DecodeString(detail.AdditionalInfo)
			if err != nil {
				log.Printf("Error decoding AdditionalInfo for %s: %v", filename, err)
				continue
			}
			// Save the decoded content to the file
			err = os.WriteFile(filename, decodedInfo, 0644)
			if err != nil {
				log.Printf("Error saving file %s: %v", filename, err)
			} else {
				log.Printf("Successfully saved file: %s", filename)
			}
		}
	}
}

func printAllErrors(pipelines map[string]interface{}, errors []UpgradeError) {
	stages, _ := getSupportedStages()
	if stages != nil {
		checkUnsupportedStages(pipelines, stages)
	}
	printResourceErrors(errors)
}

func printResourceErrors(errors []UpgradeError) error {
	for _, err := range errors {
		if !strings.Contains(err.Message, "SpinnakerStageType") {
			log.Warnf(fmt.Sprintf("  . %s\n", err.Message))
		}
	}
	return nil
}

func printCreatedEntities(resources []SuccessfullyMigratedDetail) {
	for _, detail := range resources {
		ngDetail := detail.NgEntityDetail
		log.Printf("created entity:\n  EntityType: %s\n  Identifier: %s\n  OrgIdentifier: %s\n  ProjectIdentifier: %s\n\n",
			ngDetail.EntityType, ngDetail.Identifier, ngDetail.OrgIdentifier, ngDetail.ProjectIdentifier)
	}
}
