package logger

import (
	"context"
	"log/slog"
)

type discardHandler struct {
	slog.JSONHandler
}

func (discardHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (discardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h discardHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h discardHandler) WithGroup(string) slog.Handler {
	return h
}

func SetDefaultDiscardLogger() {
	slog.SetDefault(slog.New(&discardHandler{}))
}
