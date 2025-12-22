package companies

import (
	"vivek-ray/modules/companies/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.POST("/", controller.GetCompaniesByFilter)
	router.POST("/count", controller.GetCompaniesCountByFilter)

	router.GET("/filters", controller.GetFilters)
	router.PUT("/filters", controller.UpdateActiveStatus)
	router.POST("/filters/data", controller.GetFilterData)
}
