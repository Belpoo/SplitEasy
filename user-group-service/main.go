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

func main() {
	db.InitDB()

	r := gin.Default()

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
		port = "8080"
	}

	log.Printf("Iniciando Módulo de Usuarios y Grupos en puerto %s...", port)
	r.Run(":" + port)
}
