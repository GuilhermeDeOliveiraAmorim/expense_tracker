package util

import (
	"log/slog"
	"os"
)

type Logger struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	From    string `json:"from"`
	Layer   string `json:"layer"`
	TypeLog string `json:"type_log"`
}

func NewLoggerError(code int, message, from, layer, typeLog string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Error(
		"ERROR",
		"code:", code,
		"message:", message,
		"from:", from,
		"layer:", layer,
		"type:", typeLog,
	)
}

func NewLoggerInfo(code int, message, from, layer, typeLog string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info(
		"INFO",
		"code:", code,
		"message:", message,
		"from:", from,
		"layer:", layer,
		"type:", typeLog,
	)
}

func NewLoggerWarning(code int, message, from, layer, typeLog string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Warn(
		"WARNING",
		"code:", code,
		"message:", message,
		"from:", from,
		"layer:", layer,
		"type:", typeLog,
	)
}
