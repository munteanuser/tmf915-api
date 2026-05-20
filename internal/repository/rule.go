package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type RuleRepo struct{ db *sqlx.DB }

func NewRuleRepo(db *sqlx.DB) *RuleRepo { return &RuleRepo{db: db} }

func (r *RuleRepo) List(ctx context.Context, offset, limit int) ([]models.Rule, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM [Rule]").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.Rule{}
	q := `SELECT id,href,name,base_type,schema_location,at_type FROM [Rule] ORDER BY name
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	return rows, total, nil
}

func (r *RuleRepo) Get(ctx context.Context, id string) (*models.Rule, error) {
	var row models.Rule
	if err := r.db.GetContext(ctx, &row,
		"SELECT id,href,name,base_type,schema_location,at_type FROM [Rule] WHERE id=@id",
		sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	return &row, nil
}

func (r *RuleRepo) Create(ctx context.Context, in models.RuleCreate) (*models.Rule, error) {
	id := uuid.NewString()
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/rule/%s", id)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO [Rule] (id,href,name,base_type,schema_location,at_type)
		 VALUES(@id,@href,@name,@bt,@sl,@at)`,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("name", in.Name), sql.Named("bt", in.AtBaseType),
		sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *RuleRepo) Update(ctx context.Context, id string, in models.RuleUpdate) (*models.Rule, error) {
	e, err := r.Get(ctx, id)
	if err != nil || e == nil {
		return e, err
	}
	if in.Name != nil {
		e.Name = *in.Name
	}
	if in.AtBaseType != nil {
		e.AtBaseType = in.AtBaseType
	}
	if in.AtSchemaLocation != nil {
		e.AtSchemaLocation = in.AtSchemaLocation
	}
	if in.AtType != nil {
		e.AtType = in.AtType
	}
	_, err = r.db.ExecContext(ctx,
		"UPDATE [Rule] SET name=@name,base_type=@bt,schema_location=@sl,at_type=@at WHERE id=@id",
		sql.Named("name", e.Name), sql.Named("bt", e.AtBaseType),
		sql.Named("sl", e.AtSchemaLocation), sql.Named("at", e.AtType), sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *RuleRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM [Rule] WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
