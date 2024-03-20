---
sidebar_position: 4
---

# Advance Commands

## 1. Import Spinnaker pipeline to existing project
`
harness-upgrade --load migrator-config.yml pipelines --pipeline-name pipeline_name import
`
### You should get output like this with some prompts:
```shell
INFO[2024-03-04T15:50:38-08:00] Importing the application....
INFO[2024-03-04T15:50:38-08:00]
Migration details:
  Platform: spinnaker
  Spinnaker Host: harness-eval.dynamo-staging.eng.armory.io
  App name: prasadtest
  Pipeline Name: pipeline_name
  Authentication method: basic 
  Insecure: false 
? Do you want to proceed with pipeline migration? Yes
INFO[2024-03-04T15:50:41-08:00] Spinnaker pipeline migration completed
```

