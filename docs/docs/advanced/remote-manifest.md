---
sidebar_position: 1
---

# Remote Manifests
It is a very common practice to manage K8s manifests, ECS definitions from a git remote file. These manifest files often use first gen expressions within them. Since Git remote files are outside the control of harness these files are not migrated.
In order to convert first gen expressions to next gen expressions we can use the upgrade tool to achieve the same.

The command replaces first gen expressions found in all files in current & subdirectories. By default, we only process the files with following file extensions `yml`, `yaml` & `json`.
```shell
harness-upgrade expressions 
```

Do a dry run on the files without replacing any CG expressions
```shell
harness-upgrade expressions --dry-run
```

To provide custom files extensions
```shell
harness-upgrade expressions --extensions yml,txt,xml
```

Secrets referenced in these files are converted to camel case to align with the migrator. You can provide the scope of the secrets using the `--secret-scope` flag.
```shell
harness-upgrade --secret-scope account expressions
```

To provide custom expressions or override default expressions
```shell
harness-upgrade --custom-expressions /path/to/file.yaml expressions 
```

The above command expects a yaml file whose contents are simple key values of first-gen expressions & string to replace that expression with. A sample override file would look like this -
```yaml
appName: <+org.name>
service.name: <+project.name>
```


