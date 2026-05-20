USE TMF915;
GO

CREATE TABLE dbo.Alarm (
    id                           NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href                         NVARCHAR(500)    NULL,
    ack_state                    NVARCHAR(100)    NULL,
    ack_system_id                NVARCHAR(200)    NULL,
    ack_user_id                  NVARCHAR(200)    NULL,
    alarm_changed_time           DATETIME2        NULL,
    alarm_cleared_time           DATETIME2        NULL,
    alarm_details                NVARCHAR(2000)   NULL,
    alarm_escalation             BIT              NULL,
    alarm_raised_time            DATETIME2        NULL,
    alarm_reporting_time         DATETIME2        NULL,
    alarmed_object_type          NVARCHAR(200)    NULL,
    clear_system_id              NVARCHAR(200)    NULL,
    clear_user_id                NVARCHAR(200)    NULL,
    external_alarm_id            NVARCHAR(200)    NULL,
    is_root_cause                BIT              NULL,
    planned_outage_indicator     NVARCHAR(100)    NULL,
    probable_cause               NVARCHAR(500)    NULL,
    proposed_repaired_actions    NVARCHAR(1000)   NULL,
    reporting_system_id          NVARCHAR(200)    NULL,
    service_affecting            BIT              NULL,
    source_system_id             NVARCHAR(200)    NULL,
    specific_problem             NVARCHAR(500)    NULL,
    state                        NVARCHAR(100)    NULL,
    base_type                    NVARCHAR(200)    NULL,
    schema_location              NVARCHAR(500)    NULL,
    at_type                      NVARCHAR(200)    NULL,
    alarm_type                   NVARCHAR(MAX)    NULL,   -- JSON: AlarmType enum
    alarmed_object               NVARCHAR(MAX)    NULL,   -- JSON: AlarmedObject
    perceived_severity           NVARCHAR(MAX)    NULL,   -- JSON: PerceivedSeverity enum
    crossed_threshold_information NVARCHAR(MAX)   NULL,   -- JSON: CrossedThresholdInformation
    affected_service             NVARCHAR(MAX)    NULL,   -- JSON: []AffectedService
    comments                     NVARCHAR(MAX)    NULL,   -- JSON: []Comment
    correlated_alarm             NVARCHAR(MAX)    NULL,   -- JSON: []AlarmRef
    parent_alarm                 NVARCHAR(MAX)    NULL,   -- JSON: []AlarmRef
    places                       NVARCHAR(MAX)    NULL,   -- JSON: []Place
    created_at                   DATETIME2        NOT NULL DEFAULT GETUTCDATE(),
    updated_at                   DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_Alarm_raised ON dbo.Alarm(alarm_raised_time DESC);
CREATE INDEX IX_Alarm_state  ON dbo.Alarm(state);
GO
