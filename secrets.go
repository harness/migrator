package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateSecrets(*cli.Context) (err error) {
	promptConfirm := PromptSecretDetails()
	err = MigrateEntities(promptConfirm, []string{migrationReq.SecretScope}, "secrets", Secret)
	if err != nil {
		log.Fatal("Failed to migrate secrets")
	}
	return
}
