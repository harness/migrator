---
sidebar_position: 2
slug: /install
---

# Installation
Download the latest release from GitHub releases. We support MacOS(`darwim & amd64`), Linux(`linux + (amd64/arm64)`) and Windows(`windows+amd64`). Please download the right assets. Extract the file anywhere.
We recommend that you move it to a folder that is specified in your path.

```shell
mv harness-upgrade /somepath/
harness-upgrade help
```

If you are using macOS then just do
```shell
mv harness-upgrade /usr/local/bin/
harness-upgrade help
```

If the above works successfully you should see all the commands that are supported with `harness-upgrade`

## Update the CLI
To update to a new version of the CLI
```shell
harness-upgrade update
```
