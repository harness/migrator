---
title: Introduction
description: Introduction to the Spinnaker Migration tool. Includes supported Spinnaker stages.
sidebar_position: 1
---

# Introduction

This tool was developed by Harness to assist in the migration process from Spinnaker to CD Next Gen. It can migrate specific Harness resources or multiple resources at once.

### Who is the Tool intended for?

This tool is designed for customers, CSMs, and developers who are currently using Spinnaker and wish to migrate their accounts to Harness CD Next Gen.

### Spinnaker stages that can be migrated
The following key components can be migrated using this tool:

| **Spinnaker Stages**              | **Harness**                                                                                                                                    | **Release** | **Comments**                                               |
| --------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- | ----------- | ---------------------------------------------------------- |
| Wait                              | Wait                                                                                                                                           |             |                                                            |
| Bake Manifest                     | [AMI Plugin](https://hub.docker.com/r/harnessdev/aws-bake-deploy-ami-plugin) |             |                                                            |
| Deploy Manifest                   | Deploy                                                                                                                                         |             |                                                            |
| Delete Manifest                   | Delete                                                                                                                                         |             |                                                            |
| Enable Manifest                   | Plugin                                                                                                                                         | 1.91.0      |                                                            |
| Disable Manifest                  | Plugin                                                                                                                                         | 1.91.0      |                                                            |
| Manual Judgement                  | Approval                                                                                                                                       |             |                                                            |
| Pipeline                          | Pipeline Chaining                                                                                                                              |             | Can only chain up to 2 pipelines                           |
| Find Artifacts From Resource      | Shell Script                                                                                                                                   |             |                                                            |
| Find Image From Tags              | Shell Script                                                                                                                                   |             |                                                            |
| Find Image                        | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                                                                                            |             |                                                            |
| Evaluate variables                | JEXL Expression                                                                                                                                |             |                                                            |
| Check Preconditions               | Shell Script                                                                                                                                   |             | Supports Expressions, and Multiple Check preconditions     |
| Jenkins                           | Jenkins Step                                                                                                                                   |             |                                                            |
| Bake                              | [AMI Plugin](https://hub.docker.com/r/harnessdev/aws-bake-deploy-ami-plugin)                                                                   |             | For AMI Baking, Only supports AWS cloud provider currently |
| Deploy AMI Backed with Packer     | [AMI Plugin](https://hub.docker.com/r/harnessdev/aws-bake-deploy-ami-plugin)                                                                   |             |                                                            |
| AWS.LambdaInvokeStage             | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               |             |                                                            |
| AWS.LambdaUpdateCodeStage         | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       |             |                                                            |
| Aws.LambdaDeploymentStage         | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       | 1.84.0      |                                                            |
| Aws.LambdaTrafficRoutingStage     | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               | 1.84.0      |                                                            |
| ShrinkCluster                     | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       |             |                                                            |
| ScaleDownCluster                  | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       |             |                                                            |
| ResizeServerGroup                 | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       |             |                                                            |
| DisableCluster                    | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       |             |                                                            |
| ArtifactoryPromote                | HTTP Step                                                                                                                                      |             |                                                            |
| DestroyServerGroup                | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               |             |                                                            |
| Webhook                           | HTTP Step                                                                                                                                      |             |                                                            |
| EnableServerGroup                 | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               |             |                                                            |
| RollbackCluster                   | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               |             |                                                            |
| disableServerGroup                | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               | 1.82.0      |                                                            |
| modifyAwsScalingProcess           | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               | 1.85.0      |                                                            |
| undoRolloutManifest               | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       | 1.85.0      |                                                            |
| rollingRestartManifest            | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       | 1.85.0      |                                                            |
| scaleManifest                     | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       | 1.85.0      |                                                            |
| cloneServerGroup                  | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       | 1.88.0      |                                                            |
| addJiras                          | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       | 1.89.0      |                                                            |
| runJobManifest                    | CD Stage, Kubernetes Apply Step                                                                                                                | 1.89.0      |                                                            |
| determineTargetServerGroup        | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)               |             | Implemented indirectly as a part of the CloneServerGroup   |
| applySourceServerGroupCapacity    | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)                                                                       |             | Implemented indirectly as a part of the CloneServerGroup   |
| upsertImageTags                   | Shell Script                                                                                                                           | 1.92.0      |                                                            |
| githubAddIssueLabel               | HTTP Step                                                                                                                                      | 1.92.0      |                                                            |
| githubDeleteIssueLabel            | HTTP Step                                                                                                                                      | 1.92.0      |                                                            |
| githubStatus                      | HTTP Step                                                                                                                                      | 1.92.0      |                                                            |
| githubAddIssueComment             | HTTP Step                                                                                                                                      | 1.92.0      |                                                            |
| restrictExecutionDuringTimeWindow |                                                                                                                                                | 1.93.0      |                                                            |
