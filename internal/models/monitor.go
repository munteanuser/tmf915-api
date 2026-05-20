package models

type Monitor struct {
	ID               string  `json:"id"                        db:"id"`
	Href             *string `json:"href,omitempty"            db:"href"`
	SourceHref       *string `json:"sourceHref,omitempty"      db:"source_href"`
	State            *string `json:"state,omitempty"           db:"state"`
	AtBaseType       *string `json:"@baseType,omitempty"       db:"base_type"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty" db:"schema_location"`
	AtType           *string `json:"@type,omitempty"           db:"at_type"`

	Request     *Request  `json:"request,omitempty"  db:"-"`
	RequestRaw  *string   `json:"-"                  db:"request"`
	Response    *Response `json:"response,omitempty" db:"-"`
	ResponseRaw *string   `json:"-"                  db:"response"`
}

type MonitorCreate struct {
	SourceHref       *string   `json:"sourceHref,omitempty"`
	State            *string   `json:"state,omitempty"`
	Request          *Request  `json:"request,omitempty"`
	Response         *Response `json:"response,omitempty"`
	AtBaseType       *string   `json:"@baseType,omitempty"`
	AtSchemaLocation *string   `json:"@schemaLocation,omitempty"`
	AtType           *string   `json:"@type,omitempty"`
}

type MonitorUpdate struct {
	SourceHref       *string   `json:"sourceHref,omitempty"`
	State            *string   `json:"state,omitempty"`
	Request          *Request  `json:"request,omitempty"`
	Response         *Response `json:"response,omitempty"`
	AtBaseType       *string   `json:"@baseType,omitempty"`
	AtSchemaLocation *string   `json:"@schemaLocation,omitempty"`
	AtType           *string   `json:"@type,omitempty"`
}
