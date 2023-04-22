---
sidebar_position: 20
slug: /faqs
---

# FAQs

#### What happens if I run the migration API/CLI with the same input multiple times?

Nothing. Migration is designed to be idempotent. So you will not see duplicate connectors, secrets, pipelines, etc.

#### How do we generate the identifiers?

We remove most special characters from the name & convert the expression to camel case(configurable using `--identifier-format` flag, default is camelCase). This is one of the main reasons we can achieve idempotency. Please note that identifiers are currently capped at 128 characters.

#### Will the migration on the already migrated FG entity update the corresponding entity?

No. If you wish to update the entity for some reason, then please delete the entity in NG and try again.

#### What will happen if there is an already-created entity with the same identifier?

Migration for that entity will fail. If another entity references that entity, we will reference the existing one.

Let's take an example -

Let's say that a connector in NG has an identifier test. In FG, there is a connector with the name Test. Let's say the customer is migrating an application with an artifact source in one of the services using this FG connector. When migrating, we will see an error Failed to migrate Test connector, but the service migration will go through.

#### Does migration support migration across accounts/clusters?

Migration can migrate across accounts if they reside within the same cluster. If two accounts are present in different clusters like Prod1 & Prod2 then it is not possible.

#### We use remote manifests in Git. How can we automatically migrate the expressions used in the Git remote?

The Harness Upgrade tool provides functionality to migrate expressions from a remote manifest. More details can be found [here](advanced/remote-manifest).