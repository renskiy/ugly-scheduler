package main

import (
	"github.com/renskiy/ugly-scheduler/internal/app"
	"github.com/renskiy/ugly-scheduler/internal/scheduler"
	"github.com/renskiy/ugly-scheduler/rpc/scheduler"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("addr", ":8001", "server addr:port")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	app := app.New()

	service := rpc.NewSchedulerServer(scheduler.New(app), nil)
	app.Router().PathPrefix(service.PathPrefix()).Handler(service)

	app.Run()
}
