---
sidebar_position: 20
slug: /faqs
---

# FAQs

#### If I run the migration API/CLI with the same input multiple times, what happens?

The migration process is designed to be idempotent, meaning that if you run the migration with the same input multiple times, there will be no issues or duplicate entities such as connectors, secrets, or pipelines created.

#### How are identifiers generated during the migration process?

Most special characters are removed from the entity name and the expression is converted to camel case. The `--identifier-format` flag can be used to configure the identifier format, with the default being camelCase. Identifiers are currently limited to a maximum of 128 characters. This process of generating identifiers is one of the reasons why the migration process can achieve idempotency.

#### Will the migration process update already migrated FG entities?

No, the migration process will not update already migrated FG entities. If you wish to update the entity for some reason, you must delete the entity in NG and try again.

#### What happens if there is an entity with the same identifier that has already been created?

The migration process for that entity will fail. However, if another entity references the entity that has the same identifier, we will reference the existing entity. For example, if there is a connector with the name Test in FG and a connector with the identifier test in NG, we will reference the existing NG connector when migrating an application with an artifact source in one of the services using the FG connector.

#### Is it possible to migrate entities across accounts/clusters?

Yes we support across accounts and cross cluster migration. Please refer to the [Cross Account](advanced/cross-account) section for more details. 

#### How can I migrate expressions used in Git remote manifests?

The Harness Upgrade tool provides functionality to migrate expressions from a remote manifest. Additional details on this process can be found in the [Advanced section on remote manifest](advanced/remote-manifest).

