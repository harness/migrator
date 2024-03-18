package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateApp(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if migrationReq.Platform == spinnaker {
		return migrateSpinnakerApplication()
	}

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

func migrateSpinnakerApplication() error {
	log.Info("Starting the migration of Spinnaker application")
	authMethod := authBasic
	if len(migrationReq.Cert) > 0 {
		authMethod = authx509
	}
	log.Info("Importing the application....")
	if len(migrationReq.SpinnakerHost) == 0 {
		migrationReq.SpinnakerHost = TextInput("Please provide spinnaker host : ")
	}
	if len(migrationReq.SpinnakerAppName) == 0 {
		migrationReq.SpinnakerAppName = TextInput("Please provide the Spinnaker application name : ")
	}
	logSpinnakerMigrationDetails(authMethod)
	confirm := ConfirmInput("Do you want to proceed?")
	if !confirm {
		log.Fatal("Aborting...")
	}
	if len(migrationReq.ProjectIdentifier) == 0 {
		migrationReq.ProjectIdentifier = TextInput("Name of the Project : ")
	}

	jsonBody, err := getAllPipelines(authMethod)
	if err != nil {
		return err
	}

	pipelines, err := normalizeJsonArray(jsonBody)
	if err != nil {
		return err
	}
	//first check if there are any pipeline to be migrated or not if not then don't create the project
	if len(pipelines) == 0 {
		log.Info("No pipelines found to be migrated")
		return nil
	} else {
		// Check if the project already exists for the given input project name in the given org
		projects := getProjects()
		id := findProjectIdByName(projects, migrationReq.ProjectIdentifier)

		if len(id) > 0 {
			log.Info("Project already exists with the given name")
		} else {
			log.Info("Creating project....")
			if err := createAProject(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier, formatString(migrationReq.ProjectIdentifier)); err != nil {
				log.Error(err)
			}
		}
	}

	payload := map[string][]map[string]interface{}{"pipelines": pipelines}
	_, err = createSpinnakerPipelines(payload)
	return err
}
