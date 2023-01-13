package main

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"strings"
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

func PromptDefaultInputs() bool {
	promptConfirm := false

	if len(migrationReq.Environment) == 0 {
		promptConfirm = true
		migrationReq.Environment = SelectInput("Which environment?", []string{Dev, QA, Prod, Prod3}, Dev)
	}

	// Check if auth is provided. If not provided then request for one
	migrationReq.Auth = os.Getenv("HARNESS_MIGRATOR_AUTH")
	if len(migrationReq.Auth) == 0 {
		migrationReq.Auth = TextInput("The environment variable 'HARNESS_MIGRATOR_AUTH' is not set. What is the api key?")
	}

	if migrationReq.Environment == "Dev" || migrationReq.AllowInsecureReq {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
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

	if len(migrationReq.TemplateScope) == 0 {
		promptConfirm = true
		migrationReq.TemplateScope = SelectInput("Scope for templates:", scopes, Project)
	}

	return promptConfirm
}

func PromptOrgAndProject() bool {
	promptConfirm := false
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

func migrateApp(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app that you wish to import -")
	}

	promptConfirm = PromptOrgAndProject() || promptConfirm

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

func migrateWorkflows(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app containing the workflows -")
	}

	if len(migrationReq.WorkflowIds) == 0 {
		promptConfirm = true
		migrationReq.WorkflowIds = TextInput("Provide the workflows that you wish to import as template as comma separated values(e.g. workflow1,workflow2)")
	}

	if len(migrationReq.WorkflowScope) == 0 {
		promptConfirm = true
		migrationReq.WorkflowScope = SelectInput("Scope for workflows:", scopes, Project)
	}

	promptConfirm = PromptOrgAndProject() || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with workflows migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := GetUrl(migrationReq.Environment, "save/v2", migrationReq.Account)
	// Migrating the app
	log.Info("Importing the workflows....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Workflow, Filter{
		WorkflowIds: strings.Split(migrationReq.WorkflowIds, ","),
		AppId:       migrationReq.AppId,
	}))
	log.Info("Imported the workflows.")

	return nil
}

func migratePipelines(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()
	if len(migrationReq.AppId) == 0 {
		promptConfirm = true
		migrationReq.AppId = TextInput("Please provide the application ID of the app containing the pipeline -")
	}

	if len(migrationReq.WorkflowScope) == 0 {
		promptConfirm = true
		migrationReq.WorkflowScope = SelectInput("Scope for workflow to be migrated as templates:", scopes, Project)
	}

	if len(migrationReq.PipelineIds) == 0 {
		promptConfirm = true
		migrationReq.PipelineIds = TextInput("Provide the pipelines that you wish to import as template as comma separated values(e.g. pipeline1,pipeline2)")
	}

	promptConfirm = PromptOrgAndProject() || promptConfirm

	logMigrationDetails()

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with pipeline migration?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	url := GetUrl(migrationReq.Environment, "save/v2", migrationReq.Account)
	// Migrating the app
	log.Info("Importing the pipelines....")
	CreateEntity(url, migrationReq.Auth, getReqBody(Pipeline, Filter{
		PipelineIds: strings.Split(migrationReq.PipelineIds, ","),
		AppId:       migrationReq.AppId,
	}))
	log.Info("Imported the pipelines.")

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
