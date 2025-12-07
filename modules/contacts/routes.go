package contacts

import (
	"vivek-ray/modules/contacts/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/filters", controller.GetFilters)
	router.POST("/filters/data", controller.GetFilterData)

	router.POST("/contacts", controller.GetContactsByFilter)
	router.POST("/contacts/count", controller.GetContactsCountByFilter)
}
