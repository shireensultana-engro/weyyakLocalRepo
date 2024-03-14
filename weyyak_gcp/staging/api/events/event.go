package main

import (
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Various microservices will call event to log the events.
// Based on the event type, various actions can be taken.
// Below wraps events into a struct with ID, type and message

type LogEvent struct {
	Timestamp time.Time `json:"@timestamp"`
	LogLevel  string    `json:"level"`
	Message   string    `json:"message"`
}

type UserEvent struct {
	Timestamp time.Time `json:"@timestamp"`
	UserID    string    `json:"userid"`
	EventType string    `json:"event"`
	Details   string    `json:"details"`
}

func Event(c *gin.Context) {
	log := c.MustGet("LOG").(zerolog.Logger)

	// Process the request
	evtType := c.Param("type")

	if evtType == "log" {
		var evt LogEvent
		err := c.ShouldBindJSON(&evt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		evt.Timestamp = time.Now()
		log.Info().Str("@timestamp", evt.Timestamp.Format(time.RFC3339)).Msg(evt.Message)

		// Store the event in Elastic
		esClient := c.MustGet("ELASTIC").(*elasticsearch.Client)

		// Store the log messages in Log index
		esClient.Index(
			"log",
			esutil.NewJSONReader(&evt),
			esClient.Index.WithRefresh("true"),
			esClient.Index.WithPretty(),
			esClient.Index.WithFilterPath("result", "_id"),
		)

		return
	}

	if evtType == "activity" {
		var evt UserEvent
		err := c.ShouldBindJSON(&evt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		evt.Timestamp = time.Now()
		log.Info().Str("Event", evt.EventType).
			Str("@timestamp", evt.Timestamp.Format(time.RFC3339)).
			Str("User", evt.UserID).Str("Details", evt.Details)

		// Store the event in Elastic
		esClient := c.MustGet("ELASTIC").(*elasticsearch.Client)

		// Store the log messages in Log index
		esClient.Index(
			"activity",
			esutil.NewJSONReader(&evt),
			esClient.Index.WithRefresh("true"),
			esClient.Index.WithPretty(),
			esClient.Index.WithFilterPath("result", "_id"),
		)

		return
	}
}
