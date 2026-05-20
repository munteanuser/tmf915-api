package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type TopicRepo struct{ db *sqlx.DB }

func NewTopicRepo(db *sqlx.DB) *TopicRepo { return &TopicRepo{db: db} }

func (r *TopicRepo) List(ctx context.Context, offset, limit int) ([]models.Topic, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Topic").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.Topic{}
	q := `SELECT id,href,content_query,header_query,name,base_type,schema_location,at_type
	      FROM Topic ORDER BY name
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	return rows, total, nil
}

func (r *TopicRepo) Get(ctx context.Context, id string) (*models.Topic, error) {
	var row models.Topic
	if err := r.db.GetContext(ctx, &row,
		"SELECT id,href,content_query,header_query,name,base_type,schema_location,at_type FROM Topic WHERE id=@id",
		sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	return &row, nil
}

func (r *TopicRepo) Exists(ctx context.Context, id string) (bool, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Topic WHERE id=@id", sql.Named("id", id)).Scan(&n); err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *TopicRepo) Create(ctx context.Context, in models.TopicCreate) (*models.Topic, error) {
	id := uuid.NewString()
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/topic/%s", id)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO Topic (id,href,content_query,header_query,name,base_type,schema_location,at_type)
		 VALUES(@id,@href,@cq,@hq,@name,@bt,@sl,@at)`,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("cq", in.ContentQuery), sql.Named("hq", in.HeaderQuery), sql.Named("name", in.Name),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *TopicRepo) Update(ctx context.Context, id string, in models.TopicUpdate) (*models.Topic, error) {
	e, err := r.Get(ctx, id)
	if err != nil || e == nil {
		return e, err
	}
	if in.ContentQuery != nil {
		e.ContentQuery = in.ContentQuery
	}
	if in.HeaderQuery != nil {
		e.HeaderQuery = in.HeaderQuery
	}
	if in.Name != nil {
		e.Name = *in.Name
	}
	_, err = r.db.ExecContext(ctx,
		"UPDATE Topic SET content_query=@cq,header_query=@hq,name=@name WHERE id=@id",
		sql.Named("cq", e.ContentQuery), sql.Named("hq", e.HeaderQuery),
		sql.Named("name", e.Name), sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *TopicRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM Topic WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
