package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateUserGroups(*cli.Context) error {
	_ = PromptEnvDetails()
	logMigrationDetails()
	log.Info("Importing the user groups....")
	CreateEntities(getReqBody(UserGroups, Filter{}))
	log.Info("Imported the user groups.")
	return nil
}
