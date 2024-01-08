package main

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var Version = "development"

type cliFnWrapper func(ctx *cli.Context) error

// Note: All prompt responses will be added to this
var migrationReq = struct {
	Auth                  string `survey:"auth"`
	Environment           string `survey:"environment"`
	Account               string `survey:"account"`
	SecretScope           string `survey:"secretScope"`
	ConnectorScope        string `survey:"connectorScope"`
	WorkflowScope         string `survey:"workflowScope"`
	PipelineScope         string `survey:"pipelineScope"`
	TemplateScope         string `survey:"templateScope"`
	EnvironmentScope      string `survey:"environmentScope"`
	UserGroupScope        string `survey:"userGroupScope"`
	OrgIdentifier         string `survey:"org"`
	ProjectIdentifier     string `survey:"project"`
	AppId                 string `survey:"appId"`
	AllAppEntities        bool   `survey:"all"`
	WorkflowIds           string `survey:"workflowIds"`
	PipelineIds           string `survey:"pipelineIds"`
	TriggerIds            string `survey:"triggerIds"`
	File                  string `survey:"load"`
	IdentifierCase        string `survey:"identifierCase"`
	LogLevel              string `survey:"logLevel"`
	Json                  bool   `survey:"json"`
	AllowInsecureReq      bool   `survey:"insecure"`
	ProjectName           string `survey:"projectName"`
	OrgName               string `survey:"orgName"`
	UrlNG                 string `survey:"urlNG"`
	UrlCG                 string `survey:"urlCG"`
	DryRun                bool   `survey:"dryRun"`
	FileExtensions        string `survey:"fileExtensions"`
	CustomExpressionsFile string `survey:"customExpressionsFile"`
	CustomStringsFile     string `survey:"customStringsFile"`
	OverrideFile          string `survey:"overrideFile"`
	ExportFolderPath      string `survey:"export"`
	CsvFile               string `survey:"csv"`
	Names                 string `survey:"names"`
	Identifiers           string `survey:"identifiers"`
	All                   bool   `survey:"all"`
	AsPipelines           bool   `survey:"asPipelines"`
	TargetAccount         string `survey:"targetAccount"`
	TargetAuthToken       string `survey:"targetAuth"`
	BaseUrl               string `survey:"baseUrl"`
	TargetGatewayUrl      string `survey:"targetGatewayUrl"`
	Force                 bool   `survey:"force"`
	Flags                 string `survey:"flags"`
}{}

func getReqBody(entityType EntityType, filter Filter) RequestBody {
	inputs := Inputs{
		Overrides:   LoadOverridesFromFile(migrationReq.OverrideFile),
		Replace:     LoadCustomeStringsFromFile(migrationReq.CustomStringsFile),
		Expressions: LoadYamlFromFile(migrationReq.CustomExpressionsFile),
		Settings:    LoadSettingsFromFile(migrationReq.OverrideFile),
		Defaults: Defaults{
			Secret:                EntityDefaults{Scope: getOrDefault(migrationReq.SecretScope, Project)},
			SecretManager:         EntityDefaults{Scope: getOrDefault(migrationReq.SecretScope, Project)},
			SecretManagerTemplate: EntityDefaults{Scope: getOrDefault(migrationReq.SecretScope, Project)},
			Connector:             EntityDefaults{Scope: getOrDefault(migrationReq.ConnectorScope, Project)},
			Template:              EntityDefaults{Scope: getOrDefault(migrationReq.TemplateScope, Project)},
			Environment:           EntityDefaults{Scope: getOrDefault(migrationReq.EnvironmentScope, Project)},
			Workflow:              EntityDefaults{Scope: getOrDefault(migrationReq.WorkflowScope, Project), WorkflowAsPipeline: migrationReq.AsPipelines},
			UserGroup:             EntityDefaults{Scope: getOrDefault(migrationReq.UserGroupScope, Account)},
		},
	}

	var flags []string
	if len(migrationReq.Flags) > 0 {
		flags = Split(migrationReq.Flags, ",")
	}
	flags = addIfNotExists(flags, "HELM_INFRA_WITH_STAGE_VAR")

	destination := DestinationDetails{
		ProjectIdentifier: migrationReq.ProjectIdentifier,
		OrgIdentifier:     migrationReq.OrgIdentifier,
		AccountIdentifier: migrationReq.TargetAccount,
		AuthToken:         migrationReq.TargetAuthToken,
		GatewayUrl:        migrationReq.TargetGatewayUrl,
	}
	body := RequestBody{Inputs: inputs, DestinationDetails: destination, EntityType: entityType, Filter: filter, IdentifierCaseFormat: migrationReq.IdentifierCase,
		Flags: flags}

	b, err := json.Marshal(body)
	if err != nil {
		log.Debug(err)
	}
	log.Debugf("Request details: %s", b)
	return body
}

func addIfNotExists(arr []string, strToAdd string) []string {
	for _, str := range arr {
		if str == strToAdd {
			return arr
		}
	}
	return append(arr, strToAdd)
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
		"Flags":             migrationReq.Flags,
	}).Info("Migration details")
}

func cliWrapper(fn cliFnWrapper, ctx *cli.Context) error {
	if len(migrationReq.LogLevel) > 0 {
		level, err := log.ParseLevel(migrationReq.LogLevel)
		if err != nil {
			log.Fatal("Invalid log level")
		}
		log.SetLevel(level)
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
	CheckGithubForReleases()
	globalFlags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "env",
			Usage:       "possible values - Prod, QA, Dev",
			Destination: &migrationReq.Environment,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "base-url",
			Usage:       "provide the `BASE_URL` for self managed platforms",
			Destination: &migrationReq.BaseUrl,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "destination-project",
			Usage:       "`destination-project` destination project url in next gen",
			Destination: &migrationReq.UrlNG,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "source-app",
			Usage:       "`source-app` source application url in current gen",
			Destination: &migrationReq.UrlCG,
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
			Name:        "environment-scope",
			Usage:       "`SCOPE` to create environment in. Possible values - account, org, project",
			Destination: &migrationReq.EnvironmentScope,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "user-group-scope",
			Usage:       "`SCOPE` to create user groups in. Possible values - account, org, project",
			Destination: &migrationReq.UserGroupScope,
			DefaultText: Account,
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
			Name:        "load",
			Usage:       "`FILE` to load flags from",
			Destination: &migrationReq.File,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "insecure",
			Usage:       "allow insecure API requests. This is automatically set to true if environment is Dev",
			Destination: &migrationReq.AllowInsecureReq,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "log-level",
			Usage:       "set the log level. Possible values - trace, debug, info, warn, error, fatal, panic. Default is `info`",
			Destination: &migrationReq.LogLevel,
			DefaultText: "info",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "json",
			Usage:       "log as JSON instead of standard ASCII formatter",
			Destination: &migrationReq.Json,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "identifier-format",
			Usage:       "`FORMAT` to use for generation of identifiers. Supported values as CAMEL_CASE & LOWER_CASE",
			Destination: &migrationReq.IdentifierCase,
			Value:       "CAMEL_CASE",
			DefaultText: "CAMEL_CASE",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "target-account",
			Usage:       "destination `ACCOUNT` that you wish to migrate to",
			Destination: &migrationReq.TargetAccount,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "target-api-key",
			Usage:       "`API_KEY` for the target account to authenticate & authorise the migration.",
			Destination: &migrationReq.TargetAuthToken,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "target-gateway-url",
			Usage:       "destination gateway `URL`. For Prod1 & Prod2, use https://app.harness.io/gateway, for Prod3 use https://app3.harness.io/gateway",
			Destination: &migrationReq.TargetGatewayUrl,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "custom-expressions",
			Usage:       "provide a `FILE` to load custom expressions from",
			Destination: &migrationReq.CustomExpressionsFile,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "replace",
			Usage:       "provide a `FILE` to load strings from",
			Destination: &migrationReq.CustomStringsFile,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "override",
			Usage:       "provide a `FILE` to load overrides",
			Destination: &migrationReq.OverrideFile,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "flags",
			Usage:       "provide a list of flags for custom logic",
			Destination: &migrationReq.Flags,
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
				Name:    "update",
				Aliases: []string{"upgrade"},
				Usage:   "Check for updates and upgrade the CLI",
				Action: func(context *cli.Context) error {
					return cliWrapper(Update, context)
				},
			},
			{
				Name:  "account-summary",
				Usage: "Get a summary of an account",
				Action: func(context *cli.Context) error {
					return cliWrapper(GetAccountSummary, context)
				},
			},
			{
				Name:  "application-summary",
				Usage: "Get a summary of an app",
				Action: func(context *cli.Context) error {
					return cliWrapper(GetAppSummary, context)
				},
			},
			{
				Name:  "user-groups",
				Usage: "Import user groups from First Gen to Next Gen",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if all user groups need to be migrated",
						Destination: &migrationReq.All,
					},
					&cli.StringFlag{
						Name:        "ids",
						Usage:       "`IDs` of the user groups",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the user groups",
						Destination: &migrationReq.Names,
					},
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateUserGroups, context)
				},
			},
			{
				Name:  "account",
				Usage: "Import secrets managers, secrets, connectors. This will not migrate services, environments, triggers, pipelines etc",
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateAccountLevelEntities, context)
				},
			},
			{
				Name:  "app",
				Usage: "Import an app into an existing project by providing the `appId`",
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateApp, context)
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if set will migrate all workflows & pipelines",
						Destination: &migrationReq.AllAppEntities,
					},
				},
			},
			{
				Name:    "service",
				Aliases: []string{"services"},
				Usage:   "Import services into an existing project from an application",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if all services in the app need to be migrated",
						Destination: &migrationReq.All,
					},
					&cli.StringFlag{
						Name:        "ids",
						Usage:       "`IDs` of the services",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the services",
						Destination: &migrationReq.Names,
					},
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateServices, context)
				},
			},
			{
				Name:  "secrets",
				Usage: "Import secrets",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if all secrets in the account need to be migrated",
						Destination: &migrationReq.All,
					},
					&cli.StringFlag{
						Name:        "ids",
						Usage:       "`IDs` of the secrets",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the secrets",
						Destination: &migrationReq.Names,
					},
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateSecrets, context)
				},
			},
			{
				Name:  "environments",
				Usage: "Import environments into an existing project from an application",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if all environments in the app need to be migrated",
						Destination: &migrationReq.All,
					},
					&cli.StringFlag{
						Name:        "ids",
						Usage:       "`IDs` of the environments",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the environments",
						Destination: &migrationReq.Names,
					},
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateEnvironments, context)
				},
			},
			{
				Name:  "connectors",
				Usage: "Import connectors",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if all connectors in the account need to be migrated",
						Destination: &migrationReq.All,
					},
					&cli.StringFlag{
						Name:        "ids",
						Usage:       "`IDs` of the connectors",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the connectors",
						Destination: &migrationReq.Names,
					},
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateConnectors, context)
				},
			},
			{
				Name:  "workflows",
				Usage: "Import workflows as stage or pipeline templates by providing the `appId` & `workflowIds`",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if all workflows in the app need to be migrated",
						Destination: &migrationReq.All,
					},
					&cli.BoolFlag{
						Name:        "as-pipelines",
						Usage:       "create a pipeline for the workflows, this will create stage templates where possible & reuse the template to create the pipeline",
						Destination: &migrationReq.AsPipelines,
					},
					altsrc.NewStringFlag(&cli.StringFlag{
						Name:        "workflows",
						Usage:       "workflows as comma separated values `workflowId1,workflowId2`",
						Destination: &migrationReq.WorkflowIds,
					}),
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateWorkflows, context)
				},
			},
			{
				Name:  "pipelines",
				Usage: "Import pipelines into an existing project by providing the `appId` & `pipelineIds`",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "all pipelines",
						Destination: &migrationReq.All,
					},
					altsrc.NewStringFlag(&cli.StringFlag{
						Name:        "pipelines",
						Usage:       "first gen pipeline ids as comma separated values `pipeline1,pipeline2`",
						Destination: &migrationReq.PipelineIds,
					}),
					&cli.StringFlag{
						Name:        "identifiers",
						Usage:       "`IDENTIFIERS` of the next gen pipelines",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the next gen pipeline",
						Destination: &migrationReq.Names,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "rm",
						Usage: "Remove nextgen pipelines",
						Action: func(context *cli.Context) error {
							return cliWrapper(BulkRemovePipelines, context)
						},
					},
					{
						Name:  "import",
						Usage: "import first gen pipelines to next gen",
						Action: func(context *cli.Context) error {
							return cliWrapper(migratePipelines, context)
						},
					},
				},
			},
			{
				Name:  "triggers",
				Usage: "Import triggers by providing the `appId` & `triggerIds`",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if all triggers in the app need to be migrated",
						Destination: &migrationReq.All,
					},
					altsrc.NewStringFlag(&cli.StringFlag{
						Name:        "triggers",
						Usage:       "triggers as comma separated values `triggerId1,triggerId2`",
						Destination: &migrationReq.TriggerIds,
					}),
					altsrc.NewStringFlag(&cli.StringFlag{
						Name:        "names",
						Usage:       "First Gen `NAMES` of the triggers",
						Destination: &migrationReq.Names,
					}),
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(migrateTriggers, context)
				},
			},
			{
				Name:  "expressions",
				Usage: "looks for harness CG expressions in current directory & sub directories from current folder & replaces them with equivalent NG expressions",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Usage:       "if set will only list found expressions without actually replacing anything",
						Destination: &migrationReq.DryRun,
					},
					&cli.StringFlag{
						Name:        "extensions",
						Usage:       "provide comma separated file extensions to look for expressions in. defaults to json, yaml & yml extensions",
						Value:       "json,yaml,yml",
						DefaultText: "json,yaml,yml",
						Destination: &migrationReq.FileExtensions,
					},
				},
				Action: func(context *cli.Context) error {
					return cliWrapper(ReplaceCurrentGenExpressionsWithNextGen, context)
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
					&cli.StringFlag{
						Name:        "export",
						Usage:       "`FOLDER_PATH` of where the files need to be exported to",
						Value:       ".",
						DefaultText: ".",
						Destination: &migrationReq.ExportFolderPath,
					},
					&cli.StringFlag{
						Name:        "csv",
						Usage:       "`CSV_FILE` path to the csv file",
						Destination: &migrationReq.CsvFile,
					},
					&cli.StringFlag{
						Name:        "identifiers",
						Usage:       "`IDENTIFIERS` of the projects",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the projects",
						Destination: &migrationReq.Names,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "csv-template",
						Usage: "Get a CSV with application name, project name, project identifier template for an account",
						Action: func(context *cli.Context) error {
							return cliWrapper(GetProjectCSVTemplate, context)
						},
					},
					{
						Name:  "create",
						Usage: "Create a project",
						Action: func(context *cli.Context) error {
							return cliWrapper(createProject, context)
						},
					},
					{
						Name:  "create-bulk",
						Usage: "Creates apps as projects",
						Action: func(context *cli.Context) error {
							return cliWrapper(bulkCreateProject, context)
						},
					},
					{
						Name:  "rm",
						Usage: "Remove projects",
						Action: func(context *cli.Context) error {
							return cliWrapper(bulkRemoveProject, context)
						},
					},
				},
			},
			{
				Name:  "org",
				Usage: "Org specific commands.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Usage:       "`NAME` for the org",
						Destination: &migrationReq.OrgName,
					},
					&cli.StringFlag{
						Name:        "identifier",
						Usage:       "`IDENTIFIER` for the org",
						Destination: &migrationReq.OrgIdentifier,
					},
					&cli.StringFlag{
						Name:        "identifiers",
						Usage:       "`IDENTIFIERS` of the org",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the org",
						Destination: &migrationReq.Names,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Create an org",
						Action: func(context *cli.Context) error {
							return cliWrapper(createOrg, context)
						},
					},
					{
						Name:  "rm",
						Usage: "Remove org",
						Action: func(context *cli.Context) error {
							return cliWrapper(bulkRemoveOrg, context)
						},
					},
				},
			},
			{
				Name:  "templates",
				Usage: "Template specific commands.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "identifiers",
						Usage:       "`IDENTIFIERS` of the next gen templates",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "ids",
						Usage:       "`IDS` of the first gen templates",
						Destination: &migrationReq.Identifiers,
					},
					&cli.StringFlag{
						Name:        "names",
						Usage:       "`NAMES` of the template",
						Destination: &migrationReq.Names,
					},
					&cli.BoolFlag{
						Name:        "force",
						Usage:       "to force delete template",
						Destination: &migrationReq.Force,
					},
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "if set will delete all templates",
						Destination: &migrationReq.All,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "rm",
						Usage: "Remove templates",
						Action: func(context *cli.Context) error {
							return cliWrapper(BulkRemoveTemplates, context)
						},
					},
					{
						Name:  "import",
						Usage: "import templates. pass the --app flag if you want to migrate app level templates else do not pass",
						Action: func(context *cli.Context) error {
							return cliWrapper(MigrateTemplates, context)
						},
					},
				},
			},
		},
		Before: altsrc.InitInputSourceWithContext(globalFlags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Flags:  globalFlags,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
