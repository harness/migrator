---
sidebar_position: 4
---

# Advance Commands

## 1. Import Spinnaker pipeline to existing project
`
harness-upgrade --load migrator-config.yml pipelines import
`
### You should get output like this with some prompts:
```shell
INFO[2024-03-04T15:50:38-08:00] Importing the application....
? Please provide the Spinnaker pipeline name :  Spinnaker Pipeline name
INFO[2024-03-04T15:50:38-08:00]
Migration details:
  Platform: spinnaker
  Spinnaker Host: harness-eval.dynamo-staging.eng.armory.io
  App name: Spinnaker Application Name
  Pipeline Name: Spinnaker Pipeline name
  Authentication method: basic
  Insecure: false
? Do you want to proceed with pipeline migration? Yes
INFO[2024-03-04T15:50:38-08:00] Project check details :                       Account=AccountID OrgIdentifier=default ProjectIdentifier=ProjectID
INFO[2024-03-04T15:50:38-08:00] Project with identifier ProjectID does not exist
INFO[2024-03-04T15:50:38-08:00] Project with identifier ProjectID created
INFO[2024-03-04T15:50:38-08:00] {Pipeline Logs}
INFO[2024-03-04T15:50:41-08:00] Spinnaker migration completed
```

