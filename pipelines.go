package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strings"
)

func migratePipelines(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app containing the pipeline -")
	}

	if len(migrationReq.WorkflowScope) == 0 {
		promptConfirm = true
		migrationReq.WorkflowScope = SelectInput("Scope for workflow to be migrated as templates:", scopes, Project)
	}

	if len(migrationReq.PipelineIds) == 0 {
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

	url := GetUrl(migrationReq.Environment, MIGRATOR, "save/v2", migrationReq.Account)
	// Migrating the pipelines
	log.Info("Importing the pipelines....")
	var pipelineIds []string
	if len(migrationReq.PipelineIds) > 0 {
		pipelineIds = strings.Split(migrationReq.PipelineIds, ",")
	}
	CreateEntity(url, migrationReq.Auth, getReqBody(Pipeline, Filter{
		PipelineIds: pipelineIds,
		AppId:       migrationReq.AppId,
	}))
	log.Info("Imported the pipelines.")

	return nil
}
