package contacts

import (
	"vivek-ray/modules/contacts/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.POST("/", controller.GetContactsByFilter)
	router.POST("/count", controller.GetContactsCountByFilter)

	router.GET("/filters", controller.GetFilters)
	router.PUT("/filters", controller.UpdateActiveStatus)
	router.POST("/filters/data", controller.GetFilterData)
}
