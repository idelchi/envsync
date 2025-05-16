#!/bin/sh
set -e

GITHUB_TOKEN=${ENVSYNC_GITHUB_TOKEN:-${GITHUB_TOKEN}}
DISABLE_SSL=${ENVSYNC_DISABLE_SSL:-${DISABLE_SSL}}

# Usage function
usage() {
  cat <<EOF
Usage: ${0} [OPTIONS]
Install envsync.

All arguments are passed to the installation script command. See below for details.

Environment variables:

  ENVSYNC_GITHUB_TOKEN/GITHUB_TOKEN       GitHub token to use for downloading assets from GitHub.
  ENVSYNC_DISABLE_SSL/DISABLE_SSL         Disable SSL verification when downloading assets.

Example:

    curl -sSL https://raw.githubusercontent.com/idelchi/envsync/refs/heads/dev/install.sh | GITHUB_TOKEN=<token> sh -s -- -d ~/.local/bin

EOF
  curl ${ENVSYNC_DISABLE_SSL:+-k} -sSL https://raw.githubusercontent.com/idelchi/scripts/refs/heads/dev/install.sh | INSTALLER_TOOL="envsync" sh -s -- -h

  exit 1
}

# Parse arguments
parse_args() {
  # Handle known options with getopts
  while getopts ":h" opt; do
    case "${opt}" in
      h) usage ;;
    esac
    shift $((OPTIND - 1))
    OPTIND=1
  done
}

# Download envsync
install() {
  local args="${1}"

  # Download tools using envsync
  curl ${DISABLE_SSL:+-k} -sSL https://raw.githubusercontent.com/idelchi/scripts/refs/heads/dev/install.sh | ENVSYNC_GITHUB_TOKEN=${GITHUB_TOKEN} INSTALLER_TOOL="envsync" sh -s -- ${args}
}

need_cmd() {
  if ! command -v "${1}" >/dev/null 2>&1; then
    printf "Required command '${1}' not found"
    exit 1
  fi
}

main() {
  parse_args "$@"

  # Check for required commands
  need_cmd curl

  # Install tools
  args="$@"
  install "${args}"
}

main "$@"
