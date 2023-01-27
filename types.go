package main

type EntityType string
type ImportType string

type Filter struct {
	Type        ImportType `json:"importType"`
	AppId       string     `json:"appId"`
	WorkflowIds []string   `json:"workflowIds"`
	PipelineIds []string   `json:"pipelineIds"`
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
	Workflow      EntityDefaults `json:"WORKFLOW"`
	Template      EntityDefaults `json:"TEMPLATE"`
}

type Inputs struct {
	Defaults Defaults `json:"defaults"`
}

type ProjectDetails struct {
	OrgIdentifier string   `json:"orgIdentifier"`
	Identifier    string   `json:"identifier"`
	Name          string   `json:"name"`
	Color         string   `json:"color"`
	Modules       []string `json:"modules"`
	Description   string   `json:"description"`
}

type BulkProjectResult struct {
	AppName           string       `json:"appName"`
	AppId             string       `json:"appId"`
	ProjectIdentifier string       `json:"projectIdentifier"`
	ProjectName       string       `json:"projectName"`
	Error             UpgradeError `json:"error"`
}

type BulkCreateBody struct {
	Org string `json:"orgIdentifier"`
}

type ProjectCreateBody struct {
	Project ProjectDetails `json:"project"`
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
	SuccessfullyMigrated int64 `json:"successfullyMigrated"`
	AlreadyMigrated      int64 `json:"alreadyMigrated"`
}

type Resource struct {
	RequestId       string                    `json:"requestId"`
	Stats           map[string]MigrationStats `json:"stats"`
	Errors          []UpgradeError            `json:"errors"`
	Status          string                    `json:"status"`
	ResponsePayload SummaryResponse           `json:"responsePayload"`
}

type ResponseBody struct {
	Code     string             `json:"code"`
	Message  string             `json:"message"`
	Status   string             `json:"status"`
	Data     interface{}        `json:"data"`
	Resource interface{}        `json:"resource"`
	Messages []ResponseMessages `json:"responseMessages"`
}

type ResponseMessages struct {
	Code         string      `json:"code"`
	Level        string      `json:"level"`
	Message      string      `json:"message"`
	Exception    interface{} `json:"exception"`
	FailureTypes interface{} `json:"failureTypes"`
}

type SummaryResponse struct {
	Summary map[string]EntitySummary `json:"summary"`
}

type EntitySummary struct {
	Name                     string           `json:"name"`
	Count                    int64            `json:"count"`
	TypeSummary              map[string]int64 `json:"typeSummary"`
	StepTypeSummary          map[string]int64 `json:"stepTypeSummary"`
	KindSummary              map[string]int64 `json:"kindSummary"`
	StoreSummary             map[string]int64 `json:"storeSummary"`
	DeploymentTypeSummary    map[string]int64 `json:"deploymentTypeSummary"`
	ArtifactTypeSummary      map[string]int64 `json:"artifactTypeSummary"`
	CloudProviderTypeSummary map[string]int64 `json:"cloudProviderTypeSummary"`
	Expressions              []string         `json:"expressions"`
}
