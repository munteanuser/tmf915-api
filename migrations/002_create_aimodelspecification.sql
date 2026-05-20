USE TMF915;
GO

CREATE TABLE dbo.AiModelSpecification (
    id                  NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href                NVARCHAR(500)    NULL,
    description         NVARCHAR(2000)   NULL,
    is_bundle           BIT              NULL,
    last_update         DATETIME2        NULL,
    lifecycle_status    NVARCHAR(100)    NULL,
    name                NVARCHAR(500)    NOT NULL,
    version             NVARCHAR(50)     NULL,
    base_type           NVARCHAR(200)    NULL,
    schema_location     NVARCHAR(500)    NULL,
    at_type             NVARCHAR(200)    NULL,
    valid_for           NVARCHAR(MAX)    NULL,   -- JSON: TimePeriod
    target_entity_schema NVARCHAR(MAX)   NULL,   -- JSON: TargetEntitySchema
    attachment          NVARCHAR(MAX)    NULL,   -- JSON: []AttachmentRefOrValue
    constraint_refs     NVARCHAR(MAX)    NULL,   -- JSON: []ConstraintRef
    related_party       NVARCHAR(MAX)    NULL,   -- JSON: []RelatedParty
    created_at          DATETIME2        NOT NULL DEFAULT GETUTCDATE(),
    updated_at          DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_AiModelSpecification_name ON dbo.AiModelSpecification(name);
CREATE INDEX IX_AiModelSpecification_lifecycle ON dbo.AiModelSpecification(lifecycle_status);
GO
