package dobs

import (
	"context"
	"log/slog"
)

type (
	Attrs []Attr
	Attr  = slog.Attr
	Value = slog.Value
)

func AttrsFromContext(ctx context.Context) Attrs {
	attrs, _ := ctx.Value(ContextKeyAttrs).(Attrs)
	return attrs.Copy()
}

func ContextAppendAttrs(parent context.Context, attrs ...Attr) context.Context {
	existingAttrs := AttrsFromContext(parent)
	all := append(existingAttrs, Attrs(attrs)...)

	return context.WithValue(parent, ContextKeyAttrs, all)
}

func (attrs Attrs) Copy() Attrs {
	return append(Attrs(nil), attrs...)
}

func (attrs Attrs) Find(key string) (Attr, bool) {
	var (
		result Attr
		found  bool
	)
	for _, a := range attrs {
		if a.Key != key {
			continue
		}
		result = a
		found = true
		break
	}
	return result, found
}
