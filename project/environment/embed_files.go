package environment

import _ "embed"

//go:embed .gitignore
var Gitignore string

//go:embed .golangci.yaml
var GolangCIConfig string

//go:embed autocomplete.sh
var AutocompleteShellScript string

//go:embed helper.sh
var HelperShellScript string

//go:embed README.md
var ReadmeMd string
