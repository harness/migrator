package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
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

	_, err := Post(url, migrationReq.Auth, ProjectCreateBody{
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
		f, err := os.Create(path.Join(absolutePath, result.ProjectIdentifier+".yaml"))
		if err != nil {
			log.Fatal(err)
			return err
		}

		defer f.Close()
		_, err = f.Write(yamlContent)
		if err != nil {
			log.Fatal(err)
			return err
		}

		log.Infof("Application %s was exported to file %s.yaml", result.AppName, result.ProjectIdentifier)
	}

	return nil
}
