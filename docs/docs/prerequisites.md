---
sidebar_position: 3
slug: /prerequisites
---

# Prerequisites

Before using the Harness CD First Gen to CD Next Gen migration tool, you must have the following prerequisites:

## Access Token

To use the migration tool, you need a [Service Account Token (SAT)](https://developer.harness.io/docs/platform/automation/api/add-and-manage-api-keys#create-service-account-api-keys-and-tokens) with the **Account Admin - All Resources Including Child Scopes** permissions. Ensure that you have the required permissions and create a SAT before using the tool.

## Account ID and Environment

You will need the account ID and environment information for your Harness account. You can find the account ID and environment information on the accounts overview page. To access the accounts overview page, click on **Account Settings** in the sidebar menu. The environment refers to the Harness cluster hosting your account and can be one of the following:

| Environment | Cluster                       |  
|-------------|-------------------------------|
| Prod        | Prod 1                        |  
| Prod        | Prod 2                        |
| Prod3       | Prod 3                        |
| SelfManaged | Harness self-managed platform | 

Make sure to note down the correct environment and cluster information before starting the migration process.

## SelfManaged Platform
If you are running on Harness self-managed platform please use the `--base-url https://ACME.com/gateway` flag or set the `base-url` field in the file YAML.
