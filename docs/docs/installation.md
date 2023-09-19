---
sidebar_position: 2
slug: /install
---

# Installation

## Unix systems 
For unix systems run this command to download the latest release & install - 
```shell
curl https://raw.githubusercontent.com/harness/migrator/master/install | sh
```

If the above fails, please try installing using the manual installation approach.

## Manual install
To install the Harness CD First Gen to CD Next Gen migration tool manually, follow these steps:

1. Download the latest release from the GitHub releases page. The tool supports MacOS (`darwin` + (`amd64/arm64`)), Linux (`linux` + (`amd64`/`arm64`)), and Windows (`windows`+`amd64`) platforms, so make sure to download the correct asset for your platform.
2. Extract the downloaded file to a directory of your choice. It is recommended that you move the extracted file to a folder specified in your system's path for ease of use.
3. Run the `harness-upgrade help` command to verify that the installation was successful.

If you are using macOS, you can move the harness-upgrade file to the `/usr/local/bin/` directory by running the following command:

```shell
mv harness-upgrade /usr/local/bin/
```
Then, run the harness-upgrade help command to verify that the installation was successful.

To update the CLI to a new version, run the following command:

```shell
harness-upgrade update
```
This will update the tool to the latest version.