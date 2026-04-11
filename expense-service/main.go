/**
 * Archivo: main.go
 * Propósito: Punto de entrada del microservicio Expense. Atiende tráfico de transaccionales y encola notificaciones.
 */
package main

import (
	"log"
	"os"

	"github.com/Belpoo/SplitEasy/expense-service/controllers"
	"github.com/Belpoo/SplitEasy/expense-service/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()

	expenseRoutes := r.Group("/expenses")
	{
		expenseRoutes.POST("", controllers.CreateExpense)
		expenseRoutes.GET("", controllers.ListGroupExpenses)
		expenseRoutes.DELETE("/:id", controllers.DeleteExpense)
		// Edit & Get single endpoints estarían aquí.
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // distinct from user-group default
	}

	log.Printf("Iniciando microservicio de Gastos locales en puerto %s...", port)
	r.Run(":" + port)
}
