package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func migrateSecrets(*cli.Context) (err error) {
	promptConfirm := PromptSecretDetails()
	promptConfirm = PromptOrgAndProject([]string{migrationReq.SecretScope}) || promptConfirm

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
		ids, err = GetEntityIds("secrets", migrationReq.Identifiers, migrationReq.Names)
		if err != nil {
			log.Fatal("Failed to get ids of the secrets")
		}
		if len(ids) == 0 {
			log.Fatal("No secrets found with given names/ids")
		}
	}
	log.Info("Importing the secrets....")
	CreateEntities(getReqBody(Secret, Filter{
		Type: importType,
		Ids:  ids,
	}))
	log.Info("Imported the secrets.")

	return nil
}
