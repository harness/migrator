# Harness Upgrade
CLI to help customers, CSMs and developers with migrating their current gen harness account to next gen

## Getting Started

### Installation
Download the latest release from GitHub releases. We support MacOS(`darwim & amd64`), Linux(`linux + (amd64/arm64)`) and Windows(`windows+amd64`). Please download the right assets. Extract the file anywhere.
We recommend that you move it to a folder that is specified in your path. 

```shell
mv harness-upgrade /somepath/
harness-upgrade help
```

If the above works successfully you should see all the commands that are supported with `harness-upgrade`

### Migrating using the step-by-step guide

To migrate account level entities such as secret managers, secrets & connectors
```shell
harness-upgrade
```

To migrate an application 
```shell
harness-upgrade app
```

We use Bearer auth token to make API calls. The token can be provided in the step-by-step guide in the prompt or as below

```shell
HARNESS_MIGRATOR_AUTH=token harness-upgrade
```

OR
```shell
export HARNESS_MIGRATOR_AUTH=token
harness-upgrade
```

### Migrating with a single command
```shell
HARNESS_MIGRATOR_AUTH=token harness-upgrade --project PROJECT --org ORG --account ACCOUNT_ID --secret SCOPE --connector SCOPE --env ENV
```

To migrate an application

```shell
HARNESS_MIGRATOR_AUTH=token harness-upgrade app APP_ID --project PROJECT --org ORG --account ACCOUNT_ID --secret SCOPE --connector SCOPE --env ENV
```

| Flag        | Details                                                                                 |
|-------------|-----------------------------------------------------------------------------------------|
| --env       | Your target environment. It can be either `Dev`, `QA` or `Prod`                         |
| --account   | ID of the account that you wish to migrate                                              |
| --secret    | Scope at which the secret has to be created. It can be `project`, `org` or `account`    |
| --connector | Scope at which the connector has to be created. It can be `project`, `org` or `account` |
| --org       | Identifier of the target org                                                            |
| --project   | Identifier of the target project                                                        |
| --debug     | If debug level logs need to be printed                                                  |
| --json      | Formatted the logs as JSON                                                              |

If not all the required flags are provided we will fall back to prompt based technique to capture all the required details.

## Contact
If you face any issues please reach out on `#migrate-entities-from-cg-to-ng` channel on Slack.