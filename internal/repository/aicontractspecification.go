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

type AiContractSpecificationRepo struct{ db *sqlx.DB }

func NewAiContractSpecificationRepo(db *sqlx.DB) *AiContractSpecificationRepo {
	return &AiContractSpecificationRepo{db: db}
}

func (r *AiContractSpecificationRepo) List(ctx context.Context, offset, limit int) ([]models.AiContractSpecification, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM AiContractSpecification").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.AiContractSpecification{}
	q := `SELECT id,href,description,is_bundle,last_update,lifecycle_status,name,version,
	             base_type,schema_location,at_type,valid_for,target_entity_schema,
	             attachment,constraint_refs,related_party,spec_characteristic
	      FROM AiContractSpecification ORDER BY name
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	for i := range rows {
		expandAiContractSpec(&rows[i])
	}
	return rows, total, nil
}

func (r *AiContractSpecificationRepo) Get(ctx context.Context, id string) (*models.AiContractSpecification, error) {
	var row models.AiContractSpecification
	q := `SELECT id,href,description,is_bundle,last_update,lifecycle_status,name,version,
	             base_type,schema_location,at_type,valid_for,target_entity_schema,
	             attachment,constraint_refs,related_party,spec_characteristic
	      FROM AiContractSpecification WHERE id=@id`
	if err := r.db.GetContext(ctx, &row, q, sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	expandAiContractSpec(&row)
	return &row, nil
}

func (r *AiContractSpecificationRepo) Create(ctx context.Context, in models.AiContractSpecificationCreate) (*models.AiContractSpecification, error) {
	id := uuid.NewString()
	now := time.Now().UTC()

	vfRaw, _ := marshalJSON(in.ValidFor)
	tesRaw, _ := marshalJSON(in.TargetEntitySchema)
	attRaw, _ := marshalJSON(in.Attachment)
	cRaw, _ := marshalJSON(in.Constraint)
	rpRaw, _ := marshalJSON(in.RelatedParty)
	scRaw, _ := marshalJSON(in.SpecCharacteristic)

	q := `INSERT INTO AiContractSpecification
	      (id,href,description,is_bundle,last_update,lifecycle_status,name,version,
	       base_type,schema_location,at_type,valid_for,target_entity_schema,
	       attachment,constraint_refs,related_party,spec_characteristic)
	      VALUES(@id,@href,@desc,@bundle,@lu,@ls,@name,@ver,
	             @bt,@sl,@at,@vf,@tes,@att,@cr,@rp,@sc)`
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/aiContractSpecification/%s", id)
	_, err := r.db.ExecContext(ctx, q,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("desc", in.Description), sql.Named("bundle", in.IsBundle),
		sql.Named("lu", now), sql.Named("ls", in.LifecycleStatus),
		sql.Named("name", in.Name), sql.Named("ver", in.Version),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation),
		sql.Named("at", in.AtType), sql.Named("vf", vfRaw), sql.Named("tes", tesRaw),
		sql.Named("att", attRaw), sql.Named("cr", cRaw), sql.Named("rp", rpRaw), sql.Named("sc", scRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiContractSpecificationRepo) Update(ctx context.Context, id string, in models.AiContractSpecificationUpdate) (*models.AiContractSpecification, error) {
	e, err := r.Get(ctx, id)
	if err != nil || e == nil {
		return e, err
	}
	now := time.Now().UTC()
	if in.Description != nil {
		e.Description = in.Description
	}
	if in.IsBundle != nil {
		e.IsBundle = in.IsBundle
	}
	if in.LifecycleStatus != nil {
		e.LifecycleStatus = in.LifecycleStatus
	}
	if in.Name != nil {
		e.Name = *in.Name
	}
	if in.Version != nil {
		e.Version = in.Version
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
	if in.ValidFor != nil {
		e.ValidFor = in.ValidFor
	}
	if in.TargetEntitySchema != nil {
		e.TargetEntitySchema = in.TargetEntitySchema
	}
	if in.Attachment != nil {
		e.Attachment = in.Attachment
	}
	if in.Constraint != nil {
		e.Constraint = in.Constraint
	}
	if in.RelatedParty != nil {
		e.RelatedParty = in.RelatedParty
	}
	if in.SpecCharacteristic != nil {
		e.SpecCharacteristic = in.SpecCharacteristic
	}

	vfRaw, _ := marshalJSON(e.ValidFor)
	tesRaw, _ := marshalJSON(e.TargetEntitySchema)
	attRaw, _ := marshalJSON(e.Attachment)
	cRaw, _ := marshalJSON(e.Constraint)
	rpRaw, _ := marshalJSON(e.RelatedParty)
	scRaw, _ := marshalJSON(e.SpecCharacteristic)

	q := `UPDATE AiContractSpecification SET
	        description=@desc,is_bundle=@bundle,last_update=@lu,lifecycle_status=@ls,
	        name=@name,version=@ver,base_type=@bt,schema_location=@sl,at_type=@at,
	        valid_for=@vf,target_entity_schema=@tes,attachment=@att,
	        constraint_refs=@cr,related_party=@rp,spec_characteristic=@sc
	      WHERE id=@id`
	_, err = r.db.ExecContext(ctx, q,
		sql.Named("desc", e.Description), sql.Named("bundle", e.IsBundle),
		sql.Named("lu", now), sql.Named("ls", e.LifecycleStatus),
		sql.Named("name", e.Name), sql.Named("ver", e.Version),
		sql.Named("bt", e.AtBaseType), sql.Named("sl", e.AtSchemaLocation), sql.Named("at", e.AtType),
		sql.Named("vf", vfRaw), sql.Named("tes", tesRaw), sql.Named("att", attRaw),
		sql.Named("cr", cRaw), sql.Named("rp", rpRaw), sql.Named("sc", scRaw),
		sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiContractSpecificationRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM AiContractSpecification WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func expandAiContractSpec(r *models.AiContractSpecification) {
	unmarshalJSON(r.ValidForRaw, &r.ValidFor)
	unmarshalJSON(r.TargetEntitySchemaRaw, &r.TargetEntitySchema)
	unmarshalJSON(r.AttachmentRaw, &r.Attachment)
	unmarshalJSON(r.ConstraintRaw, &r.Constraint)
	unmarshalJSON(r.RelatedPartyRaw, &r.RelatedParty)
	unmarshalJSON(r.SpecCharacteristicRaw, &r.SpecCharacteristic)
}
