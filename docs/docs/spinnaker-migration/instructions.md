 # Instructions


## Install Migrator

`curl https://raw.githubusercontent.com/harness/migrator/master/install | sh`

### Test the connection to Harness Account by executing account summary command

`harness-upgrade --api-key API_KEY --account AccountID --env Prod1 account-summary`

## Create migrator-config.yml file

```yaml
env: Prod1
api-key: API_KEY
account: AccountID
platform: spinnaker
spinnaker-host: harness-eval.dynamo-staging.eng.armory.io
project: ProjectID
org: default
secret-scope: project
connector-scope: project
template-scope: project
workflow-scope: project
auth64: AUTH_TOKEN
```

## Run command

`harness-upgrade --load migrator-config.yml app --app-name prasadtest`

## You should get output like this with some prompts:

```shell
`INFO[2024-03-04T15:50:38-08:00] Importing the application....                
INFO[2024-03-04T15:50:38-08:00] 
Migration details:
  Platform: spinnaker
  Spinnaker Host: harness-eval.dynamo-staging.eng.armory.io
  App name: prasadtest
  Pipeline Name: 
  Authentication method: basic 
  Insecure: false 
? Do you want to proceed with application migration? Yes
INFO[2024-03-04T15:50:41-08:00] Spinnaker migration completed`
```
