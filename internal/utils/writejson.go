package utils

import (
	"encoding/json"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"go.uber.org/zap"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonMarshal, err := json.Marshal(data)
	if err != nil {
		logger.Log.Info("Json encode error", zap.Error(err))
		return
	}

	w.Write(jsonMarshal)
}
