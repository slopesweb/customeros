package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	eventutils "github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

const (
	OpportunityAggregateType eventstore.AggregateType = "opportunity"
)

type OpportunityAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Opportunity *model.Opportunity
}

func NewOpportunityAggregateWithTenantAndID(tenant, id string) *OpportunityAggregate {
	oppAggregate := OpportunityAggregate{}
	oppAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(OpportunityAggregateType, tenant, id)
	oppAggregate.SetWhen(oppAggregate.When)
	oppAggregate.Opportunity = &model.Opportunity{}
	oppAggregate.Tenant = tenant

	return &oppAggregate
}

func (a *OpportunityAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *opportunitypb.CreateOpportunityGrpcRequest:
		return nil, a.createOpportunity(ctx, r)
	case *opportunitypb.UpdateOpportunityGrpcRequest:
		return nil, a.updateOpportunity(ctx, r)
	case *opportunitypb.CreateRenewalOpportunityGrpcRequest:
		return nil, a.createRenewalOpportunity(ctx, r)
	case *opportunitypb.UpdateRenewalOpportunityGrpcRequest:
		return nil, a.updateRenewalOpportunity(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OpportunityAggregate) createOpportunity(ctx context.Context, request *opportunitypb.CreateOpportunityGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.createOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), createdAtNotNil)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	dataFields := model.OpportunityDataFields{
		Name:              request.Name,
		MaxAmount:         request.MaxAmount,
		InternalType:      model.OpportunityInternalType(request.InternalType),
		ExternalType:      request.ExternalType,
		InternalStage:     model.OpportunityInternalStage(request.InternalStage),
		ExternalStage:     request.ExternalStage,
		EstimatedClosedAt: utils.TimestampProtoToTimePtr(request.EstimatedCloseDate),
		OwnerUserId:       request.OwnerUserId,
		CreatedByUserId:   utils.StringFirstNonEmpty(request.CreatedByUserId, request.LoggedInUserId),
		GeneralNotes:      request.GeneralNotes,
		NextSteps:         request.NextSteps,
		OrganizationId:    request.OrganizationId,
		Currency:          request.Currency,
		LikelihoodRate:    request.LikelihoodRate,
	}

	createEvent, err := events.NewOpportunityCreateEvent(a, dataFields, sourceFields, externalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&createEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *OpportunityAggregate) updateOpportunity(ctx context.Context, request *opportunitypb.UpdateOpportunityGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.updateOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	dataFields := model.OpportunityDataFields{
		Name:              request.Name,
		Amount:            request.Amount,
		MaxAmount:         request.MaxAmount,
		ExternalStage:     request.ExternalStage,
		ExternalType:      request.ExternalType,
		EstimatedClosedAt: utils.TimestampProtoToTimePtr(request.EstimatedCloseDate),
		OwnerUserId:       request.OwnerUserId,
		GeneralNotes:      request.GeneralNotes,
		NextSteps:         request.NextSteps,
		InternalStage:     model.OpportunityInternalStage(request.InternalStage),
		Currency:          request.Currency,
		LikelihoodRate:    request.LikelihoodRate,
	}

	fieldsMask := extractFieldsMask(request.FieldsMask)

	updateEvent, err := events.NewOpportunityUpdateEvent(a, dataFields, sourceFields.Source, externalSystem, updatedAtNotNil, fieldsMask)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityUpdateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *OpportunityAggregate) createRenewalOpportunity(ctx context.Context, request *opportunitypb.CreateRenewalOpportunityGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.createRenewalOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), createdAtNotNil)
	renewedAt := utils.TimestampProtoToTimePtr(request.RenewedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	renewalLikelihood := model.RenewalLikelihood(request.RenewalLikelihood).StringEnumValue()
	adjustedRate := request.RenewalAdjustedRate
	if string(renewalLikelihood) == "" {
		renewalLikelihood = neo4jenum.RenewalLikelihoodHigh
	}
	if renewalLikelihood == neo4jenum.RenewalLikelihoodHigh && adjustedRate == 0 {
		adjustedRate = 100
	}

	if adjustedRate < 0 {
		adjustedRate = 0
	} else if adjustedRate > 100 {
		adjustedRate = 100
	}

	createRenewalEvent, err := opportunityevent.NewOpportunityCreateRenewalEvent(a, request.ContractId, string(renewalLikelihood), request.RenewalApproved, sourceFields, createdAtNotNil, updatedAtNotNil, renewedAt, adjustedRate)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCreateRenewalEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&createRenewalEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return a.Apply(createRenewalEvent)
}

func (a *OpportunityAggregate) updateRenewalOpportunity(ctx context.Context, request *opportunitypb.UpdateRenewalOpportunityGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.updateRenewalOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())
	renewedAt := utils.TimestampProtoToTimePtr(request.RenewedAt)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	renewalLikelihood := model.RenewalLikelihood(request.RenewalLikelihood).StringEnumValue()
	adjustedRate := request.RenewalAdjustedRate
	if string(renewalLikelihood) == "" {
		adjustedRate = 100
		renewalLikelihood = neo4jenum.RenewalLikelihoodHigh
	}

	if adjustedRate < 0 {
		adjustedRate = 0
	} else if adjustedRate > 100 {
		adjustedRate = 100
	}

	fieldsMask := extractFieldsMask(request.FieldsMask)

	updateRenewalEvent, err := events.NewOpportunityUpdateRenewalEvent(a, string(renewalLikelihood), request.Comments, request.LoggedInUserId, sourceFields.Source, request.Amount, request.RenewalApproved, updatedAtNotNil, fieldsMask, request.OwnerUserId, renewedAt, adjustedRate)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityUpdateRenewalEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateRenewalEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return a.Apply(updateRenewalEvent)
}

func (a *OpportunityAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case opportunityevent.OpportunityCreateV1:
		return a.onOpportunityCreate(evt)
	case opportunityevent.OpportunityUpdateV1:
		return a.onOpportunityUpdate(evt)
	case opportunityevent.OpportunityCreateRenewalV1:
		return a.onRenewalOpportunityCreate(evt)
	case opportunityevent.OpportunityUpdateRenewalV1:
		return a.onRenewalOpportunityUpdate(evt)
	case opportunityevent.OpportunityUpdateNextCycleDateV1:
		return a.onOpportunityUpdateNextCycleDate(evt)
	case opportunityevent.OpportunityCloseWinV1:
		return a.onOpportunityCloseWin(evt)
	case opportunityevent.OpportunityCloseLooseV1:
		return a.onOpportunityCloseLoose(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), eventutils.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *OpportunityAggregate) onOpportunityCreate(evt eventstore.Event) error {
	var eventData events.OpportunityCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.ID = a.ID
	a.Opportunity.Tenant = a.Tenant
	a.Opportunity.OrganizationId = eventData.OrganizationId
	a.Opportunity.Name = eventData.Name
	a.Opportunity.MaxAmount = eventData.MaxAmount
	a.Opportunity.Currency = eventData.Currency
	a.Opportunity.LikelihoodRate = eventData.LikelihoodRate
	a.Opportunity.InternalType = eventData.InternalType
	a.Opportunity.ExternalType = eventData.ExternalType
	a.Opportunity.InternalStage = eventData.InternalStage
	a.Opportunity.ExternalStage = eventData.ExternalStage
	a.Opportunity.EstimatedClosedAt = eventData.EstimatedClosedAt
	a.Opportunity.OwnerUserId = eventData.OwnerUserId
	a.Opportunity.CreatedByUserId = eventData.CreatedByUserId
	a.Opportunity.GeneralNotes = eventData.GeneralNotes
	a.Opportunity.NextSteps = eventData.NextSteps
	a.Opportunity.CreatedAt = eventData.CreatedAt
	a.Opportunity.UpdatedAt = eventData.UpdatedAt
	a.Opportunity.Source = eventData.Source
	if eventData.ExternalSystem.Available() {
		a.Opportunity.ExternalSystems = []commonmodel.ExternalSystem{eventData.ExternalSystem}
	}

	return nil
}

func (a *OpportunityAggregate) onRenewalOpportunityCreate(evt eventstore.Event) error {
	var eventData opportunityevent.OpportunityCreateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.ID = a.ID
	a.Opportunity.Tenant = a.Tenant
	a.Opportunity.ContractId = eventData.ContractId
	a.Opportunity.InternalType = neo4jenum.OpportunityInternalTypeRenewal.String()
	a.Opportunity.InternalStage = eventData.InternalStage
	a.Opportunity.CreatedAt = eventData.CreatedAt
	a.Opportunity.UpdatedAt = eventData.UpdatedAt
	a.Opportunity.Source = eventData.Source
	a.Opportunity.RenewalDetails.RenewalLikelihood = eventData.RenewalLikelihood
	a.Opportunity.RenewalDetails.RenewalApproved = eventData.RenewalApproved
	a.Opportunity.RenewalDetails.RenewedAt = eventData.RenewedAt
	a.Opportunity.RenewalDetails.RenewalAdjustedRate = eventData.RenewalAdjustedRate

	return nil
}

func (a *OpportunityAggregate) onOpportunityUpdateNextCycleDate(evt eventstore.Event) error {
	var eventData opportunityevent.OpportunityUpdateNextCycleDateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.RenewalDetails.RenewedAt = eventData.RenewedAt

	return nil
}

func (a *OpportunityAggregate) onOpportunityUpdate(evt eventstore.Event) error {
	var eventData events.OpportunityUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// Update only if the source of truth is 'openline' or the new source matches the source of truth
	if eventData.Source == eventutils.SourceOpenline {
		a.Opportunity.Source.SourceOfTruth = eventData.Source
	}

	if eventData.Source != a.Opportunity.Source.SourceOfTruth && a.Opportunity.Source.SourceOfTruth == eventutils.SourceOpenline {
		// Update fields only if they are empty
		if a.Opportunity.Name == "" && eventData.UpdateName() {
			a.Opportunity.Name = eventData.Name
		}
	} else {
		if eventData.UpdateName() {
			a.Opportunity.Name = eventData.Name
		}
		if eventData.UpdateAmount() {
			a.Opportunity.Amount = eventData.Amount
		}
		if eventData.UpdateMaxAmount() {
			a.Opportunity.MaxAmount = eventData.MaxAmount
		}
		if eventData.UpdateExternalStage() {
			a.Opportunity.ExternalStage = eventData.ExternalStage
		}
		if eventData.UpdateExternalType() {
			a.Opportunity.ExternalType = eventData.ExternalType
		}
		if eventData.UpdateEstimatedClosedAt() {
			a.Opportunity.EstimatedClosedAt = eventData.EstimatedClosedAt
		}
		if eventData.UpdateOwnerUserId() {
			a.Opportunity.OwnerUserId = eventData.OwnerUserId
		}
		if eventData.UpdateInternalStage() {
			a.Opportunity.InternalStage = eventData.InternalStage
		}
		if eventData.UpdateCurrency() {
			a.Opportunity.Currency = eventData.Currency
		}
		if eventData.UpdateNextSteps() {
			a.Opportunity.NextSteps = eventData.NextSteps
		}
		if eventData.UpdateLikelihoodRate() {
			a.Opportunity.LikelihoodRate = eventData.LikelihoodRate
		}
	}
	a.Opportunity.UpdatedAt = eventData.UpdatedAt

	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.Opportunity.ExternalSystems {
			if externalSystem.ExternalSystemId == eventData.ExternalSystem.ExternalSystemId && externalSystem.ExternalId == eventData.ExternalSystem.ExternalId {
				found = true
				externalSystem.ExternalUrl = eventData.ExternalSystem.ExternalUrl
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
				if eventData.ExternalSystem.ExternalIdSecond != "" {
					externalSystem.ExternalIdSecond = eventData.ExternalSystem.ExternalIdSecond
				}
			}
		}
		if !found {
			a.Opportunity.ExternalSystems = append(a.Opportunity.ExternalSystems, eventData.ExternalSystem)
		}
	}

	return nil
}

func (a *OpportunityAggregate) onRenewalOpportunityUpdate(evt eventstore.Event) error {
	var eventData events.OpportunityUpdateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.UpdatedAt = eventData.UpdatedAt
	if eventData.UpdateRenewalLikelihood() {
		a.Opportunity.RenewalDetails.RenewalLikelihood = eventData.RenewalLikelihood
	}
	if eventData.RenewalApproved {
		a.Opportunity.RenewalDetails.RenewalApproved = eventData.RenewalApproved
	}
	if eventData.UpdatedByUserId != "" &&
		(eventData.Amount != a.Opportunity.Amount || eventData.RenewalLikelihood != a.Opportunity.RenewalDetails.RenewalLikelihood) {
		a.Opportunity.RenewalDetails.RenewalUpdatedByUserAt = &eventData.UpdatedAt
		a.Opportunity.RenewalDetails.RenewalUpdatedByUserId = eventData.UpdatedByUserId
	}
	if eventData.UpdateComments() {
		a.Opportunity.Comments = eventData.Comments
	}
	if eventData.UpdateAmount() {
		a.Opportunity.Amount = eventData.Amount
	}
	if eventData.Source == eventutils.SourceOpenline {
		a.Opportunity.Source.SourceOfTruth = eventData.Source
	}
	if eventData.OwnerUserId != "" {
		a.Opportunity.OwnerUserId = eventData.OwnerUserId
	}
	if eventData.UpdateRenewedAt() {
		a.Opportunity.RenewalDetails.RenewedAt = eventData.RenewedAt
	}
	if eventData.UpdateRenewalAdjustedRate() {
		a.Opportunity.RenewalDetails.RenewalAdjustedRate = eventData.RenewalAdjustedRate
	}

	return nil
}

func (a *OpportunityAggregate) onOpportunityCloseWin(evt eventstore.Event) error {
	var eventData opportunityevent.OpportunityCloseWinEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Opportunity.InternalStage = neo4jenum.OpportunityInternalStageClosedWon.String()
	a.Opportunity.ClosedAt = &eventData.ClosedAt
	return nil
}

func (a *OpportunityAggregate) onOpportunityCloseLoose(evt eventstore.Event) error {
	var eventData opportunityevent.OpportunityCloseLooseEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Opportunity.InternalStage = neo4jenum.OpportunityInternalStageClosedLost.String()
	a.Opportunity.ClosedAt = &eventData.ClosedAt
	return nil
}

func extractFieldsMask(requestMaskFields []opportunitypb.OpportunityMaskField) []string {
	maskFields := make([]string, 0)
	if requestMaskFields == nil || len(requestMaskFields) == 0 {
		return maskFields
	}
	for _, field := range requestMaskFields {
		switch field {
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_NAME:
			maskFields = append(maskFields, opportunityevent.FieldMaskName)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT:
			maskFields = append(maskFields, opportunityevent.FieldMaskAmount)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_COMMENTS:
			maskFields = append(maskFields, opportunityevent.FieldMaskComments)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD:
			maskFields = append(maskFields, opportunityevent.FieldMaskRenewalLikelihood)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT:
			maskFields = append(maskFields, opportunityevent.FieldMaskMaxAmount)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWED_AT:
			maskFields = append(maskFields, opportunityevent.FieldMaskRenewedAt)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ADJUSTED_RATE:
			maskFields = append(maskFields, opportunityevent.FieldMaskAdjustedRate)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_EXTERNAL_TYPE:
			maskFields = append(maskFields, opportunityevent.FieldMaskExternalType)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_EXTERNAL_STAGE:
			maskFields = append(maskFields, opportunityevent.FieldMaskExternalStage)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_INTERNAL_STAGE:
			maskFields = append(maskFields, opportunityevent.FieldMaskInternalStage)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ESTIMATED_CLOSE_DATE:
			maskFields = append(maskFields, opportunityevent.FieldMaskEstimatedClosedAt)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_OWNER_USER_ID:
			maskFields = append(maskFields, opportunityevent.FieldMaskOwnerUserId)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_CURRENCY:
			maskFields = append(maskFields, opportunityevent.FieldMaskCurrency)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_NEXT_STEPS:
			maskFields = append(maskFields, opportunityevent.FieldMaskNextSteps)
		case opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_LIKELIHOOD_RATE:
			maskFields = append(maskFields, opportunityevent.FieldMaskLikelihoodRate)
		}
	}
	return utils.RemoveDuplicates(maskFields)
}
