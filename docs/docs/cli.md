---
sidebar_position: 5
slug: /cli
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
| secrets             | Import secrets                                                                                                                             |  
| environments        | Import environments into an existing project from an application                                                                           |  
| connectors          | Import connectors                                                                                                                          |
| workflows           | Import workflows as stage or pipeline templates by providing the `appId` & `workflowIds`                                                   |  
| pipelines           | Import pipelines into an existing project by providing the `appId` & `pipelineIds`                                                         |  
| triggers            | Import triggers by providing the `appId` & `triggerIds`                                                                                    |  
| expressions         | looks for harness CG expressions in current directory & sub directories from current folder & replaces them with equivalent NG expressions |  
| project             | Project specific commands like create, delete, list etc.                                                                                   |  
| org                 | Org specific commands.                                                                                                                     |  
| templates           | Template specific commands.                                                                                                                |
| help, h             | Shows a list of commands or help for one command                                                                                           |  


## Global Options

| Flag                         | Description                                                                                                                     |   
|------------------------------|---------------------------------------------------------------------------------------------------------------------------------|
| --env `ENV`                  | possible values - `Prod`, `Prod3`, `SelfManaged`, `QA`, `Dev`                                                                   |  
| --base-url `BASE_URL`        | provide the `BASE_URL` for self managed platforms                                                                               |  
| --account `ACCOUNT`          | `ACCOUNT` that you wish to migrate                                                                                              |  
| --api-key `API_KEY`          | `API_KEY` to authenticate & authorise the migration                                                                             |  
| --secret-scope `SCOPE`       | `SCOPE` to create secrets in. Possible values - `account`, `org`, `project`                                                     |  
| --connector-scope `SCOPE`    | `SCOPE` to create connectors in. Possible values - `account`, `org`, `project`                                                  |  
| --workflow-scope `SCOPE`     | `SCOPE` to create stage templates in. Possible values - `account`, `org`, `project`                                             |  
| --template-scope `SCOPE`     | `SCOPE` to create templates in. Possible values - `account`, `org`, `project`                                                   |  
| --user-group-scope `SCOPE`   | `SCOPE` to create user groups in. Possible values - `account`, `org`, `project`                                                 |
| --org `IDENTIFIER`           | organisation `IDENTIFIER` in next gen                                                                                           |  
| --project `IDENTIFIER`       | project `IDENTIFIER` in next gen                                                                                                |  
| --app `APP_ID`               | `APP_ID` in current gen                                                                                                         |  
| --identifier-format `FORMAT` | `FORMAT` to use for generation of identifiers. Supported values as `CAMEL_CASE` & `LOWER_CASE` (default: `CAMEL_CASE`)          |  
| --custom-expressions `FILE`  | provide a `FILE` to load custom expressions from                                                                                |  
| --override `FILE`            | provide a `FILE` to load overrides                                                                                              |
| --target-account `ACCOUNT`   | destination `ACCOUNT` that you wish to migrate to                                                                               |  
| --target-api-key `API_KEY`   | `API_KEY` for the target account to authenticate & authorise the migration.                                                     |
| --target-gateway-url `URL`   | destination gateway `URL`. For Prod1 & Prod2, use https://app.harness.io/gateway, for Prod3 use https://app3.harness.io/gateway |
| --load `FILE`                | `FILE` to load flags from                                                                                                       |
| --insecure                   | allow insecure API requests. This is automatically set to true if environment is Dev (default: false)                           |
| --log-level                  | set the log level. Possible values - trace, debug, info, warn, error, fatal, panic. Default is `info`                           |
| --flags value                | provide a list of flags for custom logic. Please refer [here](advanced/flags).                                                  |
| --json                       | log as JSON instead of standard ASCII formatter (default: false).                                                               |
| --help, -h                   | show help.                                                                                                                      |
| --version, -v                | print the version                                                                                                               |

If not all the required flags are provided we will fall back to prompt based technique to capture all the required details.  
               
