package main

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the HTTP method and request URL

		start := time.Now()

		traceID := uuid.New().String()
		
		c.Request.Header.Set("X-Trace-ID", traceID)

		c.Next()

		end := time.Now()

		duration := end.Sub(start)

		method := c.Request.Method

		httpstatus := c.Writer.Status()
		fmt.Printf("Status: %d\n", httpstatus)
		var logType string

		switch httpstatus {
		case 200:
			logType = "info"
		case 400:
			logType = "alert"
		case 404:
			logType = "fatal"
		case 500:
			logType = "warning"
		}
		fmt.Println("duration",duration , "method",method, "logType",logType)

		//  Loki(c, logType, method, traceID, duration)
	}
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func Loki(c *gin.Context, logType string, msg string, TraceId string, duration time.Duration) {
	// config, _ := AppConfig(".env")
	v := strconv.FormatInt(time.Now().UnixNano(), 10)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	msg = msg + " " + "traceId=" + TraceId + " " + c.Request.URL.Path

	fmt.Println(string(v))

	Lokijson := map[string]interface{}{
		"streams": []interface{}{
			map[string]interface{}{
				"stream": map[string]interface{}{
					"status":     c.Writer.Status(),
					"message":    msg,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"referer":    c.Request.Referer(),
					"duration":   duration.String(),
					"hostname":   hostname,
					"client_ip":  GetOutboundIP(),
					"user_agent": c.Request.UserAgent(),
					"level":      logType,
					"job":        "Event",
				},
				"values": []interface{}{
					[]interface{}{
						string(v),
						msg,
					},
				},
			},
		},
	}
	fmt.Println("lokijson",Lokijson)

	// jsonValue, _ := json.Marshal(Lokijson)

	// url := "http://3.110.118.98:3100/loki/api/v1/push"

	// // resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	// log.Println(resp, err)
}