---
sidebar_position: 2
---

# Provide custom string replacement

You can provide custom string replacement information using the `--replace` flag. This flag takes a YAML file as input. 

The YAML file should have the following structure:

```yaml
- old: "this is old"
  new: "this is new string"
- old: "another old string"
  new: "another new string"
```

The usage of the flag is as follows:

```shell
harness-upgrade --load file.yaml --replace replace.yaml app --all
```

