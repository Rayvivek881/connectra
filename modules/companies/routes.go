package companies

import (
	"vivek-ray/modules/companies/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/:uuid", controller.GetCompanyByUUID)
	router.POST("/", controller.GetCompaniesByFilter)
	router.POST("/count", controller.GetCompaniesCountByFilter)
	router.POST("/create", controller.CreateCompany)
	router.PUT("/:uuid", controller.UpdateCompany)
	router.DELETE("/:uuid", controller.DeleteCompany)
	router.POST("/upsert", controller.UpsertCompany)
	router.POST("/batch-upsert", controller.BatchUpsert)
}
