package entity

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

// Deprecated, use neo4j module instead
type MeetingEntity struct {
	Id                 string
	Name               *string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	StartedAt          *time.Time
	EndedAt            *time.Time
	ConferenceUrl      *string
	MeetingExternalUrl *string
	AppSource          string
	Agenda             *string
	AgendaContentType  *string
	Source             neo4jentity.DataSource
	SourceOfTruth      neo4jentity.DataSource
	Recording          *string
	DataloaderKey      string
	Status             *MeetingStatus
}

func (meeting MeetingEntity) ToString() string {
	return fmt.Sprintf("id: %s\n", meeting.Id)
}

type MeetingEntities []MeetingEntity

func (MeetingEntity) IsTimelineEvent() {
}

func (MeetingEntity) TimelineEventLabel() string {
	return model.NodeLabelMeeting
}

func (MeetingEntity) IsAnalysisDescribe() {
}

func (MeetingEntity) AnalysisDescribeLabel() string {
	return model.NodeLabelMeeting
}

func (meeting *MeetingEntity) SetDataloaderKey(key string) {
	meeting.DataloaderKey = key
}

func (meeting MeetingEntity) GetDataloaderKey() string {
	return meeting.DataloaderKey
}

func (meeting MeetingEntity) Labels(tenant string) []string {
	return []string{"Meeting", "TimelineEvent", "Meeting_" + tenant, "TimelineEvent_" + tenant}
}
