USE TMF915;
GO

CREATE TABLE dbo.Event (
    id               NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href             NVARCHAR(500)    NULL,
    correlation_id   NVARCHAR(200)    NULL,
    description      NVARCHAR(2000)   NULL,
    domain           NVARCHAR(200)    NULL,
    event_id         NVARCHAR(200)    NULL,
    event_time       DATETIME2        NULL,
    event_type       NVARCHAR(200)    NULL,
    priority         NVARCHAR(100)    NULL,
    time_occurred    DATETIME2        NULL,
    title            NVARCHAR(500)    NULL,
    topic_id         NVARCHAR(36)     NOT NULL,
    base_type        NVARCHAR(200)    NULL,
    schema_location  NVARCHAR(500)    NULL,
    at_type          NVARCHAR(200)    NULL,
    event_payload    NVARCHAR(MAX)    NULL,   -- JSON: any
    related_party    NVARCHAR(MAX)    NULL,   -- JSON: []RelatedParty
    reporting_system NVARCHAR(MAX)    NULL,   -- JSON: EntityRef
    source           NVARCHAR(MAX)    NULL,   -- JSON: EntityRef
    created_at       DATETIME2        NOT NULL DEFAULT GETUTCDATE(),

    CONSTRAINT FK_Event_Topic FOREIGN KEY (topic_id) REFERENCES dbo.Topic(id) ON DELETE CASCADE
);
GO

CREATE INDEX IX_Event_topic_time ON dbo.Event(topic_id, event_time DESC);
CREATE INDEX IX_Event_type       ON dbo.Event(event_type);
GO
