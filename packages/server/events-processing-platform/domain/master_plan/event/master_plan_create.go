package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type MasterPlanCreateEvent struct {
	Tenant       string        `json:"tenant" validate:"required"`
	Name         string        `json:"name"`
	CreatedAt    time.Time     `json:"createdAt"`
	SourceFields common.Source `json:"sourceFields"`
}

func NewMasterPlanCreateEvent(aggregate eventstore.Aggregate, name string, sourceFields common.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := MasterPlanCreateEvent{
		Tenant:       aggregate.GetTenant(),
		Name:         name,
		CreatedAt:    createdAt,
		SourceFields: sourceFields,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate MasterPlanCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, MasterPlanCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for MasterPlanCreateEvent")
	}

	return event, nil
}
