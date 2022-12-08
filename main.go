package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var Version = "development"

// Note: All prompt responses will be added to this
var migrationReq = struct {
	Auth              string `survey:"auth"`
	Environment       string `survey:"environment"`
	Account           string `survey:"account"`
	SecretScope       string `survey:"secretScope"`
	ConnectorScope    string `survey:"connectorScope"`
	OrgIdentifier     string `survey:"org"`
	ProjectIdentifier string `survey:"project"`
	AppId             string `survey:"appId"`
}{}

func getReqBody(entityType EntityType, filter Filter) RequestBody {
	inputs := Inputs{
		Defaults: Defaults{
			Secret:        EntityDefaults{Scope: migrationReq.SecretScope},
			SecretManager: EntityDefaults{Scope: migrationReq.SecretScope},
			Connector:     EntityDefaults{Scope: migrationReq.ConnectorScope},
		},
	}
	destination := DestinationDetails{ProjectIdentifier: migrationReq.ProjectIdentifier, OrgIdentifier: migrationReq.OrgIdentifier}
	return RequestBody{Inputs: inputs, DestinationDetails: destination, EntityType: entityType, Filter: filter}
}

func PromptDefaultInputs() bool {
	promptConfirm := false

	// Check if auth is provided. If not provided then request for one
	migrationReq.Auth = os.Getenv("HARNESS_MIGRATOR_AUTH")
	if len(migrationReq.Auth) == 0 {
		migrationReq.Auth = TextInput("The environment variable 'HARNESS_MIGRATOR_AUTH' is not set. What is the auth token?")
	}

	if len(migrationReq.Environment) == 0 {
		promptConfirm = true
		migrationReq.Environment = SelectInput("Which environment?", []string{"Dev", "QA", "Prod"}, Dev)
	}

	if len(migrationReq.Account) == 0 {
		promptConfirm = true
		migrationReq.Account = TextInput("Account that you wish to migrateAccountLevelEntities:")
	}

	if len(migrationReq.SecretScope) == 0 {
		promptConfirm = true
		migrationReq.SecretScope = SelectInput("Scope for secrets & secret managers:", scopes, Project)
	}

	if len(migrationReq.ConnectorScope) == 0 {
		promptConfirm = true
		migrationReq.ConnectorScope = SelectInput("Scope for connectors:", scopes, Project)
	}
	return promptConfirm
}

func logMigrationDetails() {
	log.Printf("Env - %s\n"+
		"Account - %s\n"+
		"SecretScope - %s\n"+
		"ConnectorScope - %s\n"+
		"App ID - %s\n"+
		"Org Identifier - %s\n"+
		"Project Identifier - %s\n",
		migrationReq.Environment,
		migrationReq.Account,
		migrationReq.SecretScope,
		migrationReq.ConnectorScope,
		migrationReq.AppId,
		migrationReq.OrgIdentifier,
		migrationReq.ProjectIdentifier)
}

func migrateAccountLevelEntities(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()

	promptOrg := false
	promptProject := false
	// Based on the scopes of entities determine the destination details
	if migrationReq.SecretScope == Project || migrationReq.ConnectorScope == Project {
		promptOrg = len(migrationReq.OrgIdentifier) == 0
		promptProject = len(migrationReq.ProjectIdentifier) == 0
	} else if migrationReq.SecretScope == Org || migrationReq.ConnectorScope == Org {
		promptOrg = len(migrationReq.OrgIdentifier) == 0
	}

	if promptOrg {
		promptConfirm = true
		migrationReq.OrgIdentifier = TextInput("Which Org?")
	}
	if promptProject {
		promptConfirm = true
		migrationReq.ProjectIdentifier = TextInput("Which Project?")
	}

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
	log.Println("Importing all secret managers from CG to NG...")
	CreateEntity(url, migrationReq.Auth, getReqBody(SecretManager, Filter{
		Type: All,
	}))
	log.Println("Imported all secret managers.")

	// Create Secrets
	log.Println("Importing all secrets from CG to NG...")
	CreateEntity(url, migrationReq.Auth, getReqBody(Secret, Filter{
		Type: All,
	}))
	log.Println("Imported all secrets.")

	// Create Connectors
	log.Println("Importing all connectors from CG to NG....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Connector, Filter{
		Type: All,
	}))
	log.Println("Imported all connectors.")

	return nil
}

func migrateApp(ctx *cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	migrationReq.AppId = ctx.Args().Get(0)
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app that you wish to import -")
	}

	promptOrg := len(migrationReq.OrgIdentifier) == 0
	promptProject := len(migrationReq.ProjectIdentifier) == 0

	if promptOrg {
		promptConfirm = true
		migrationReq.OrgIdentifier = TextInput("Which Org?")
	}
	if promptProject {
		promptConfirm = true
		migrationReq.ProjectIdentifier = TextInput("Which Project?")
	}

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with app migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := GetUrl(migrationReq.Environment, "save/v2", migrationReq.Account)
	// Migrating the app
	log.Println("Importing the application....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Application, Filter{
		AppId: migrationReq.AppId,
	}))
	log.Println("Imported the application.")

	return nil
}

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Println(cCtx.App.Version)
	}
	app := &cli.App{
		Name:                 "harness-upgrade",
		Version:              Version,
		Usage:                "Upgrade Harness CD from Current Gen to Next Gen!",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:   "app",
				Usage:  "Import an app into a existing project by providing the `appId`",
				Action: migrateApp,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env",
				Usage:       "possible values - Prod, QA, Dev",
				Destination: &migrationReq.Environment,
			},
			&cli.StringFlag{
				Name:        "account",
				Usage:       "`ID` of the account that you wish to migrateAccountLevelEntities",
				Destination: &migrationReq.Account,
			},
			&cli.StringFlag{
				Name:        "secret",
				Usage:       "`scope` to create secrets in. Possible values - account, org, project",
				Destination: &migrationReq.SecretScope,
			},
			&cli.StringFlag{
				Name:        "connector",
				Usage:       "`scope` to create connectors in. Possible values - account, org, project",
				Destination: &migrationReq.ConnectorScope,
			},
			&cli.StringFlag{
				Name:        "org",
				Usage:       "organisation `identifier` in next gen",
				Destination: &migrationReq.OrgIdentifier,
			},
			&cli.StringFlag{
				Name:        "project",
				Usage:       "project `identifier` in next gen",
				Destination: &migrationReq.ProjectIdentifier,
			},
		},
		Action: migrateAccountLevelEntities,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
