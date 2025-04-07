package routes

import (
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/controller"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/middleware"
	"github.com/gin-gonic/gin"
)

func PatientRoutes(router *gin.Engine) {
	router.GET("/patient/search/:id", controller.GetPatient)

	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/patient/search", controller.SearchPatients)
	}
}
