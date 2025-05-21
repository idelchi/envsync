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

`envprof` is a CLI tool for managing named environment profiles in `YAML` or `TOML`.

Supports profile inheritance (layering) and importing of `.env` files.

## Features

- Define multiple environment profiles in a single YAML or TOML file, with inheritance and dotenv support
- List profiles, export to `.env` files or the current shell, or spawn a subshell with the selected environment

## Installation

For a quick installation, you can use the provided installation script:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/envprof/refs/heads/dev/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
# list all profiles
envprof list
```

```sh
# list all variables in a profile with inheritance information
envprof list dev -v
```

```sh
# list a specific variable
envprof list dev HOST
```

```sh
# export profile to a file
envprof export dev .env
```

```sh
# spawn a subshell with the environment loaded
envprof shell dev
```

```sh
# export to current shell
eval "$(envprof export dev)"
```

## Format

Complex types (arrays, maps) are serialized as JSON; all other values are simple strings.

### YAML

```yaml
dev:
  dotenv:
    - secrets.env
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
dotenv = ['secrets.env']
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

Inheritance is resolved in order: later profiles override earlier ones.
`dotenv` files have the lowest priority and load first, before applying profile layers.

As an example, running `envprof export dev .env` with the previous YAML definition
as well as a sample `secrets.env`:

```sh
TOKEN=secret
```

produces the following `.env` file:

```sh
# Active profile: "dev"
DEBUG=true
HOST=localhost
PORT=80
TOKEN=secret
```

Inspecting with `envprof list dev -v` would show the variables and their origins:

```sh
DEBUG=true              (inherited from "staging")
HOST=localhost
PORT=80                 (inherited from "prod")
TOKEN=secret            (inherited from "secrets.env")
```

The layering chain here is:

```sh
secrets.env -> prod -> staging -> dev
```

from lowest to highest priority (left to right).

## Flags

All commands accept the following flag:

```sh
--file, -f
```

which can be used to specify a file (or a list of fallback files) to load.

Defaults to the first found among `envprof.yaml`, `envprof.yml`, or `envprof.toml`, unless `ENVPROF_FILE` is set.

## Subcommands

<details>
<summary><strong>list / ls</strong> — List profiles or variables</summary>

- **Usage:**

  - `envprof list [--verbose/-v] [profile] [variable]`

- **Flags:**
  - `--verbose`, `-v` – Show variable origins

</details>

<details>
<summary><strong>export / x</strong> — Export profile to file or stdout</summary>

- **Usage:**

  - `envprof export [--prefix <string>] <profile> [file]`

- **Flags:**
    <!-- markdownlint-disable MD038 -->
  - `--prefix` – String to prefix variables (default: `export `)
    <!-- markdownlint-enable MD038 -->
    </details>

<details>
<summary><strong>shell / sh</strong> — Spawn a subshell with profile</summary>

- **Usage:**

  - `envprof shell [--isolate/-i] [--shell <string>] <profile>`

- **Flags:**
  - `--isolate`, `-i` – Prevent inheriting current shell variables
  - `--shell` – Force shell (default empty string -> detected)

</details>

## Shell integration

When using the `shell` subcommand, `envprof` sets `ENVPROF_ACTIVE_SHELL` in the environment -
use it for customizing your prompt.

An example for `starship.toml` would be:

```toml
format = """${env_var.envprof}$all"""

[env_var.envprof]
variable = "ENVPROF_ACTIVE_SHELL"
format = "\\[envprof: $env_value\\]($style) "
```

This variable is also used to detect if you’re already in an `envprof` subshell, preventing nested sessions.

For convenience, define a shell function to quickly switch profiles:

```sh
envprof-activate() {
  local output
  if output="$(envprof export "${1}" 2>&1)"; then
    eval "${output}"
  else
    echo "${output}" >&2
  fi
}
```

Use `envprof-activate dev` to switch to the `dev` profile.

> [!NOTE]
> This will export variables into your current shell, potentially overwriting existing ones.
> Repeated use will also mix the variables from different profiles.

## Demo

![Demo](assets/gifs/envprof.gif)
