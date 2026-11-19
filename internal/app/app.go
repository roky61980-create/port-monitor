package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/roky61980-create/port-monitor/internal/alerts"
	"github.com/roky61980-create/port-monitor/internal/api"
	"github.com/roky61980-create/port-monitor/internal/config"
	"github.com/roky61980-create/port-monitor/internal/events"
	"github.com/roky61980-create/port-monitor/internal/snmptrap"
	"github.com/roky61980-create/port-monitor/internal/store"
	"github.com/roky61980-create/port-monitor/internal/telemetry"
)

func Run(ctx context.Context) error {
	cfg,err:=config.Load("");if err!=nil{return err};db,err:=store.Open(cfg.Database.Path);if err!=nil{return err};defer db.Close()
	hub:=events.NewHub();engine:=alerts.New(db,hub);tel,err:=telemetry.Start(ctx,cfg.OTel.Endpoint);if err!=nil{slog.Warn("OpenTelemetry exporter unavailable", "error",err)}else{defer func(){c,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel();if e:=tel.Shutdown(c);e!=nil{slog.Error("telemetry shutdown failed","error",e)}}()}
	trap,err:=snmptrap.Start(cfg.SNMP.TrapListen,snmptrap.Normalizer{Devices:db},func(ev model.NetworkEvent)error{if tel!=nil{tel.Event(context.Background(),ev.EventType,ev.Device,ev.Vendor,ev.Model,ev.Interface,ev.Severity,"")};return engine.Process(ev)});if err!=nil{slog.Warn("SNMP trap listener not started", "error",err)}else{defer trap.Close()}
	server:=&http.Server{Addr:cfg.Server.Listen,Handler:api.New(db,hub,engine,tel).Handler(),ReadHeaderTimeout:5*time.Second}
	errCh:=make(chan error,1);go func(){slog.Info("HTTP server listening","address",cfg.Server.Listen);errCh<-server.ListenAndServe()}()
	select{case <-ctx.Done():shutdown,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel();if err:=server.Shutdown(shutdown);err!=nil{return err};return nil;case err:=<-errCh:if errors.Is(err,http.ErrServerClosed){return nil};return err}
}
