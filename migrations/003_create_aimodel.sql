USE TMF915;
GO

CREATE TABLE dbo.AiModel (
    id                      NVARCHAR(36)     NOT NULL PRIMARY KEY,
    href                    NVARCHAR(500)    NULL,
    category                NVARCHAR(200)    NULL,
    description             NVARCHAR(2000)   NULL,
    end_date                DATETIME2        NULL,
    has_started             BIT              NULL,
    is_bundle               BIT              NULL,
    is_service_enabled      BIT              NULL,
    is_stateful             BIT              NULL,
    name                    NVARCHAR(500)    NOT NULL,
    service_date            DATETIME2        NULL,
    service_type            NVARCHAR(200)    NULL,
    start_date              DATETIME2        NULL,
    start_mode              NVARCHAR(100)    NULL,
    state                   NVARCHAR(100)    NULL,
    base_type               NVARCHAR(200)    NULL,
    schema_location         NVARCHAR(500)    NULL,
    at_type                 NVARCHAR(200)    NULL,
    ai_model_specification  NVARCHAR(MAX)    NULL,   -- JSON: EntityRef
    service_specification   NVARCHAR(MAX)    NULL,   -- JSON: EntityRef
    gpu                     NVARCHAR(MAX)    NULL,   -- JSON: ResourceRef
    training_data           NVARCHAR(MAX)    NULL,   -- JSON: EntityRef
    notes                   NVARCHAR(MAX)    NULL,   -- JSON: []Note
    places                  NVARCHAR(MAX)    NULL,   -- JSON: []any
    related_entity          NVARCHAR(MAX)    NULL,   -- JSON: []any
    related_party           NVARCHAR(MAX)    NULL,   -- JSON: []RelatedParty
    service_characteristic  NVARCHAR(MAX)    NULL,   -- JSON: []Characteristic
    created_at              DATETIME2        NOT NULL DEFAULT GETUTCDATE(),
    updated_at              DATETIME2        NOT NULL DEFAULT GETUTCDATE()
);
GO

CREATE INDEX IX_AiModel_name  ON dbo.AiModel(name);
CREATE INDEX IX_AiModel_state ON dbo.AiModel(state);
GO
