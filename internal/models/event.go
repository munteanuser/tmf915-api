package models

import "time"

type Event struct {
	ID               string     `json:"id"                        db:"id"`
	Href             *string    `json:"href,omitempty"            db:"href"`
	CorrelationID    *string    `json:"correlationId,omitempty"   db:"correlation_id"`
	Description      *string    `json:"description,omitempty"     db:"description"`
	Domain           *string    `json:"domain,omitempty"          db:"domain"`
	EventID          *string    `json:"eventId,omitempty"         db:"event_id"`
	EventTime        *time.Time `json:"eventTime,omitempty"       db:"event_time"`
	EventType        *string    `json:"eventType,omitempty"       db:"event_type"`
	Priority         *string    `json:"priority,omitempty"        db:"priority"`
	TimeOccurred     *time.Time `json:"timeOccurred,omitempty"    db:"time_occurred"`
	Title            *string    `json:"title,omitempty"           db:"title"`
	TopicID          string     `json:"-"                         db:"topic_id"`
	AtBaseType       *string    `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string    `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string    `json:"@type,omitempty"           db:"at_type"`

	// JSON blobs
	EventPayload        interface{}    `json:"event,omitempty"          db:"-"`
	EventPayloadRaw     *string        `json:"-"                        db:"event_payload"`
	RelatedParty        []RelatedParty `json:"relatedParty,omitempty"   db:"-"`
	RelatedPartyRaw     *string        `json:"-"                        db:"related_party"`
	ReportingSystem     *EntityRef     `json:"reportingSystem,omitempty" db:"-"`
	ReportingSystemRaw  *string        `json:"-"                        db:"reporting_system"`
	Source              *EntityRef     `json:"source,omitempty"         db:"-"`
	SourceRaw           *string        `json:"-"                        db:"source"`
}

type EventCreate struct {
	CorrelationID    *string        `json:"correlationId,omitempty"`
	Description      *string        `json:"description,omitempty"`
	Domain           *string        `json:"domain,omitempty"`
	EventID          *string        `json:"eventId,omitempty"`
	EventTime        *time.Time     `json:"eventTime,omitempty"`
	EventType        *string        `json:"eventType,omitempty"`
	Priority         *string        `json:"priority,omitempty"`
	TimeOccurred     *time.Time     `json:"timeOccurred,omitempty"`
	Title            *string        `json:"title,omitempty"`
	EventPayload     interface{}    `json:"event,omitempty"`
	RelatedParty     []RelatedParty `json:"relatedParty,omitempty"`
	ReportingSystem  *EntityRef     `json:"reportingSystem,omitempty"`
	Source           *EntityRef     `json:"source,omitempty"`
	AtBaseType       *string        `json:"@baseType,omitempty"`
	AtSchemaLocation *string        `json:"@schemaLocation,omitempty"`
	AtType           *string        `json:"@type,omitempty"`
}
