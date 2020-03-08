package scheduler

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	rpc "github.com/renskiy/ugly-scheduler/rpc/scheduler"
)

type newEventForm struct {
	*rpc.Event
}

func (form *newEventForm) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Delay, validation.Min(0)),
		validation.Field(&form.Message, validation.Required, validation.Length(0, 1024)),
	)
}
