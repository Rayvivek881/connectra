package companies

import (
	"vivek-ray/modules/companies/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	// Read endpoints
	router.POST("/", controller.GetCompaniesByFilter)
	router.POST("/count", controller.GetCompaniesCountByFilter)

	// Filter endpoints
	router.GET("/filters", controller.GetFilters)
	router.PUT("/filters", controller.UpdateActiveStatus)
	router.POST("/filters/data", controller.GetFilterData)

	// Write endpoints
	router.POST("/create", controller.CreateCompany)
	router.PUT("/:uuid", controller.UpdateCompany)
	router.DELETE("/:uuid", controller.DeleteCompany)
	router.POST("/upsert", controller.UpsertCompany)
	router.POST("/bulk", controller.BulkUpsertCompanies)
}
