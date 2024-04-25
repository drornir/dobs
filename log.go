package dobs

import (
	"context"
	"log/slog"
)

func NewSlogHandler(wrapped slog.Handler) *SlogHandler {
	return &SlogHandler{
		wrapped: wrapped,
	}
}

type SlogHandler struct {
	wrapped slog.Handler
}

func (slh *SlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return slh.Enabled(ctx, level)
}

func (slh *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	r = r.Clone()
	attrs := AttrsFromContext(ctx)
	r.AddAttrs(attrs...)

	return slh.Handle(ctx, r)
}

func (slh *SlogHandler) WithAttrs(attrs []Attr) slog.Handler {
	return &SlogHandler{
		wrapped: slh.wrapped.WithAttrs(attrs),
	}
}

func (slh *SlogHandler) WithGroup(name string) slog.Handler {
	return &SlogHandler{
		wrapped: slh.wrapped.WithGroup(name),
	}
}
