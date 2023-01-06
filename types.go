package main

type EntityType string
type ImportType string

type Filter struct {
	Type        ImportType `json:"importType"`
	AppId       string     `json:"appId"`
	WorkflowIds []string   `json:"workflowIds"`
}

type DestinationDetails struct {
	ProjectIdentifier string `json:"projectIdentifier"`
	OrgIdentifier     string `json:"orgIdentifier"`
}

type EntityDefaults struct {
	Scope string `json:"scope"`
}

type Defaults struct {
	SecretManager EntityDefaults `json:"SECRET_MANAGER"`
	Secret        EntityDefaults `json:"SECRET"`
	Connector     EntityDefaults `json:"CONNECTOR"`
}

type Inputs struct {
	Defaults Defaults `json:"defaults"`
}

type RequestBody struct {
	DestinationDetails DestinationDetails `json:"destinationDetails"`
	EntityType         EntityType         `json:"entityType"`
	Filter             Filter             `json:"filter"`
	Inputs             Inputs             `json:"inputs"`
}

type CurrentGenEntity struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	AppId string `json:"appId"`
}

type UpgradeError struct {
	Message string           `json:"message"`
	Entity  CurrentGenEntity `json:"entity"`
}

type MigrationStats struct {
}

type Resource struct {
	Stats  MigrationStats `json:"stats"`
	Errors []UpgradeError `json:"errors"`
}

type MigrationResponseBody struct {
	Resource Resource    `json:"resource"`
	Messages interface{} `json:"responseMessages"`
}
