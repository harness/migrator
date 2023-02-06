package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateApp(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app that you wish to import -")
	}

	if len(migrationReq.WorkflowScope) == 0 {
		promptConfirm = true
		migrationReq.WorkflowScope = SelectInput("Scope for workflows:", scopes, Project)
	}

	promptConfirm = PromptOrgAndProject([]string{Project}) || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with app migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	// Migrating the app
	log.Info("Importing the application....")
	log.Info("Importing the services, environments, infra, manifests...")
	CreateEntities(getReqBody(Application, Filter{
		AppId: migrationReq.AppId,
	}))
	if migrationReq.AllAppEntities {
		log.Info("Importing all the workflows...")
		CreateEntities(getReqBody(Workflow, Filter{
			AppId: migrationReq.AppId,
		}))
		log.Info("Importing all the pipelines...")
		CreateEntities(getReqBody(Pipeline, Filter{
			AppId: migrationReq.AppId,
		}))
	}
	log.Info("Imported the application.")

	return nil
}
