package common

import (
	"vivek-ray/modules/common/controller"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.POST("/batch-upsert", controller.BatchUpsert)

	// Upload
	router.GET("/upload-url", controller.GetUploadURL)
	router.GET("/download-url", controller.GetDownloadURL)

	// Jobs
	router.POST("/jobs", controller.ListJobs)
	router.POST("/jobs/create", controller.CreateJob)
	router.GET("/jobs/:job_uuid", controller.GetJobByUuid)

	// Filters
	router.GET("/:service/filters", controller.GetFilters)
	router.POST("/:service/filters/data", controller.GetFilterData)
}
