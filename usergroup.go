package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateUserGroups(*cli.Context) (err error) {
	promptConfirm := PromptEnvDetails()
	log.Info("Importing the user groups....")
	return MigrateEntities(promptConfirm, []string{migrationReq.UserGroupScope}, "usergroups", UserGroups)
}
