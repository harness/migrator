---
sidebar_position: 1
---

# Overrides & Settings

## Override Names and Identifiers

During the upgrade process from FirstGen to NextGen in Harness, we usually keep the names of entities the same as in FirstGen (sometimes removing certain characters such as `(`,`)`,`{`,`}`,`[`,`]`). For identifiers, we convert the names to camel case (e.g., `Hello World` -> `helloWorld`).
This approach generally works well, but in some cases, there may be a need to customize the identifier generation logic to a different format, such as using all lowercase letters. In such situations, you can use the `--identifier-format` flag to override the default logic.

However, there are scenarios where more granular control over the generation of names and identifiers is required. For example, in FirstGen, a connector may be named `Harness Docker Registry`, but in NextGen, we want to name it `Docker Registry` with the identifier `docker_registry`. To achieve this level of customization, we can utilize the `--override` flag. This flag takes a file as input, and the file's contents should follow this structure:

```yaml
overrides:
  - name: Fetch Release Information    # Name of the entity in NextGen 
    identifier: fetch_release_info     # Identifier of the entity in NextGen
    type: TEMPLATE                     # Type of entity in FirstGen
    firstGenName: Release_Fetch_Info   # Name of entity in FirstGen
  - name: HarnessDockerRegistry
    identifier: Harness_Registry
    type: CONNECTOR
    id: pwrySmUISimoRrJL6Nsvbw         # ID of the entity in FirstGen
```

The type field allows the following values: `SECRET`, `CONNECTOR`, `SERVICE`, `ENVIRONMENT`, `WORKFLOW`, `PIPELINE`, `TEMPLATE`.

Here's an example of how to use the override flag with the upgrade tool:

```shell
# Importing app from FirstGen to NextGen
harness-upgrade --load file.yaml --override overrides.yaml app --all
```

You can use the `--override` flag with any command that imports entities, such as services, workflows, pipelines, environments, connectors, templates, secrets, triggers, user-groups etc.

:::caution

It's important to use the `--override` flag with all entity import commands. If it's not used, the default logic for name and identifier generation will be applied, potentially resulting in duplicate and inconsistent entities in NextGen.

:::

## Settings

Often times when we upgrade from FirstGen to NextGen we have few default behaviours that are considered. 
For example, we default to allow deployment to same infra as `false`. But in some cases you may want to override this behaviour & you would like to set that as `true`.

We have various different settings to control different behaviours. You leverage the same by using `--override` flag. You add a settings section in the overrides file & add the settings that you want to override. 

```yaml  
overrides:
  ...
settings:
  - type: SETTING_1
    value: VALUE
  - type: SETTING_2
    value: VALUE
```

The possible values for the setting types are  - 

| Setting Type                            | Description                                                                       | Allowed Values    |   
|-----------------------------------------|-----------------------------------------------------------------------------------|-------------------|
| `SIMULTANEOUS_DEPLOYMENT_ON_SAME_INFRA` | Used to control if you would like to allow simultaneous deployments to same infra | `true` or `false` |  
