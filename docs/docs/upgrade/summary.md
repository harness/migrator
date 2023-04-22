---
sidebar_position: 1
---

# Summary Report

Before migrating entities from FirstGen to NextGen, it is important to have a summary report of the entities present in the account or application.

## Account Summary
The account summary report provides an overview of all the entities present in an account. For some entities, it also summarizes the effort required for upgrading them.
To generate the account summary, use the following command `harness-upgrade account-summary` 

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV account-summary
```

## Application Summary
Similar to the account summary, the application summary generates a report summarizing all the entities being used by an application.
To generate the application summary, use the following command `harness-upgrade application-summary`

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --app APP_ID --env ENV application-summary
```

:::info
Please note that these commands only generate a summary report of entities present in the account or application. They do not create any entities in NextGen.
:::