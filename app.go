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

	promptConfirm = PromptOrgAndProject([]string{Project}) || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with app migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := GetUrl(migrationReq.Environment, MIGRATOR, "save/v2", migrationReq.Account)
	// Migrating the app
	log.Info("Importing the application....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Application, Filter{
		AppId: migrationReq.AppId,
	}))
	log.Info("Imported the application.")

	return nil
}
