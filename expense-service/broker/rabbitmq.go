/**
 * Archivo: broker/rabbitmq.go
 * Propósito: Emisor de eventos al sistema desacoplado (Asincronicidad del punto Ajuste 2).
 * Decisiones Tácticas:
 *  Se lee la variable USE_BROKER = "true"/"false" de forma dinámica.
 *  Si es false: Actúa como simulador para no colapsar testing/entregas básicas, 
 *  solo emitiendo Logs.
 *  Si es true: Conectará y tirará el evento JSON "expense.created"
 */
package broker

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type SplitPayload struct {
	UserID     string  `json:"user_id"`
	AmountOwed float64 `json:"amount_owed"`
}

type ExpenseCreatedEvent struct {
	ExpenseID string         `json:"expense_id"`
	GroupID   string         `json:"group_id"`
	PaidBy    string         `json:"paid_by"`
	Amount    float64        `json:"amount"`
	Splits    []SplitPayload `json:"splits"`
}

func PublishExpenseEvent(event ExpenseCreatedEvent) {
	useBroker := os.Getenv("USE_BROKER")
	
	payloadBytes, _ := json.Marshal(event)

	if useBroker != "true" {
		log.Printf("📨 [BROKER SIMULADO] Evento expense.created: %s\n", string(payloadBytes))
		return
	}

	// Lógica genuina a RabbitMQ
	brokerURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		log.Printf("Fallo conexión a RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Fallo obteniendo canal broker: %v", err)
		return
	}
	defer ch.Close()

	// Declarar Cola
	q, err := ch.QueueDeclare(
		"expense.created", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Printf("Fallo declarando queue: %v", err)
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payloadBytes,
		})
	if err != nil {
		log.Printf("Error empujando mensaje al exchange: %v", err)
	} else {
		log.Printf("📨 [BROKER ACTIVO] Evento Publicado a RabbitMQ correctamente.")
	}
}
