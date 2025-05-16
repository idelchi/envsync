<p align="center">
  <h3 align="center"><code>envsync</code></h3>
  <p align="center">Profile-based environment variable manager</p>
</p>

---

[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/envsync.svg)](https://pkg.go.dev/github.com/idelchi/envsync)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/envsync)](https://goreportcard.com/report/github.com/idelchi/envsync)
[![Build Status](https://github.com/idelchi/envsync/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/envsync/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`envsync` is a CLI tool for managing named environment profiles in YAML, TOML, or JSON. It supports layering (inheritance) of profiles and is designed for clarity and shell integration.

## Features

- Define multiple environment profiles in one file
- Supports YAML, TOML, and JSON
- Profile inheritance with conflict resolution
- Export to shell using `eval "$(envsync export <profile>)"`
- Apply profiles to current environment or write to `.env` files

## Usage

```sh
# set a variable
envsync set dev DB_URL postgres://localhost:5432

# read a variable
envsync get dev DB_URL

# list all variables (after inheritance)
envsync list dev

# export for shell eval
eval "$(envsync export dev)"

# apply profile and write to .env
envsync apply dev .env

# remove a key
envsync remove dev DB_URL

# delete entire profile
envsync delete dev

# list all profiles
envsync profiles
```

## Format

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

### JSON

```json
{
  "base": {
    "env": {
      "DB_HOST": "localhost",
      "DB_PORT": "5432"
    }
  },
  "dev": {
    "extends": ["base"],
    "env": {
      "DB_NAME": "devdb",
      "DEBUG": "true"
    }
  }
}
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

- `base.env` applied
- then `common.env` (overrides base on conflict)
- then `staging.env` (final override)
