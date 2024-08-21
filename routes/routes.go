package routes

import (
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	router.POST("/import", controllers.ImportExcel)
	router.GET("/records", controllers.ViewRecords)
	router.PUT("/record/:id", controllers.EditRecord)
	router.DELETE("/records/:id", controllers.DeleteRecord)
	router.DELETE("/records/delete-all", controllers.DeleteAllRecords)

}
