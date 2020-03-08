package scheduler

import (
	"github.com/renskiy/ugly-scheduler/internal/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SchedulerTestSuite struct {
	test.Suite
}

func TestToken(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}

func (s *SchedulerTestSuite) TestSchedule() {

}
