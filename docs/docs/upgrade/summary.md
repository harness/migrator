---
sidebar_position: 1
---

# Summary Report

Before we move entities from first gen to next gen it is important to get a gist of entities that are present in the account.

## Account Summary
This generate a report that summarises all the entities that are present in an account. For few entities we try and summarise the effort required with the upgrade. 
We can use the `harness-upgrade account-summary` command to generate the account summary 

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV account-summary
```

## Application Summary
Similar to account summary we generate an application summary. A report gets generated that summarises all the entities that are being used by an application.
We can use the `harness-upgrade application-summary` command to generate the application summary.

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --app APP_ID --env ENV application-summary
```

:::info

These commands are only generating the summary of entities present in the account or application. They are not creating entities in Next Gen.

:::