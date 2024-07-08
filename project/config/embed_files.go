package config

import _ "embed"

//go:embed aws/config.ex
var AwsConfigExample string

//go:embed aws/credentials.ex
var AwsCredentialsExample string

//go:embed .gitignore
var Gitignore string

//go:embed main.yaml
var MainConfig string

//go:embed main-local.yaml.ex
var MainLocalConfigExample string
