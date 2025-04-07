package main

import (
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/config"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config.ConnectDB()
	routes.PatientRoutes(router)
	routes.StaffRoutes(router)

	router.Run() // listen and serve on 0.0.0.0:8080
}
