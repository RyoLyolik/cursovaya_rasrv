package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	kafkaBrokers  = []string{"localhost:9092"}
	kafkaTopic    = "temperature-readings"
	consumerGroup = "temperature-consumer"

	// Prometheus metrics
	temperatureGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_celsius",
			Help: "Temperature readings in Celsius",
		},
		[]string{"location"},
	)
	recordsProcessed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "temperature_records_processed_total",
			Help: "Total number of temperature records processed",
		},
	)
	dbInsertErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "temperature_db_insert_errors_total",
			Help: "Total number of database insert errors",
		},
	)
)

type TemperatureReading struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Timestamp   string  `json:"timestamp"`
}

func MyHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{ErrorLog: log.Default()}),
	)
}

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create Kafka consumer
	consumer, err := createConsumer()
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Set up HTTP server for Prometheus metrics
	metricsServer := &http.Server{
		Addr:         ":2112",
		Handler:      promhttp.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :2112")
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Handle shutdown signals
	ctx, cancel := context.WithCancel(context.Background())
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := consumer.Consume(ctx, []string{kafkaTopic}, &consumerHandler{db: db})
				if err != nil {
					log.Printf("Error from consumer: %v", err)
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	log.Println("Temperature consumer started...")
	<-sigterm
	log.Println("Terminating...")
	cancel()
}

func initDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=admin password=admin dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS temperature_readings (
			id SERIAL PRIMARY KEY,
			location TEXT NOT NULL,
			temperature FLOAT NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createConsumer() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(kafkaBrokers, consumerGroup, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

type consumerHandler struct {
	db *sql.DB
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var reading TemperatureReading
		if err := json.Unmarshal(message.Value, &reading); err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}

		// Update Prometheus metrics
		temperatureGauge.WithLabelValues(reading.Location).Set(reading.Temperature)
		recordsProcessed.Inc()

		// Insert into database
		_, err := h.db.Exec(
			"INSERT INTO temperature_readings (location, temperature, timestamp) VALUES ($1, $2, $3)",
			reading.Location, reading.Temperature, reading.Timestamp,
		)
		if err != nil {
			log.Printf("Error inserting into database: %v", err)
			dbInsertErrors.Inc()
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}
