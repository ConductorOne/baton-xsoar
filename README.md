![Baton Logo](./docs/images/baton-logo.png)

# `baton-demisto` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-demisto.svg)](https://pkg.go.dev/github.com/conductorone/baton-demisto) ![main ci](https://github.com/conductorone/baton-demisto/actions/workflows/main.yaml/badge.svg)

`baton-demisto` is a connector for Demisto (Cortex XSOAR) built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with created instance of Cortex XSOAR API to sync data about users and their roles.

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Prerequisites

To use this connector, you need to install Cortex XSOAR (Demisto) instance and create API key to access the API.

More information about installation of Cortex XSOAR can be found in [Cortex XSOAR administration guide](https://docs-cortex.paloaltonetworks.com/r/Cortex-XSOAR/6.6/Cortex-XSOAR-Administrator-Guide) for example [this guide](https://docs-cortex.paloaltonetworks.com/r/Cortex-XSOAR/6.6/Cortex-XSOAR-Administrator-Guide/Install-a-Server-for-a-Single-Server-Deployment) for installing server for single server deployment.

To obtain API key, login to webview of Cortex XSOAR instance and go to `Settings` -> `Integrations` -> `API Keys` and click on `Get Your Key` button. You can find more information about API keys [here](https://docs-cortex.paloaltonetworks.com/r/Cortex-XSOAR/6.6/Cortex-XSOAR-Administrator-Guide/API-Keys).

# Getting Started

The instance comes by default with invalid SSL certificate, so to bypass validation you have to set `BATON_UNSAFE` environment variable to `true` or use `--unsafe` flag.

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-demisto

BATON_TOKEN=token baton-demisto
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_TOKEN=token BATON_UNSAFE=true ghcr.io/conductorone/baton-demisto:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-demisto/cmd/baton-demisto@main

BATON_TOKEN=token baton-demisto
baton resources
```

# Data Model

`baton-demisto` will fetch information about the following resources:

- Users
- Roles

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-demisto` Command Line Usage

```
baton-demisto

Usage:
  baton-demisto [flags]
  baton-demisto [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string              The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string          The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                   The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                          help for baton-demisto
      --log-format string             The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string              The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --token string                  Access token used to connect to the Cortex XSOAR API. ($BATON_TOKEN)
      --unsafe                        Allow insecure TLS connections to Cortex XSOAR instance. ($BATON_UNSAFE)
  -v, --version                       version for baton-demisto

Use "baton-demisto [command] --help" for more information about a command.

```
