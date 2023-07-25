---
sidebar_position: 4
---

# User groups
When migrating from FirstGen to NextGen, it is possible to create user groups and add users to them. However, due to significant differences in permissions between the two platforms, the permissions will not be automatically migrated.
Use the following command - 

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV user-groups --all
```

By default, the user groups will be migrated to account. Use the below command if you want to scope the user groups to a specific org or project.

```shell
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV --org ORG --project PROJECT --user-group-scope SCOPE user-groups --all
```

If you wish to migrate specific user groups you can use either the `--ids` or `--names` flag. The `--ids` flag takes a comma separated list of user group ids. The `--names` flag takes a comma separated list of user group names. If both flags are provided, the `--ids` flag will be used.

```shell
# By ID
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV user-groups --ids "user-group-id-1,user-group-id-2"

# By name
harness-upgrade --api-key SAT_API_KEY --account ACCOUNT_ID --env ENV user-groups --names "User Group 1,User Group 2"
```


:::tip
Instead of providing the flags inline, you can save all the flags into a YAML file and load them using `--load FILE`
```yaml
env: Prod
api-key: sat.kmpySmUISimoRrJL6NL73w....
account: kmpySmUISimoRrJL6NL73w
app: APP_ID
project: demo
org: default
user-group-scope: project
```

We can then do this -
```shell
harness-upgrade --load file.yaml user-groups --all
```
:::