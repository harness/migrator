package main

import "github.com/urfave/cli/v2"

func getAccountSummary(*cli.Context) error {
	_ = PromptDefaultInputs()
	logMigrationDetails()

	_ = GetUrl(migrationReq.Environment, MIGRATOR, "save/v2", migrationReq.Account)
	return nil
}

func getAppSummary(*cli.Context) error {
	_ = PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		migrationReq.AppId = TextInput("Please provide the application ID - ")
	}

	_ = GetUrl(migrationReq.Environment, MIGRATOR, "save/v2", migrationReq.Account)
	return nil
}
