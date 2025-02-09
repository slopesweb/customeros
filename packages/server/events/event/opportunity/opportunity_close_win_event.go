package opportunity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OpportunityCloseWinEvent struct {
	Tenant   string    `json:"tenant" validate:"required"`
	ClosedAt time.Time `json:"closedAt" validate:"required"`
}

func NewOpportunityCloseWinEvent(aggregate eventstore.Aggregate, closedAt time.Time) (eventstore.Event, error) {
	eventData := OpportunityCloseWinEvent{
		Tenant:   aggregate.GetTenant(),
		ClosedAt: closedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityCloseWinEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityCloseWinV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityCloseWinEvent")
	}
	return event, nil
}
