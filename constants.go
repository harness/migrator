package main

const (
	SecretManager       EntityType = "SECRET_MANAGER"
	Secret                         = "SECRET"
	Service                        = "SERVICE"
	Connector                      = "CONNECTOR"
	Application                    = "APPLICATION"
	Workflow                       = "WORKFLOW"
	Trigger                        = "TRIGGER"
	WorkflowExecution              = "WORKFLOW_EXECUTION"
	Pipeline                       = "PIPELINE"
	Infrastructure                 = "INFRA"
	Environment                    = "ENVIRONMENT"
	ApplicationManifest            = "MANIFEST"
	Template                       = "TEMPLATE"
)

const (
	Project string = "project"
	Org            = "org"
	Account        = "account"
)

const (
	All ImportType = "ALL"
)

var scopes = []string{"project", "org", "account"}
