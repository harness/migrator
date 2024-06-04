---
sidebar_position: 3
redirect_from: /migrator/spinnaker-migration/advance
---

# Instructions

## Install Migrator

`curl https://raw.githubusercontent.com/harness/migrator/master/install | sh`

### Test the connection to Harness Account by executing account summary command

`harness-upgrade --api-key API_KEY --account AccountID --env Prod1 account-summary`

## Create migrator-config.yml file

```yaml
env: Prod1
api-key: NG_API_KEY
account: NGAccountID
platform: spinnaker
spinnaker-host: harness-eval.dynamo-staging.eng.armory.io
project: DestinationProjectID
org: DestinationOrgID
secret-scope: project
connector-scope: project
template-scope: project
workflow-scope: project
app-name: SpinnakerAppName
auth64: AUTH_TOKEN
```

## Run command

`harness-upgrade --load migrator-config.yml app`

### You should get output like this with some prompts:

```shell
INFO[2024-03-04T15:50:38-08:00] Starting the migration of Spinnaker application                
INFO[2024-03-04T15:50:38-08:00] Importing the application....                
INFO[2024-03-04T15:50:38-08:00] 
Migration details:
  Platform: spinnaker
  Spinnaker Host: harness-eval.dynamo-staging.eng.armory.io
  App name: prasadtest
  Pipeline Name: 
  Authentication method: basic 
  Insecure: false 
? Do you want to proceed? Yes
INFO[2024-03-04T15:50:41-08:00] Spinnaker migration completed
```

## Import Spinnaker pipeline to existing project
`
harness-upgrade --load migrator-config.yml pipelines import
`
### You should get output like this with some prompts:
```shell
INFO[2024-03-04T15:50:38-08:00] Importing the application....
? Please provide the Spinnaker pipeline name :  Spinnaker Pipeline name
INFO[2024-03-04T15:50:38-08:00]
Migration details:
  Platform: spinnaker
  Spinnaker Host: harness-eval.dynamo-staging.eng.armory.io
  App name: Spinnaker Application Name
  Pipeline Name: Spinnaker Pipeline name
  Authentication method: basic
  Insecure: false
? Do you want to proceed with pipeline migration? Yes
INFO[2024-03-04T15:50:38-08:00] Project check details :                       Account=AccountID OrgIdentifier=default ProjectIdentifier=ProjectID
INFO[2024-03-04T15:50:38-08:00] Project with identifier ProjectID does not exist
INFO[2024-03-04T15:50:38-08:00] Project with identifier ProjectID created
INFO[2024-03-04T15:50:38-08:00] {Pipeline Logs}
INFO[2024-03-04T15:50:41-08:00] Spinnaker migration completed
```

## Video Demo

<iframe width="560" height="315" src="https://www.youtube.com/embed/gkWK3BGpDEU?si=ExGLw9B5Q7h2VMdn" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>