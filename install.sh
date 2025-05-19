#!/bin/sh
set -e

GITHUB_TOKEN=${ENVPROF_GITHUB_TOKEN:-${GITHUB_TOKEN}}
DISABLE_SSL=${ENVPROF_DISABLE_SSL:-${DISABLE_SSL}}

# Usage function
usage() {
  cat <<EOF
Usage: ${0} [OPTIONS]
Install envprof.

All arguments are passed to the installation script command. See below for details.

Environment variables:

  ENVPROF_GITHUB_TOKEN/GITHUB_TOKEN       GitHub token to use for downloading assets from GitHub.
  ENVPROF_DISABLE_SSL/DISABLE_SSL         Disable SSL verification when downloading assets.

Example:

    curl -sSL https://raw.githubusercontent.com/idelchi/envprof/refs/heads/main/install.sh | GITHUB_TOKEN=<token> sh -s -- -d ~/.local/bin

EOF
  curl ${ENVPROF_DISABLE_SSL:+-k} -sSL https://raw.githubusercontent.com/idelchi/scripts/refs/heads/main/install.sh | INSTALLER_TOOL="envprof" sh -s -- -h

  exit 1
}

# Parse arguments
parse_args() {
  # Handle known options with getopts
  while getopts ":h" opt; do
    case "${opt}" in
      h) usage ;;
      *) : ;; # Do nothing for other options
    esac
  done

  # Reset option parsing
  shift $((OPTIND - 1))
  OPTIND=1
}

# Download envprof
install() {
  # Download tools using envprof
  curl ${DISABLE_SSL:+-k} -sSL https://raw.githubusercontent.com/idelchi/scripts/refs/heads/main/install.sh | ENVPROF_GITHUB_TOKEN=${GITHUB_TOKEN} INSTALLER_TOOL="envprof" sh -s -- "$@"
}

need_cmd() {
  if ! command -v "${1}" >/dev/null 2>&1; then
    printf "Required command '%s' not found\n" "${1}"
    exit 1
  fi
}

main() {
  parse_args "$@"

  # Check for required commands
  need_cmd curl

  # Install tools
  install "$@"
}

main "$@"
