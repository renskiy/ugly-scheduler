package scheduler

import (
	"context"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/renskiy/ugly-scheduler/internal/app"
	rpc "github.com/renskiy/ugly-scheduler/rpc/scheduler"
)

func New(app app.App) rpc.Scheduler {
	return &scheduler{App: app}
}

type scheduler struct {
	app.App
}

func (app *scheduler) Schedule(ctx context.Context, eent *rpc.Event) (*google_protobuf.Empty, error) {
	return nil, nil
}
