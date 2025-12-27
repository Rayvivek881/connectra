package main

import (
	"vivek-ray/cmd"
	"vivek-ray/conf"
	"vivek-ray/connections"
)

func main() {
	v := conf.Viper{}
	v.Init()

	connections.InitDB()
	defer connections.CloseDB()

	connections.InitSearchEngine()
	defer connections.CloseSearchEngine()

	connections.InitS3()

	cmd.Execute()
}
