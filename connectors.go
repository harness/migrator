package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateConnectors(*cli.Context) error {
	promptConfirm := PromptConnectorDetails()
	promptConfirm = PromptOrgAndProject([]string{migrationReq.ConnectorScope, migrationReq.SecretScope}) || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	importType := ImportType("ALL")
	ids, err := GetEntityIds("connectors", migrationReq.Identifiers, migrationReq.Names)
	if err != nil {
		log.Fatal("Failed to get ids of the connectors")
	}
	if len(ids) > 0 {
		importType = "SPECIFIC"
	}

	log.Info("Importing the connectors....")
	CreateEntities(getReqBody(Connector, Filter{
		Type: importType,
		Ids:  ids,
	}))
	log.Info("Imported the connectors.")
	return nil
}
