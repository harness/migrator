# Harness Upgrade

Harness has developed a tool to help user's migrate from Harness CD Current Gen to CD Next Generation. The tool can migrate specific Harness resources or bulk migration of resources. 

### Who is the Tool intended for?

- CLI to help customers, CSMs and developers with migrating their current gen harness account to next gen

### Key Components that can be migrated 

1. Pipelines
2. Workflows
3. Service (Service Definition - Manifests, variables included)
4. Environments (Overrides included)
5. Infrastructure Definition
6. Connectors
7. Secrets & Secret Managers
8. Templates

## Installation
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

If the above works successfully you should see all the commands that are supported with `harness-upgrade`. [Complete documentation can be found here](https://harness.github.io/migrator/)

## To update the CLI
To update to a new version of the CLI
```shell
harness-upgrade update
```

## Contact
If you face any issues please reach out to us or feel free to create a GitHub issue.