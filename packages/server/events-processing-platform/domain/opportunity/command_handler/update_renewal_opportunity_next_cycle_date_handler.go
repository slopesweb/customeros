package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type UpdateRenewalOpportunityNextCycleDateCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpdateRenewalOpportunityNextCycleDateCommand) error
}

type updateRenewalOpportunityNextCycleDateCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewUpdateRenewalOpportunityNextCycleDateCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) UpdateRenewalOpportunityNextCycleDateCommandHandler {
	return &updateRenewalOpportunityNextCycleDateCommandHandler{
		log: log,
		es:  es,
		cfg: cfg,
	}
}

func (h *updateRenewalOpportunityNextCycleDateCommandHandler) Handle(ctx context.Context, cmd *command.UpdateRenewalOpportunityNextCycleDateCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateRenewalOpportunityNextCycleDateCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	// Validate the command fields
	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		opportunityAggregate, err := aggregate.LoadOpportunityAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		// Apply the command to the aggregate
		if err = opportunityAggregate.HandleCommand(ctx, cmd); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		err = h.es.Save(ctx, opportunityAggregate)
		if err == nil {
			return nil // Save successful
		}

		if eventstore.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			// Handle concurrency error
			if attempt == h.cfg.RetriesOnOptimisticLockException-1 {
				// If we have reached the maximum number of retries, return an error
				tracing.TraceErr(span, err)
				return err
			}
			span.LogFields(log.Int("retryAttempt", attempt+1))
			time.Sleep(utils.BackOffExponentialDelay(attempt)) // backoffDelay is a function that increases the delay with each attempt
			continue                                           // Retry
		} else {
			// Some other error occurred
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
