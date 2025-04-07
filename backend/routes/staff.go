package routes

import (
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/controller"
	"github.com/gin-gonic/gin"
)

func StaffRoutes(router *gin.Engine) {
	router.POST("/staff/create", controller.CreateStaff)

	router.POST("/staff/login", controller.LoginStaff)
}
