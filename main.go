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
	WorkflowScope     string `survey:"workflowScope"`
	TemplateScope     string `survey:"templateScope"`
	OrgIdentifier     string `survey:"org"`
	ProjectIdentifier string `survey:"project"`
	AppId             string `survey:"appId"`
	WorkflowIds       string `survey:"workflowIds"`
	PipelineIds       string `survey:"pipelineIds"`
	Debug             bool   `survey:"debug"`
	Json              bool   `survey:"json"`
	AllowInsecureReq  bool   `survey:"insecure"`
}{}

func getReqBody(entityType EntityType, filter Filter) RequestBody {
	inputs := Inputs{
		Defaults: Defaults{
			Secret:        EntityDefaults{Scope: getOrDefault(migrationReq.SecretScope, Project)},
			SecretManager: EntityDefaults{Scope: getOrDefault(migrationReq.SecretScope, Project)},
			Connector:     EntityDefaults{Scope: getOrDefault(migrationReq.ConnectorScope, Project)},
			Template:      EntityDefaults{Scope: getOrDefault(migrationReq.TemplateScope, Project)},
			Workflow:      EntityDefaults{Scope: getOrDefault(migrationReq.WorkflowScope, Project)},
		},
	}
	destination := DestinationDetails{ProjectIdentifier: migrationReq.ProjectIdentifier, OrgIdentifier: migrationReq.OrgIdentifier}
	return RequestBody{Inputs: inputs, DestinationDetails: destination, EntityType: entityType, Filter: filter}
}

func logMigrationDetails() {
	log.WithFields(log.Fields{
		"Account":           migrationReq.Account,
		"SecretScope":       migrationReq.SecretScope,
		"ConnectorScope":    migrationReq.ConnectorScope,
		"TemplateScope":     migrationReq.TemplateScope,
		"AppID":             migrationReq.AppId,
		"OrgIdentifier":     migrationReq.OrgIdentifier,
		"ProjectIdentifier": migrationReq.ProjectIdentifier,
	}).Info("Migration details")
}

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

func cliWrapper(fn cliFnWrapper, ctx *cli.Context) error {
	if migrationReq.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if migrationReq.Json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	return fn(ctx)
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
		Suggest:              true,
		Commands: []*cli.Command{
			{
				Name:  "app",
				Usage: "Import an app into an existing project by providing the `appId`",
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateApp, context)
				},
			},
			{
				Name:  "workflows",
				Usage: "Import workflows as stage or pipeline templates by providing the `appId` & `workflowIds`",
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateWorkflows, context)
				},
			},
			{
				Name:  "pipelines",
				Usage: "Import pipelines into an existing project by providing the `appId` & `pipelineIds`",
				Action: func(context *cli.Context) error {
					return cliWrapper(migratePipelines, context)
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
				Usage:       "`id` of the account that you wish to migrate",
				Destination: &migrationReq.Account,
			},
			&cli.StringFlag{
				Name:        "secret-scope",
				Usage:       "`scope` to create secrets in. Possible values - account, org, project",
				Destination: &migrationReq.SecretScope,
			},
			&cli.StringFlag{
				Name:        "connector-scope",
				Usage:       "`scope` to create connectors in. Possible values - account, org, project",
				Destination: &migrationReq.ConnectorScope,
			},
			&cli.StringFlag{
				Name:        "workflow-scope",
				Usage:       "`scope` to create workflows in. Possible values - account, org, project",
				Destination: &migrationReq.WorkflowScope,
			},
			&cli.StringFlag{
				Name:        "template-scope",
				Usage:       "`scope` to create templates in. Possible values - account, org, project",
				Destination: &migrationReq.TemplateScope,
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
			&cli.StringFlag{
				Name:        "app",
				Usage:       "application `id` in current gen",
				Destination: &migrationReq.AppId,
			},
			&cli.StringFlag{
				Name:        "workflows",
				Usage:       "workflows as comma separated values `workflowId1,workflowId2`",
				Destination: &migrationReq.WorkflowIds,
			},
			&cli.StringFlag{
				Name:        "pipelines",
				Usage:       "pipelines as comma separated values `pipeline1,pipeline2`",
				Destination: &migrationReq.WorkflowIds,
			},
			&cli.BoolFlag{
				Name:        "insecure",
				Usage:       "allow insecure API requests. This is automatically set to true if environment is Dev",
				Destination: &migrationReq.AllowInsecureReq,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "print debug level logs",
				Destination: &migrationReq.Debug,
			},
			&cli.BoolFlag{
				Name:        "json",
				Usage:       "log as JSON instead of standard ASCII formatter",
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
