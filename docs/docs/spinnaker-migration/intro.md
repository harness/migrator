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

|    | Stages                        | Harness                                                                      | Can we Migrate? | Comments                                                   |
|----|-------------------------------|------------------------------------------------------------------------------|-----------------|------------------------------------------------------------|
| 1  | Wait                          | Wait                                                                         | Yes             |                                                            |
| 2  | Bake Manifest                 | [AMI Plugin](https://hub.docker.com/r/harnessdev/aws-bake-deploy-ami-plugin) | Yes             |                                                            |
| 3  | Deploy Manifest               | Deploy                                                                       | Yes             |                                                            |
| 4  | Delete Manifest               | Delete                                                                       | Yes             |                                                            |
| 5  | Manual Judgement              | Approval                                                                     | Yes             |                                                            |
| 6  | Pipeline                      | Pipeline Chaining                                                            | Yes             | Can only chain up to 2 pipelines                           |
| 7  | Find Artifacts From Resource  | Shell Script                                                                 | Yes             |                                                            |
| 8  | Find Image From Tags          | Shell Script                                                                 | Yes             |                                                            |
| 9  | Find Image                    | Shell Script                                                                 | Yes             |                                                            |
| 10 | Evaluate Variables            | JEXL Expression                                                              | Yes             |                                                            |
| 11 | Check Preconditions           | Shell Script                                                                 | Yes             | Supports Expressions, and Multiple Check preconditions     |
| 12 | Jenkins                       | Jenkins Step                                                                 | Yes             |                                                            |
| 13 | Bake                          | [AMI Plugin](https://hub.docker.com/r/harnessdev/aws-bake-deploy-ami-plugin) | Yes             | For AMI Baking, Only supports AWS cloud provider currently |
| 14 | Deploy AMI Backed with Packer | [AMI Plugin](https://hub.docker.com/r/harnessdev/aws-bake-deploy-ami-plugin) | Yes             |                                                            |
| 15 | AWS.LambdaInvokeStage         | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 16 | AWS.LambdaUpdateCodeStage     | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 17 | ShrinkCluster                 | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 18 | ScaleDownCluster              | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 19 | DisableCluster                | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 20 | RollbackCluster               | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 21 | ArtifactoryPromote            | HTTP Step                                                                    | Yes             |                                                            |
| 22 | ResizeServerGroup             | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 23 | DestroyServerGroup            | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 24 | EnableServerGroup             | [AWS Drone Plugin](https://hub.docker.com/r/harnessdev/aws-drone-plugin)     | Yes             |                                                            |
| 25 | Webhook                       | HTTP step                                                                    | Yes             |                                                            |
