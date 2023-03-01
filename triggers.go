package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strings"
)

func migrateTriggers(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app containing the triggers -")
	}

	if len(migrationReq.WorkflowScope) == 0 {
		promptConfirm = true
		migrationReq.WorkflowScope = SelectInput("Scope for workflows:", scopes, Project)
	}

	if len(migrationReq.TriggerIds) == 0 && !migrationReq.All {
		allTriggerConfirm := ConfirmInput("No triggers provided. This defaults to migrating all triggers within the application. Do you want to proceed?")
		if !allTriggerConfirm {
			promptConfirm = true
			migrationReq.TriggerIds = TextInput("Provide the triggers that you wish to import as comma separated values(e.g. trigger1,trigger2)")
		}
	}

	promptConfirm = PromptOrgAndProject([]string{Project}) || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with triggers migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	// Migrating the triggers
	var triggerIds []string
	if len(migrationReq.TriggerIds) > 0 {
		triggerIds = strings.Split(migrationReq.TriggerIds, ",")
	}
	log.Info("Importing the triggers....")
	CreateEntities(getReqBody(Trigger, Filter{
		TriggerIds: triggerIds,
		AppId:      migrationReq.AppId,
	}))
	log.Info("Imported the triggers.")

	return nil
}
