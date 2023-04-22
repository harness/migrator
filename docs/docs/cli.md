---
sidebar_position: 5
slug: /
---

# CLI

Summary of all CLI commands and options

## Usage
```shell
harness-upgrade [global options] command [command options] [arguments...]
```

## Commands

| Command             | Description                                                                                                                                |   
|---------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| update, upgrade     | Check for updates and upgrade the CLI                                                                                                      |  
| account-summary     | Get a summary of the account                                                                                                               |  
| application-summary | Get a summary of an app                                                                                                                    |
| user-groups         | Import user groups from First Gen to Next Gen                                                                                              |  
| account             | Import secrets managers, secrets, connectors. This will not migrate services, environments, triggers, pipelines etc                        |  
| app                 | Import an app into an existing project by providing the `appId`                                                                            |  
| service, services   | Import services into an existing project from an application                                                                               |  
| workflows           | Import workflows as stage or pipeline templates by providing the `appId` & `workflowIds`                                                   |  
| pipelines           | Import pipelines into an existing project by providing the `appId` & `pipelineIds`                                                         |  
| triggers            | Import triggers by providing the `appId` & `triggerIds`                                                                                    |  
| expressions         | looks for harness CG expressions in current directory & sub directories from current folder & replaces them with equivalent NG expressions |  
| project             | Project specific commands like create, delete, list etc.                                                                                   |  
| org                 | Org specific commands.                                                                                                                     |  
| help, h             | Shows a list of commands or help for one command                                                                                           |  


## Global Options

| Flag                         | Description                                                                                                            |   
|------------------------------|------------------------------------------------------------------------------------------------------------------------|
| --env `ENV`                  | possible values - `Prod`, `Prod3`, `SelfManaged`, `QA`, `Dev`                                                          |  
| --base-url `BASE_URL`        | provide the `BASE_URL` for self managed platforms                                                                      |  
| --account `ACCOUNT`          | `ACCOUNT` that you wish to migrate                                                                                     |  
| --api-key `API_KEY`          | `API_KEY` to authenticate & authorise the migration                                                                    |  
| --secret-scope `SCOPE`       | `SCOPE` to create secrets in. Possible values - `account`, `org`, `project`                                            |  
| --connector-scope `SCOPE`    | `SCOPE` to create connectors in. Possible values - `account`, `org`, `project`                                         |  
| --workflow-scope `SCOPE`     | `SCOPE` to create stage templates in. Possible values - `account`, `org`, `project`                                    |  
| --template-scope `SCOPE`     | `SCOPE` to create templates in. Possible values - `account`, `org`, `project`                                          |  
| --org `IDENTIFIER`           | organisation `IDENTIFIER` in next gen                                                                                  |  
| --project `IDENTIFIER`       | project `IDENTIFIER` in next gen                                                                                       |  
| --app `APP_ID`               | `APP_ID` in current gen                                                                                                |  
| --identifier-format `FORMAT` | `FORMAT` to use for generation of identifiers. Supported values as `CAMEL_CASE` & `LOWER_CASE` (default: `CAMEL_CASE`) |  
| --target-account `ACCOUNT`   | destination `ACCOUNT` that you wish to migrate to                                                                      |  
| --target-api-key `API_KEY`   | `API_KEY` for the target account to authenticate & authorise the migration.                                            |
| --load `FILE`                | `FILE` to load flags from                                                                                              |
| --insecure                   | allow insecure API requests. This is automatically set to true if environment is Dev (default: false)                  |
| --debug                      | print debug level logs (default: false)                                                                                |
| --json                       | log as JSON instead of standard ASCII formatter (default: false).                                                      |
| --help, -h                   | show help.                                                                                                             |
| --version, -v                | print the version                                                                                                      |

If not all the required flags are provided we will fall back to prompt based technique to capture all the required details.  
               
