USE TMF915;
GO

CREATE TABLE dbo.Monitor (
    id              NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href            NVARCHAR(500)    NULL,
    source_href     NVARCHAR(500)    NULL,
    state           NVARCHAR(100)    NULL,
    base_type       NVARCHAR(200)    NULL,
    schema_location NVARCHAR(500)    NULL,
    at_type         NVARCHAR(200)    NULL,
    request         NVARCHAR(MAX)    NULL,   -- JSON: Request
    response        NVARCHAR(MAX)    NULL,   -- JSON: Response
    created_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE(),
    updated_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_Monitor_state ON dbo.Monitor(state);
GO
