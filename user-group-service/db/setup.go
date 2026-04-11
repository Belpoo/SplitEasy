/**
 * Archivo: db/setup.go
 * Propósito: Inicializar conexión con PostgreSQL.
 * Decisiones: Crear tablas profiles, groups y group_members garantizando
 * el enfoque en ids del tipo UUID compatibles con Auth Service, pero
 * manteniendo estos separados per se.
 */
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://usr:pwd@localhost:5432/usergroupdb?sslmode=disable"
	}

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error abriendo PostgreSQL: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Advertencia: No se pudo hacer ping a la base de datos (ignorar en testing aislado): %v", err)
	} else {
		log.Println("Conexión exitosa a la Base de Datos User Group")
		createTables()
	}
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS profiles (
		id UUID PRIMARY KEY,
		display_name VARCHAR(100),
		avatar_url TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS groups (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(100) NOT NULL,
		description TEXT,
		created_by UUID NOT NULL REFERENCES profiles(id),
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS group_members (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
		role VARCHAR(10) CHECK (role IN ('admin', 'member')) DEFAULT 'member',
		joined_at TIMESTAMP DEFAULT NOW(),
		UNIQUE(group_id, user_id)
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("Error ejecutando esquema DB: %v", err)
	}
}
