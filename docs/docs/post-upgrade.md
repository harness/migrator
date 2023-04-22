---
sidebar_position: 9
slug: /post-upgrade
---

# Post Upgrade

- Expressions

  - Update expressions in remote manifests using our tool

  - Review summary of expressions not migrated

  - Secret expressions using secrets.getValue need to use secret ID instead of name and a scoping prefix

- Review pipelines and stages

  - Adjust runtime inputs vs expressions and fixed values, as desired

  - Factor out common stage variables to pipeline variables

- Execute a cut-over deployment with same version

- Enable NextGen triggers, disable FirstGen triggers

- Disable FirstGen access