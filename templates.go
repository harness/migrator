package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strconv"
)

func BulkRemoveTemplates(*cli.Context) error {
	promptConfirm := PromptEnvDetails()
	names := Split(migrationReq.Names, ",")
	identifiers := Split(migrationReq.Identifiers, ",")

	if migrationReq.All {
		identifiers = []string{}
		templates := getTemplates(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier, []string{})
		for _, template := range templates {
			identifiers = append(identifiers, template.Identifier)
		}
	}

	if len(names) == 0 && len(identifiers) == 0 {
		log.Fatal("No names or identifiers for the templates provided. Aborting")
	}
	if len(names) > 0 && len(identifiers) > 0 {
		log.Fatal("Both names and identifiers for the templates provided. Aborting")
	}

	n := len(identifiers)
	if len(names) > 0 {
		n = len(names)
	}
	if promptConfirm {
		confirm := ConfirmInput("Are you sure you want to proceed with deletion of " + strconv.Itoa(n) + " templates?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	if len(names) > 0 {
		templates := getTemplates(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier, []string{})
		for _, name := range names {
			id := findTemplateIdByName(templates, name)
			if len(id) > 0 {
				identifiers = append(identifiers, id)
			}
		}
		log.Debugf("Valid identifiers for the given names are - %s", identifiers)
	}

	for _, identifier := range identifiers {
		deleteTemplate(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier, identifier, migrationReq.Force)
	}
	log.Info("Finished operation for all given templates")
	return nil
}

func deleteTemplate(orgId string, projectId string, templateId string, force bool) {
	templates := getTemplates(orgId, projectId, []string{templateId})
	var versions []string
	for _, template := range templates {
		versions = append(versions, template.VersionLabel)
	}
	queryParams := map[string]string{
		AccountIdentifier: migrationReq.Account,
		"forceDelete":     strconv.FormatBool(force),
	}
	if len(orgId) > 0 {
		queryParams[OrgIdentifier] = orgId
	}
	if len(projectId) > 0 {
		queryParams[ProjectIdentifier] = projectId
	}
	url := GetUrlWithQueryParams(migrationReq.Environment, TemplateService, fmt.Sprintf("api/templates/%s", templateId), queryParams)

	log.Infof("Deleting the template with identifier %s", templateId)

	_, err := Delete(url, migrationReq.Auth, TemplateDeleteBody{TemplateVersionLabels: versions})

	if err == nil {
		log.Infof("Successfully deleted the template - %s", templateId)
	} else {
		log.Errorf("Failed to delete the template - %s", templateId)
	}
}

func getTemplates(orgId string, projectId string, templateIdentifiers []string) []TemplateDetails {
	queryParams := map[string]string{
		AccountIdentifier:  migrationReq.Account,
		"size":             "1000",
		"templateListType": "LastUpdated",
	}
	if len(orgId) > 0 {
		queryParams[OrgIdentifier] = orgId
	}
	if len(projectId) > 0 {
		queryParams[ProjectIdentifier] = projectId
	}
	url := GetUrlWithQueryParams(migrationReq.Environment, TemplateService, "api/templates/list-metadata", queryParams)
	resp, err := Post(url, migrationReq.Auth, FilterRequestBody{FilterType: TemplateService, TemplateIdentifiers: templateIdentifiers})
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to fetch templates", err)
	}
	byteData, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatal("Failed to fetch templates", err)
	}
	var templateListBody TemplateListBody
	err = json.Unmarshal(byteData, &templateListBody)
	if err != nil {
		log.Fatal("Failed to fetch templates", err)
	}

	return templateListBody.Templates
}

func findTemplateIdByName(templates []TemplateDetails, templateName string) string {
	for _, o := range templates {
		if o.Name == templateName {
			return o.Identifier
		}
	}
	return ""
}

func MigrateTemplates(*cli.Context) (err error) {
	promptConfirm := PromptDefaultInputs()
	migrationReq.All = true
	err = MigrateEntities(promptConfirm, []string{Project}, "templates", Service)
	if err != nil {
		log.Fatal("Failed to migrate templates")
	}
	return
}
