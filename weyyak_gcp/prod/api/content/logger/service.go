package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func JaegerMiddleware(c *gin.Context) {
	if c.Request.URL.Path != "/health" {
		start := time.Now()
		ctx := c.Request.Context()

		_, span := otel.Tracer(c.Request.URL.Path).Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		span.SetAttributes(attribute.String("name", c.Request.URL.Path))
		span.SetAttributes(attribute.String("IP", c.ClientIP()))
		// span.SetAttributes(attribute.String("collection", c.Params("collection")))
		span.SetAttributes(attribute.String("latency", time.Since(start).String()))
		span.SetAttributes(attribute.String("user-agent", c.Request.UserAgent()))

		ctx = trace.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)
	}

	c.Next()
	// Update the span with the response status code
}

func JSON(c *gin.Context, httpCode int, ginH interface{}) {
	if c.Request.URL.Path != "/health" {
		span := trace.SpanFromContext(c.Request.Context())
		defer span.End()

		span.SetAttributes(attribute.Int("http.status_code", httpCode))

		bytes, err := json.Marshal(ginH)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		message := string(bytes)

		switch httpCode {
		case 200:
			span.AddEvent(fmt.Sprintf("%s", message))
		case 400, 404, 500:
			span.SetStatus(codes.Error, fmt.Sprintf("%s", message))
			span.RecordError(errors.New(fmt.Sprintf("%s", message)))
		}
	}

	c.JSON(httpCode, ginH)
	return
}
