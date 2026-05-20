package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type AiModelSpecificationRepo struct{ db *sqlx.DB }

func NewAiModelSpecificationRepo(db *sqlx.DB) *AiModelSpecificationRepo {
	return &AiModelSpecificationRepo{db: db}
}

func (r *AiModelSpecificationRepo) List(ctx context.Context, offset, limit int) ([]models.AiModelSpecification, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM AiModelSpecification").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	rows := []models.AiModelSpecification{}
	q := `SELECT id,href,description,is_bundle,last_update,lifecycle_status,name,version,
	             base_type,schema_location,at_type,valid_for,target_entity_schema,
	             attachment,constraint_refs,related_party
	      FROM AiModelSpecification
	      ORDER BY name
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q,
		sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	for i := range rows {
		if err := expandAiModelSpec(&rows[i]); err != nil {
			return nil, 0, err
		}
	}
	return rows, total, nil
}

func (r *AiModelSpecificationRepo) Get(ctx context.Context, id string) (*models.AiModelSpecification, error) {
	var row models.AiModelSpecification
	q := `SELECT id,href,description,is_bundle,last_update,lifecycle_status,name,version,
	             base_type,schema_location,at_type,valid_for,target_entity_schema,
	             attachment,constraint_refs,related_party
	      FROM AiModelSpecification WHERE id = @id`
	if err := r.db.GetContext(ctx, &row, q, sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	if err := expandAiModelSpec(&row); err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AiModelSpecificationRepo) Create(ctx context.Context, in models.AiModelSpecificationCreate) (*models.AiModelSpecification, error) {
	id := uuid.NewString()
	now := time.Now().UTC()

	validForRaw, _ := marshalJSON(in.ValidFor)
	tesRaw, _ := marshalJSON(in.TargetEntitySchema)
	attachRaw, _ := marshalJSON(in.Attachment)
	constRaw, _ := marshalJSON(in.Constraint)
	rpRaw, _ := marshalJSON(in.RelatedParty)

	q := `INSERT INTO AiModelSpecification
	      (id,href,description,is_bundle,last_update,lifecycle_status,name,version,
	       base_type,schema_location,at_type,valid_for,target_entity_schema,attachment,constraint_refs,related_party)
	      VALUES(@id,@href,@description,@is_bundle,@last_update,@lifecycle_status,@name,@version,
	             @base_type,@schema_location,@at_type,@valid_for,@tes,@attachment,@constraint_refs,@related_party)`

	basePath := fmt.Sprintf("/tmf-api/AiM/v4/aiModelSpecification/%s", id)
	_, err := r.db.ExecContext(ctx, q,
		sql.Named("id", id),
		sql.Named("href", basePath),
		sql.Named("description", in.Description),
		sql.Named("is_bundle", in.IsBundle),
		sql.Named("last_update", now),
		sql.Named("lifecycle_status", in.LifecycleStatus),
		sql.Named("name", in.Name),
		sql.Named("version", in.Version),
		sql.Named("base_type", in.AtBaseType),
		sql.Named("schema_location", in.AtSchemaLocation),
		sql.Named("at_type", in.AtType),
		sql.Named("valid_for", validForRaw),
		sql.Named("tes", tesRaw),
		sql.Named("attachment", attachRaw),
		sql.Named("constraint_refs", constRaw),
		sql.Named("related_party", rpRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiModelSpecificationRepo) Update(ctx context.Context, id string, in models.AiModelSpecificationUpdate) (*models.AiModelSpecification, error) {
	existing, err := r.Get(ctx, id)
	if err != nil || existing == nil {
		return existing, err
	}

	// Apply patch fields
	if in.Description != nil {
		existing.Description = in.Description
	}
	if in.IsBundle != nil {
		existing.IsBundle = in.IsBundle
	}
	if in.LifecycleStatus != nil {
		existing.LifecycleStatus = in.LifecycleStatus
	}
	if in.Name != nil {
		existing.Name = *in.Name
	}
	if in.Version != nil {
		existing.Version = in.Version
	}
	if in.AtBaseType != nil {
		existing.AtBaseType = in.AtBaseType
	}
	if in.AtSchemaLocation != nil {
		existing.AtSchemaLocation = in.AtSchemaLocation
	}
	if in.AtType != nil {
		existing.AtType = in.AtType
	}
	if in.ValidFor != nil {
		existing.ValidFor = in.ValidFor
	}
	if in.TargetEntitySchema != nil {
		existing.TargetEntitySchema = in.TargetEntitySchema
	}
	if in.Attachment != nil {
		existing.Attachment = in.Attachment
	}
	if in.Constraint != nil {
		existing.Constraint = in.Constraint
	}
	if in.RelatedParty != nil {
		existing.RelatedParty = in.RelatedParty
	}

	now := time.Now().UTC()
	validForRaw, _ := marshalJSON(existing.ValidFor)
	tesRaw, _ := marshalJSON(existing.TargetEntitySchema)
	attachRaw, _ := marshalJSON(existing.Attachment)
	constRaw, _ := marshalJSON(existing.Constraint)
	rpRaw, _ := marshalJSON(existing.RelatedParty)

	q := `UPDATE AiModelSpecification SET
	        description=@description, is_bundle=@is_bundle, last_update=@last_update,
	        lifecycle_status=@lifecycle_status, name=@name, version=@version,
	        base_type=@base_type, schema_location=@schema_location, at_type=@at_type,
	        valid_for=@valid_for, target_entity_schema=@tes,
	        attachment=@attachment, constraint_refs=@constraint_refs, related_party=@related_party
	      WHERE id=@id`
	_, err = r.db.ExecContext(ctx, q,
		sql.Named("description", existing.Description),
		sql.Named("is_bundle", existing.IsBundle),
		sql.Named("last_update", now),
		sql.Named("lifecycle_status", existing.LifecycleStatus),
		sql.Named("name", existing.Name),
		sql.Named("version", existing.Version),
		sql.Named("base_type", existing.AtBaseType),
		sql.Named("schema_location", existing.AtSchemaLocation),
		sql.Named("at_type", existing.AtType),
		sql.Named("valid_for", validForRaw),
		sql.Named("tes", tesRaw),
		sql.Named("attachment", attachRaw),
		sql.Named("constraint_refs", constRaw),
		sql.Named("related_party", rpRaw),
		sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiModelSpecificationRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM AiModelSpecification WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func expandAiModelSpec(r *models.AiModelSpecification) error {
	if err := unmarshalJSON(r.ValidForRaw, &r.ValidFor); err != nil {
		return err
	}
	if err := unmarshalJSON(r.TargetEntitySchemaRaw, &r.TargetEntitySchema); err != nil {
		return err
	}
	if err := unmarshalJSON(r.AttachmentRaw, &r.Attachment); err != nil {
		return err
	}
	if err := unmarshalJSON(r.ConstraintRaw, &r.Constraint); err != nil {
		return err
	}
	return unmarshalJSON(r.RelatedPartyRaw, &r.RelatedParty)
}
