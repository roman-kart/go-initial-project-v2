#!/bin/bash

SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")
GO_MODULE_URL="github.com/roman-kart/go-initial-project/v2"

function rand-str {
    # Return random alpha-numeric string of given LENGTH
    #
    # Usage: VALUE=$(rand-str $LENGTH)
    #    or: VALUE=$(rand-str)

    local DEFAULT_LENGTH=64
    local LENGTH=${1:-$DEFAULT_LENGTH}

    LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c $LENGTH
    # LC_ALL=C: required for Mac OS X - https://unix.stackexchange.com/a/363194/403075
    # -dc: delete complementary set == delete all except given set
}

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
FIX="FALSE"

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
    --fix)
      FIX="TRUE"
      shift  1
      ;;
    *)
      echo "Error: Invalid argument: '$1'"
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
  --fix     - Use automatic issue fixing
EOF
    else
      arguments=""
          echo "Performing testing"
          if [ "$IS_VERBOSE" != "FALSE" ]
          then
            arguments="$arguments --verbose"
          fi

          if [ "$FIX" != "FALSE" ]
          then
            arguments="$arguments --fix"
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
  lint-check-autofix)
    if [ "$IS_HELP" != "FALSE" ]; then
      cat << EOF
Usage: ./helper.sh lint (<argument>...)

Helps to check which of autofix linters breaks down the code.

IMPORTANT: commit all your changes before executing this command - it will erase all your changes
EOF
    fi

    if [[ $(git status --porcelain) ]]; then
      echo "Commit all your changes before executing this command"
      exit 1
    fi

    autofix_linters=(
      "gci"
      "gocritic"
      "godot"
      "gofmt"
      "gofumpt"
      "goheader"
      "goimports"
      "mirror"
      "misspell"
      "nolintlint"
      "protogetter"
      "tagalign"
      "whitespace"
    )

    current_git_branch=$(git rev-parse --abbrev-ref HEAD)
    temporal_git_branch=$(git rev-parse  --abbrev-ref HEAD)
    while [[ $(git rev-parse --verify "$temporal_git_branch" 2>/dev/null) ]]; do
      temporal_git_branch="lint_check_autofix_$(rand-str 32)"
    done

    echo "Current branch: $current_git_branch"
    echo "Temporal branch: $temporal_git_branch"

    git checkout -b  "$temporal_git_branch"

    lintersWithChanges=()

    for autofix_linter in "${autofix_linters[@]}"; do
        echo "Start checking $autofix_linter"
        golangci-lint.exe run --enable "$autofix_linter" --fix
        if [[ $(git status --porcelain) ]]; then
          echo "Find new changes after linting!"
          echo "Check is these changes are suitable. Use: golangci-lint.exe run --enable \"$autofix_linter\" --fix"
          lintersWithChanges+=( "$autofix_linter" )

          git add . && git commit -m 'tmp: lint-check-autofix temporal commit' && previous_git_commit=$(git rev-parse HEAD^) && echo "Previous commit: $previous_git_commit" && git reset --hard "$previous_git_commit"
        fi
    done

    git checkout "$current_git_branch"
    git branch -D "$temporal_git_branch"

    if ! [ ${#lintersWithChanges[@]} -eq 0 ]; then
      echo "Linters which made changes:"
      for  autofix_linter in "${lintersWithChanges[@]}"; do
        echo "$autofix_linter"
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
  gotest - Run all GO-tests
EOF
    ;;
  gotest)
    if [ "$IS_HELP" != "FALSE" ]; then
          cat << EOF
Usage: ./helper.sh test

Arguments:
  --verbose - Enable verbose output (default FALSE)
EOF
    else
      args=""

      if  [ "$IS_VERBOSE" != "FALSE" ]; then
        args="$args -v"
      fi

      # shellcheck disable=SC2086
      go test $args ./project/tests/
    fi
    ;;
  *)
    echo "Error: Invalid action"
    exit 1
esac

exit 0
