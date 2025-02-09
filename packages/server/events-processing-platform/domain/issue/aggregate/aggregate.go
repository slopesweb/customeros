package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/pkg/errors"
	"strings"
)

const (
	IssueAggregateType eventstore.AggregateType = "issue"
)

type IssueAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Issue *model.Issue
}

func NewIssueAggregateWithTenantAndID(tenant, id string) *IssueAggregate {
	issueAggregate := IssueAggregate{}
	issueAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(IssueAggregateType, tenant, id)
	issueAggregate.SetWhen(issueAggregate.When)
	issueAggregate.Issue = &model.Issue{}
	issueAggregate.Tenant = tenant

	return &issueAggregate
}

func (a *IssueAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.IssueCreateV1:
		return a.onIssueCreate(evt)
	case event.IssueUpdateV1:
		return a.onIssueUpdate(evt)
	case event.IssueAddUserAssigneeV1:
		return a.onIssueAddUserAssignee(evt)
	case event.IssueRemoveUserAssigneeV1:
		return a.onIssueRemoveUserAssignee(evt)
	case event.IssueAddUserFollowerV1:
		return a.onIssueAddUserFollower(evt)
	case event.IssueRemoveUserFollowerV1:
		return a.onIssueRemoveUserFollower(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), events2.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *IssueAggregate) onIssueCreate(evt eventstore.Event) error {
	var eventData event.IssueCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.ID = a.ID
	a.Issue.Tenant = a.Tenant
	a.Issue.Subject = eventData.Subject
	a.Issue.Description = eventData.Description
	a.Issue.Status = eventData.Status
	a.Issue.Priority = eventData.Priority
	a.Issue.ReportedByOrganizationId = eventData.ReportedByOrganizationId
	a.Issue.SubmittedByOrganizationId = eventData.SubmittedByOrganizationId
	a.Issue.SubmittedByUserId = eventData.SubmittedByUserId
	a.Issue.Source = cmnmod.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.Source,
		AppSource:     eventData.AppSource,
	}
	a.Issue.CreatedAt = eventData.CreatedAt
	a.Issue.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Issue.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *IssueAggregate) onIssueUpdate(evt eventstore.Event) error {
	var eventData event.IssueUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == events2.SourceOpenline {
		a.Issue.Source.SourceOfTruth = eventData.Source
	}
	if eventData.Source != a.Issue.Source.SourceOfTruth && a.Issue.Source.SourceOfTruth == events2.SourceOpenline {
		if a.Issue.Subject == "" {
			a.Issue.Subject = eventData.Subject
		}
		if a.Issue.Description == "" {
			a.Issue.Description = eventData.Description
		}
		if a.Issue.Status == "" {
			a.Issue.Status = eventData.Status
		}
		if a.Issue.Priority == "" {
			a.Issue.Priority = eventData.Priority
		}
	} else {
		a.Issue.Subject = eventData.Subject
		a.Issue.Description = eventData.Description
		a.Issue.Status = eventData.Status
		a.Issue.Priority = eventData.Priority
	}
	a.Issue.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *IssueAggregate) onIssueAddUserAssignee(evt eventstore.Event) error {
	var eventData event.IssueAddUserAssigneeEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.AddAssignedToUserId(eventData.UserId)
	return nil
}

func (a *IssueAggregate) onIssueRemoveUserAssignee(evt eventstore.Event) error {
	var eventData event.IssueRemoveUserAssigneeEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.RemoveAssignedToUserId(eventData.UserId)
	return nil
}

func (a *IssueAggregate) onIssueAddUserFollower(evt eventstore.Event) error {
	var eventData event.IssueAddUserFollowerEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.AddFollowedByUserId(eventData.UserId)
	return nil
}

func (a *IssueAggregate) onIssueRemoveUserFollower(evt eventstore.Event) error {
	var eventData event.IssueRemoveUserFollowerEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.RemoveFollowedByUserId(eventData.UserId)
	return nil
}
