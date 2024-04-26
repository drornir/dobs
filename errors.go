package dobs

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

var ErrorRecordGroup = "error"

type Error struct {
	Err   error
	Attrs Attrs
	pc    uintptr
	t     time.Time
}

func NewError(err error, addPCDepth int) Error {
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function]
	runtime.Callers(2+addPCDepth, pcs[:])
	pc = pcs[0]
	return Error{
		Err:   err,
		Attrs: nil,
		pc:    pc,
		t:     time.Now(),
	}
}

func Errorf(f string, a ...any) Error {
	err := fmt.Errorf(f, a...)
	return NewError(err, 1)
}

func (e Error) WithAttrs(attrs ...Attr) Error {
	newAttrs := append(e.Attrs.Copy(), attrs...)

	e.Attrs = newAttrs
	return e
}

func (e Error) WithContextAttrs(ctx context.Context) Error {
	attrs := AttrsFromContext(ctx)
	newAttrs := append(e.Attrs.Copy(), attrs...)

	e.Attrs = newAttrs
	return e
}

func (e Error) Record() slog.Record {
	r := slog.NewRecord(e.t, slog.LevelError, e.Error(), e.pc)
	r.AddAttrs(
		slog.Attr{
			Key:   ErrorRecordGroup,
			Value: slog.GroupValue(e.WrappedAttrs()...),
		},
	)
	return r
}

func (e Error) LogTo(ctx context.Context, logger *slog.Logger) {
	h := logger.Handler()
	if !h.Enabled(ctx, slog.LevelError) {
		return
	}

	r := e.Record()
	_ = h.Handle(ctx, r)
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

func UnwrapAttrs(err error) Attrs {
	return Error{Err: err}.WrappedAttrs()
}

func (e Error) WrappedAttrs() Attrs {
	attrs := e.Attrs.Copy()

	err := e.Err
	for err != nil {
		switch terr := err.(type) {
		case Error:
			attrs = append(attrs, terr.Attrs...)
			err = terr.Err
		case interface{ Unwrap() error }:
			err = terr.Unwrap()
		case interface{ Unwrap() []error }:
			for _, err2 := range terr.Unwrap() {
				attrs = append(attrs, UnwrapAttrs(err2)...)
			}
			err = nil
		default:
			err = nil
		}
	}

	return attrs
}
