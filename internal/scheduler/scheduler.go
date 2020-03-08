package scheduler

import (
	"context"
	"database/sql"
	"github.com/renskiy/ugly-scheduler/internal/app"
	"github.com/renskiy/ugly-scheduler/internal/models"
	rpc "github.com/renskiy/ugly-scheduler/rpc/scheduler"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"os"
	"time"
)

var host string

func init() {
	host, _ = os.Hostname()
}

// New returns new implementation of rpc.Scheduler
func New(app app.App) rpc.Scheduler {
	scheduler := &scheduler{App: app}
	scheduler.RunTask(scheduler.Worker)
	return scheduler
}

type scheduler struct {
	app.App
}

func (app *scheduler) Schedule(ctx context.Context, event *rpc.Event) (*rpc.Empty, error) {
	if err := (&newEventForm{event}).Validate(); err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	delay := time.Duration(event.Delay) * time.Second
	app.DB().MustExec(app.DB().Rebind(insertNewEventSQL), time.Now().Add(delay), event.Message)

	return &rpc.Empty{}, nil
}

func (app *scheduler) Worker(shutdown <-chan struct{}) error {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-shutdown:
			return nil
		case <-timer.C:
			if err := app.processEvents(); err != nil {
				app.Logger().Error("event process error", zap.Error(err))
			}
		}
	}
}

func (app *scheduler) processEvents() error {
	tx := app.DB().MustBegin()
	defer tx.Rollback()

	event := new(models.Event)
	if err := tx.Get(
		event,
		tx.Rebind(selectEventToPrecessSQL),
		time.Now(),
	); err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	app.Logger().Info(event.Message, zap.Int("id", event.ID), zap.String("host", host))

	tx.MustExec(tx.Rebind(markEventAsDoneSQL), event.ID)

	return tx.Commit()
}
