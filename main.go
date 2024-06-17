// Package main - contains usage example of go-initial-project
package main

import (
	"os"

	"github.com/roman-kart/go-initial-project/project"
	"github.com/roman-kart/go-initial-project/project/tools"
)

func main() {
	test()
}

func test() {
	rootPath := tools.GetRootPath()
	configFolder := rootPath + string(os.PathSeparator) + "project" + string(os.PathSeparator) + "config"
	app, cleanup, err := project.InitializeApplication(configFolder)

	defer cleanup()

	tools.PanicOnError(err)

	app.Logger.Logger.Info("Starting application")
}
