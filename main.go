package main

import (
	"fmt"
	"os"

	"github.com/roman-kart/go-initial-project/project/config"

	"github.com/roman-kart/go-initial-project/project/managers"
	"github.com/roman-kart/go-initial-project/project/utils"
)

func main() {
	test()
}

func test() {
	fmt.Println("test")
	fields, err := utils.RetrieveClickhouseTags(managers.ApplicationStatsModel{})
	utils.PanicOnError(err)
	options, err := utils.BuildOptionsFromClickhouseTags(fields)
	utils.PanicOnError(err)
	fmt.Println(fmt.Sprintf("%+v", options))

	rootPath := utils.GetRootPath()
	configPath := rootPath + string(os.PathSeparator) + "project" + string(os.PathSeparator) + "config"
	cfg, err := config.NewConfig(configPath)
	utils.PanicOnError(err)
	fmt.Println(fmt.Sprintf("%+v", cfg))
}
