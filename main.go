package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
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
	File              string `survey:"load"`
	Debug             bool   `survey:"debug"`
	Json              bool   `survey:"json"`
	AllowInsecureReq  bool   `survey:"insecure"`
	ProjectName       string `survey:"projectName"`
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
	globalFlags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "env",
			Usage:       "possible values - Prod, QA, Dev",
			Destination: &migrationReq.Environment,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "account",
			Usage:       "`ACCOUNT` that you wish to migrate",
			Destination: &migrationReq.Account,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "api-key",
			Usage:       "`API_KEY` to authenticate & authorise the migration.",
			Destination: &migrationReq.Auth,
			EnvVars:     []string{"HARNESS_MIGRATOR_AUTH"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "secret-scope",
			Usage:       "`SCOPE` to create secrets in. Possible values - account, org, project",
			Destination: &migrationReq.SecretScope,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "connector-scope",
			Usage:       "`SCOPE` to create connectors in. Possible values - account, org, project",
			Destination: &migrationReq.ConnectorScope,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "workflow-scope",
			Usage:       "`SCOPE` to create workflows in. Possible values - account, org, project",
			Destination: &migrationReq.WorkflowScope,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "template-scope",
			Usage:       "`SCOPE` to create templates in. Possible values - account, org, project",
			Destination: &migrationReq.TemplateScope,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "org",
			Usage:       "organisation `IDENTIFIER` in next gen",
			Destination: &migrationReq.OrgIdentifier,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "project",
			Usage:       "project `IDENTIFIER` in next gen",
			Destination: &migrationReq.ProjectIdentifier,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "app",
			Usage:       "`APP_ID` in current gen",
			Destination: &migrationReq.AppId,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "workflows",
			Usage:       "workflows as comma separated values `workflowId1,workflowId2`",
			Destination: &migrationReq.WorkflowIds,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "pipelines",
			Usage:       "pipelines as comma separated values `pipeline1,pipeline2`",
			Destination: &migrationReq.WorkflowIds,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "load",
			Usage:       "`FILE` to load flags from",
			Destination: &migrationReq.File,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "insecure",
			Usage:       "allow insecure API requests. This is automatically set to true if environment is Dev",
			Destination: &migrationReq.AllowInsecureReq,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "debug",
			Usage:       "print debug level logs",
			Destination: &migrationReq.Debug,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "json",
			Usage:       "log as JSON instead of standard ASCII formatter",
			Destination: &migrationReq.Json,
		}),
	}
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
			{
				Name:  "project",
				Usage: "Project specific commands like create, delete, list etc.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Usage:       "`NAME` for the project",
						Destination: &migrationReq.ProjectName,
					},
					&cli.StringFlag{
						Name:        "identifier",
						Usage:       "`IDENTIFIER` for the project",
						Destination: &migrationReq.ProjectIdentifier,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Create a project",
						Action: func(context *cli.Context) error {
							return cliWrapper(createProject, context)
						},
					},
				},
			},
		},
		Before: altsrc.InitInputSourceWithContext(globalFlags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Flags:  globalFlags,
		Action: func(context *cli.Context) error {
			return cliWrapper(migrateAccountLevelEntities, context)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
