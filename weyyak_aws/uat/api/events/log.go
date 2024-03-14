package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	log := zerolog.New(os.Stdout).With().Logger()

	logMode := os.Getenv("LOG_MODE") // FILE | ELASTIC | PROMETHEUS etc.

	if logMode != "FILE" {
		return log
		// Control returns log with Stdout, no need of else
	}

	// below code handles FILE option, no need of else as control returns before this
	tempFile, err := ioutil.TempFile(os.TempDir(), "app.log")
	if err != nil {
		return log
	}
	log = zerolog.New(tempFile).With().Logger()

	return log
}

func Log(c *gin.Context) {
	log := c.MustGet("LOG").(zerolog.Logger)
	var evt LogEvent
	err := c.ShouldBindJSON(&evt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	log.Info().Str("@timestamp", time.Now().Format(time.RFC3339)).Msg(evt.Message)
}
