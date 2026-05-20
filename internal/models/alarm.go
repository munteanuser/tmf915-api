package models

import "time"

type Alarm struct {
	ID                      string     `json:"id"                             db:"id"`
	Href                    *string    `json:"href,omitempty"                 db:"href"`
	AckState                *string    `json:"ackState,omitempty"             db:"ack_state"`
	AckSystemId             *string    `json:"ackSystemId,omitempty"          db:"ack_system_id"`
	AckUserId               *string    `json:"ackUserId,omitempty"            db:"ack_user_id"`
	AlarmChangedTime        *time.Time `json:"alarmChangedTime,omitempty"     db:"alarm_changed_time"`
	AlarmClearedTime        *time.Time `json:"alarmClearedTime,omitempty"     db:"alarm_cleared_time"`
	AlarmDetails            *string    `json:"alarmDetails,omitempty"         db:"alarm_details"`
	AlarmEscalation         *bool      `json:"alarmEscalation,omitempty"      db:"alarm_escalation"`
	AlarmRaisedTime         *time.Time `json:"alarmRaisedTime,omitempty"      db:"alarm_raised_time"`
	AlarmReportingTime      *time.Time `json:"alarmReportingTime,omitempty"   db:"alarm_reporting_time"`
	AlarmedObjectType       *string    `json:"alarmedObjectType,omitempty"    db:"alarmed_object_type"`
	ClearSystemId           *string    `json:"clearSystemId,omitempty"        db:"clear_system_id"`
	ClearUserId             *string    `json:"clearUserId,omitempty"          db:"clear_user_id"`
	ExternalAlarmId         *string    `json:"externalAlarmId,omitempty"      db:"external_alarm_id"`
	IsRootCause             *bool      `json:"isRootCause,omitempty"          db:"is_root_cause"`
	PlannedOutageIndicator  *string    `json:"plannedOutageIndicator,omitempty" db:"planned_outage_indicator"`
	ProbableCause           *string    `json:"probableCause,omitempty"        db:"probable_cause"`
	ProposedRepairedActions *string    `json:"proposedRepairedActions,omitempty" db:"proposed_repaired_actions"`
	ReportingSystemId       *string    `json:"reportingSystemId,omitempty"    db:"reporting_system_id"`
	ServiceAffecting        *bool      `json:"serviceAffecting,omitempty"     db:"service_affecting"`
	SourceSystemId          *string    `json:"sourceSystemId,omitempty"       db:"source_system_id"`
	SpecificProblem         *string    `json:"specificProblem,omitempty"      db:"specific_problem"`
	State                   *string    `json:"state,omitempty"                db:"state"`
	AtBaseType              *string    `json:"@baseType,omitempty"            db:"base_type"`
	AtSchemaLocation        *string    `json:"@schemaLocation,omitempty"      db:"schema_location"`
	AtType                  *string    `json:"@type,omitempty"                db:"at_type"`

	// JSON blobs
	AlarmType                    *string `json:"alarmType,omitempty"                    db:"alarm_type"`
	AlarmedObject                *string `json:"alarmedObject,omitempty"                db:"alarmed_object"`
	PerceivedSeverity            *string `json:"perceivedSeverity,omitempty"            db:"perceived_severity"`
	CrossedThresholdInformation  *string `json:"crossedThresholdInformation,omitempty"  db:"crossed_threshold_information"`
	AffectedService              *string `json:"affectedService,omitempty"              db:"affected_service"`
	Comment                      *string `json:"comment,omitempty"                      db:"comments"`
	CorrelatedAlarm              *string `json:"correlatedAlarm,omitempty"              db:"correlated_alarm"`
	ParentAlarm                  *string `json:"parentAlarm,omitempty"                  db:"parent_alarm"`
	Place                        *string `json:"place,omitempty"                        db:"places"`
}

type AlarmCreate struct {
	AckState                *string    `json:"ackState,omitempty"`
	AckSystemId             *string    `json:"ackSystemId,omitempty"`
	AckUserId               *string    `json:"ackUserId,omitempty"`
	AlarmChangedTime        *time.Time `json:"alarmChangedTime,omitempty"`
	AlarmClearedTime        *time.Time `json:"alarmClearedTime,omitempty"`
	AlarmDetails            *string    `json:"alarmDetails,omitempty"`
	AlarmEscalation         *bool      `json:"alarmEscalation,omitempty"`
	AlarmRaisedTime         *time.Time `json:"alarmRaisedTime,omitempty"`
	AlarmReportingTime      *time.Time `json:"alarmReportingTime,omitempty"`
	AlarmedObjectType       *string    `json:"alarmedObjectType,omitempty"`
	ExternalAlarmId         *string    `json:"externalAlarmId,omitempty"`
	IsRootCause             *bool      `json:"isRootCause,omitempty"`
	ProbableCause           *string    `json:"probableCause,omitempty"`
	SpecificProblem         *string    `json:"specificProblem,omitempty"`
	State                   *string    `json:"state,omitempty"`
	AlarmType               interface{} `json:"alarmType,omitempty"`
	AlarmedObject           interface{} `json:"alarmedObject,omitempty"`
	PerceivedSeverity       interface{} `json:"perceivedSeverity,omitempty"`
	AtBaseType              *string    `json:"@baseType,omitempty"`
	AtSchemaLocation        *string    `json:"@schemaLocation,omitempty"`
	AtType                  *string    `json:"@type,omitempty"`
}

type AlarmUpdate struct {
	AckState                *string    `json:"ackState,omitempty"`
	AckSystemId             *string    `json:"ackSystemId,omitempty"`
	AckUserId               *string    `json:"ackUserId,omitempty"`
	AlarmChangedTime        *time.Time `json:"alarmChangedTime,omitempty"`
	AlarmClearedTime        *time.Time `json:"alarmClearedTime,omitempty"`
	AlarmDetails            *string    `json:"alarmDetails,omitempty"`
	AlarmEscalation         *bool      `json:"alarmEscalation,omitempty"`
	PlannedOutageIndicator  *string    `json:"plannedOutageIndicator,omitempty"`
	ProbableCause           *string    `json:"probableCause,omitempty"`
	ProposedRepairedActions *string    `json:"proposedRepairedActions,omitempty"`
	ServiceAffecting        *bool      `json:"serviceAffecting,omitempty"`
	SpecificProblem         *string    `json:"specificProblem,omitempty"`
	State                   *string    `json:"state,omitempty"`
	AtBaseType              *string    `json:"@baseType,omitempty"`
	AtSchemaLocation        *string    `json:"@schemaLocation,omitempty"`
	AtType                  *string    `json:"@type,omitempty"`
}
