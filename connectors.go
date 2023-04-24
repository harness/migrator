package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateConnectors(*cli.Context) (err error) {
	promptConfirm := PromptConnectorDetails()
	promptConfirm = PromptOrgAndProject([]string{migrationReq.ConnectorScope, migrationReq.SecretScope}) || promptConfirm

	err = MigrateEntities(promptConfirm, []string{migrationReq.ConnectorScope, migrationReq.SecretScope}, "connectors", Connector)
	if err != nil {
		log.Fatal("Failed to migrate connectors")
	}
	return
}
