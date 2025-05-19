<h1 align="center">envprof</h1>

<p align="center">
  <img alt="envprof logo" src="assets/images/envprof.png" height="150" />
  <p align="center">Profile-based environment variable manager</p>
</p>

---

[![GitHub release](https://img.shields.io/github/v/release/idelchi/envprof)](https://github.com/idelchi/envprof/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/envprof.svg)](https://pkg.go.dev/github.com/idelchi/envprof)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/envprof)](https://goreportcard.com/report/github.com/idelchi/envprof)
[![Build Status](https://github.com/idelchi/envprof/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/envprof/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`envprof` is a CLI tool for managing named environment profiles in YAML or TOML.
Supports layering (inheritance) of profiles and imports of `.env` files.

## Features

- Define multiple environment profiles in one file in YAML or TOML, with profile inheritance and dotenv imports
- List profiles, export to a `.env` file, the current shell, or enter a subshell with the environment loaded

## Installation

For a quick installation, you can use the provided installation script:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/envprof/refs/heads/main/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
# list all profiles
envprof list

# list all variables in a profile with inheritance information
envprof list dev -v

# list a specific variable
envprof list dev DB_URL

# export profile to a file
envprof export dev .env

# enter a subshell with the environment loaded,
# optionally inheriting environment variables from the parent shell
envprof shell [--inherit] [--shell <shell|detected>] dev

# or export to current shell
eval "$(envprof export [--prefix=export] dev)"
```

## Format

Complex types (arrays, maps) are serialized as JSON representations, everything else is stringified.

### YAML

```yaml
dev:
  dotenv:
    - .env
  extends:
    - staging
  env:
    HOST: localhost

staging:
  extends:
    - prod
  env:
    HOST: staging.example.com
    DEBUG: true

prod:
  env:
    HOST: prod.example.com
    PORT: 80
    DEBUG: false
```

### TOML

```toml
[dev]
extends = ['staging']
dotenv = ['.env']
[dev.env]
HOST = 'localhost'

[staging]
extends = ['prod']
[staging.env]
DEBUG = true
HOST = 'staging.example.com'

[prod.env]
DEBUG = false
HOST = 'prod.example.com'
PORT = 80
```

## Inheritance Behavior

Inheritance is resolved in order. Later profiles override earlier ones. `dotenv` entries
have least priority and are loaded in the order they are defined, before layering the profiles on top.

As an example, running `envprof export dev dev.env` with the previous YAML definition
as well as a sample `.env`:

```sh
TOKEN=secret
```

produces the following `dev.env` file:

```sh
# Active profile: "dev"
DEBUG=true
TOKEN=secret
HOST=localhost
PORT=80
```

Inspecting with `envprof list dev -v` would show the inheritance chain:

```sh
DEBUG=true              (inherited from "staging")
TOKEN=secret            (inherited from ".env")
HOST=localhost
PORT=80                 (inherited from "prod")
```

The inheritance chain is:

```sh
.env -> prod -> staging -> dev
```

from lowest to highest priority.

## Subcommands

```sh
list/ls [--verbose/-v] [profile] [variable]
```

List all profiles, all variables in a profile, or a specific variable in a profile.
`--verbose` shows from which source each variable is inherited.

```sh
export/x [--prefix] <profile> [file]
```

Print out the environment variables of a profile to stdout (with `--prefix` defaulting to `export`) or export them to a file.

```sh
shell/sh [--inherit] [--shell <shell|detected>] <profile>
```

Enter a subshell with the environment loaded, optionally inheriting environment variables from the parent shell.
The shell can be specified or detected automatically.

## Flags

The following flag is available on all commands:

```sh
--file, -f
```

Specify a file (or list of fallbacks) to load.

Defaults to the first found of `envprof.yaml`, `envprof.yml`, or `envprof.toml`, unless `ENVPROF_FILE` is set.

## Shell integration

When using the `shell` subcommand, `envprof` exports `ENVPROF_ACTIVE_SHELL` to the active shell session,
which can be used to customize the prompt.

An example for `starship.toml` would be:

```toml
format = """${env_var.envprof}$all"""

[env_var.envprof]
variable = "ENVPROF_ACTIVE_SHELL"
format = "\\[envprof: $env_value\\]($style) "
```

The same variable is used to detect if `envprof` is running in a subshell, and prevent entering a new one (before exiting the current one).

You can also define a function to quickly switch profiles in your shell:

```sh
envprofx() {
  local output
  if output="$(envprof export "${1}" 2>&1)"; then
    eval "${output}"
  else
    echo "${output}" >&2
  fi
}
```

and use as `envprofx dev` to switch to the `dev` profile.

## Demo

![Demo](assets/gifs/envprof.gif)
