---
sidebar_position: 1
---

# Override Names and Identifiers

When Harness upgrades entities from FirstGen to NextGen we try to keep the name same as that of First Gen(sometime we strip-out character such as `(,),{,},[,]` etc.) and for identifiers we convert the name to camel case(Hello World -> helloWorld).
This works out best for most cases. In few cases there may be a need to control the identifier generation logic to a different format(e.g all lower case) and in such cases we provide `--identifier-format` flag to override the default logic.

There may be times when we want a much more granular control over the generation of name & identifier. For example in first gen a connector may be named as `Harness Docker Registry` and in next gen we want to name it as `Docker Registry` and the identifier should be `docker_registry`. 
If you would like to achieve this then we can leverage the `--override` flag. This flag takes a file as an input and the contents of the file look as following - 

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

The type allows the following values - `SECRET`, `CONNECTOR`, `SERVICE`, `ENVIRONMENT`, `WORKFLOW`, `PIPELINE`, `TEMPLATE`.

Here is an example on how to use it with the upgrade tool - 

```shell
# Import app from FirstGen to NextGen
harness-upgrade --load file.yaml --override overrides.yaml app --all
```

You can use the flag with any command that are importing entities such as services, workflows, pipelines, environments, connectors, templates, secrets, triggers etc.

:::caution

It is important that the flag is used with all commands that are importing entities. If is not used then we revert to using the default logic for name and identifier generation. This may be lead to duplicate and inconsistent entities in NextGen.

:::
