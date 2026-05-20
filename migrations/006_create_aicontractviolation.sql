USE TMF915;
GO

CREATE TABLE dbo.AiContractViolation (
    id              NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href            NVARCHAR(500)    NULL,
    date            DATETIME2        NULL,
    base_type       NVARCHAR(200)    NULL,
    schema_location NVARCHAR(500)    NULL,
    at_type         NVARCHAR(200)    NULL,
    ai_contract     NVARCHAR(MAX)    NULL,   -- JSON: EntityRef
    violation       NVARCHAR(MAX)    NULL,   -- JSON: Violation
    related_party   NVARCHAR(MAX)    NULL,   -- JSON: []RelatedParty
    created_at      DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_AiContractViolation_date ON dbo.AiContractViolation(date DESC);
GO
