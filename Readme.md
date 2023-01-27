# Harness Upgrade

Harness has developed a tool to help user's migrate from Harness CD Current Gen to CD Next Generation. The tool can migrate specific Harness resources or bulk migration of resources. 

### Who is the Tool intended for?

- CLI to help customers, CSMs and developers with migrating their current gen harness account to next gen

### Key Components that can be migrated 

1. Pipelines
2. Workflows
3. Service (Service Definition - Manifests, variables included)
5. Environments (Overrides included)
6. Infrastructure Definition
7. Connectors
8. Secrets
9. Templates 


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

To migrate workflows
```shell
harness-upgrade workflows
```

To migrate pipelines
```shell
harness-upgrade pipelines
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
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --env ENV
```

### To migrate an application

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --app APP_ID --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --env ENV app
```

### To migrate workflows

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --app APP_ID --workflows WORKFLOW_IDS --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --workflow-scope SCOPE --env ENV workflows 
```
> If you do not provide use `--workflows` flag it will migrate all workflows in the app

### To migrate pipelines

```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --app APP_ID --pipelines PIPELINE_IDS --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --workflow-scope SCOPE --env ENV pipelines 
```
> If you do not provide use `--pipelines` flag it will migrate all pipelines in the app

### To create a project
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --env ENV --org ORG project --name PROJECT_NAME --identifier PROJECT_IDENTIFIER create
```

### To create projects in bulk
The following command creates a corresponding project for every app in the account. It then exports a YAML file for every project to the specified export path(defaults to current dir) specified. 
```shell
HARNESS_MIGRATOR_AUTH=apiKey harness-upgrade --account ACCOUNT_ID --env ENV --org ORG --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE project --export FOLDER_PATH create-bulk
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
harness-upgrade --load 'templates/account.yaml'

# To migrate the app
harness-upgrade app --load 'templates/app.yaml'

# To migrate the workflows
harness-upgrade workflows --load 'templates/workflows.yaml'

# To migrate the pipelines
harness-upgrade pipelines --load 'templates/pipelines.yaml'
```

## Global Flags

| Flag                  | Details                                                                                                                |
|-----------------------|------------------------------------------------------------------------------------------------------------------------|
| --env                 | Your target environment. It can be either `Dev`, `QA`, `Prod` or `Prod3`                                               |
| --account             | `ACCOUNT_ID` of the account that you wish to migrate                                                                   |
| --api-key             | `API_KEY` to authenticate & authorise the migration. You may also use the `HARNESS_MIGRATOR_AUTH` env variable instead |
| --secret-scope        | Scope at which the secret has to be created. It can be `project`, `org` or `account`                                   |
| --connector-scope     | Scope at which the connector has to be created. It can be `project`, `org` or `account`                                |
| --template-scope      | Scope at which the templates has to be created. It can be `project`, `org` or `account`                                |
| --workflow-scope      | Scope at which the workflow as template has to be created. It can be `project`, `org` or `account`                     |
| --org                 | Identifier of the target org                                                                                           |
| --project             | Identifier of the target project                                                                                       |
| --app                 | Application ID from current gen                                                                                        |
| --workflows           | Workflow Ids as comma separated values(ex. `workflow1,workflow2,workflow3`)                                            |
| --pipelines           | Pipeline Ids as comma separated values(ex. `pipeline1,pipeline2,pipeline3`)                                            |
| --destination-project | URL of the project where we want to migrate                                                                            |
| --source-app          | URL of the application from which we will migrate the entities                                                         |
| --debug               | If debug level logs need to be printed                                                                                 |
| --json                | Formatted the logs as JSON                                                                                             |

If not all the required flags are provided we will fall back to prompt based technique to capture all the required details.

## Contact
If you face any issues please reach out to us or feel free to create a GitHub issue.
