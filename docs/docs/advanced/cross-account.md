---
sidebar_position: 3
---

# Cross Account

For most accounts, the upgrade needs to happen within an account. But for a few accounts, there may be a need to create entities from one account to another. This could be the case when there are different accounts for different organizations of a company. In NextGen, it could be an effort to bring them into a single account and manage them as different organizations.

We can achieve this by using the `--target-account` & `--target-api-key` flags. All the commands like `account`, `app`, `workflows`, `pipelines` etc. can leverage these flags. The `--account` & `--api-key` will refer to the FirstGen account ID and API key.

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
We currently do not support cross-cluster upgrade.

:::