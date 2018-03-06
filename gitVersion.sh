#!/bin/bash

# Define paths relative to location of this script.

ORIGINAL_DIR=$PWD
SCRIPT_DIR=$( cd "$(dirname "${BASH_SOURCE}")" ; pwd -P )

BOLD=$(tput bold)
BLACK=$(tput setaf 0)
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
YELLOW=$(tput setaf 3)
BLUE=$(tput setaf 4)
MAGENTA=$(tput setaf 5)
CYAN=$(tput setaf 6)
WHITE=$(tput setaf 7)
NORM=$(tput sgr0)  # Important to have normal in last position.


function onerror {
  echo "[$(date)] `basename $0`: ${RED}ERROR${NORM}"
}

function usage_and_clean_exit() {
  trap - EXIT
  usage
  cd "${ORIGINAL_DIR}" || true
  exit 0
}


function usage_and_error_exit() {
  trap - EXIT
  usage
  cd "${ORIGINAL_DIR}" || true
  exit 1
}


function clean_exit() {
  trap - EXIT
  cd "${ORIGINAL_DIR}" || true
  exit 0
}


function error_exit() {
  trap - EXIT
  usage
  echo "[$(date)] `basename $0`: ${RED}Error: $1 ${NORM}"
  cd "${ORIGINAL_DIR}" || true
  exit 1
}

# Define default command-line option values.
# Environment variables, if set, become default.

if [ -z "$GITVERSION_APPEND" ]; then
  GITVERSION_APPEND=0
fi
if [ -z "$GITVERSION_OUTPUT_FILE" ]; then
  GITVERSION_OUTPUT_FILE="/dev/stdout"
fi
if [ -z "$GITVERSION_OUTPUT_FORMAT" ]; then
  GITVERSION_OUTPUT_FORMAT="sem"
fi
if [ -z "$GITVERSION_REVISION" ]; then
  GITVERSION_REVISION=0
fi
if [ -z "$GITVERSION_VERBOSE" ]; then
  GITVERSION_VERBOSE=0
fi

# Set error handling.

trap onerror EXIT

# After these commands, any unset variables will trigger an error.
# See https://vaneyckt.io/posts/safer_bash_scripts_with_set_euxo_pipefail/

# set -x
set -euo pipefail

# =============================================================================
#   Help
# =============================================================================

function usage() {

  cat << EOF
usage: `basename $0` [-h] [-a] [-v] [-f output-format] [-r revision-number] [-o output-file]

Get version information from git in the current working directory.

OPTIONS:
   -h   Show this message.
   -a   Append to file (versus overwrite file)
   -f   Output format.  Examples:
          hpp - Creates C++ *.hpp file for VERSION and BIX_VERSION_STR
          rev - Current revision number
          sem - Current semantic version
          sha - SHA of a revision specified by -r parameter
   -o   Output file. Example /tmp/output.txt
   -v   Verbose output
   -r   Revision number. Used with '-f sha'

EXAMPLES:
   `basename $0`
       Returns "M.m.P[-iter.<sha>]" semantic version to STDOUT.
   `basename $0` -o test.txt
       Creates / overwrites test.txt with "M.m.P" contents.
   `basename $0` -f hpp -o src/common/version.hpp
       Creates / overwrites src/common/version.hpp with "static..." contents.
   `basename $0` -f sha -r 10
       Returns the SHA of the 10th revision.
EOF
}

# =============================================================================
#   Git queries
# =============================================================================

# Returns M.m.P-B-xxxxxxxx
function git_describe_always_tags() {
  echo $(git describe --always --tags)
}


# Returns non-zero if the workspace is dirty, or "0" if it is clean.
function git_is_workspace_dirty() {
  expr `git status --porcelain 2>/dev/null| wc -l`
}


# Returns number of builds.
function git_max_revision_number() {
  echo $(git rev-list HEAD | wc -l)
}


# Given a build number, return the SHA of the version.
function git_sha_for_revision_number() {
  echo $(git rev-list --reverse HEAD | nl | awk "{ if(\$1 == "$1") { print \$2 }}")
}


# Returns M.m.P
function git_tag() {
  read TAG COMMIT_SINCE_TAG SHA1 <<< $(echo $(git_describe_always_tags) | awk -F "-" '{print $1 " " $2 " " substr($3, 2, length($3))}')
  echo "${TAG}"
}

# =============================================================================
#   Output functions
# =============================================================================

# Returns 'static const char *BIX_VERSION_STR="<SemVer>";'
function hpp_output() {
  VERSION=$(semver_output)
  echo "static const char *BIX_VERSION_STR=\"$VERSION\";\n"
}


# Returns the number of revisions.
function max_revision_number_output() {
  echo "$(git_max_revision_number)"
}


# Returns SemVer.
function semver_output() {
  TAG=$(git log --oneline --decorate=short | grep 'tag:' | awk -F"tag: " '{sub(/[ |,|\)].*/,"",$2); print $2}' | grep '^[[:digit:]].[[:digit:]].[[:digit:]]' | sort --reverse | head -n 1)
  VERSION=$TAG
  ITERATION=$(git log ${TAG}..HEAD --oneline | wc -l | tr -d '[:space:]')
  if [ "${ITERATION}" != "0" ]; then
    SHA=$(git describe --always)
    if [[ $VERSION = *"-"* ]]; then
      VERSION=${VERSION}.${SHA}
    else
      VERSION=${VERSION}-${ITERATION}.${SHA}
    fi
    VERSION=$(echo "${VERSION}" | sed "s/dirty\-/dirty\./")
  fi
  if [ "$(git_is_workspace_dirty)" != "0" ]; then
    VERSION=${VERSION}.dirty
  fi
  echo "${VERSION}"
}


# If possible, return absolute path of file.
function absolute_path() {
  UNAME_STR=$(uname)
  if [[ "${UNAME_STR}" == 'Linux' ]]; then
    echo $(readlink --canonicalize $1);
  elif [[ "${UNAME_STR}" == 'Darwin' ]]; then
    echo $1;
  else
    echo $1;
  fi
 }


# Given a revision number, returns the SHA of the version.
function sha_for_revision_output() {
  echo "$(git_sha_for_revision_number $1)"
}


# Returns ???? - this is just for testing.
function test_output() {
  echo "$(git_sha_for_revision_number $1)"
}

# =============================================================================
#   MAIN
# =============================================================================

# Test for prerequisites.

if ! type "git" > /dev/null; then
  error_exit "ERROR: 'git' is not installed."
fi

# Parse command-line options.

OPTIND=1  # A POSIX variable. Reset in case getopts has been used previously in the shell.
while getopts "af:h?o:r:v" OPTION
do
  case ${OPTION} in
    h) usage_and_clean_exit ;;
    a) GITVERSION_APPEND=1 ;;
    f) GITVERSION_OUTPUT_FORMAT="$OPTARG" ;;
    o) GITVERSION_OUTPUT_FILE=$(absolute_path $OPTARG) ;;
    r) GITVERSION_REVISION="$OPTARG" ;;
    v) GITVERSION_VERBOSE=1 ;;
    *) error_exit "ERROR: Incorrect option specified." ;;
  esac
done

# Based on output format request (-o), call the correct sub-routine.

case ${GITVERSION_OUTPUT_FORMAT} in
  bix)
    OUTPUT="$(bix_output)"
    ;;
  hpp)
    OUTPUT="$(hpp_output)"
    ;;
  rev)
    OUTPUT="$(max_revision_number_output)"
    ;;
  sem)
    OUTPUT="$(semver_output)"
    ;;
  sha)
    OUTPUT="$(sha_for_revision_output ${GITVERSION_REVISION})"
    ;;
  test)
    OUTPUT="$(test_output ${GITVERSION_REVISION})"
    ;;
  *)
    error_exit "Option '-f ${GITVERSION_OUTPUT_FORMAT}' not supported."
    ;;
esac

# Print verbose output.

if [ ${GITVERSION_VERBOSE} != 0 ]; then
   echo ${GREEN}Output format: ${BOLD}${GITVERSION_OUTPUT_FORMAT}${NORM}
   echo ${GREEN}Output file: ${BOLD}${GITVERSION_OUTPUT_FILE}${NORM}
   echo ${GREEN}Output: ${BOLD}${OUTPUT}${NORM}
fi

# Create or append to output file. Output file may be /dev/stdout.

if [ ${GITVERSION_APPEND} == 0 ]; then
  printf "${OUTPUT}" > "${GITVERSION_OUTPUT_FILE}"
else
  printf "${OUTPUT}" >> "${GITVERSION_OUTPUT_FILE}"
fi

# Epilog

if [ ${GITVERSION_OUTPUT_FILE} != "/dev/stdout" ]; then
  echo "Output: ${GITVERSION_OUTPUT_FILE}"
fi

clean_exit
