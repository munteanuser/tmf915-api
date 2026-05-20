package models

import "time"

type AiContract struct {
	ID               string     `json:"id"                        db:"id"`
	Href             *string    `json:"href,omitempty"            db:"href"`
	ApprovalDate     *time.Time `json:"approvalDate,omitempty"    db:"approval_date"`
	Approved         *bool      `json:"approved,omitempty"        db:"approved"`
	Description      *string    `json:"description,omitempty"     db:"description"`
	Name             string     `json:"name"                      db:"name"`
	State            *string    `json:"state,omitempty"           db:"state"`
	Version          *string    `json:"version,omitempty"         db:"version"`
	AtBaseType       *string    `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string    `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string    `json:"@type,omitempty"           db:"at_type"`

	AiContractSpecification    *EntityRef     `json:"aiContractSpecification,omitempty" db:"-"`
	AiContractSpecificationRaw *string        `json:"-"                                 db:"ai_contract_specification"`
	AiModel                    *EntityRef     `json:"aiModel,omitempty"                 db:"-"`
	AiModelRaw                 *string        `json:"-"                                 db:"ai_model"`
	Template                   *TemplateRef   `json:"template,omitempty"                db:"-"`
	TemplateRaw                *string        `json:"-"                                 db:"template_ref"`
	ValidFor                   *TimePeriod    `json:"validFor,omitempty"                db:"-"`
	ValidForRaw                *string        `json:"-"                                 db:"valid_for"`
	Characteristic             []Characteristic `json:"characteristic,omitempty"        db:"-"`
	CharacteristicRaw          *string        `json:"-"                                 db:"characteristics"`
	RelatedParty               []RelatedParty `json:"relatedParty,omitempty"            db:"-"`
	RelatedPartyRaw            *string        `json:"-"                                 db:"related_party"`
	Rule                       []RuleRef      `json:"rule,omitempty"                    db:"-"`
	RuleRaw                    *string        `json:"-"                                 db:"rules"`
}

type AiContractCreate struct {
	ApprovalDate            *time.Time     `json:"approvalDate,omitempty"`
	Approved                *bool          `json:"approved,omitempty"`
	Description             *string        `json:"description,omitempty"`
	Name                    string         `json:"name"`
	State                   *string        `json:"state,omitempty"`
	Version                 *string        `json:"version,omitempty"`
	AiContractSpecification *EntityRef     `json:"aiContractSpecification,omitempty"`
	AiModel                 *EntityRef     `json:"aiModel,omitempty"`
	Template                *TemplateRef   `json:"template,omitempty"`
	ValidFor                *TimePeriod    `json:"validFor,omitempty"`
	Characteristic          []Characteristic `json:"characteristic,omitempty"`
	RelatedParty            []RelatedParty `json:"relatedParty,omitempty"`
	Rule                    []RuleRef      `json:"rule,omitempty"`
	AtBaseType              *string        `json:"@baseType,omitempty"`
	AtSchemaLocation        *string        `json:"@schemaLocation,omitempty"`
	AtType                  *string        `json:"@type,omitempty"`
}

type AiContractUpdate struct {
	ApprovalDate            *time.Time     `json:"approvalDate,omitempty"`
	Approved                *bool          `json:"approved,omitempty"`
	Description             *string        `json:"description,omitempty"`
	Name                    *string        `json:"name,omitempty"`
	State                   *string        `json:"state,omitempty"`
	Version                 *string        `json:"version,omitempty"`
	AiContractSpecification *EntityRef     `json:"aiContractSpecification,omitempty"`
	AiModel                 *EntityRef     `json:"aiModel,omitempty"`
	Template                *TemplateRef   `json:"template,omitempty"`
	ValidFor                *TimePeriod    `json:"validFor,omitempty"`
	Characteristic          []Characteristic `json:"characteristic,omitempty"`
	RelatedParty            []RelatedParty `json:"relatedParty,omitempty"`
	Rule                    []RuleRef      `json:"rule,omitempty"`
	AtBaseType              *string        `json:"@baseType,omitempty"`
	AtSchemaLocation        *string        `json:"@schemaLocation,omitempty"`
	AtType                  *string        `json:"@type,omitempty"`
}
