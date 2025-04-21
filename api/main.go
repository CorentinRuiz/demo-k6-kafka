package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/IBM/sarama"
)

type Event struct {
    UserID int    `json:"user_id"`
    Action string `json:"action"`
}

var producer sarama.SyncProducer

func initProducer() {
    brokers := []string{"kafka:9092"}

    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Partitioner = sarama.NewRandomPartitioner

    var err error
    producer, err = sarama.NewSyncProducer(brokers, config)
    if err != nil {
        log.Fatalf("Failed to create Kafka producer: %v", err)
    }
}

func main() {
    initProducer()
    defer producer.Close()

    router := gin.Default()

    router.POST("/event", func(c *gin.Context) {
        var event Event
        if err := c.ShouldBindJSON(&event); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
            return
        }

        msgBytes, _ := json.Marshal(event)
        msg := &sarama.ProducerMessage{
            Topic: "events",
            Value: sarama.ByteEncoder(msgBytes),
        }

        _, _, err := producer.SendMessage(msg)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send to Kafka"})
            return
        }

        c.JSON(http.StatusAccepted, gin.H{"status": "Event sent"})
    })

    router.GET("/health", func(c *gin.Context) {
        c.String(http.StatusOK, "OK")
    })

    log.Println("Gin API running on :8080")
    router.Run(":8080")
}
