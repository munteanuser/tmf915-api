package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type AiContractViolationRepo struct{ db *sqlx.DB }

func NewAiContractViolationRepo(db *sqlx.DB) *AiContractViolationRepo {
	return &AiContractViolationRepo{db: db}
}

func (r *AiContractViolationRepo) List(ctx context.Context, offset, limit int) ([]models.AiContractViolation, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM AiContractViolation").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.AiContractViolation{}
	q := `SELECT id,href,date,base_type,schema_location,at_type,ai_contract,violation,related_party
	      FROM AiContractViolation ORDER BY date DESC
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	for i := range rows {
		unmarshalJSON(rows[i].AiContractRaw, &rows[i].AiContract)
		unmarshalJSON(rows[i].ViolationRaw, &rows[i].Violation)
		unmarshalJSON(rows[i].RelatedPartyRaw, &rows[i].RelatedParty)
	}
	return rows, total, nil
}

func (r *AiContractViolationRepo) Get(ctx context.Context, id string) (*models.AiContractViolation, error) {
	var row models.AiContractViolation
	q := `SELECT id,href,date,base_type,schema_location,at_type,ai_contract,violation,related_party
	      FROM AiContractViolation WHERE id=@id`
	if err := r.db.GetContext(ctx, &row, q, sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	unmarshalJSON(row.AiContractRaw, &row.AiContract)
	unmarshalJSON(row.ViolationRaw, &row.Violation)
	unmarshalJSON(row.RelatedPartyRaw, &row.RelatedParty)
	return &row, nil
}

func (r *AiContractViolationRepo) Create(ctx context.Context, in models.AiContractViolationCreate) (*models.AiContractViolation, error) {
	id := uuid.NewString()
	acRaw, _ := marshalJSON(in.AiContract)
	vRaw, _ := marshalJSON(in.Violation)
	rpRaw, _ := marshalJSON(in.RelatedParty)

	q := `INSERT INTO AiContractViolation (id,href,date,base_type,schema_location,at_type,ai_contract,violation,related_party)
	      VALUES(@id,@href,@date,@bt,@sl,@at,@ac,@viol,@rp)`
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/aiContractViolation/%s", id)
	_, err := r.db.ExecContext(ctx, q,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("date", in.Date),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
		sql.Named("ac", acRaw), sql.Named("viol", vRaw), sql.Named("rp", rpRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}
