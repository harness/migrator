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

	err = MigrateEntities(promptConfirm, []string{Project}, "environments", Environment)
	if err != nil {
		log.Fatal("Failed to migrate environments")
	}
	return
}
