---
sidebar_position: 2
---

# Utility Commands
During the upgrade process it sometimes requires managing next gen entities like project & organisations. 

## Org Management

### Create an org
We can use the `harness-upgrade org create` command to create an org.
```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  org --name ORG_NAME --identifier ORG_IDENTIFIER create  
```  

### Remove organisations
The following command removes organisations from an account. You can provide the names or identifiers of the organisations.

To remove organisations by name use `--names`
```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \ 
  org --names name1,name2 rm  
```  

To remove organisations by identifier use `--identifiers`
```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  org --identifiers identifier1,identifier2 rm  
```  

## Project Management

### Create a project
We can use the `harness-upgrade project create` command to create project.
```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  --org ORG \
  project --name PROJECT_NAME --identifier PROJECT_IDENTIFIER create  
```

### Bulk create projects for every app in First Gen
Projects in NextGen are counterparts to applications from FirstGen. So a common requirement is to create projects in NextGen with the same name as application from FirstGen.
The following command creates a corresponding project for every app in the account. It then exports a YAML file for every project to the specified export path(defaults to current dir) specified.

```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  --org ORG \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  project --export FOLDER_PATH create-bulk  
```

:::info Note
The above command does not migrate entities it only creates projects and export YAML. To create entities in NextGen please run the application migration
:::

:::tip Automate
If you wish to migrate all the above apps then run this -
```shell  
for f in /path/to/folder/*.yaml; do harness-upgrade --load $f app --all; done  
```
:::

### Bulk create projects using CSV
When migrating from FirstGen to NextGen there may be a requirements to create projects for a FirstGen app in different organisations rather than the same org. 
The bulk command earlier let us create all apps as projects into a single org. The CSV provides more control as to which org fo projects get created in.

You can use a CSV containing the mapping for first gen applications to a next gen projects.

Generate a template that contains application name, project name, project identifier & org identifier. We default project name to application name, project identifier defaults to camelCase format of the application name & org is default. You can modify the csv if you want to customize them.

```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  project --csv PATH_TO_CSV csv-template  
```  

Sample contents of csv file
```text  
AppName,ProjectName,ProjectIdentifier,OrgIdentifier  
Demo,Demo,demo,default  
Test App,Test App,testApp,default  
```  

You can then create & generate yaml file based on the above CSV

```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  --secret-scope SCOPE \
  --connector-scope SCOPE \
  --template-scope SCOPE \
  project --csv PATH_TO_CSV --export FOLDER_PATH create-bulk  
```  

### Remove projects
The following command removes projects from a given org in an account. You can provide the names or identifiers of the projects.

To remove projects by name use `--names`
```shell  
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  --org ORG \
  project --names name1,name2 rm  
```  

To remove projects by identifier use `--identifiers`
```shell
harness-upgrade --api-key SAT_API_KEY \
  --account ACCOUNT_ID \
  --env ENV \
  --org ORG \
  project --identifiers identifier1,identifier2 rm  
```  
