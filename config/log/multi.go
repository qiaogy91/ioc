package log

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type MultiHandler struct {
	hs []slog.Handler
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}

	for _, h := range m.hs {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.hs {
		// 如果有一个 handler 不启用当前级别，则返回 false
		if !h.Enabled(ctx, level) {
			return false
		}
	}
	// 只有当所有 handler 都启用时，才返回 true
	return true
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hs := make([]slog.Handler, len(m.hs))
	for i, h := range m.hs {
		hs[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{hs: hs}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	hs := make([]slog.Handler, len(m.hs))
	for i, h := range m.hs {
		hs[i] = h.WithGroup(name)
	}
	return &MultiHandler{hs: hs}
}
