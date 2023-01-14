package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strings"
)

func migrateWorkflows(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app containing the workflows -")
	}

	if len(migrationReq.WorkflowIds) == 0 {
		promptConfirm = true
		migrationReq.WorkflowIds = TextInput("Provide the workflows that you wish to import as template as comma separated values(e.g. workflow1,workflow2)")
	}

	if len(migrationReq.WorkflowScope) == 0 {
		promptConfirm = true
		migrationReq.WorkflowScope = SelectInput("Scope for workflows:", scopes, Project)
	}

	promptConfirm = PromptOrgAndProject([]string{migrationReq.WorkflowScope, migrationReq.SecretScope, migrationReq.ConnectorScope, migrationReq.TemplateScope}) || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with workflows migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := GetUrl(migrationReq.Environment, MIGRATOR, "save/v2", migrationReq.Account)
	// Migrating the app
	log.Info("Importing the workflows....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Workflow, Filter{
		WorkflowIds: strings.Split(migrationReq.WorkflowIds, ","),
		AppId:       migrationReq.AppId,
	}))
	log.Info("Imported the workflows.")

	return nil
}
