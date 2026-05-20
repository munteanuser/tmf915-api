package models

import "time"

type AiContractSpecification struct {
	ID               string     `json:"id"                        db:"id"`
	Href             *string    `json:"href,omitempty"            db:"href"`
	Description      *string    `json:"description,omitempty"     db:"description"`
	IsBundle         *bool      `json:"isBundle,omitempty"        db:"is_bundle"`
	LastUpdate       *time.Time `json:"lastUpdate,omitempty"      db:"last_update"`
	LifecycleStatus  *string    `json:"lifecycleStatus,omitempty" db:"lifecycle_status"`
	Name             string     `json:"name"                      db:"name"`
	Version          *string    `json:"version,omitempty"         db:"version"`
	AtBaseType       *string    `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string    `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string    `json:"@type,omitempty"           db:"at_type"`

	ValidFor              *TimePeriod            `json:"validFor,omitempty"           db:"-"`
	ValidForRaw           *string                `json:"-"                            db:"valid_for"`
	TargetEntitySchema    *TargetEntitySchema    `json:"targetEntitySchema,omitempty" db:"-"`
	TargetEntitySchemaRaw *string                `json:"-"                            db:"target_entity_schema"`
	Attachment            []AttachmentRefOrValue `json:"attachment,omitempty"         db:"-"`
	AttachmentRaw         *string                `json:"-"                            db:"attachment"`
	Constraint            []ConstraintRef        `json:"constraint,omitempty"         db:"-"`
	ConstraintRaw         *string                `json:"-"                            db:"constraint_refs"`
	RelatedParty          []RelatedParty         `json:"relatedParty,omitempty"       db:"-"`
	RelatedPartyRaw       *string                `json:"-"                            db:"related_party"`
	SpecCharacteristic    []Characteristic       `json:"specCharacteristic,omitempty" db:"-"`
	SpecCharacteristicRaw *string                `json:"-"                            db:"spec_characteristic"`
}

type AiContractSpecificationCreate struct {
	Description        *string                `json:"description,omitempty"`
	IsBundle           *bool                  `json:"isBundle,omitempty"`
	LastUpdate         *time.Time             `json:"lastUpdate,omitempty"`
	LifecycleStatus    *string                `json:"lifecycleStatus,omitempty"`
	Name               string                 `json:"name"`
	Version            *string                `json:"version,omitempty"`
	ValidFor           *TimePeriod            `json:"validFor,omitempty"`
	TargetEntitySchema *TargetEntitySchema    `json:"targetEntitySchema,omitempty"`
	Attachment         []AttachmentRefOrValue `json:"attachment,omitempty"`
	Constraint         []ConstraintRef        `json:"constraint,omitempty"`
	RelatedParty       []RelatedParty         `json:"relatedParty,omitempty"`
	SpecCharacteristic []Characteristic       `json:"specCharacteristic,omitempty"`
	AtBaseType         *string                `json:"@baseType,omitempty"`
	AtSchemaLocation   *string                `json:"@schemaLocation,omitempty"`
	AtType             *string                `json:"@type,omitempty"`
}

type AiContractSpecificationUpdate struct {
	Description        *string                `json:"description,omitempty"`
	IsBundle           *bool                  `json:"isBundle,omitempty"`
	LastUpdate         *time.Time             `json:"lastUpdate,omitempty"`
	LifecycleStatus    *string                `json:"lifecycleStatus,omitempty"`
	Name               *string                `json:"name,omitempty"`
	Version            *string                `json:"version,omitempty"`
	ValidFor           *TimePeriod            `json:"validFor,omitempty"`
	TargetEntitySchema *TargetEntitySchema    `json:"targetEntitySchema,omitempty"`
	Attachment         []AttachmentRefOrValue `json:"attachment,omitempty"`
	Constraint         []ConstraintRef        `json:"constraint,omitempty"`
	RelatedParty       []RelatedParty         `json:"relatedParty,omitempty"`
	SpecCharacteristic []Characteristic       `json:"specCharacteristic,omitempty"`
	AtBaseType         *string                `json:"@baseType,omitempty"`
	AtSchemaLocation   *string                `json:"@schemaLocation,omitempty"`
	AtType             *string                `json:"@type,omitempty"`
}
