package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type AiContractRepo struct{ db *sqlx.DB }

func NewAiContractRepo(db *sqlx.DB) *AiContractRepo { return &AiContractRepo{db: db} }

func (r *AiContractRepo) List(ctx context.Context, offset, limit int) ([]models.AiContract, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM AiContract").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.AiContract{}
	q := `SELECT id,href,approval_date,approved,description,name,state,version,
	             base_type,schema_location,at_type,
	             ai_contract_specification,ai_model,template_ref,valid_for,
	             characteristics,related_party,rules
	      FROM AiContract ORDER BY name
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	for i := range rows {
		expandAiContract(&rows[i])
	}
	return rows, total, nil
}

func (r *AiContractRepo) Get(ctx context.Context, id string) (*models.AiContract, error) {
	var row models.AiContract
	q := `SELECT id,href,approval_date,approved,description,name,state,version,
	             base_type,schema_location,at_type,
	             ai_contract_specification,ai_model,template_ref,valid_for,
	             characteristics,related_party,rules
	      FROM AiContract WHERE id=@id`
	if err := r.db.GetContext(ctx, &row, q, sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	expandAiContract(&row)
	return &row, nil
}

func (r *AiContractRepo) Create(ctx context.Context, in models.AiContractCreate) (*models.AiContract, error) {
	id := uuid.NewString()
	acsRaw, _ := marshalJSON(in.AiContractSpecification)
	amRaw, _ := marshalJSON(in.AiModel)
	tplRaw, _ := marshalJSON(in.Template)
	vfRaw, _ := marshalJSON(in.ValidFor)
	charRaw, _ := marshalJSON(in.Characteristic)
	rpRaw, _ := marshalJSON(in.RelatedParty)
	ruleRaw, _ := marshalJSON(in.Rule)

	q := `INSERT INTO AiContract
	      (id,href,approval_date,approved,description,name,state,version,
	       base_type,schema_location,at_type,
	       ai_contract_specification,ai_model,template_ref,valid_for,
	       characteristics,related_party,rules)
	      VALUES(@id,@href,@ad,@appr,@desc,@name,@state,@ver,
	             @bt,@sl,@at,@acs,@am,@tpl,@vf,@chars,@rp,@rules)`
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/aiContract/%s", id)
	_, err := r.db.ExecContext(ctx, q,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("ad", in.ApprovalDate), sql.Named("appr", in.Approved),
		sql.Named("desc", in.Description), sql.Named("name", in.Name),
		sql.Named("state", in.State), sql.Named("ver", in.Version),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
		sql.Named("acs", acsRaw), sql.Named("am", amRaw), sql.Named("tpl", tplRaw), sql.Named("vf", vfRaw),
		sql.Named("chars", charRaw), sql.Named("rp", rpRaw), sql.Named("rules", ruleRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiContractRepo) Update(ctx context.Context, id string, in models.AiContractUpdate) (*models.AiContract, error) {
	e, err := r.Get(ctx, id)
	if err != nil || e == nil {
		return e, err
	}
	if in.ApprovalDate != nil {
		e.ApprovalDate = in.ApprovalDate
	}
	if in.Approved != nil {
		e.Approved = in.Approved
	}
	if in.Description != nil {
		e.Description = in.Description
	}
	if in.Name != nil {
		e.Name = *in.Name
	}
	if in.State != nil {
		e.State = in.State
	}
	if in.Version != nil {
		e.Version = in.Version
	}
	if in.AiContractSpecification != nil {
		e.AiContractSpecification = in.AiContractSpecification
	}
	if in.AiModel != nil {
		e.AiModel = in.AiModel
	}
	if in.Template != nil {
		e.Template = in.Template
	}
	if in.ValidFor != nil {
		e.ValidFor = in.ValidFor
	}
	if in.Characteristic != nil {
		e.Characteristic = in.Characteristic
	}
	if in.RelatedParty != nil {
		e.RelatedParty = in.RelatedParty
	}
	if in.Rule != nil {
		e.Rule = in.Rule
	}

	acsRaw, _ := marshalJSON(e.AiContractSpecification)
	amRaw, _ := marshalJSON(e.AiModel)
	tplRaw, _ := marshalJSON(e.Template)
	vfRaw, _ := marshalJSON(e.ValidFor)
	charRaw, _ := marshalJSON(e.Characteristic)
	rpRaw, _ := marshalJSON(e.RelatedParty)
	ruleRaw, _ := marshalJSON(e.Rule)

	q := `UPDATE AiContract SET
	        approval_date=@ad,approved=@appr,description=@desc,name=@name,state=@state,version=@ver,
	        base_type=@bt,schema_location=@sl,at_type=@at,
	        ai_contract_specification=@acs,ai_model=@am,template_ref=@tpl,valid_for=@vf,
	        characteristics=@chars,related_party=@rp,rules=@rules
	      WHERE id=@id`
	_, err = r.db.ExecContext(ctx, q,
		sql.Named("ad", e.ApprovalDate), sql.Named("appr", e.Approved),
		sql.Named("desc", e.Description), sql.Named("name", e.Name),
		sql.Named("state", e.State), sql.Named("ver", e.Version),
		sql.Named("bt", e.AtBaseType), sql.Named("sl", e.AtSchemaLocation), sql.Named("at", e.AtType),
		sql.Named("acs", acsRaw), sql.Named("am", amRaw), sql.Named("tpl", tplRaw), sql.Named("vf", vfRaw),
		sql.Named("chars", charRaw), sql.Named("rp", rpRaw), sql.Named("rules", ruleRaw),
		sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiContractRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM AiContract WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func expandAiContract(r *models.AiContract) {
	unmarshalJSON(r.AiContractSpecificationRaw, &r.AiContractSpecification)
	unmarshalJSON(r.AiModelRaw, &r.AiModel)
	unmarshalJSON(r.TemplateRaw, &r.Template)
	unmarshalJSON(r.ValidForRaw, &r.ValidFor)
	unmarshalJSON(r.CharacteristicRaw, &r.Characteristic)
	unmarshalJSON(r.RelatedPartyRaw, &r.RelatedParty)
	unmarshalJSON(r.RuleRaw, &r.Rule)
}
