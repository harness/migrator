package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func migrateAccountLevelEntities(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	// Based on the scopes of entities determine the destination details
	promptConfirm = PromptOrgAndProject([]string{migrationReq.SecretScope, migrationReq.ConnectorScope}) || promptConfirm
	logMigrationDetails()

	// We confirm if they wish to proceed or not
	if promptConfirm {
		confirm := ConfirmInput("Do you wish to proceed importing all secret managers, secrets & connectors?")
		if !confirm {
			os.Exit(1)
		}
	}

	// Finally Make the API calls to create all entities
	url := GetUrl(migrationReq.Environment, "save/v2", migrationReq.Account)

	// Create Secret Managers
	log.Info("Importing all secret managers from CG to NG...")
	CreateEntity(url, migrationReq.Auth, getReqBody(SecretManager, Filter{
		Type: All,
	}))
	log.Info("Imported all secret managers.")

	// Create Secrets
	log.Info("Importing all secrets from CG to NG...")
	CreateEntity(url, migrationReq.Auth, getReqBody(Secret, Filter{
		Type: All,
	}))
	log.Info("Imported all secrets.")

	// Create Connectors
	log.Info("Importing all connectors from CG to NG....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Connector, Filter{
		Type: All,
	}))
	log.Info("Imported all connectors.")

	return nil
}
