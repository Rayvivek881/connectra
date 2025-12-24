package contacts

import (
	"vivek-ray/modules/contacts/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	// Read endpoints
	router.POST("/", controller.GetContactsByFilter)
	router.POST("/count", controller.GetContactsCountByFilter)

	// Filter endpoints
	router.GET("/filters", controller.GetFilters)
	router.PUT("/filters", controller.UpdateActiveStatus)
	router.POST("/filters/data", controller.GetFilterData)

	// Write endpoints
	router.POST("/create", controller.CreateContact)
	router.PUT("/:uuid", controller.UpdateContact)
	router.DELETE("/:uuid", controller.DeleteContact)
	router.POST("/upsert", controller.UpsertContact)
	router.POST("/bulk", controller.BulkUpsertContacts)
}
