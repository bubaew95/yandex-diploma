package middleware

import (
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
	}

	loggerResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (w *loggerResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.responseData.status = statusCode
}

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
		}

		lw := loggerResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)
		
		duration := time.Since(start)

		logger.Log.Info("HTTP request&response",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.status),
		)
	})
}
