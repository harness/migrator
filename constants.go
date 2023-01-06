package main

const (
	SecretManager EntityType = "SECRET_MANAGER"
	Secret                   = "SECRET"
	Connector                = "CONNECTOR"
	Application              = "APPLICATION"
	Workflow                 = "WORKFLOW"
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
