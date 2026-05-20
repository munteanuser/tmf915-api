USE TMF915;
GO

CREATE TABLE dbo.Hub (
    id              NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href            NVARCHAR(500)    NULL,
    callback        NVARCHAR(2000)   NOT NULL,
    query           NVARCHAR(1000)   NULL,
    base_type       NVARCHAR(200)    NULL,
    schema_location NVARCHAR(500)    NULL,
    at_type         NVARCHAR(200)    NULL,
    created_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE TABLE dbo.Listener (
    id         NVARCHAR(36)     NOT NULL PRIMARY KEY,
    callback   NVARCHAR(2000)   NOT NULL,
    query      NVARCHAR(1000)   NULL,
    created_at DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO
-- No index on callback (NVARCHAR(2000) exceeds 1700-byte index key limit)
