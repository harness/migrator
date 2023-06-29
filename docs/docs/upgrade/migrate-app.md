---
sidebar_position: 3
---

# Migrating Application

Before we move entities from first gen to next gen it is important to get a gist of entities that are present in the app. So ensure to run application summary

## Migrating an application
To migrate all services, environments, application defaults, referenced secrets, referenced connectors we use the `harness-upgrade app`
```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --env ENV \
  app
```

:::info
The above command does not create projects or org. The above command does not migrate workflows, pipelines & triggers
:::

## Migrating application with workflows and pipelines
To migrate the app  `harness-upgrade app --all`

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  app --all
```

## Migrating workflows

### Migrating all workflows
Workflows are created as stage templates in NextGen by default. The only exception to this rule are multi-service workflows that are created as Pipelines.
To migrate all workflows in the app  `harness-upgrade workflows --all`

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  workflows --all
```

### Migrating specific workflows
To migrate specific workflows in the app  `harness-upgrade workflows --workflows workflow1,workflow2`. Pass comma separated workflow ids.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  workflows --workflows WORKFLOW_IDS
```

### Migrating workflows as pipelines

If we want to create Pipelines for workflows we can use the `--as-pipelines` flag.
```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  workflows --all --as-pipelines
```

:::note
When creating workflows as pipelines, workflows are first migrated as stage templates and these stage templates are used to create the pipelines.
:::

## Migrating pipelines

### Migrating all pipelines
To migrate all pipelines in the app  `harness-upgrade pipelines --all`

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  pipelines --all import
```

### Migrating specific pipelines
To migrate specific pipelines in the app  `harness-upgrade pipelines --pipelines pipeline1,pipeline2`. Pass comma separated pipeline ids.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  pipelines --pipelines PIPELINE_IDS import
```

## Migrating triggers

### Migrating all triggers
To migrate all triggers in the app  `harness-upgrade triggers --all`

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  triggers --all
```

### Migrating specific triggers
To migrate specific triggers in the app  `harness-upgrade triggers --triggers trigger1,trigger2`. Pass comma separated trigger ids.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  triggers --triggers TRIGGER_IDS
```

:::caution
Migrating Triggers is WIP. The current results requires considerable manual effort post migration
:::


## Migrating services

### Migrating all services
To migrate all services in the app  `harness-upgrade services --all`

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  services --all
```

### Migrating services based on names
To migrate specific services in the app  `harness-upgrade services --names name1,name2`. Pass comma separated service names.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  services --names SERVICE_NAMES
```

### Migrating services based on ids
To migrate specific services in the app  `harness-upgrade services --ids id1,id2`. Pass comma separated service ids.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  services --ids SERVIE_IDS
```

## Migrating environments

### Migrating all environments
To migrate all environments in the app  `harness-upgrade environments --all`

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  environments --all
```

### Migrating environments based on names
To migrate specific environments in the app  `harness-upgrade environments --names name1,name2`. Pass comma separated environment names.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  environments --names ENVIRONMENT_NAMES
```

### Migrating environments based on ids
To migrate specific environments in the app  `harness-upgrade environments --ids id1,id2`. Pass comma separated environment ids.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  environments --ids ENVIRONMENT_IDS
```

## Migrating app level templates
To migrate templates from the app, use the following command `harness-upgrade --app APP_ID templates --all import`

### Migrating all templates
To migrate all templates on app level, use the following command

```shell
harness-upgrade --api-key SAT_API_KEY \ 
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --env ENV \
  templates --all import
```

### Migrating templates based on names
To migrate specific templates in the app  `harness-upgrade --app APP_ID templates --names name1,name2 import`. Pass comma separated template names.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  templates --names TEMPLATE_NAMES import
```

### Migrating templates based on ids
To migrate specific templates in the app  `harness-upgrade --app APP_ID templates --ids id1,id2 import`. Pass comma separated template ids.

```shell
harness-upgrade --api-key SAT_API_KEY \
  --project PROJECT \
  --org ORG \
  --account ACCOUNT_ID \
  --app APP_ID \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  --workflow-scope SCOPE \
  --env ENV \
  templates --ids TEMPLATE_IDS import
```

:::tip
Instead of providing the flags, you can save all the flags into a YAML file and load them using `--load FILE`
```yaml
env: Prod
api-key: sat.kmpySmUISimoRrJL6NL73w....
account: kmpySmUISimoRrJL6NL73w
app: APP_ID
project: demo
org: default
secret-scope: account
connector-scope: account
template-scope: account
workflow-scope: project
```

We can then do this -
```shell
harness-upgrade --load file.yaml app --all
```

```shell
harness-upgrade --load file.yaml workflows --all
```
:::
