# Harness Upgrade

Harness has developed a tool to help user's migrate from Harness CD Current Gen to CD Next Generation. The tool can migrate specific Harness resources or bulk migration of resources. 

### Who is the Tool intended for?

- CLI to help customers, CSMs and developers with migrating their current gen harness account to next gen

### Key Components that can be migrated 

1. Pipelines
2. Workflows
3. Service (Service Definition - Manifests, variables included)
4. Environments (Overrides included)
5. Infrastructure Definition
6. Connectors
7. Secrets & Secret Managers
8. Templates 


## Installation
Download the latest release from GitHub releases. We support MacOS(`darwim & amd64`), Linux(`linux + (amd64/arm64)`) and Windows(`windows+amd64`). Please download the right assets. Extract the file anywhere.
We recommend that you move it to a folder that is specified in your path. 

```shell
mv harness-upgrade /somepath/
harness-upgrade help
```

If you are using macOS then just do
```shell
mv harness-upgrade /usr/local/bin/
harness-upgrade help
```

If the above works successfully you should see all the commands that are supported with `harness-upgrade`

## Migrating using the step-by-step guide

To migrate account level entities such as secret managers, secrets & connectors
```shell
harness-upgrade
```

To migrate an application 
```shell
harness-upgrade app
```

To migrate an application and all the workflows & pipelines
```shell
harness-upgrade app --all
```

To migrate workflows
```shell
harness-upgrade workflows
```

> All workflows are created as stage templates in NextGen except for multi-service workflows

To create pipelines out of workflows
```shell
harness-upgrade workflows --as-pipelines 
```

> Workflows will continue to be migrated to stage templates & the pipelines are created using these stage templates.

To migrate pipelines
```shell
harness-upgrade pipelines
```

To migrate triggers
```shell
harness-upgrade triggers
```

To create project
```shell
harness-upgrade project create
```

We use API keys created in NextGen to make API calls. The token can be provided in the step-by-step guide in the prompt or as below

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade
```

OR
```shell
export HARNESS_MIGRATOR_AUTH=apiKey
harness-upgrade
```

OR
```shell
harness-upgrade --api-key apiKey
```

## Migrating with a single command
Using the step-by-step guide is the recommended way to get started with upgrade, but filling the prompts everytime can be tedious. If you wish to provide all or a few inputs you can pass them using the flags. If required arguments are not provided we will prompt for the inputs.  

### To migrate all account level entities

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --env ENV account
```

### To migrate an application

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --app APP_ID --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --env ENV app
```

### To migrate workflows

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --app APP_ID --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --workflow-scope SCOPE --env ENV workflows --workflows WORKFLOW_IDS
```
> If you do not provide use `--workflows` flag it will migrate all workflows in the app

### To migrate pipelines

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --app APP_ID --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --workflow-scope SCOPE --env ENV pipelines --pipelines PIPELINE_IDS 
```
> If you do not provide use `--pipelines` flag it will migrate all pipelines in the app

### To migrate triggers

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --app APP_ID --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --workflow-scope SCOPE --env ENV triggers --triggers TRIGGER_IDS 
```
> If you do not provide use `--triggers` flag it will migrate all triggers in the app

> Migrating Triggers is WIP. The current results requires considerable manual effort post migration

### To create a project
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --env ENV --org ORG project --name PROJECT_NAME --identifier PROJECT_IDENTIFIER create
```

### To create projects in bulk
The following command creates a corresponding project for every app in the account. It then exports a YAML file for every project to the specified export path(defaults to current dir) specified. 
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --env ENV --org ORG --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE project --export FOLDER_PATH create-bulk
```

If you wish to migrate all the above apps then run this - 
```shell
for f in /path/to/folder/*.yaml; do harness-upgrade --load $f app --all; done
```

### To remove projects
The following command removes projects from a given org in an account. You can provide the names or identifiers of the projects.

To remove projects by name use `--names`
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --env ENV --org ORG project --names name1,name2 rm
```

To remove projects by identifier use `--identifiers`
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --env ENV --org ORG project --identifiers identifier1,identifier2 rm
```

### To get account summary
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --env ENV account-summary
```

### To get application summary
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --app APP_ID --env ENV application-summary
```
### Using URL to parse details
```shell
harness-upgrade --source-app <CURRENT_GEN_APP_URL> --destination-project <NG_PROJECT_URL>
```

## Migrating by providing the flags from a file
If you wish to provide the flags from a file you can use the `--load` to load flags from a file. You can find templates for various options in the `templates/` directory.

```shell
# To migrate the account level entities
harness-upgrade --load templates/account.yaml

# To migrate the app
harness-upgrade --load templates/app.yaml app

# To migrate the app along with all workflows & pipelines
harness-upgrade --load templates/app.yaml app --all

# To migrate the workflows
harness-upgrade --load templates/workflows.yaml workflows

# To migrate the pipelines
harness-upgrade --load templates/pipelines.yaml pipelines
```

## Replace First Gen expressions with Next Gen

The command replaces first gen expressions found in all files in current & subdirectories. By default, we only process the files with following file extensions `yml`, `yaml` & `json`.  
```shell
harness-upgrade expressions 
```

To provide custom files extensions
```shell
harness-upgrade expressions --extensions yml,txt,xml
```

Do a dry run on the files without replacing any CG expressions
```shell
harness-upgrade expressions --dry-run
```

Secrets referenced in these files are converted to camel case to align with the migrator. You can provide the scope of the secrets using the `--secret-scope` flag.
```shell
harness-upgrade --secret-scope account expressions
```

To provide custom expressions or override default expressions
```shell
harness-upgrade expressions --override /path/to/file.yaml
```

The above command expects a yaml file whose contents are simple key values of first-gen expressions & string to replace that expression with. A sample override file would look like this - 
```yaml
appName: <+org.name>
service.name: <+project.name>
```

## Global Flags

| Flag                  | Details                                                                                                                          |
|-----------------------|----------------------------------------------------------------------------------------------------------------------------------|
| --env                 | Your target environment. It can be either `Dev`, `QA`, `Prod` or `Prod3`                                                         |
| --account             | `ACCOUNT_ID` of the account that you wish to migrate                                                                             |
| --api-key             | `API_KEY` to authenticate & authorise the migration. You may also use the `HARNESS_MIGRATOR_AUTH` env variable instead           |
| --secret-scope        | Scope at which the secret has to be created. It can be `project`, `org` or `account`                                             |
| --connector-scope     | Scope at which the connector has to be created. It can be `project`, `org` or `account`                                          |
| --template-scope      | Scope at which the templates has to be created. It can be `project`, `org` or `account`                                          |
| --workflow-scope      | Scope at which the workflow as template has to be created. It can be `project`, `org` or `account`                               |
| --org                 | Identifier of the target org                                                                                                     |
| --project             | Identifier of the target project                                                                                                 |
| --app                 | Application ID from current gen                                                                                                  |
| --workflows           | Workflow Ids as comma separated values(ex. `workflow1,workflow2,workflow3`)                                                      |
| --pipelines           | Pipeline Ids as comma separated values(ex. `pipeline1,pipeline2,pipeline3`)                                                      |
| --destination-project | URL of the project where we want to migrate                                                                                      |
| --source-app          | URL of the application from which we will migrate the entities                                                                   |
| --identifier-format   | To control the format of the identifier generated. It defaults to `CAMEL_CASE`, we currently support `CAMEL_CASE` & `LOWER_CASE` |
| --debug               | If debug level logs need to be printed                                                                                           |
| --json                | Formatted the logs as JSON                                                                                                       |

If not all the required flags are provided we will fall back to prompt based technique to capture all the required details.

## Contact
If you face any issues please reach out to us or feel free to create a GitHub issue.
