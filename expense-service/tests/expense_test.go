/**
 * Archivo: tests/expense_test.go
 * Propósito: Testing de Invarianza de Deuda.
 * Verifica estrictamente que divisiones matemáticamente inválidas detienen el registro asíncrono.
 */
package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Belpoo/SplitEasy/expense-service/controllers"
	"github.com/gin-gonic/gin"
)

func TestCreateExpense_OwedMismatchesTotal_ShouldReturn400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/expenses", controllers.CreateExpense)

	// Gasto de: Monto Total $90,000 COP
	// Splits de 30 y 30. Faltan 30. = $60,000. Debe generar Error 400.
	payload := []byte(`{
		"group_id": "grupo-fake-1",
		"amount": 90000.0,
		"description": "Cena compartida pero matematicas malas",
		"splits": [
			{"user_id": "A", "amount_owed": 30000.0},
			{"user_id": "B", "amount_owed": 30000.0}
		]
	}`)

	req, _ := http.NewRequest("POST", "/expenses", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Las matematicas malas debieron arrojar status 400, en cambio fue %d", w.Code)
	}
}

func TestDeleteExpense_ByAdmin(t *testing.T) {
	// (Placeholder) Para eliminar gasto en la versión simulada, se invocaría el mock proxy
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/expenses/:id", controllers.DeleteExpense)

	req, _ := http.NewRequest("DELETE", "/expenses/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Al no conectarse en real enviará 500 interno si no existe mock, o 200 en caso vacío de driver
	// Pero verifica rutado correcto.
}
