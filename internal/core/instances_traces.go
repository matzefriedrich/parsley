package core

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	SpanAttrResolverInstanceFound             = "resolver.instance.found"
	SpanAttrResolverInstanceId                = "resolver.instance.id"
	SpanAttrResolverInstanceSource            = "resolver.instance.source"
	SpanAttrResolverInstanceStorage           = "resolver.instance.storage"
	SpanAttrResolverRegistrationId            = "resolver.registration.id"
	SpanAttrResolverRegistrationServiceType   = "resolver.registration.service-type"
	SpanAttrResolverRegistrationLifetimeScope = "resolver.registration.lifetime"
)

const (
	SourceLocal                    = "local"
	SourceNotFound                 = "not-found"
	SourceScope                    = "scope"
	StorageLocationLocalSingleton  = "local-singleton"
	StorageLocationParentSingleton = "parent-singleton"
	StorageLocationScope           = "scope"
	StorageLocationTransient       = "transient"
)

type tryResolveInstanceSpan struct {
	TypedSpan
}

type TryResolveInstanceSpan interface {
	trace.Span
	InstanceFound(id uint64, source string)
	InstanceNotFound()
}

var _ TryResolveInstanceSpan = (*tryResolveInstanceSpan)(nil)

func (t tryResolveInstanceSpan) InstanceFound(id uint64, source string) {
	t.SetAttributes(
		attribute.Int64(SpanAttrResolverInstanceId, int64(id)),
		attribute.String(SpanAttrResolverInstanceSource, source),
		attribute.Bool(SpanAttrResolverInstanceFound, true),
	)
}

func (t tryResolveInstanceSpan) InstanceNotFound() {
	t.SetAttributes(
		attribute.String(SpanAttrResolverInstanceSource, SourceNotFound),
		attribute.Bool(SpanAttrResolverInstanceFound, false),
	)
}

func newTryResolveInstanceSpan(scope context.Context) (context.Context, TryResolveInstanceSpan) {

	tp := otel.GetTracerProvider()
	tracer := tp.Tracer(TracerName)
	ctx, span := tracer.Start(scope, "instances.resolve.TryResolveInstance")

	return ctx, &tryResolveInstanceSpan{
		TypedSpan: NewTypedSpan("instances", "resolve", "TryResolveInstance", span),
	}
}

type keepInstanceSpan struct {
	TypedSpan
	registration types.ServiceRegistration
}

type KeepInstanceSpan interface {
	trace.Span
	InstanceStorage(location string)
}

var _ KeepInstanceSpan = (*keepInstanceSpan)(nil)

func (k keepInstanceSpan) InstanceStorage(location string) {
	k.SetAttributes(attribute.String(SpanAttrResolverInstanceStorage, location))
}

func newKeepInstanceSpan(scope context.Context, registration types.ServiceRegistration) (context.Context, KeepInstanceSpan) {

	tp := otel.GetTracerProvider()
	tracer := tp.Tracer(TracerName)
	ctx, span := tracer.Start(scope, "instances.resolve.TryKeepInstance")

	lifetimeScope := registration.LifetimeScope()
	registrationId := int64(registration.Id())
	serviceTypeName := registration.ServiceType().Name()

	span.SetAttributes(
		attribute.Int64(SpanAttrResolverRegistrationId, registrationId),
		attribute.String(SpanAttrResolverRegistrationServiceType, serviceTypeName),
		attribute.String(SpanAttrResolverRegistrationLifetimeScope, lifetimeScope.String()),
	)

	return ctx, &keepInstanceSpan{
		TypedSpan:    NewTypedSpan("instances", "resolve", "TryKeepInstance", span),
		registration: registration,
	}
}
