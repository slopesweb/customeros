package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	event2 "github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ContactEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewContactEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *ContactEventHandler {
	return &ContactEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

func (h *ContactEventHandler) OnContactCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *h.services.CommonServices.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error

		data := neo4jrepository.ContactCreateFields{
			AggregateVersion: evt.Version,
			FirstName:        eventData.FirstName,
			LastName:         eventData.LastName,
			Prefix:           eventData.Prefix,
			Description:      eventData.Description,
			Timezone:         eventData.Timezone,
			ProfilePhotoUrl:  eventData.ProfilePhotoUrl,
			Username:         eventData.Username,
			Name:             eventData.Name,
			CreatedAt:        eventData.CreatedAt,
			SourceFields: neo4jmodel.Source{
				Source:        helper.GetSource(eventData.Source),
				SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceOfTruth),
				AppSource:     helper.GetAppSource(eventData.AppSource),
			},
		}
		err = h.services.CommonServices.Neo4jRepositories.ContactWriteRepository.CreateContactInTx(ctx, tx, eventData.Tenant, contactId, data)
		if err != nil {
			h.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
			return nil, err
		}
		if eventData.ExternalSystem.Available() {
			externalSystemData := neo4jmodel.ExternalSystem{
				ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
				ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
				ExternalId:       eventData.ExternalSystem.ExternalId,
				ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
				ExternalSource:   eventData.ExternalSystem.ExternalSource,
				SyncDate:         eventData.ExternalSystem.SyncDate,
			}
			err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, contactId, model.NodeLabelContact, externalSystemData)
			if err != nil {
				h.log.Errorf("Error while link contact %s with external system %s: %s", contactId, eventData.ExternalSystem.ExternalSystemId, err.Error())
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if eventData.SocialUrl != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
			return h.grpcClients.ContactClient.AddSocial(ctx, &contactpb.ContactAddSocialGrpcRequest{
				Tenant:    eventData.Tenant,
				ContactId: contactId,
				Url:       eventData.SocialUrl,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("AddSocial failed: %v", err.Error())
		}
	}

	return nil
}

func (h *ContactEventHandler) OnContactUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	data := neo4jrepository.ContactUpdateFields{
		AggregateVersion:      evt.Version,
		FirstName:             eventData.FirstName,
		LastName:              eventData.LastName,
		Prefix:                eventData.Prefix,
		Description:           eventData.Description,
		Timezone:              eventData.Timezone,
		ProfilePhotoUrl:       eventData.ProfilePhotoUrl,
		Username:              eventData.Username,
		Name:                  eventData.Name,
		Source:                eventData.Source,
		UpdateFirstName:       eventData.UpdateFirstName(),
		UpdateLastName:        eventData.UpdateLastName(),
		UpdateName:            eventData.UpdateName(),
		UpdatePrefix:          eventData.UpdatePrefix(),
		UpdateDescription:     eventData.UpdateDescription(),
		UpdateTimezone:        eventData.UpdateTimezone(),
		UpdateProfilePhotoUrl: eventData.UpdateProfilePhotoUrl(),
		UpdateUsername:        eventData.UpdateUsername(),
	}
	err := h.services.CommonServices.Neo4jRepositories.ContactWriteRepository.UpdateContact(ctx, eventData.Tenant, contactId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		session := utils.NewNeo4jWriteSession(ctx, *h.services.CommonServices.Neo4jRepositories.Neo4jDriver)
		defer session.Close(ctx)

		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			//var err error
			if eventData.ExternalSystem.Available() {
				externalSystemData := neo4jmodel.ExternalSystem{
					ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
					ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
					ExternalId:       eventData.ExternalSystem.ExternalId,
					ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
					ExternalSource:   eventData.ExternalSystem.ExternalSource,
					SyncDate:         eventData.ExternalSystem.SyncDate,
				}
				innerErr := h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, contactId, model.NodeLabelContact, externalSystemData)
				if innerErr != nil {
					h.log.Errorf("Error while link contact %s with external system %s: %s", contactId, eventData.ExternalSystem.ExternalSystemId, err.Error())
					return nil, innerErr
				}
			}
			return nil, nil
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (h *ContactEventHandler) OnPhoneNumberLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnPhoneNumberLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.PhoneNumberId, eventData.Label, eventData.Primary)

	return err
}

func (h *ContactEventHandler) OnEmailLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnEmailLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.EmailWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.EmailId, eventData.Label, eventData.Primary)

	return err
}

func (h *ContactEventHandler) OnLocationLinkToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnLocationLinkToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.LocationId)

	return err
}

func (h *ContactEventHandler) OnContactLinkToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactLinkToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactLinkWithOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.JobRoleCreateFields{
		Description: eventData.Description,
		JobTitle:    eventData.JobTitle,
		Primary:     eventData.Primary,
		CreatedAt:   eventData.CreatedAt,
		StartedAt:   eventData.StartedAt,
		EndedAt:     eventData.EndedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.SourceFields.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceFields.SourceOfTruth),
			AppSource:     helper.GetAppSource(eventData.SourceFields.AppSource),
		},
	}
	err := h.services.CommonServices.Neo4jRepositories.JobRoleWriteRepository.LinkContactWithOrganization(ctx, eventData.Tenant, contactId, eventData.OrganizationId, data)

	return err
}

func (h *ContactEventHandler) OnSocialAddedToContactV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnSocialAddedToContactV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	data := neo4jrepository.SocialFields{
		SocialId:       eventData.SocialId,
		Url:            eventData.Url,
		Alias:          eventData.Alias,
		ExternalId:     eventData.ExternalId,
		FollowersCount: eventData.FollowersCount,
		CreatedAt:      eventData.CreatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source.Source),
			SourceOfTruth: helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetSource(eventData.Source.AppSource),
		},
	}
	err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.MergeSocialForEntity(ctx, eventData.Tenant, contactId, model.NodeLabelContact, data)

	return err
}

func (h *ContactEventHandler) OnSocialRemovedFromContactV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnSocialRemovedFromContactV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactRemoveSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	if eventData.SocialId != "" {
		err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.RemoveSocialForEntityById(ctx, eventData.Tenant, contactId, model.NodeLabelContact, eventData.SocialId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}
	} else {
		err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.RemoveSocialForEntityByUrl(ctx, eventData.Tenant, contactId, model.NodeLabelContact, eventData.Url)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}
	}
	return nil
}

func (h *ContactEventHandler) OnAddTag(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnAddTag")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactAddTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	err := h.services.CommonServices.Neo4jRepositories.TagWriteRepository.LinkTagByIdToEntity(ctx, eventData.Tenant, eventData.TagId, contactId, model.NodeLabelContact, eventData.TaggedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while adding tag %s to contact %s: %s", eventData.TagId, contactId, err.Error())
	}

	return err
}

func (h *ContactEventHandler) OnRemoveTag(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnRemoveTag")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactRemoveTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	err := h.services.CommonServices.Neo4jRepositories.TagWriteRepository.UnlinkTagByIdFromEntity(ctx, eventData.Tenant, eventData.TagId, contactId, model.NodeLabelContact)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while removing tag %s to contact %s: %s", eventData.TagId, contactId, err.Error())
	}

	return err
}

func (h *ContactEventHandler) OnLocationAddedToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnLocationAddedToContact")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactAddLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	data := neo4jrepository.LocationCreateFields{
		RawAddress: eventData.RawAddress,
		Name:       eventData.Name,
		CreatedAt:  eventData.CreatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source),
			SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
			AppSource:     helper.GetSource(eventData.AppSource),
		},
		AddressDetails: neo4jrepository.AddressDetails{
			Latitude:      eventData.Latitude,
			Longitude:     eventData.Longitude,
			Country:       eventData.Country,
			CountryCodeA2: eventData.CountryCodeA2,
			CountryCodeA3: eventData.CountryCodeA3,
			Region:        eventData.Region,
			District:      eventData.District,
			Locality:      eventData.Locality,
			Street:        eventData.Street,
			Address:       eventData.AddressLine1,
			Address2:      eventData.AddressLine2,
			Zip:           eventData.ZipCode,
			AddressType:   eventData.AddressType,
			HouseNumber:   eventData.HouseNumber,
			PostalCode:    eventData.PostalCode,
			PlusFour:      eventData.PlusFour,
			Commercial:    eventData.Commercial,
			Predirection:  eventData.Predirection,
			TimeZone:      eventData.TimeZone,
			UtcOffset:     eventData.UtcOffset,
		},
	}

	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.CreateLocation(ctx, eventData.Tenant, eventData.LocationId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while creating location %s: %s", eventData.LocationId, err.Error())
		return err
	}
	err = h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while linking location %s to contact %s: %s", eventData.LocationId, contactId, err.Error())
		return err
	}

	return nil
}

func (h *ContactEventHandler) OnContactHide(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactHide")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactHideEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	err := h.services.CommonServices.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, eventData.Tenant, contactId, neo4jentity.ContactPropertyHide, true)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while hiding contact %s: %s", contactId, err.Error())
	}

	return err
}

func (h *ContactEventHandler) OnContactShow(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactEventHandler.OnContactShow")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event2.ContactShowEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contactId := contact.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	err := h.services.CommonServices.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, eventData.Tenant, contactId, neo4jentity.ContactPropertyHide, false)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while showing contact %s: %s", contactId, err.Error())
	}

	return err
}
