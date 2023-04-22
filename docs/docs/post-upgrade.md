---
sidebar_position: 9
slug: /post-upgrade
---

# Post Upgrade

After the upgrade, some additional steps need to be taken to ensure a smooth transition. Here are the steps that need to be taken:

### Expressions

Expressions are an essential part of Harness, and they need to be updated post-upgrade. Here's what you need to do:

- Update expressions in remote manifests using our tool. This tool migrates expressions from remote manifests and provides a summary of expressions not migrated.
- Please go over the summary post upgrade to replace the FirstGen expressions with equivalent next-gen expression

### Review pipelines and stages

Review your pipelines and stages to ensure they are running smoothly post-upgrade. Here's what you need to do:

- Adjust runtime inputs vs. expressions and fixed values, as desired.
- Factor out common stage variables to pipeline variables.

### Execute a cut-over deployment with the same version

Execute a deployment with the same versions as the one you have currently deployed. This ensures that everything is working as expected.

### Enable NextGen triggers, disable FirstGen triggers

After the upgrade, you need to enable NextGen triggers and disable FirstGen triggers. This is essential to ensure that your Harness pipelines are running smoothly.

### Disable FirstGen access

Finally, disable FirstGen access. This ensures that all new workflows are created in NextGen and not in FirstGen.
