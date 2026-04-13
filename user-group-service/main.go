/**
 * Archivo: main.go
 * Propósito: Router universal GIN y puerto expuesto 8080 del User Group Service.
 */
package main

import (
	"log"
	"os"

	"github.com/Belpoo/SplitEasy/user-group-service/controllers"
	"github.com/Belpoo/SplitEasy/user-group-service/db"
	"github.com/Belpoo/SplitEasy/user-group-service/middleware"
	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5177") // Port for user-group-frontend
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	db.InitDB()

	r := gin.Default()
	r.Use(corsMiddleware())

	// Rutas protegidas genéricas
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// Usuarios
		protected.GET("/users/me", controllers.GetMe)
		protected.PUT("/users/me", controllers.UpdateMe)

		// Grupos
		protected.POST("/groups", controllers.CreateGroup)
		protected.GET("/groups", controllers.ListGroups)
		protected.GET("/groups/:id", controllers.GetGroup)
		protected.POST("/groups/:id/members", controllers.InviteMember)
		protected.DELETE("/groups/:id/members/:userId", controllers.RemoveMember)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Iniciando Módulo de Usuarios y Grupos en puerto %s...", port)
	r.Run(":" + port)
}

