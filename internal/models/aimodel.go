package models

import "time"

type AiModel struct {
	ID               string     `json:"id"                        db:"id"`
	Href             *string    `json:"href,omitempty"            db:"href"`
	Category         *string    `json:"category,omitempty"        db:"category"`
	Description      *string    `json:"description,omitempty"     db:"description"`
	EndDate          *time.Time `json:"endDate,omitempty"         db:"end_date"`
	HasStarted       *bool      `json:"hasStarted,omitempty"      db:"has_started"`
	IsBundle         *bool      `json:"isBundle,omitempty"        db:"is_bundle"`
	IsServiceEnabled *bool      `json:"isServiceEnabled,omitempty" db:"is_service_enabled"`
	IsStateful       *bool      `json:"isStateful,omitempty"      db:"is_stateful"`
	Name             string     `json:"name"                      db:"name"`
	ServiceDate      *time.Time `json:"serviceDate,omitempty"     db:"service_date"`
	ServiceType      *string    `json:"serviceType,omitempty"     db:"service_type"`
	StartDate        *time.Time `json:"startDate,omitempty"       db:"start_date"`
	StartMode        *string    `json:"startMode,omitempty"       db:"start_mode"`
	State            *string    `json:"state,omitempty"           db:"state"`
	AtBaseType       *string    `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string    `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string    `json:"@type,omitempty"           db:"at_type"`

	// Stored as JSON blobs
	AiModelSpecification    *EntityRef     `json:"aiModelSpecification,omitempty"  db:"-"`
	AiModelSpecificationRaw *string        `json:"-"                               db:"ai_model_specification"`
	ServiceSpecification    *EntityRef     `json:"serviceSpecification,omitempty"  db:"-"`
	ServiceSpecificationRaw *string        `json:"-"                               db:"service_specification"`
	GPU                     *ResourceRef   `json:"gpu,omitempty"                   db:"-"`
	GPURaw                  *string        `json:"-"                               db:"gpu"`
	TrainingData            *EntityRef     `json:"trainingData,omitempty"          db:"-"`
	TrainingDataRaw         *string        `json:"-"                               db:"training_data"`
	Note                    []Note         `json:"note,omitempty"                  db:"-"`
	NoteRaw                 *string        `json:"-"                               db:"notes"`
	Place                   []interface{}  `json:"place,omitempty"                 db:"-"`
	PlaceRaw                *string        `json:"-"                               db:"places"`
	RelatedEntity           []interface{}  `json:"relatedEntity,omitempty"         db:"-"`
	RelatedEntityRaw        *string        `json:"-"                               db:"related_entity"`
	RelatedParty            []RelatedParty `json:"relatedParty,omitempty"          db:"-"`
	RelatedPartyRaw         *string        `json:"-"                               db:"related_party"`
	ServiceCharacteristic   []Characteristic `json:"serviceCharacteristic,omitempty" db:"-"`
	ServiceCharacteristicRaw *string        `json:"-"                               db:"service_characteristic"`
}

type AiModelCreate struct {
	Category             *string          `json:"category,omitempty"`
	Description          *string          `json:"description,omitempty"`
	EndDate              *time.Time       `json:"endDate,omitempty"`
	HasStarted           *bool            `json:"hasStarted,omitempty"`
	IsBundle             *bool            `json:"isBundle,omitempty"`
	IsServiceEnabled     *bool            `json:"isServiceEnabled,omitempty"`
	IsStateful           *bool            `json:"isStateful,omitempty"`
	Name                 string           `json:"name"`
	ServiceDate          *time.Time       `json:"serviceDate,omitempty"`
	ServiceType          *string          `json:"serviceType,omitempty"`
	StartDate            *time.Time       `json:"startDate,omitempty"`
	StartMode            *string          `json:"startMode,omitempty"`
	State                *string          `json:"state,omitempty"`
	AiModelSpecification *EntityRef       `json:"aiModelSpecification,omitempty"`
	ServiceSpecification *EntityRef       `json:"serviceSpecification,omitempty"`
	GPU                  *ResourceRef     `json:"gpu,omitempty"`
	TrainingData         *EntityRef       `json:"trainingData,omitempty"`
	Note                 []Note           `json:"note,omitempty"`
	RelatedParty         []RelatedParty   `json:"relatedParty,omitempty"`
	ServiceCharacteristic []Characteristic `json:"serviceCharacteristic,omitempty"`
	AtBaseType           *string          `json:"@baseType,omitempty"`
	AtSchemaLocation     *string          `json:"@schemaLocation,omitempty"`
	AtType               *string          `json:"@type,omitempty"`
}

type AiModelUpdate struct {
	Category             *string          `json:"category,omitempty"`
	Description          *string          `json:"description,omitempty"`
	EndDate              *time.Time       `json:"endDate,omitempty"`
	HasStarted           *bool            `json:"hasStarted,omitempty"`
	IsBundle             *bool            `json:"isBundle,omitempty"`
	IsServiceEnabled     *bool            `json:"isServiceEnabled,omitempty"`
	IsStateful           *bool            `json:"isStateful,omitempty"`
	Name                 *string          `json:"name,omitempty"`
	ServiceDate          *time.Time       `json:"serviceDate,omitempty"`
	ServiceType          *string          `json:"serviceType,omitempty"`
	StartDate            *time.Time       `json:"startDate,omitempty"`
	StartMode            *string          `json:"startMode,omitempty"`
	State                *string          `json:"state,omitempty"`
	AiModelSpecification *EntityRef       `json:"aiModelSpecification,omitempty"`
	ServiceSpecification *EntityRef       `json:"serviceSpecification,omitempty"`
	GPU                  *ResourceRef     `json:"gpu,omitempty"`
	TrainingData         *EntityRef       `json:"trainingData,omitempty"`
	Note                 []Note           `json:"note,omitempty"`
	RelatedParty         []RelatedParty   `json:"relatedParty,omitempty"`
	ServiceCharacteristic []Characteristic `json:"serviceCharacteristic,omitempty"`
	AtBaseType           *string          `json:"@baseType,omitempty"`
	AtSchemaLocation     *string          `json:"@schemaLocation,omitempty"`
	AtType               *string          `json:"@type,omitempty"`
}
