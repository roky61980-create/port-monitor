package telemetry

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct { provider *sdktrace.TracerProvider; tracer trace.Tracer }
func Start(ctx context.Context,endpoint string)(*Telemetry,error){ exp,err:=otlptracegrpc.New(ctx,otlptracegrpc.WithEndpoint(endpoint),otlptracegrpc.WithInsecure());if err!=nil{return nil,err};res,err:=resource.New(ctx,resource.WithAttributes(attribute.String("service.name","em-network-port-monitor")));if err!=nil{return nil,err};tp:=sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp),sdktrace.WithResource(res));return &Telemetry{provider:tp,tracer:tp.Tracer("port-monitor/events")},nil }
func (t *Telemetry) Event(ctx context.Context,kind,device,vendor,model,iface,severity,state string){if t==nil||t.tracer==nil{return};_,span:=t.tracer.Start(ctx,kind,trace.WithAttributes(attribute.String("event_type",kind),attribute.String("device",device),attribute.String("vendor",vendor),attribute.String("model",model),attribute.String("interface",iface),attribute.String("alert_severity",severity),attribute.String("alert_state",state)));span.End()}
func (t *Telemetry) Shutdown(ctx context.Context)error{if t==nil||t.provider==nil{return nil};return t.provider.Shutdown(ctx)}
var ErrDisabled=errors.New("telemetry disabled")
