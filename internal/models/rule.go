package models

type Rule struct {
	ID               string  `json:"id"                        db:"id"`
	Href             *string `json:"href,omitempty"            db:"href"`
	Name             string  `json:"name"                      db:"name"`
	AtBaseType       *string `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string `json:"@type,omitempty"           db:"at_type"`
}

type RuleCreate struct {
	Name             string  `json:"name"`
	AtBaseType       *string `json:"@baseType,omitempty"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty"`
	AtType           *string `json:"@type,omitempty"`
}

type RuleUpdate struct {
	Name             *string `json:"name,omitempty"`
	AtBaseType       *string `json:"@baseType,omitempty"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty"`
	AtType           *string `json:"@type,omitempty"`
}
