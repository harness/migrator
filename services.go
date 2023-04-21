package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateServices(*cli.Context) error {
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
	ids, err := GetEntityIds("services", migrationReq.Identifiers, migrationReq.Names)
	if err != nil {
		log.Fatal("Failed to get ids of the services")
	}
	if len(ids) > 0 {
		importType = "SPECIFIC"
	}

	log.Info("Importing the services....")
	CreateEntities(getReqBody(Service, Filter{
		AppId: migrationReq.AppId,
		Type:  importType,
		Ids:   ids,
	}))
	log.Info("Imported the services.")

	return nil
}
