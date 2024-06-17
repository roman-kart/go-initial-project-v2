#!/bin/bash

# This is autocomplete script for helper.sh.
# To use it, add "source /path/to/autocomplete.sh" to your.bashrc file

_autocomplete() {
  helpArgs="--help -h"

  latest="${COMP_WORDS[$COMP_CWORD]}"
  prev="${COMP_WORDS[$COMP_CWORD - 1]}"
  words=""
  case "${prev}" in
    ./helper.sh)
      words="$helpArgs godoc lint"
      ;;
    lint)
      words="$helpArgs --verbose --fix"
      ;;
    godoc)
      words="$helpArgs --port"
      ;;
    *)
      ;;
  esac
  COMPREPLY=($(compgen -W "$words" -- $latest))
  return 0
}

complete -F _autocomplete ./helper.sh
