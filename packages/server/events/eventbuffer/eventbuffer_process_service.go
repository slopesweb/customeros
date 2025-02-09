package eventbuffer

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
)

type EventBufferProcessService struct {
	eventBufferRepository postgresRepository.EventBufferRepository
	logger                logger.Logger
	grpc_clients          *grpc_client.Clients
	signalChannel         chan os.Signal
	ticker                *time.Ticker
}

func NewEventBufferProcessService(ebr postgresRepository.EventBufferRepository, logger logger.Logger, grpc_clients *grpc_client.Clients) *EventBufferProcessService {
	return &EventBufferProcessService{eventBufferRepository: ebr, logger: logger, grpc_clients: grpc_clients}
}

func (eb *EventBufferProcessService) Start(ctx context.Context) {
	eb.logger.Info("EventBufferProcessService started")

	eb.ticker = time.NewTicker(time.Second * 30)
	eb.signalChannel = make(chan os.Signal, 1)
	signal.Notify(eb.signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				// Run dispatch logic every n seconds
				eb.logger.Info("EventBufferWatcher.Dispatch: dispatch buffered events")
				err := eb.Dispatch(ctx)
				if err != nil {
					eb.logger.Errorf("EventBufferWatcher.Dispatch: error dispatching events: %s", err.Error())
				}
			case <-eb.signalChannel:
				// Shutdown goroutine
				eb.logger.Info("EventBufferWatcher.Dispatch: Got signal, exiting...")
				runtime.Goexit()
			}
		}
	}(eb.ticker)
}

// Stop stops the EventBufferWatcher
func (eb *EventBufferProcessService) Stop() {
	eb.signalChannel <- syscall.SIGTERM // TODO get the signal from the caller
	eb.ticker.Stop()
	eb.logger.Info("EventBufferWatcher stopped")
	close(eb.signalChannel)
	eb.signalChannel = nil
}

// Dispatch dispatches all expired events from event_buffer table, and delete them after dispatching
func (eb *EventBufferProcessService) Dispatch(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferProcessService.Dispatch")
	defer span.Finish()

	now := utils.Now()
	eventBuffers, err := eb.eventBufferRepository.GetByExpired(now)
	if err != nil {
		return err
	}
	if len(eventBuffers) == 0 {
		return nil
	}
	tracing.LogObjectAsJson(span, "expiredEvents", eventBuffers)
	for _, eventBuffer := range eventBuffers {
		if err := eb.HandleEvent(ctx, eventBuffer); err != nil {
			tracing.TraceErr(span, err)
			eb.logger.Errorf("EventBufferWatcher.Dispatch: error handling event: %s", err.Error())
			continue
		}
		err = eb.eventBufferRepository.Delete(&eventBuffer)
		if err != nil {
			return err
		}
	}
	return err
}

func (eb *EventBufferProcessService) HandleEvent(ctx context.Context, eventBuffer postgresEntity.EventBuffer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferProcessService.HandleEvent")
	defer span.Finish()

	//skip these 2 events that are handled by subscribers until we migrate and test them
	if eventBuffer.EventType == "V1_ORGANIZATION_UPDATE_OWNER_NOTIFICATION" {
		return errors.New("Event type not supported")
	}

	_, err := eb.grpc_clients.EventStoreClient.StoreEvent(context.Background(), &eventstorepb.StoreEventGrpcRequest{
		EventDataBytes: eventBuffer.EventData,
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
