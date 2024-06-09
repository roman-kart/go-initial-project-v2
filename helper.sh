#!/bin/bash

SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")
GO_MODULE_URL="github.com/roman-kart/go-initial-project"

# Check if no arguments were passed
# The number of arguments is stored in $#
if [ $# -eq 0 ]
then
  echo "No arguments were passed"
  exit 1
fi

# The first argument is the type of action
action=$1

# Shift the arguments so we can iterate over the rest
shift

IS_VERBOSE="FALSE"
PORT="8080"
IS_HELP="FALSE"

# Parse the remaining arguments
while (( "$#" )); do
  case "$1" in
    --verbose|-v)
      IS_VERBOSE="TRUE"
      shift 1
      ;;
    --port|-p)
      PORT="$1"
      shift  1
      ;;
    --help|-h)
      IS_HELP="TRUE"
      shift  1
      ;;
    *)
      echo "Error: Invalid argument"
      exit 1
  esac
done

# Perform the action
case $action in
  lint)
    if [ "$IS_HELP" !=  "FALSE" ]; then
      cat << EOF
Usage: ./helper.sh lint (<argument>...)

Arguments:
  --verbose - Enable verbose output
EOF
    else
      arguments=""
          echo "Performing testing"
          if [ "$IS_VERBOSE" != "FALSE" ]
          then
            arguments="$arguments --verbose"
          fi

          testingOutput=$(golangci-lint.exe run --enable-all $arguments)

          # shellcheck disable=SC2016
          : '
            Regexp explanation with example `gip\tools.go:67:10: error returned from external`
            s# - Replace compart after
            (^.*) - part before backslash (gip)
            (\\) - backslash (\)
            (.*\.go:[0-9]+:[0-9]*:?) - part after backslash and before error message text (tools.go:67:10:)
            #\1/\3 - change backslash to slash between first and third parts
            #g - globalk replacement
          '
          pattern='^(.*)(\\)(.*\.go:[0-9]+:[0-9]*:?)'
          # infinite loop because global replacement not working (IDK why)
          for (( ; ; ))
          do
              if echo "$testingOutput" | grep -Eq "$pattern"; then
                testingOutput=$(echo "$testingOutput" | sed -E "s#$pattern#\1/\3#g")
              else
                echo "$testingOutput"
                exit 0
              fi
          done
    fi
    ;;
  godoc)
    if [ "$IS_HELP" != "FALSE" ]; then
      cat << EOF
Usage: ./helper.sh godoc (<argument>...)

Arguments:
  --port  - provide custom port number (default 8080)
EOF
    else
      echo "Performing godoc. Port $PORT"
      echo "Url: http://localhost:$PORT/"
      echo "Url current go module: http://localhost:$PORT/pkg/$GO_MODULE_URL"
      godoc -http ":$PORT"
    fi
    ;;
  --help|-h)
    cat << EOF
Usage: ./helper.sh (<command>|--help|-h)

Documentation for command: ./helper.sh <command> (--help|-h)

Commands:
  lint   - Perform linting with golangci-lint
  godoc  - Start godoc server
EOF
    ;;
  *)
    echo "Error: Invalid action"
    exit 1
esac
