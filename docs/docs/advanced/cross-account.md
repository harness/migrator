---
sidebar_position: 3
---

# Cross Account

To support different scenarios where entities need to be migrated across accounts or clusters, we provide cross-account functionality. This is particularly useful when dealing with multiple organizations within a company or when consolidating entities into a single account in NextGen.

There are two common use cases for cross-account migration:

1. Migrating entities from one account to another within the same cluster, such as migrating from `Account1` in `Prod1` to `Account2` in `Prod1`.
2. Migrating entities across clusters or installations, for example, migrating from `Account1` in `Prod1` to `Account2` in `Prod3`, or migrating from `Account1` in `SelfManaged` to `Account2` in `Prod2`.

## Same Cluster

When both accounts are within the same installation/cluster, you can use the `--target-account` and `--target-api-key` flags to achieve cross-account migration. These flags can be used with commands like `account`, `app`, `workflows`, `pipelines`, etc. The `--account` and `--api-key` will refer to the FirstGen account ID and API key.

**Example:**

```shell
harness-upgrade --account FG_ACCOUNT --api-key FG_API_KEY --target-account NG_ACCOUNT --target-api-key NG_API_KEY
```

To make this seamless please create a yaml file & load it using the `--load FILE` flag.

```yaml
env: Prod
api-key: sat.kmpySmUISimoRrJL6NL73w...
account: kmpySmUISimoRrJL6NL73w
target-account: px7xd_BFRCi-pfWPYXVjvw
target-api-key: sat.px7xd_BFRCi-pfWPYXVjvw...
app: APP_ID
project: demo
org: default
secret-scope: org
connector-scope: org
template-scope: org
workflow-scope: project
```

:::caution

For this technique to work, both accounts must be in the same cluster. For example, If the FirstGen account is in `Prod1` and the NextGen target account is in `Prod2`, this method will not work.

:::

## Different Cluster
There are certain scenarios where migration is needed between clusters or installations. Here are a few examples:

1. Moving from SelfManaged FirstGen to SaaS NextGen.
2. Moving from Prod1 FirstGen cluster to Prod3 NextGen cluster.
3. Merging two FirstGen SelfManaged installations into a single NextGen SelfManaged installation.

This technique can also be used to have a single account for multiple organizations across different SaaS clusters.
To achieve this, you can use the `--target-account`, `--target-api-key`, and `--target-gateway-url` flags with commands like `account`, `app`, `workflows`, `pipelines`, etc. The `--account` and `--api-key` will refer to the FirstGen account ID and API key.

**Example -**

```shell
harness-upgrade --account FG_ACCOUNT --api-key FG_API_KEY --target-account NG_ACCOUNT --target-api-key NG_API_KEY --target-gateway-url TARGET_GATEWAY_URL
```

To make this seamless please create a yaml file & load it using the `--load FILE` flag.

```yaml
env: Prod
api-key: sat.kmpySmUISimoRrJL6NL73w...
account: kmpySmUISimoRrJL6NL73w
target-account: px7xd_BFRCi-pfWPYXVjvw
target-api-key: sat.px7xd_BFRCi-pfWPYXVjvw...
target-gateway-url: https://app3.harness.io/gateway
app: APP_ID
project: demo
org: default
secret-scope: org
connector-scope: org
template-scope: org
workflow-scope: project
```

:::caution

This is a beta feature. Not all utility commands support this. We are working on adding support for all commands.

:::
