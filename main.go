package main

import (
	"fmt"
	"os"

	"github.com/roman-kart/go-initial-project/project"

	"github.com/roman-kart/go-initial-project/project/utils"
)

func main() {
	test()
}

func test() {
	rootPath := utils.GetRootPath()
	configFolder := rootPath + string(os.PathSeparator) + "project" + string(os.PathSeparator) + "config"
	app, cleanup, err := project.InitializeApplication(configFolder)
	defer cleanup()
	utils.PanicOnError(err)

	fmt.Println(fmt.Sprintf("%+v", app.Config))
}
