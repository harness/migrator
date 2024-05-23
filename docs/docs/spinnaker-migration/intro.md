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

|    | Stages                        | Harness           | Can we Migrate? | Comments                                               |
|----|-------------------------------|-------------------|-----------------|--------------------------------------------------------|
| 1  | pipeline                      | Pipeline Chaining | Yes             | Can only chain up to 2 pipelines                       |
| 2  | manualJudgment                | Approval          | Yes             |                                                        |
| 3  | checkPreconditions            | Shell Script      | Yes             | Supports Expressions, and Multiple Check preconditions |
| 4  | deleteManifest                | Delete            | Yes             |                                                        |
| 5  | evaluateVariables             | JEXL Expression   | Yes             |                                                        |
| 6  | wait                          | Wait              | Yes             |                                                        |
| 7  | deployManifest                | Deploy            | Yes             |                                                        |
| 8  | webhook                       | HTTP Step         | Yes             |                                                        |
| 9  | findImageFromTags             | Shell Script      | Yes             |                                                        |
| 10 | bake                          | Bake Plugin       | Yes             | Only supports AWS cloud provider currently             |
| 11 | jenkins                       | Jenkins Step      | Yes             |                                                        |
| 12 | findArtifactsFromResource     | Shell Script      | Yes             |                                                        |
| 13 | AWS: invokeLambda             | Plugin            | Yes             |                                                        |
| 14 | ShrinkCluster                 | Plugin            | Yes             |                                                        |
| 15 | ScaleDownCluster              | Plugin            | Yes             |                                                        |
| 16 | ResizeServerGroup             | Plugin            | Yes             |                                                        |
| 17 | DisableCluster                | Plugin            | Yes             |                                                        |
| 18 | ArtifactoryPromote            | HTTP Step         | Yes             |                                                        |
| 19 | DestroyServerGroup            | Plugin            | Yes             |                                                        |
| 20 | Deploy AMI Backed with Packer | Plugin            | Yes             |                                                        |
