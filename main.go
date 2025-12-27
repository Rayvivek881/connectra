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
	connections.InitS3()

	cmd.Execute()
}
