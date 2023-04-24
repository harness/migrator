package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateEnvironments(*cli.Context) (err error) {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID -")
	}

	promptConfirm = PromptOrgAndProject([]string{Project}) || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	importType := ImportType("ALL")
	var ids []string
	if !migrationReq.All {
		importType = "SPECIFIC"
		ids, err = GetEntityIds("environments", migrationReq.Identifiers, migrationReq.Names)
		if err != nil {
			log.Fatal("Failed to get ids of the environments")
		}
		if len(ids) == 0 {
			log.Fatal("No environments found with given names/ids")
		}
	}
	log.Info("Importing the environments....")
	CreateEntities(getReqBody(Environment, Filter{
		AppId: migrationReq.AppId,
		Type:  importType,
		Ids:   ids,
	}))
	log.Info("Imported the environments.")

	return nil
}
