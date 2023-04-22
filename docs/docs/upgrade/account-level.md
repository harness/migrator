---
sidebar_position: 2
---

# Migrating Account level entities

Before migrating entities from FirstGen to NextGen, it is important to understand the entities present in the account. Run the account level summary to get a summary report.

## Migrating all account level entities
To migrate all secrets, connectors & cloud providers from an account we use the following command `harness-upgrade account`
```shell
harness-upgrade --api-key SAT_API_KEY --project PROJECT --org ORG --account ACCOUNT_ID --secret-scope SCOPE --connector-scope SCOPE --template-scope SCOPE --env ENV account
```

## Migrating secrets
To migrate secrets from an account, use the following command `harness-upgrade secrets`

### Migrating all secrets
To migrate all secrets in an account, use the following command:

```shell
harness-upgrade --api-key SAT_API_KEY \ 
--project PROJECT \
--org ORG \
--account ACCOUNT_ID \
--secret-scope SCOPE \
--connector-scope SCOPE \
--template-scope SCOPE \
--env ENV secrets
```

### Migrating specific secrets by names
To migrate specific secrets by names, pass the comma-separated names of the secrets using the `--names` flag

```shell
harness-upgrade --api-key SAT_API_KEY \ 
--project PROJECT \
--org ORG \
--account ACCOUNT_ID \
--secret-scope SCOPE \
--connector-scope SCOPE \
--template-scope SCOPE \
--env ENV secrets --names name1,name2
```

### Migrating specific secrets by ids
To migrate specific secrets by ids, pass the comma-separated ids of the secrets using the `--ids` flag

```shell
harness-upgrade --api-key SAT_API_KEY \ 
--project PROJECT \
--org ORG \
--account ACCOUNT_ID \
--secret-scope SCOPE \
--connector-scope SCOPE \
--template-scope SCOPE \
--env ENV secrets --ids id1,id2
```

## Migrating connectors & cloud providers
To migrate connectors from an account, use the following command `harness-upgrade connectors`

### Migrating all connectors
To migrate all connectors in an account, use the following command

```shell
harness-upgrade --api-key SAT_API_KEY \ 
--project PROJECT \
--org ORG \
--account ACCOUNT_ID \
--secret-scope SCOPE \
--connector-scope SCOPE \
--template-scope SCOPE \
--env ENV connectors
```

### Migrating specific connectors by names
To migrate specific connectors by names, pass the comma-separated names of the connectors using the `--names` flag:

```shell
harness-upgrade --api-key SAT_API_KEY \ 
--project PROJECT \
--org ORG \
--account ACCOUNT_ID \
--secret-scope SCOPE \
--connector-scope SCOPE \
--template-scope SCOPE \
--env ENV connectors --names name1,name2
```

### Migrating specific connectors by ids
To migrate specific connectors by ids, pass the comma-separated ids of the connectors using the `--ids` flag:

```shell
harness-upgrade --api-key SAT_API_KEY \ 
--project PROJECT \
--org ORG \
--account ACCOUNT_ID \
--secret-scope SCOPE \
--connector-scope SCOPE \
--template-scope SCOPE \
--env ENV connectors --ids id1,id2
```

:::tip
Instead of providing the flag, you can save all the flags into a YAML file and load them using `--load FILE`
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
harness-upgrade --load file.yaml account
```
:::