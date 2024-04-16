package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

// reconcilePipeline is a function that takes a ResponseBody and a map of query parameters.
// It processes the response to extract pipeline details and checks if reconciliation is needed.
// If reconciliation is needed, it fetches the pipeline YAML, refreshes it and updates the pipeline with the refreshed YAML.
func reconcilePipeline(resp ResponseBody, queryParams map[string]string) {
	result := extractMigratedDetails(resp)
	pipelineID := getPipelineID(result)

	if pipelineID == "" {
		log.Fatalf("Pipeline ID not found in response")
	}

	uuid, err := getPipelineUUID(pipelineID, queryParams)
	if err != nil {
		log.Fatalf("Error getting pipeline UUID: %v", err)
	}

	if reconcileNeeded, _ := checkReconcileNeeded(uuid, queryParams); reconcileNeeded {
		log.Info("Pipeline Reconciliation is needed")
		performReconciliation(pipelineID, queryParams)
	}
}

func extractMigratedDetails(resp ResponseBody) map[string]interface{} {
	var result map[string]interface{}
	jsonData, _ := json.Marshal(resp)
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		return nil
	}
	return result
}

func getPipelineID(result map[string]interface{}) string {
	var pipelineID string
	successfullyMigratedDetails := result["resource"].(map[string]interface{})["successfullyMigratedDetails"].([]interface{})
	for _, detail := range successfullyMigratedDetails {
		detailMap := detail.(map[string]interface{})
		ngEntityDetail := detailMap["ngEntityDetail"].(map[string]interface{})
		if ngEntityDetail["entityType"].(string) == "PIPELINE" {
			pipelineID = ngEntityDetail["identifier"].(string)
			break
		}
	}
	return pipelineID
}

func performReconciliation(pipelineID string, queryParams map[string]string) {
	pipelineYaml, _ := getPipelineYaml(pipelineID, queryParams)
	refreshedYaml, _ := refreshPipelineYaml(pipelineYaml, queryParams)
	success, _ := updatePipelineYaml(pipelineID, refreshedYaml, queryParams)
	if !success {
		log.Fatalf("Failed to update pipeline")
	}
	log.Info("Pipeline Reconciliation completed successfully")
}

// getPipelineUUID is a function that takes a pipeline identifier and a map of query parameters.
// It makes a request to fetch the UUID of the pipeline.
func getPipelineUUID(identifier string, queryParams map[string]string) (string, error) {
	queryParams["getDefaultFromOtherRepo"] = "true"
	url := GetUrlWithQueryParams(migrationReq.Environment, PipelineService, "api/pipelines/"+identifier+"/validate", queryParams)
	respBody, err := Post(url, migrationReq.Auth, nil)
	if err != nil {
		return "", err
	}
	uuid, ok := respBody.Data.(map[string]interface{})["uuid"].(string)
	if !ok {
		return "", errors.New("UUID not found in response")
	}
	return uuid, nil
}

// checkReconcileNeeded is a function that takes a pipeline UUID and a map of query parameters.
// It makes a request to check if reconciliation is needed for the pipeline.
func checkReconcileNeeded(uuid string, queryParams map[string]string) (bool, error) {
	url := GetUrlWithQueryParams(migrationReq.Environment, PipelineService, "api/pipelines/validate/"+uuid, queryParams)
	respBodyObj, err := Get(url, migrationReq.Auth)
	if err != nil {
		return false, err
	}
	validateResp, ok := respBodyObj.Data.(map[string]interface{})["validateTemplateReconcileResponseDTO"].(map[string]interface{})
	if !ok {
		return false, errors.New("validateTemplateReconcileResponseDTO not found in response")
	}
	reconcileNeeded, ok := validateResp["reconcileNeeded"].(bool)
	if !ok {
		return false, errors.New("reconcileNeeded not found in response")
	}
	return reconcileNeeded, nil
}

// getPipelineYaml is a function that takes a pipeline identifier and a map of query parameters.
// It makes a request to fetch the YAML of the pipeline.
func getPipelineYaml(identifier string, queryParams map[string]string) (string, error) {
	queryParams["validateAsync"] = "true"
	url := GetUrlWithQueryParams(migrationReq.Environment, PipelineService, "api/pipelines/"+identifier, queryParams)
	respBodyObj, err := Get(url, migrationReq.Auth)
	if err != nil {
		return "", err
	}
	yaml, ok := respBodyObj.Data.(map[string]interface{})["yamlPipeline"].(string)
	if !ok {
		return "", errors.New("yaml not found in response")
	}
	return yaml, nil
}

// refreshPipelineYaml is a function that takes a pipeline YAML and a map of query parameters.
// It makes a request to refresh the pipeline YAML.
func refreshPipelineYaml(yaml string, queryParams map[string]string) (string, error) {
	url := GetUrlWithQueryParams(migrationReq.Environment, TemplateService, "api/refresh-template/refreshed-yaml", queryParams)
	respBodyObj, err := Post(url, migrationReq.Auth, map[string]string{"yaml": yaml})
	if err != nil {
		return "", err
	}
	refreshedYaml, ok := respBodyObj.Data.(map[string]interface{})["refreshedYaml"].(string)
	if !ok {
		return "", errors.New("refreshedYaml not found in response")
	}
	return refreshedYaml, nil
}

// updatePipelineYaml updates the pipeline with the provided YAML content.
// It returns true if the update was successful, false otherwise.
func updatePipelineYaml(pipelineID, yamlContent string, params map[string]string) (bool, error) {
	url := GetUrlWithQueryParams(migrationReq.Environment, PipelineService, "api/pipelines/v2/"+pipelineID, params)
	respBodyObj, err := Put(url, migrationReq.Auth, strings.NewReader(yamlContent))
	if err != nil || respBodyObj.Status != "SUCCESS" {
		return false, fmt.Errorf("failed to update pipeline: %v", err)
	}
	return true, nil
}
