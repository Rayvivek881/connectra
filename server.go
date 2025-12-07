package main

import (
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/modules/companies"
	"vivek-ray/modules/contacts"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {

	v := conf.Viper{}
	v.Init()

	connections.InitDB()
	connections.InitSearchEngine()

	router := gin.Default()
	router.Use(gin.Recovery())

	companies.Routes(router.Group("/companies"))
	contacts.Routes(router.Group("/contacts"))

	if err := router.Run(":8000"); err != nil {
		log.Error().Err(err).Msgf("Error starting server %v", err.Error())
		return
	}
}
