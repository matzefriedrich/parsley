package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	ModuleName                  = "resolver"
	SpanAttrResolverServiceType = "resolver.resolve.service-type"
)

type resolverTraces struct {
	core.TypedSpan
}

type ResolverTraces interface {
	trace.Span
}

var _ ResolverTraces = (*resolverTraces)(nil)

func newResolveWithOptionsSpan(scope context.Context, serviceType types.ServiceType) (context.Context, trace.Span) {

	tp := otel.GetTracerProvider()
	tracer := tp.Tracer(core.TracerName)
	ctx, span := tracer.Start(scope, "resolver.resolve.ResolveWithOptions")

	span.SetAttributes(
		attribute.String(SpanAttrResolverServiceType, serviceType.Name()),
	)
	return ctx, &resolverTraces{
		TypedSpan: core.NewTypedSpan(ModuleName, "resolve", "ResolveWithOptions", span),
	}
}
