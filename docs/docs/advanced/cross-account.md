---
sidebar_position: 3
---

# Cross Account

For most accounts, the upgrade needs to happen within an account. But for a few accounts, there may be a need to create entities from one account to another. This could be the case when there are different accounts for different organizations of a company. In NextGen, it could be an effort to bring them into a single account and manage them as different organizations.
There are commonly two use cases. 
1. Migrate entities from one account to another within same cluster. For example `Account1` in `Prod1` to `Account2` in `Prod1`.
2. Migrate entities across clusters/installation. For example `Account1` in `Prod1` to `Account2` in `Prod3` or `Account1` in `SelfManaged` to `Account2` in `Prod2`.

## Same Cluster
If the two accounts are in within same installation/cluster this technique can be used. We can achieve this by using the `--target-account` & `--target-api-key` flags. All the commands like `account`, `app`, `workflows`, `pipelines` etc. can leverage these flags. The `--account` & `--api-key` will refer to the FirstGen account ID and API key.

**Example -**

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

:::note

The two accounts need to be in the same cluster for this to work. For example, if the FirstGen account is in Prod1 & the NextGen target account is in Prod2, this will not work.

:::

## Different Cluster
There are certain scenarios that would require migration from one cluster to another or one installation to another. Here are few scenarios - 
1. You would like to move from SelfManaged FirstGen to SaaS Nextgen.
2. You would like to move from Prod1 FirstGen cluster to Prod3 Nextgen cluster.
3. You would like to merge two FirstGen SelfManaged installations to a single NextGen SelfManaged installation.

We can also use this technique to have a single accounts for multiple organizations that are spread across different SaaS clusters.
We can achieve this by using the `--target-account`, `--target-api-key` & `--target-gateway-url` flags. All the commands like `account`, `app`, `workflows`, `pipelines` etc. can leverage these flags. The `--account` & `--api-key` will refer to the FirstGen account ID and API key.

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

:::warn

This is a beta feature. Not all utility commands support this. We are working on adding support for all commands.

:::
