package scheduler

import (
	"context"
	"github.com/renskiy/ugly-scheduler/internal/models"
	"github.com/renskiy/ugly-scheduler/internal/test"
	rpc "github.com/renskiy/ugly-scheduler/rpc/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/twitchtv/twirp"
	"testing"
	"time"
)

type SchedulerTestSuite struct {
	test.Suite
}

func TestToken(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}

func (s *SchedulerTestSuite) TestSchedule() {
	scheduler := &scheduler{s.App}

	now := time.Now()
	var delay int64 = 10
	message := "message"

	_, err := scheduler.Schedule(context.Background(), &rpc.Event{Delay: delay, Message: message})
	s.NoError(err)

	event := new(models.Event)
	s.NoError(s.App.DB().Get(event, `SELECT * FROM events LIMIT 1`))

	s.Equal(now.Add(time.Duration(delay)*time.Second).Round(time.Second), event.When.Round(time.Second))
	s.Equal(message, event.Message)
}

func (s *SchedulerTestSuite) TestScheduleErrors() {
	cases := []struct {
		name string
		form *rpc.Event
		err  error
	}{
		{
			name: "empty_form",
			form: &rpc.Event{},
			err:  twirp.NewError(twirp.InvalidArgument, "message: cannot be blank."),
		},
		{
			name: "negative_delay",
			form: &rpc.Event{
				Delay:   -1,
				Message: "message",
			},
			err: twirp.NewError(twirp.InvalidArgument, "delay: must be no less than 0."),
		},
	}

	scheduler := &scheduler{s.App}

	for _, testCase := range cases {
		s.T().Run(testCase.name, func(t *testing.T) {
			_, err := scheduler.Schedule(context.Background(), testCase.form)
			assert.Equal(t, testCase.err, err)
		})
	}
}

func (s *SchedulerTestSuite) TestProcessEvents() {
	now := time.Now()
	s.App.DB().MustExec(
		s.App.DB().Rebind(`INSERT INTO events ("when", message) VALUES (?, ?), (?, ?), (?, ?)`),
		now.Add(-time.Minute), "message 1",
		now, "message 2",
		now.Add(time.Minute), "message 3",
	)

	scheduler := &scheduler{s.App}
	s.NoError(scheduler.processEvents())

	event := new(models.Event)
	s.NoError(s.App.DB().Get(event, `SELECT * FROM events WHERE done`))
	s.Equal("message 1", event.Message)

	s.NoError(scheduler.processEvents())
	s.NoError(scheduler.processEvents())

	s.NoError(s.App.DB().Get(event, `SELECT * FROM events WHERE NOT done`))
	s.Equal("message 3", event.Message)
}
