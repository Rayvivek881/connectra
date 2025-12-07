package companies

import (
	"vivek-ray/modules/companies/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/filters", controller.GetFilters)
	router.POST("/filters/data", controller.GetFilterData)

	router.POST("/companies", controller.GetCompaniesByFilter)
	router.POST("/companies/count", controller.GetCompaniesCountByFilter)
}
