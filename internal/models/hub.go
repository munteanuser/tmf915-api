package models

type Hub struct {
	ID               string  `json:"id"                        db:"id"`
	Href             *string `json:"href,omitempty"            db:"href"`
	Callback         string  `json:"callback"                  db:"callback"`
	Query            *string `json:"query,omitempty"           db:"query"`
	AtBaseType       *string `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string `json:"@type,omitempty"           db:"at_type"`
}

type HubCreate struct {
	Callback         string  `json:"callback"`
	Query            *string `json:"query,omitempty"`
	AtBaseType       *string `json:"@baseType,omitempty"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty"`
	AtType           *string `json:"@type,omitempty"`
}

// EventSubscription is used for the /listener endpoints.
type EventSubscription struct {
	ID       string  `json:"id"               db:"id"`
	Callback string  `json:"callback"         db:"callback"`
	Query    *string `json:"query,omitempty"  db:"query"`
}

type EventSubscriptionInput struct {
	Callback string  `json:"callback"`
	Query    *string `json:"query,omitempty"`
}
