---
sidebar_position: 4
slug: /concepts
---

# Key Concepts
Before we get started with the upgrade process, it is important to understand the what, why & how behind the upgrade tool.

## Idempotent
The upgrade tool is idempotent. Once an entity is created in next gen the upgrade tool will not create or update another entity in the same scope. It is also deterministic as it always creates the entities with same name & identifier always.

:::info

If you wish to update one or more entities then delete them in next gen and rerun the upgrade. 

:::

:::caution Exception
Exception to this are Service & Environment overrides. As they get updated everytime.
:::

## Migrates referenced entities
The upgrade can migrate one or more entities at once. It ensure to migrate an entities it is referencing and it's subsequent entities. Take for example we want to migrate a Pipeline. A pipeline references workflows, the workflows internally may be referencing services, environments etc. 
The services reference connectors via artifact sources, the connectors reference secrets and secrets reference secret manager. So a simple migration of Pipeline will not only create a pipeline but it will also create connectors, secrets etc.

![/images/mig-graph.png](images/mig-graph.png)

## Scopes
This is more of next generation concept than upgrade. But it is a key concept to be aware of. In first gen connectors & secrets were always account level entities. It was not possible to create them at application level. 
In next gen we have org & project scopes that give more control on how these entities are managed. The upgrade tool provides the ability to scope entities to account, org & project so that one can leverage these capabilities early on. 


| CG Entity        | NG Entity                                | Default Scope | Recommended Scope | Configurable |   
|------------------|------------------------------------------|---------------|-------------------|--------------|
| Secret Manager   | Connectors                               | Project       | Account           | Yes          |  
| Cloud Providers  | Connectors                               | Project       | Account           | Yes          |
| Connectors       | Connectors                               | Project       | Account           | Yes          |
| Templates        | Step Templates, Artifact Source Template | Project       | Account           | Yes          |
 | User Groups      | User Groups                              | Account       | Account           | No           |
| Service          | Service                                  | Project       | Project           | No           |
| Environment      | Environment                              | Project       | Project           | No           |
| Infrastructure   | Infrastructure                           | Project       | Project           | No           |
| Workflows        | Stage Templates or Pipelines             | Project       | Project           | Yes          |
| Pipeline         | Pipeline                                 | Project       | Project           | No           |
| Inline Manifests | File Store                               | Project       | Project           | No           |
| Triggers         | Triggers                                 | Project       | Project           | No           |
| App Defaults     | Project level Variables                  | Project       | Project           | No           |

:::tip
Before you get started take some time to think of how you would like to structure you account in next gen. 
:::

:::caution
Once the scopes of entities during the migration is set it should not be changed after that. If different scopes are provided each time, you will end up with same First Gen entity in multiple NG scopes.
:::