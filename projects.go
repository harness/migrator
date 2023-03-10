package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"path"
	"path/filepath"
	"strconv"
)

func createProject(*cli.Context) error {
	promptConfirm := PromptEnvDetails()
	promptConfirm = PromptOrgAndProject([]string{Org}) || promptConfirm
	if len(migrationReq.ProjectName) == 0 {
		promptConfirm = true
		migrationReq.ProjectName = TextInput("Name of the Project - ")
	}
	if len(migrationReq.ProjectIdentifier) == 0 {
		promptConfirm = true
		migrationReq.ProjectIdentifier = TextInput("Identifier for the Project - ")
	}

	log.WithFields(log.Fields{
		"Account":           migrationReq.Account,
		"OrgIdentifier":     migrationReq.OrgIdentifier,
		"ProjectName":       migrationReq.ProjectName,
		"ProjectIdentifier": migrationReq.ProjectIdentifier,
	}).Info("Project creation details")

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with project creation?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := fmt.Sprintf("%s/api/projects?accountIdentifier=%s&orgIdentifier=%s", urlMap[migrationReq.Environment][NG], migrationReq.Account, migrationReq.OrgIdentifier)

	log.Info("Creating the project....")

	_, err := Post(url, migrationReq.Auth, ProjectBody{
		Project: ProjectDetails{
			OrgIdentifier: migrationReq.OrgIdentifier,
			Identifier:    migrationReq.ProjectIdentifier,
			Name:          migrationReq.ProjectName,
			Color:         "#0063f7",
			Modules:       []string{"CD"},
			Description:   "",
		}})

	if err == nil {
		log.Info("Created the project!")
	} else {
		log.Error("There was an error creating the project")
		return err
	}

	return nil
}

func bulkCreateProject(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	promptConfirm = PromptOrgAndProject([]string{Org}) || promptConfirm

	log.WithFields(log.Fields{
		"Account":       migrationReq.Account,
		"OrgIdentifier": migrationReq.OrgIdentifier,
	}).Info("Bulk Project creation details")

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with projects creation?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := GetUrlWithQueryParams(migrationReq.Environment, MIGRATOR, "projects/bulk", map[string]string{
		"accountIdentifier": migrationReq.Account,
	})

	log.Info("Creating the projects....")

	resp, err := Post(url, migrationReq.Auth, BulkCreateBody{
		Org: migrationReq.OrgIdentifier,
	})

	if err != nil {
		log.Fatal("There was an error creating the projects")
		return err
	}

	byteData, err := json.Marshal(resp.Resource)
	if err != nil {
		return err
	}
	var data []BulkProjectResult
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return err
	}

	for _, result := range data {
		if len(result.Error.Message) > 0 {
			log.Errorf("When creating a project for application '%s' there was an error %s", result.AppName, result.Error.Message)
			continue
		}

		yamlData := map[string]string{
			"env":             migrationReq.Environment,
			"account":         migrationReq.Account,
			"api-key":         migrationReq.Auth,
			"app":             result.AppId,
			"org":             migrationReq.OrgIdentifier,
			"project":         result.ProjectIdentifier,
			"secret-scope":    migrationReq.SecretScope,
			"connector-scope": migrationReq.ConnectorScope,
			"template-scope":  migrationReq.TemplateScope,
			"workflow-scope":  Project,
		}

		yamlContent, err := yaml.Marshal(&yamlData)
		if err != nil {
			log.Fatalf("error: %v", err)
			return err
		}

		absolutePath, err := filepath.Abs(migrationReq.ExportFolderPath)
		if err != nil {
			return err
		}
		err = MkDir(absolutePath)
		if err != nil {
			return err
		}
		err = WriteToFile(path.Join(absolutePath, result.ProjectIdentifier+".yaml"), yamlContent)
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Infof("Application %s was exported to file %s.yaml", result.AppName, result.ProjectIdentifier)
	}

	return nil
}

func bulkRemoveProject(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	promptConfirm = PromptOrgAndProject([]string{Org}) || promptConfirm
	names := Split(migrationReq.Names, ",")
	identifiers := Split(migrationReq.Identifiers, ",")
	if len(names) == 0 && len(identifiers) == 0 {
		log.Fatal("No names or identifiers for the projects provided. Aborting")
	}
	if len(names) > 0 && len(identifiers) > 0 {
		log.Fatal("Both names and identifiers for the projects provided. Aborting")
	}

	n := len(identifiers)
	if len(names) > 0 {
		n = len(names)
	}
	if promptConfirm {
		confirm := ConfirmInput("Are you sure you want to proceed with deletion of " + strconv.Itoa(n) + " projects?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	if len(names) > 0 {
		projects := getProjects()
		for _, name := range names {
			id := findProjectIdByName(projects, name)
			if len(id) > 0 {
				identifiers = append(identifiers, id)
			}
		}
		log.Debugf("Valid identifers for the given names are - %s", identifiers)
	}

	for _, identifier := range identifiers {
		deleteProject(identifier)
	}
	log.Info("Finished operation for all given projects")
	return nil
}

func deleteProject(projectId string) {
	url := fmt.Sprintf("%s/api/projects/%s?accountIdentifier=%s&orgIdentifier=%s", urlMap[migrationReq.Environment][NG], projectId, migrationReq.Account, migrationReq.OrgIdentifier)

	log.Infof("Deleting the project with identifier %s", projectId)

	_, err := Delete(url, migrationReq.Auth)

	if err == nil {
		log.Infof("Successfully deleted the project - %s", projectId)
	} else {
		log.Errorf("Failed to delete the project - %s", projectId)
	}
}

func getProjects() []ProjectDetails {
	url := fmt.Sprintf("%s/api/projects?accountIdentifier=%s&orgIdentifier=%s&pageSize=1000", urlMap[migrationReq.Environment][NG], migrationReq.Account, migrationReq.OrgIdentifier)
	resp, err := Get(url, migrationReq.Auth)
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to fetch projects", err)
	}
	byteData, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatal("Failed to fetch projects", err)
	}
	var projects ProjectListBody
	err = json.Unmarshal(byteData, &projects)
	if err != nil {
		log.Fatal("Failed to fetch projects", err)
	}
	var projectDetails []ProjectDetails

	for _, p := range projects.Projects {
		projectDetails = append(projectDetails, p.Project)
	}
	return projectDetails
}

func findProjectIdByName(projects []ProjectDetails, projectName string) string {
	for _, p := range projects {
		if p.Name == projectName {
			return p.Identifier
		}
	}
	return ""
}
