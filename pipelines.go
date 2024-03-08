package main

import (
	"encoding/json"
	"fmt"
	"strconv"

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

	if len(migrationReq.SpinnakerHost) == 0 {
		migrationReq.SpinnakerHost = TextInput("Please provide spinnaker host")
	}
	if len(migrationReq.SpinnakerAppName) == 0 {
		migrationReq.SpinnakerAppName = TextInput("Please provide the Spinnaker application name")
	}

	if !migrationReq.All {
		migrationReq.PipelineName = TextInput("Please provide the Spinnaker pipeline name")
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
		jsonBody, err = getAllPipelines(authMethod)
	}
	if err != nil {
		return err
	}
	pipelines, err = normalizeJsonArray(jsonBody)
	if err != nil {
		return err
	}
	payload := map[string][]map[string]interface{}{"pipelines": pipelines}
	_, err = createSpinnakerPipelines(payload)
	return err
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

func getAllPipelines(authMethod string) ([]byte, error) {
	return GetWithAuth(migrationReq.SpinnakerHost, "applications/"+migrationReq.SpinnakerAppName+"/pipelineConfigs", authMethod, migrationReq.Auth64, migrationReq.Cert, migrationReq.Key, migrationReq.Insecure)
}

func getSinglePipeline(authMethod string, name string) ([]byte, error) {
	return GetWithAuth(migrationReq.SpinnakerHost, "applications/"+migrationReq.SpinnakerAppName+"/pipelineConfigs/"+name, authMethod, migrationReq.Auth64, migrationReq.Cert, migrationReq.Key, migrationReq.Insecure)
}

func createSpinnakerPipelines(pipelines interface{}) (reqId string, err error) {
	queryParams := map[string]string{
		ProjectIdentifier: migrationReq.ProjectIdentifier,
		OrgIdentifier:     migrationReq.OrgIdentifier,
		AccountIdentifier: migrationReq.Account,
	}
	url := GetUrlWithQueryParams(migrationReq.Environment, MigratorService, "spinnaker/pipelines", queryParams)
	resp, err := Post(url, migrationReq.Auth, pipelines)
	if err != nil {
		log.Fatal("Failed to create pipelines", err)
	}
	resource, err := getResource(resp.Resource)
	if err == nil && resource.Errors != nil && len(resource.Errors) > 0 {
		// Convert the data to JSON
		jsonData, err := json.MarshalIndent(resource.Errors, "", "    ")
		if err != nil {
			// Handle the error
			log.Error(err)
		}
		// Convert bytes to string and print
		jsonString := string(jsonData)
		log.Warnf(jsonString)
	}
	if len(resource.RequestId) != 0 {
		reqId = resource.RequestId
		log.Infof("The request id is - %s", reqId)
	}
	log.Info("Spinnaker migration completed")
	return
}
