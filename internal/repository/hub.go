package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type HubRepo struct{ db *sqlx.DB }

func NewHubRepo(db *sqlx.DB) *HubRepo { return &HubRepo{db: db} }

func (r *HubRepo) List(ctx context.Context) ([]models.Hub, error) {
	rows := []models.Hub{}
	if err := r.db.SelectContext(ctx, &rows,
		"SELECT id,href,callback,query,base_type,schema_location,at_type FROM Hub ORDER BY id"); err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	return rows, nil
}

func (r *HubRepo) Get(ctx context.Context, id string) (*models.Hub, error) {
	var row models.Hub
	if err := r.db.GetContext(ctx, &row,
		"SELECT id,href,callback,query,base_type,schema_location,at_type FROM Hub WHERE id=@id",
		sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	return &row, nil
}

func (r *HubRepo) Create(ctx context.Context, in models.HubCreate) (*models.Hub, error) {
	id := uuid.NewString()
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/hub/%s", id)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO Hub (id,href,callback,query,base_type,schema_location,at_type)
		 VALUES(@id,@href,@cb,@query,@bt,@sl,@at)`,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("cb", in.Callback), sql.Named("query", in.Query),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *HubRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM Hub WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Listener (EventSubscription) operations

type ListenerRepo struct{ db *sqlx.DB }

func NewListenerRepo(db *sqlx.DB) *ListenerRepo { return &ListenerRepo{db: db} }

func (r *ListenerRepo) Create(ctx context.Context, in models.EventSubscriptionInput) (*models.EventSubscription, error) {
	id := uuid.NewString()
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO Listener (id,callback,query) VALUES(@id,@cb,@query)",
		sql.Named("id", id), sql.Named("cb", in.Callback), sql.Named("query", in.Query),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &models.EventSubscription{ID: id, Callback: in.Callback, Query: in.Query}, nil
}

func (r *ListenerRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM Listener WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *HubRepo) AllCallbacks(ctx context.Context) ([]string, error) {
	var callbacks []string
	if err := r.db.SelectContext(ctx, &callbacks, "SELECT callback FROM Hub"); err != nil {
		return nil, err
	}
	return callbacks, nil
}
