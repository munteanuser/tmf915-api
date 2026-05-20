package models

import "time"

type AiContractViolation struct {
	ID               string     `json:"id"                        db:"id"`
	Href             *string    `json:"href,omitempty"            db:"href"`
	Date             *time.Time `json:"date,omitempty"            db:"date"`
	AtBaseType       *string    `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string    `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string    `json:"@type,omitempty"           db:"at_type"`

	AiContract      *EntityRef     `json:"aiContract,omitempty"  db:"-"`
	AiContractRaw   *string        `json:"-"                     db:"ai_contract"`
	Violation       *Violation     `json:"violation,omitempty"   db:"-"`
	ViolationRaw    *string        `json:"-"                     db:"violation"`
	RelatedParty    []RelatedParty `json:"relatedParty,omitempty" db:"-"`
	RelatedPartyRaw *string        `json:"-"                     db:"related_party"`
}

type AiContractViolationCreate struct {
	Date         *time.Time     `json:"date,omitempty"`
	AiContract   *EntityRef     `json:"aiContract,omitempty"`
	Violation    *Violation     `json:"violation,omitempty"`
	RelatedParty []RelatedParty `json:"relatedParty,omitempty"`
	AtBaseType   *string        `json:"@baseType,omitempty"`
	AtSchemaLocation *string    `json:"@schemaLocation,omitempty"`
	AtType       *string        `json:"@type,omitempty"`
}
