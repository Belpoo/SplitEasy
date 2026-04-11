/**
 * Archivo: db/setup.go
 * Propósito: Configurar base de datos PostgreSQL de gastos.
 * Nota: Los Microservicios son "Shared Nothing", lo cual indica que Expenses DB
 * es independiente al Users DB para evitar acoples de estado, solo referencian UUIDs.
 */
package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://exp_user:exp_pwd@localhost:5434/exp_db?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error de inicio en driver DB de Gastos: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Aviso: Fallo ping inicial PG Exp. DB offline: %v", err)
	} else {
		createTables()
	}
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS expenses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		group_id UUID NOT NULL,
		paid_by UUID NOT NULL,
		amount DECIMAL(12,2) NOT NULL,
		description VARCHAR(255),
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS expense_splits (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
		user_id UUID NOT NULL,
		amount_owed DECIMAL(12,2) NOT NULL,
		paid BOOLEAN DEFAULT FALSE,
		settled_at TIMESTAMP
	);`

	if _, err := DB.Exec(query); err != nil {
		log.Printf("Fallo montando tablas en BD Expense: %v", err)
	} else {
		log.Println("Tablas DB Expense iniciadas (expenses, expense_splits)")
	}
}
