USE TMF915;
GO

CREATE TABLE dbo.Topic (
    id              NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href            NVARCHAR(500)    NULL,
    content_query   NVARCHAR(1000)   NULL,
    header_query    NVARCHAR(1000)   NULL,
    name            NVARCHAR(500)    NOT NULL,
    base_type       NVARCHAR(200)    NULL,
    schema_location NVARCHAR(500)    NULL,
    at_type         NVARCHAR(200)    NULL,
    created_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE(),
    updated_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_Topic_name ON dbo.Topic(name);
GO
