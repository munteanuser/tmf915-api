package models

import "time"

// TimePeriod represents a time range.
type TimePeriod struct {
	StartDateTime *time.Time `json:"startDateTime,omitempty" db:"start_date_time"`
	EndDateTime   *time.Time `json:"endDateTime,omitempty"   db:"end_date_time"`
}

// EntityRef is a generic reference to another entity.
type EntityRef struct {
	ID              string  `json:"id"                        db:"ref_id"`
	Href            *string `json:"href,omitempty"            db:"ref_href"`
	Name            *string `json:"name,omitempty"            db:"ref_name"`
	AtReferredType  *string `json:"@referredType,omitempty"   db:"ref_referred_type"`
	AtBaseType      *string `json:"@baseType,omitempty"       db:"ref_base_type"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty" db:"ref_schema_location"`
	AtType          *string `json:"@type,omitempty"           db:"ref_type"`
}

// RelatedParty links a party to a resource.
type RelatedParty struct {
	ID             string  `json:"id"                       db:"party_id"`
	Href           *string `json:"href,omitempty"           db:"party_href"`
	Name           *string `json:"name,omitempty"           db:"party_name"`
	Role           *string `json:"role,omitempty"           db:"party_role"`
	AtReferredType *string `json:"@referredType,omitempty"  db:"party_referred_type"`
	AtBaseType     *string `json:"@baseType,omitempty"      db:"party_base_type"`
	AtType         *string `json:"@type,omitempty"          db:"party_type"`
}

// Note represents a free-text note.
type Note struct {
	ID       *string    `json:"id,omitempty"      db:"note_id"`
	Author   *string    `json:"author,omitempty"  db:"note_author"`
	Date     *time.Time `json:"date,omitempty"    db:"note_date"`
	Text     string     `json:"text"              db:"note_text"`
	AtBaseType *string  `json:"@baseType,omitempty" db:"note_base_type"`
	AtType     *string  `json:"@type,omitempty"     db:"note_type"`
}

// Characteristic is a name-value pair.
type Characteristic struct {
	ID             *string     `json:"id,omitempty"             db:"char_id"`
	Name           string      `json:"name"                     db:"char_name"`
	ValueType      *string     `json:"valueType,omitempty"      db:"char_value_type"`
	Value          interface{} `json:"value,omitempty"          db:"-"`
	ValueRaw       *string     `json:"-"                        db:"char_value"`
	AtBaseType     *string     `json:"@baseType,omitempty"      db:"char_base_type"`
	AtSchemaLocation *string   `json:"@schemaLocation,omitempty" db:"char_schema_location"`
	AtType         *string     `json:"@type,omitempty"          db:"char_type"`
}

// TargetEntitySchema references an external schema.
type TargetEntitySchema struct {
	AtSchemaLocation string `json:"@schemaLocation" db:"target_schema_location"`
	AtType           string `json:"@type"           db:"target_type"`
}

// Error is the standard TMF error response.
type Error struct {
	Code           string  `json:"code"`
	Reason         string  `json:"reason"`
	Message        *string `json:"message,omitempty"`
	Status         *string `json:"status,omitempty"`
	ReferenceError *string `json:"referenceError,omitempty"`
	AtBaseType     *string `json:"@baseType,omitempty"`
	AtSchemaLocation *string `json:"@schemaLocation,omitempty"`
	AtType         *string `json:"@type,omitempty"`
}

// AttachmentRef or Value
type AttachmentRefOrValue struct {
	ID          *string `json:"id,omitempty"          db:"attach_id"`
	Href        *string `json:"href,omitempty"        db:"attach_href"`
	Name        *string `json:"name,omitempty"        db:"attach_name"`
	Description *string `json:"description,omitempty" db:"attach_description"`
	MimeType    *string `json:"mimeType,omitempty"    db:"attach_mime_type"`
	URL         *string `json:"url,omitempty"         db:"attach_url"`
	AtBaseType  *string `json:"@baseType,omitempty"   db:"attach_base_type"`
	AtType      *string `json:"@type,omitempty"       db:"attach_type"`
}

// ConstraintRef references a constraint.
type ConstraintRef struct {
	ID             string  `json:"id"                       db:"const_id"`
	Href           *string `json:"href,omitempty"           db:"const_href"`
	Name           *string `json:"name,omitempty"           db:"const_name"`
	Version        *string `json:"version,omitempty"        db:"const_version"`
	AtReferredType *string `json:"@referredType,omitempty"  db:"const_referred_type"`
	AtBaseType     *string `json:"@baseType,omitempty"      db:"const_base_type"`
	AtType         *string `json:"@type,omitempty"          db:"const_type"`
}

// RuleRef references a Rule.
type RuleRef struct {
	ID             string  `json:"id"                       db:"rule_ref_id"`
	Href           *string `json:"href,omitempty"           db:"rule_ref_href"`
	Name           *string `json:"name,omitempty"           db:"rule_ref_name"`
	AtReferredType *string `json:"@referredType,omitempty"  db:"rule_ref_referred_type"`
	AtBaseType     *string `json:"@baseType,omitempty"      db:"rule_ref_base_type"`
	AtType         *string `json:"@type,omitempty"          db:"rule_ref_type"`
}

// TemplateRef references a template.
type TemplateRef struct {
	ID             string  `json:"id"                       db:"tmpl_id"`
	Href           *string `json:"href,omitempty"           db:"tmpl_href"`
	Name           *string `json:"name,omitempty"           db:"tmpl_name"`
	AtReferredType *string `json:"@referredType,omitempty"  db:"tmpl_referred_type"`
	AtBaseType     *string `json:"@baseType,omitempty"      db:"tmpl_base_type"`
	AtType         *string `json:"@type,omitempty"          db:"tmpl_type"`
}

// Violation describes the nature of an AI contract violation.
type Violation struct {
	ID          *string `json:"id,omitempty"          db:"viol_id"`
	Description *string `json:"description,omitempty" db:"viol_description"`
	Type        *string `json:"type,omitempty"        db:"viol_type"`
	AtBaseType  *string `json:"@baseType,omitempty"   db:"viol_base_type"`
	AtType      *string `json:"@type,omitempty"       db:"viol_type_ext"`
}

// Request / Response used by Monitor
type Request struct {
	Method  *string              `json:"method,omitempty"  db:"req_method"`
	URL     *string              `json:"url,omitempty"     db:"req_url"`
	Headers []HeaderItem         `json:"header,omitempty"  db:"-"`
	Body    *string              `json:"body,omitempty"    db:"req_body"`
}

type Response struct {
	StatusCode *int         `json:"statusCode,omitempty" db:"resp_status_code"`
	Headers    []HeaderItem `json:"header,omitempty"     db:"-"`
	Body       *string      `json:"body,omitempty"       db:"resp_body"`
}

type HeaderItem struct {
	Name  string `json:"name"  db:"header_name"`
	Value string `json:"value" db:"header_value"`
}

// ServiceSpecificationRef
type ServiceSpecificationRef struct {
	ID             string  `json:"id"                       db:"svcspec_id"`
	Href           *string `json:"href,omitempty"           db:"svcspec_href"`
	Name           *string `json:"name,omitempty"           db:"svcspec_name"`
	Version        *string `json:"version,omitempty"        db:"svcspec_version"`
	AtReferredType *string `json:"@referredType,omitempty"  db:"svcspec_referred_type"`
	AtBaseType     *string `json:"@baseType,omitempty"      db:"svcspec_base_type"`
	AtType         *string `json:"@type,omitempty"          db:"svcspec_type"`
}

// ResourceRef
type ResourceRef struct {
	ID             string  `json:"id"                       db:"res_id"`
	Href           *string `json:"href,omitempty"           db:"res_href"`
	Name           *string `json:"name,omitempty"           db:"res_name"`
	AtReferredType *string `json:"@referredType,omitempty"  db:"res_referred_type"`
	AtBaseType     *string `json:"@baseType,omitempty"      db:"res_base_type"`
	AtType         *string `json:"@type,omitempty"          db:"res_type"`
}
