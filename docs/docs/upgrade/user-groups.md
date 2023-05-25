---
sidebar_position: 4
---

# User groups
When migrating from FirstGen to NextGen, it is possible to create user groups and add users to them. However, due to significant differences in permissions between the two platforms, the permissions will not be automatically migrated.
Use the following command - 

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV user-groups
```

Use the below command if you want to scope the user groups to a specific org or project.

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV --org ORG --project PROJECT --user-group-scope SCOPE user-groups
```
