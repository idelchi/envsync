<p align="center">
  <img alt="envprof logo" src="assets/envprof.png" height="150" />
  <h3 align="center"><code>envprof</code></h3>
  <p align="center">Profile-based environment variable manager</p>
</p>

---

[![GitHub release](https://img.shields.io/github/v/release/idelchi/envprof)](https://github.com/idelchi/envprof/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/envprof.svg)](https://pkg.go.dev/github.com/idelchi/envprof)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/envprof)](https://goreportcard.com/report/github.com/idelchi/envprof)
[![Build Status](https://github.com/idelchi/envprof/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/envprof/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`envprof` is a CLI tool for managing named environment profiles in YAML or TOML.
It supports layering (inheritance) of profiles and is designed for clarity and shell integration.

## Features

- Define multiple environment profiles in one file
- Supports YAML and TOML formats
- Profile inheritance with conflict resolution
- Export profiles to shell using `eval "$(envprof export <profile>)"`
- Export profiles to `.env` files

## Installation

For a quick installation, you can use the provided installation script:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/envprof/refs/heads/main/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
# list all profiles
envprof list

# list all variables in a profile
envprof list dev

# list a specific variable
envprof list dev DB_URL

# export for shell eval
eval "$(envprof export dev)"

# export profile to a file
envprof export dev .env
```

The following flags are available:

- `--file`, `-f`: Specify a file (or list of fallbacks) to load.
  Defaults to `ENVPROF_FILE` or `envprof.yaml, envprof.yml, envprof.toml`.
- `--verbose`, `-v`: Enable verbose output to trace inheritance for variables.

## Format

Complex types (arrays, maps) are serialized as JSON representations.

### YAML

```yaml
base:
  env:
    DB_HOST: localhost
    DB_PORT: "5432"

dev:
  extends: [base]
  env:
    DB_NAME: devdb
    DEBUG: "true"
```

### TOML

```toml
[base.env]
DB_HOST = "localhost"
DB_PORT = "5432"

[dev]
extends = ["base"]

[dev.env]
DB_NAME = "devdb"
DEBUG = "true"
```

## Inheritance Behavior

Inheritance is resolved in order. Later profiles override earlier ones. For example:

```yaml
staging:
  extends:
    - base
    - common
  env:
    DEBUG: "false"
```

Results in:

- `base` applied
- then `common` (overrides base on conflict)
- then `staging` (final override)

## Destructive commands

Below commands are destructive (reformat the profiles file)
and should just be used to set up the initial file if you care about comments and formatting.

```sh
# import values from a dotenv file into a profile (creating it if it doesn't exist),
# with the option whether to keep existing values in case of conflicts or to overwrite them
envprof import dev [--keep] .env

# convert current file to another format
envprof convert [yaml|toml]
```
