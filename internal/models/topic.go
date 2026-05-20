package models

type Topic struct {
	ID               string  `json:"id"                        db:"id"`
	Href             *string `json:"href,omitempty"            db:"href"`
	ContentQuery     *string `json:"contentQuery,omitempty"    db:"content_query"`
	HeaderQuery      *string `json:"headerQuery,omitempty"     db:"header_query"`
	Name             string  `json:"name"                      db:"name"`
	AtBaseType       *string `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string `json:"@type,omitempty"           db:"at_type"`
}

type TopicCreate struct {
	ContentQuery     *string `json:"contentQuery,omitempty"`
	HeaderQuery      *string `json:"headerQuery,omitempty"`
	Name             string  `json:"name"`
	AtBaseType       *string `json:"@baseType,omitempty"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty"`
	AtType           *string `json:"@type,omitempty"`
}

type TopicUpdate struct {
	ContentQuery     *string `json:"contentQuery,omitempty"`
	HeaderQuery      *string `json:"headerQuery,omitempty"`
	Name             *string `json:"name,omitempty"`
	AtBaseType       *string `json:"@baseType,omitempty"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty"`
	AtType           *string `json:"@type,omitempty"`
}
