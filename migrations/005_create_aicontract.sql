USE TMF915;
GO

CREATE TABLE dbo.AiContract (
    id                       NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href                     NVARCHAR(500)    NULL,
    approval_date            DATETIME2        NULL,
    approved                 BIT              NULL,
    description              NVARCHAR(2000)   NULL,
    name                     NVARCHAR(500)    NOT NULL,
    state                    NVARCHAR(100)    NULL,
    version                  NVARCHAR(50)     NULL,
    base_type                NVARCHAR(200)    NULL,
    schema_location          NVARCHAR(500)    NULL,
    at_type                  NVARCHAR(200)    NULL,
    ai_contract_specification NVARCHAR(MAX)   NULL,   -- JSON: EntityRef
    ai_model                 NVARCHAR(MAX)    NULL,   -- JSON: EntityRef
    template_ref             NVARCHAR(MAX)    NULL,   -- JSON: TemplateRef
    valid_for                NVARCHAR(MAX)    NULL,   -- JSON: TimePeriod
    characteristics          NVARCHAR(MAX)    NULL,   -- JSON: []Characteristic
    related_party            NVARCHAR(MAX)    NULL,   -- JSON: []RelatedParty
    rules                    NVARCHAR(MAX)    NULL,   -- JSON: []RuleRef
    created_at               DATETIME2        NOT NULL DEFAULT GETUTCDATE(),
    updated_at               DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_AiContract_name  ON dbo.AiContract(name);
CREATE INDEX IX_AiContract_state ON dbo.AiContract(state);
GO
