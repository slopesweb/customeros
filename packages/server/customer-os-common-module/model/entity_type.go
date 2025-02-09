package model

type EntityType string

const (
	CONTACT           EntityType = "CONTACT"
	USER              EntityType = "USER"
	ORGANIZATION      EntityType = "ORGANIZATION"
	EMAIL             EntityType = "EMAIL"
	MEETING           EntityType = "MEETING"
	CONTRACT          EntityType = "CONTRACT"
	INVOICE           EntityType = "INVOICE"
	INTERACTION_EVENT EntityType = "INTERACTION_EVENT"
	COMMENT           EntityType = "COMMENT"
	ISSUE             EntityType = "ISSUE"
	LOG_ENTRY         EntityType = "LOG_ENTRY"
	OPPORTUNITY       EntityType = "OPPORTUNITY"
	SERVICE_LINE_ITEM EntityType = "SERVICE_LINE_ITEM"
	REMINDER          EntityType = "REMINDER"
)

func (entityType EntityType) String() string {
	return string(entityType)
}

func (entityType EntityType) Neo4jLabel() string {
	switch entityType {
	case CONTACT:
		return NodeLabelContact
	case USER:
		return NodeLabelUser
	case ORGANIZATION:
		return NodeLabelOrganization
	case EMAIL:
		return NodeLabelEmail
	case MEETING:
		return NodeLabelMeeting
	case CONTRACT:
		return NodeLabelContract
	case INVOICE:
		return NodeLabelInvoice
	case INTERACTION_EVENT:
		return NodeLabelInteractionEvent
	case COMMENT:
		return NodeLabelComment
	case ISSUE:
		return NodeLabelIssue
	case LOG_ENTRY:
		return NodeLabelLogEntry
	case REMINDER:
		return NodeLabelReminder
	}
	return "Unknown"
}
