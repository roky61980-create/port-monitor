package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/roky61980-create/port-monitor/internal/app"
)
func main(){slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout,nil)));ctx,stop:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM);defer stop();if err:=app.Run(ctx);err!=nil{slog.Error("application stopped","error",err);os.Exit(1)}}
