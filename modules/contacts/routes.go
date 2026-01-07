package contacts

import (
	"vivek-ray/modules/contacts/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/:uuid", controller.GetContactByUUID)
	router.POST("/", controller.GetContactsByFilter)
	router.POST("/count", controller.GetContactsCountByFilter)
	router.POST("/create", controller.CreateContact)
	router.PUT("/:uuid", controller.UpdateContact)
	router.DELETE("/:uuid", controller.DeleteContact)
	router.POST("/upsert", controller.UpsertContact)
	router.POST("/batch-upsert", controller.BatchUpsert)
}
