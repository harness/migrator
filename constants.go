package main

const (
	SecretManager       = "SECRET_MANAGER"
	Secret              = "SECRET"
	Service             = "SERVICE"
	Connector           = "CONNECTOR"
	Application         = "APPLICATION"
	Workflow            = "WORKFLOW"
	Trigger             = "TRIGGER"
	WorkflowExecution   = "WORKFLOW_EXECUTION"
	Pipeline            = "PIPELINE"
	Infrastructure      = "INFRA"
	Environment         = "ENVIRONMENT"
	ApplicationManifest = "MANIFEST"
	Template            = "TEMPLATE"
	UserGroups          = "USER_GROUP"
)

const (
	Project = "project"
	Org     = "org"
	Account = "account"
)

const (
	All ImportType = "ALL"
)

var scopes = []string{"project", "org", "account"}

const (
	AccountIdentifier = "accountIdentifier"
	OrgIdentifier     = "orgIdentifier"
	ProjectIdentifier = "projectIdentifier"
)

const (
	AccountScope TemplateScope = "ACCOUNT"
)

const (
	AppScope = "APP"
)
