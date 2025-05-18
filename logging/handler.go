package logging

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"strings"
)

const _TIME_FORMAT string = "2006-01-02 15:04:05.999"

var _STRING_BY_LEVEL = map[slog.Level]string{
	slog.LevelDebug: "DBUG",
	slog.LevelInfo: "INFO",
	slog.LevelWarn: "WARN",
	slog.LevelError: "ERRO",
}

type _Handler struct {
	level slog.Level
	attrs []slog.Attr
	groups []string
}

func (h _Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func buildAttrString(builder *strings.Builder, prefix string, attr slog.Attr) {
	builder.WriteString(prefix)
	switch attr.Value.Kind() {
	case slog.KindGroup:
		for i, a := range attr.Value.Group() {
			if i > 0 {
				builder.WriteString(" ")
			}
			buildAttrString(builder, prefix + attr.Key + ".", a)
		}
	default:
		builder.WriteString(attr.String())
	}
}

func (h _Handler) Handle(ctx context.Context, record slog.Record) error {
	var pcs = []uintptr{record.PC}
	var frame, _ = runtime.CallersFrames(pcs).Next()
	var filename = frame.File
	if idx := strings.LastIndex(filename, "/"); idx > -1 {
		filename = frame.File[idx + 1:]
	}

	var builder strings.Builder
	for _, attr := range h.attrs {
		buildAttrString(&builder, "", attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		builder.WriteString(" ")
		var prefix string
		if len(h.groups) > 0 {
			prefix = strings.Join(h.groups, ".") + "."
		}
		buildAttrString(&builder, prefix, attr)
		return true
	})

	fmt.Printf(
		"[%s][%s] %s (%s:%d) - %s%s\n",
		record.Time.Format(_TIME_FORMAT),
		_STRING_BY_LEVEL[record.Level],
		frame.Function,
		filename,
		frame.Line,
		record.Message,
		builder.String(),
	)
	return nil
}

func (h _Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(h.groups) == 0 {
		return _Handler{
			level: h.level,
			attrs: append(h.attrs, attrs...),
		}
	}

	var groupArgs = make([]any, len(attrs))
	for _, attr := range attrs {
		groupArgs = append(groupArgs, attr)
	}
	var group = slog.Group(h.groups[len(h.groups) - 1], groupArgs...)
	for i := len(h.groups) - 2; i >= 0; i-- {
		group = slog.Group(h.groups[i], group)
	}

	return _Handler{
		level: h.level,
		attrs: append(slices.Clone(h.attrs), group),
		groups: slices.Clone(h.groups),
	}
}

func (h _Handler) WithGroup(name string) slog.Handler {
	return _Handler{
		level: h.level,
		attrs: slices.Clone(h.attrs),
		groups: append(slices.Clone(h.groups), name),
	}
}
