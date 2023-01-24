package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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
