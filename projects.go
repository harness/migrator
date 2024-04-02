package main

import (
	"encoding/json"
	"fmt"
	"github.com/jszwec/csvutil"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"path"
	"path/filepath"
	"strconv"
)

func createProject(*cli.Context) error {
	promptConfirm := PromptEnvDetails()
	promptConfirm = PromptOrgAndProject([]string{Org}) || promptConfirm
	if len(migrationReq.ProjectName) == 0 {
		promptConfirm = true
		migrationReq.ProjectName = TextInput("Name of the Project - ")
	}
	if len(migrationReq.ProjectIdentifier) == 0 {
		promptConfirm = true
		migrationReq.ProjectIdentifier = TextInput("Identifier for the Project - ")
	}

	log.WithFields(log.Fields{
		"Account":           migrationReq.Account,
		"OrgIdentifier":     migrationReq.OrgIdentifier,
		"ProjectName":       migrationReq.ProjectName,
		"ProjectIdentifier": migrationReq.ProjectIdentifier,
	}).Info("Project creation details")

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with project creation?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	log.Info("Creating the project....")

	err := createAProject(migrationReq.OrgIdentifier, migrationReq.ProjectName, migrationReq.ProjectIdentifier)

	if err == nil {
		log.Info("Created the project!")
	} else {
		log.Error("There was an error creating the project")
		return err
	}

	return nil
}

func createAProject(orgIdentifier string, name string, identifier string) error {
	url := fmt.Sprintf("%s/api/projects?accountIdentifier=%s&orgIdentifier=%s", GetBaseUrl(migrationReq.Environment, NextGenService), migrationReq.Account, orgIdentifier)
	_, err := Post(url, migrationReq.Auth, ProjectBody{
		Project: ProjectDetails{
			OrgIdentifier: orgIdentifier,
			Identifier:    identifier,
			Name:          name,
			Color:         "#0063f7",
			Modules:       []string{"CD"},
			Description:   "",
		}})
	return err
}

func bulkCreateProject(*cli.Context) error {
	promptConfirm := PromptDefaultInputs()

	if len(migrationReq.CsvFile) != 0 {
		fmt.Printf("Importing from - %s\n", migrationReq.CsvFile)
	} else {
		promptConfirm = PromptOrgAndProject([]string{Org}) || promptConfirm
	}

	if len(migrationReq.ExportFolderPath) == 0 {
		migrationReq.ExportFolderPath = TextInput("Where would you like to export the generated files?")
		promptConfirm = true
	}

	log.WithFields(log.Fields{
		"Account":       migrationReq.Account,
		"OrgIdentifier": migrationReq.OrgIdentifier,
		"Export folder": migrationReq.ExportFolderPath,
	}).Info("Bulk Project creation details")

	if promptConfirm {
		confirm := ConfirmInput("Do you want to proceed with projects creation?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	if len(migrationReq.CsvFile) != 0 {
		return CreateProjectsUsingCSV()
	}

	url := GetUrlWithQueryParams(migrationReq.Environment, MigratorService, "projects/bulk", map[string]string{
		AccountIdentifier: migrationReq.Account,
	})

	log.Info("Creating the projects....")

	resp, err := Post(url, migrationReq.Auth, BulkCreateBody{
		DestinationAccountIdentifier: migrationReq.TargetAccount,
		DestinationAuthToken:         migrationReq.TargetAuthToken,
		DestinationGatewayUrl:        migrationReq.TargetGatewayUrl,
		Org:                          migrationReq.OrgIdentifier,
		IdentifierCaseFormat:         migrationReq.IdentifierCase,
	})

	if err != nil {
		log.Fatal("There was an error creating the projects", err)
		return err
	}

	byteData, err := json.Marshal(resp.Resource)
	if err != nil {
		return err
	}
	var data []BulkProjectResult
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return err
	}

	for _, result := range data {
		if len(result.Error.Message) > 0 {
			log.Errorf("When creating a project for application '%s' there was an error %s", result.AppName, result.Error.Message)
			continue
		}
		err = writeYamlToFile(result.AppId, result.AppName, migrationReq.OrgIdentifier, result.ProjectIdentifier)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func writeYamlToFile(appId string, appName string, orgIdentifier string, projectIdentifier string) error {
	yamlData := map[string]string{
		"env":                  migrationReq.Environment,
		"account":              migrationReq.Account,
		"api-key":              migrationReq.Auth,
		"app":                  appId,
		"org":                  orgIdentifier,
		"project":              projectIdentifier,
		"secret-scope":         migrationReq.SecretScope,
		"secret-manager-scope": migrationReq.SecretManagerScope,
		"connector-scope":      migrationReq.ConnectorScope,
		"template-scope":       migrationReq.TemplateScope,
		"workflow-scope":       Project,
	}

	yamlContent, err := yaml.Marshal(&yamlData)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	absolutePath, err := filepath.Abs(migrationReq.ExportFolderPath)
	if err != nil {
		return err
	}
	err = MkDir(absolutePath)
	if err != nil {
		return err
	}
	err = WriteToFile(path.Join(absolutePath, projectIdentifier+".yaml"), yamlContent)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Infof("Application %s was exported to file %s.yaml", appName, projectIdentifier)
	return nil
}

func bulkRemoveProject(*cli.Context) error {
	promptConfirm := PromptEnvDetails()
	promptConfirm = PromptOrgAndProject([]string{Org}) || promptConfirm
	names := Split(migrationReq.Names, ",")
	identifiers := Split(migrationReq.Identifiers, ",")
	if len(names) == 0 && len(identifiers) == 0 {
		log.Fatal("No names or identifiers for the projects provided. Aborting")
	}
	if len(names) > 0 && len(identifiers) > 0 {
		log.Fatal("Both names and identifiers for the projects provided. Aborting")
	}

	n := len(identifiers)
	if len(names) > 0 {
		n = len(names)
	}
	if promptConfirm {
		confirm := ConfirmInput("Are you sure you want to proceed with deletion of " + strconv.Itoa(n) + " projects?")
		if !confirm {
			log.Fatal("Aborting...")
		}
	}

	if len(names) > 0 {
		projects := getProjects()
		for _, name := range names {
			id := findProjectIdByName(projects, name)
			if len(id) > 0 {
				identifiers = append(identifiers, id)
			}
		}
		log.Debugf("Valid identifiers for the given names are - %s", identifiers)
	}

	for _, identifier := range identifiers {
		deleteProject(identifier)
	}
	log.Info("Finished operation for all given projects")
	return nil
}

func deleteProject(projectId string) {
	url := fmt.Sprintf("%s/api/projects/%s?accountIdentifier=%s&orgIdentifier=%s", GetBaseUrl(migrationReq.Environment, NextGenService), projectId, migrationReq.Account, migrationReq.OrgIdentifier)

	log.Infof("Deleting the project with identifier %s", projectId)

	_, err := Delete(url, migrationReq.Auth, nil)

	if err == nil {
		log.Infof("Successfully deleted the project - %s", projectId)
	} else {
		log.Errorf("Failed to delete the project - %s", projectId)
	}
}

func getProjects() []ProjectDetails {
	url := fmt.Sprintf("%s/api/projects?accountIdentifier=%s&orgIdentifier=%s&pageSize=1000", GetBaseUrl(migrationReq.Environment, NextGenService), migrationReq.Account, migrationReq.OrgIdentifier)
	resp, err := Get(url, migrationReq.Auth)
	if err != nil || resp.Status != "SUCCESS" {
		log.Fatal("Failed to fetch projects", err)
	}
	byteData, err := json.Marshal(resp.Data)
	if err != nil {
		log.Fatal("Failed to fetch projects", err)
	}
	var projects ProjectListBody
	err = json.Unmarshal(byteData, &projects)
	if err != nil {
		log.Fatal("Failed to fetch projects", err)
	}
	var projectDetails []ProjectDetails

	for _, p := range projects.Projects {
		projectDetails = append(projectDetails, p.Project)
	}
	return projectDetails
}

func findProjectIdByName(projects []ProjectDetails, projectName string) string {
	for _, p := range projects {
		if p.Name == projectName {
			return p.Identifier
		}
	}
	return ""
}

func GetProjectCSVTemplate(*cli.Context) (err error) {
	_ = PromptEnvDetails()

	if len(migrationReq.CsvFile) == 0 {
		migrationReq.CsvFile = TextInput("File to export the csv to - ")
	}

	apps, err := listEntities("apps")
	if err != nil {
		return
	}
	if len(apps) == 0 {
		return
	}

	var records []ProjectCSV
	for _, app := range apps {
		records = append(records, ProjectCSV{
			ProjectIdentifier: ToCamelCase(app.Name),
			AppName:           app.Name,
			OrgIdentifier:     "default",
			ProjectName:       app.Name,
		})
	}

	content, err := csvutil.Marshal(records)
	if err != nil {
		return
	}

	err = WriteToFile(migrationReq.CsvFile, content)
	if err != nil {
		return err
	}
	return nil
}

func CreateProjectsUsingCSV() (err error) {
	data, err := ReadFile(migrationReq.CsvFile)
	if err != nil {
		return
	}

	var records []ProjectCSV
	if err = csvutil.Unmarshal([]byte(data), &records); err != nil {
		return
	}

	apps, err := listEntities("apps")
	if err != nil {
		return err
	}
	var appsMap = make(map[string]string)
	for _, app := range apps {
		appsMap[app.Name] = app.Id
	}

	for _, record := range records {
		appId, ok := appsMap[record.AppName]
		if !ok {
			continue
		}
		err = createAProject(record.OrgIdentifier, record.ProjectName, record.ProjectIdentifier)
		if err != nil {
			log.Error(err)
			continue
		}
		err = writeYamlToFile(appId, record.AppName, record.OrgIdentifier, record.ProjectIdentifier)
		if err != nil {
			return
		}
	}
	return
}

func CheckProjectExistsAndCreate() error {
	if len(migrationReq.ProjectIdentifier) == 0 {
		migrationReq.ProjectIdentifier = TextInput("Identifier for the Project - ")
	}

	log.WithFields(log.Fields{
		"Account":           migrationReq.Account,
		"OrgIdentifier":     migrationReq.OrgIdentifier,
		"ProjectIdentifier": migrationReq.ProjectIdentifier,
	}).Info("Project check details : ")

	projects := getProjects()
	id := findProjectIdByName(projects, migrationReq.ProjectIdentifier)

	if len(id) > 0 {
		log.Infof("Project with identifier %s exists", migrationReq.ProjectIdentifier)
	} else {
		log.Infof("Project with identifier %s does not exist", migrationReq.ProjectIdentifier)
		if err := createAProject(migrationReq.OrgIdentifier, migrationReq.ProjectIdentifier, formatString(migrationReq.ProjectIdentifier)); err != nil {
			log.Error(err)
		}
		log.Infof("Project with identifier %s created", migrationReq.ProjectIdentifier)
	}
	return nil
}
