USE TMF915;
GO

-- RULE is a reserved keyword in SQL Server, so we use brackets
CREATE TABLE dbo.[Rule] (
    id              NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href            NVARCHAR(500)    NULL,
    name            NVARCHAR(500)    NOT NULL,
    base_type       NVARCHAR(200)    NULL,
    schema_location NVARCHAR(500)    NULL,
    at_type         NVARCHAR(200)    NULL,
    created_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE(),
    updated_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_Rule_name ON dbo.[Rule](name);
GO
