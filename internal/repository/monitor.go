package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type MonitorRepo struct{ db *sqlx.DB }

func NewMonitorRepo(db *sqlx.DB) *MonitorRepo { return &MonitorRepo{db: db} }

func (r *MonitorRepo) List(ctx context.Context, offset, limit int) ([]models.Monitor, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Monitor").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.Monitor{}
	q := `SELECT id,href,source_href,state,base_type,schema_location,at_type,request,response
	      FROM Monitor ORDER BY id
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	for i := range rows {
		unmarshalJSON(rows[i].RequestRaw, &rows[i].Request)
		unmarshalJSON(rows[i].ResponseRaw, &rows[i].Response)
	}
	return rows, total, nil
}

func (r *MonitorRepo) Get(ctx context.Context, id string) (*models.Monitor, error) {
	var row models.Monitor
	if err := r.db.GetContext(ctx, &row,
		"SELECT id,href,source_href,state,base_type,schema_location,at_type,request,response FROM Monitor WHERE id=@id",
		sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	unmarshalJSON(row.RequestRaw, &row.Request)
	unmarshalJSON(row.ResponseRaw, &row.Response)
	return &row, nil
}

func (r *MonitorRepo) Create(ctx context.Context, in models.MonitorCreate) (*models.Monitor, error) {
	id := uuid.NewString()
	reqRaw, _ := marshalJSON(in.Request)
	respRaw, _ := marshalJSON(in.Response)
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/monitor/%s", id)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO Monitor (id,href,source_href,state,base_type,schema_location,at_type,request,response)
		 VALUES(@id,@href,@sh,@state,@bt,@sl,@at,@req,@resp)`,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("sh", in.SourceHref), sql.Named("state", in.State),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
		sql.Named("req", reqRaw), sql.Named("resp", respRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *MonitorRepo) Update(ctx context.Context, id string, in models.MonitorUpdate) (*models.Monitor, error) {
	e, err := r.Get(ctx, id)
	if err != nil || e == nil {
		return e, err
	}
	if in.SourceHref != nil {
		e.SourceHref = in.SourceHref
	}
	if in.State != nil {
		e.State = in.State
	}
	if in.Request != nil {
		e.Request = in.Request
	}
	if in.Response != nil {
		e.Response = in.Response
	}
	reqRaw, _ := marshalJSON(e.Request)
	respRaw, _ := marshalJSON(e.Response)
	_, err = r.db.ExecContext(ctx,
		"UPDATE Monitor SET source_href=@sh,state=@state,request=@req,response=@resp WHERE id=@id",
		sql.Named("sh", e.SourceHref), sql.Named("state", e.State),
		sql.Named("req", reqRaw), sql.Named("resp", respRaw), sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *MonitorRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM Monitor WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
