package core

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

type TypedSpan struct {
	embedded.Span
	Inner trace.Span
}

const (
	SpanAttrModule    = "module"
	SpanAttrOperation = "operation"
	SpanAttrMethod    = "method"
)

func (t TypedSpan) End(options ...trace.SpanEndOption) {
	t.Inner.End(options...)
}

func (t TypedSpan) AddEvent(name string, options ...trace.EventOption) {
	t.Inner.AddEvent(name, options...)
}

func (t TypedSpan) AddLink(link trace.Link) {
	t.Inner.AddLink(link)
}

func (t TypedSpan) IsRecording() bool {
	return t.Inner.IsRecording()
}

func (t TypedSpan) RecordError(err error, options ...trace.EventOption) {
	t.Inner.RecordError(err, options...)
}

func (t TypedSpan) SpanContext() trace.SpanContext {
	return t.Inner.SpanContext()
}

func (t TypedSpan) SetStatus(code codes.Code, description string) {
	t.Inner.SetStatus(code, description)
}

func (t TypedSpan) SetName(name string) {
	t.Inner.SetName(name)
}

func (t TypedSpan) SetAttributes(kv ...attribute.KeyValue) {
	t.Inner.SetAttributes(kv...)
}

func (t TypedSpan) TracerProvider() trace.TracerProvider {
	return t.Inner.TracerProvider()
}

var _ trace.Span = (*TypedSpan)(nil)

func NewTypedSpan(module string, operation string, method string, span trace.Span) TypedSpan {
	span.SetAttributes(
		attribute.String(SpanAttrModule, module),
		attribute.String(SpanAttrOperation, operation),
		attribute.String(SpanAttrMethod, method),
	)
	return TypedSpan{
		Inner: span,
	}
}
