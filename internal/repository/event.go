package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type EventRepo struct{ db *sqlx.DB }

func NewEventRepo(db *sqlx.DB) *EventRepo { return &EventRepo{db: db} }

func (r *EventRepo) List(ctx context.Context, topicID string, offset, limit int) ([]models.Event, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM Event WHERE topic_id=@tid", sql.Named("tid", topicID)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.Event{}
	q := `SELECT id,href,correlation_id,description,domain,event_id,event_time,event_type,priority,
	             time_occurred,title,topic_id,base_type,schema_location,at_type,
	             event_payload,related_party,reporting_system,source
	      FROM Event WHERE topic_id=@tid ORDER BY event_time DESC
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q,
		sql.Named("tid", topicID), sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	for i := range rows {
		expandEvent(&rows[i])
	}
	return rows, total, nil
}

func (r *EventRepo) Get(ctx context.Context, topicID, id string) (*models.Event, error) {
	var row models.Event
	q := `SELECT id,href,correlation_id,description,domain,event_id,event_time,event_type,priority,
	             time_occurred,title,topic_id,base_type,schema_location,at_type,
	             event_payload,related_party,reporting_system,source
	      FROM Event WHERE id=@id AND topic_id=@tid`
	if err := r.db.GetContext(ctx, &row, q, sql.Named("id", id), sql.Named("tid", topicID)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	expandEvent(&row)
	return &row, nil
}

func (r *EventRepo) Create(ctx context.Context, topicID string, in models.EventCreate) (*models.Event, error) {
	id := uuid.NewString()
	payloadRaw, _ := marshalJSON(in.EventPayload)
	rpRaw, _ := marshalJSON(in.RelatedParty)
	rsRaw, _ := marshalJSON(in.ReportingSystem)
	srcRaw, _ := marshalJSON(in.Source)

	q := `INSERT INTO Event
	      (id,href,correlation_id,description,domain,event_id,event_time,event_type,priority,
	       time_occurred,title,topic_id,base_type,schema_location,at_type,
	       event_payload,related_party,reporting_system,source)
	      VALUES(@id,@href,@cid,@desc,@domain,@eid,@et,@etype,@prio,@to,@title,@tid,
	             @bt,@sl,@at,@payload,@rp,@rs,@src)`
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/topic/%s/event/%s", topicID, id)
	_, err := r.db.ExecContext(ctx, q,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("cid", in.CorrelationID), sql.Named("desc", in.Description),
		sql.Named("domain", in.Domain), sql.Named("eid", in.EventID),
		sql.Named("et", in.EventTime), sql.Named("etype", in.EventType),
		sql.Named("prio", in.Priority), sql.Named("to", in.TimeOccurred),
		sql.Named("title", in.Title), sql.Named("tid", topicID),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
		sql.Named("payload", payloadRaw), sql.Named("rp", rpRaw),
		sql.Named("rs", rsRaw), sql.Named("src", srcRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, topicID, id)
}

func expandEvent(r *models.Event) {
	unmarshalJSON(r.EventPayloadRaw, &r.EventPayload)
	unmarshalJSON(r.RelatedPartyRaw, &r.RelatedParty)
	unmarshalJSON(r.ReportingSystemRaw, &r.ReportingSystem)
	unmarshalJSON(r.SourceRaw, &r.Source)
}
