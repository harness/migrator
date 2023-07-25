---
sidebar_position: 2
---

# Provide custom expressions

First Gen & Next Gen have different expression syntaxes. For example, in First Gen, we use `${variable}` to reference a variable, 
but in Next Gen, we use `<+variable>`. For most expressions, the CLI can automatically convert the syntax from First Gen to Next Gen. 
This ensures that your pipelines will work as expected after the migration. 
In some cases, however, we cannot determine the equivalent value for an expression, in such cases the expressions are left untouched. 
Additional to this the auto generated expression may not be the right expression in your setup. In such cases there is a need to provide custom expressions.
You can provide custom expressions using the `--custom-expressions` flag. This flag takes a YAML file as input. 

The YAML file should have the following structure:

```yaml
appName: <+org.name>
service.name: <+project.name>
context.variable_exported.download_url: <+stage.variables.download_url>
```

The usage of the flag is as follows:

```shell
harness-upgrade --load file.yaml --custom-expressions custom.yaml app --all
```

You can use the `--custom-expressions` flag with any command that imports entities, such as services, workflows, pipelines, environments, connectors, templates, secrets, triggers, etc.
