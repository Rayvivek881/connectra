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
	connections.InitS3()
	defer connections.CloseDB()

	connections.InitSearchEngine()

	cmd.Execute()
}
