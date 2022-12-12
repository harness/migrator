package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

var Version = "development"

type cliFnWrapper func(ctx *cli.Context) error

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
	Debug             bool   `survey:"debug"`
	Json              bool   `survey:"json"`
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
		migrationReq.Account = TextInput("Account that you wish to migrate:")
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
	log.WithFields(log.Fields{
		"Account":           migrationReq.Account,
		"SecretScope":       migrationReq.SecretScope,
		"ConnectorScope":    migrationReq.ConnectorScope,
		"AppID":             migrationReq.AppId,
		"OrgIdentifier":     migrationReq.OrgIdentifier,
		"ProjectIdentifier": migrationReq.ProjectIdentifier,
	}).Info("Migration details")
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

func cliWrapper(fn cliFnWrapper, ctx *cli.Context) error {
	if migrationReq.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if migrationReq.Json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	return fn(ctx)
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
	log.Info("Importing the application....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Application, Filter{
		AppId: migrationReq.AppId,
	}))
	log.Info("Imported the application.")

	return nil
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Println(cCtx.App.Version)
	}
}

func main() {
	app := &cli.App{
		Name:                 "harness-upgrade",
		Version:              Version,
		Usage:                "Upgrade Harness CD from Current Gen to Next Gen!",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "app",
				Usage: "Import an app into a existing project by providing the `appId`",
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateApp, context)
				},
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
				Usage:       "`ID` of the account that you wish to migrate",
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
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "If debug level logs need to be printed",
				Destination: &migrationReq.Debug,
			},
			&cli.BoolFlag{
				Name:        "json",
				Usage:       "If debug level logs need to be printed",
				Destination: &migrationReq.Json,
			},
		},
		Action: func(context *cli.Context) error {
			return cliWrapper(migrateAccountLevelEntities, context)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
