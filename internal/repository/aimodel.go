package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type AiModelRepo struct{ db *sqlx.DB }

func NewAiModelRepo(db *sqlx.DB) *AiModelRepo { return &AiModelRepo{db: db} }

func (r *AiModelRepo) List(ctx context.Context, offset, limit int) ([]models.AiModel, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM AiModel").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.AiModel{}
	q := `SELECT id,href,category,description,end_date,has_started,is_bundle,is_service_enabled,
	             is_stateful,name,service_date,service_type,start_date,start_mode,state,
	             base_type,schema_location,at_type,
	             ai_model_specification,service_specification,gpu,training_data,
	             notes,places,related_entity,related_party,service_characteristic
	      FROM AiModel ORDER BY name
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	for i := range rows {
		expandAiModel(&rows[i])
	}
	return rows, total, nil
}

func (r *AiModelRepo) Get(ctx context.Context, id string) (*models.AiModel, error) {
	var row models.AiModel
	q := `SELECT id,href,category,description,end_date,has_started,is_bundle,is_service_enabled,
	             is_stateful,name,service_date,service_type,start_date,start_mode,state,
	             base_type,schema_location,at_type,
	             ai_model_specification,service_specification,gpu,training_data,
	             notes,places,related_entity,related_party,service_characteristic
	      FROM AiModel WHERE id=@id`
	if err := r.db.GetContext(ctx, &row, q, sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	expandAiModel(&row)
	return &row, nil
}

func (r *AiModelRepo) Create(ctx context.Context, in models.AiModelCreate) (*models.AiModel, error) {
	id := uuid.NewString()
	amsRaw, _ := marshalJSON(in.AiModelSpecification)
	svcspecRaw, _ := marshalJSON(in.ServiceSpecification)
	gpuRaw, _ := marshalJSON(in.GPU)
	tdRaw, _ := marshalJSON(in.TrainingData)
	noteRaw, _ := marshalJSON(in.Note)
	rpRaw, _ := marshalJSON(in.RelatedParty)
	scRaw, _ := marshalJSON(in.ServiceCharacteristic)

	q := `INSERT INTO AiModel
	      (id,href,category,description,end_date,has_started,is_bundle,is_service_enabled,
	       is_stateful,name,service_date,service_type,start_date,start_mode,state,
	       base_type,schema_location,at_type,
	       ai_model_specification,service_specification,gpu,training_data,
	       notes,related_party,service_characteristic)
	      VALUES(@id,@href,@category,@description,@end_date,@has_started,@is_bundle,@is_service_enabled,
	             @is_stateful,@name,@service_date,@service_type,@start_date,@start_mode,@state,
	             @base_type,@schema_location,@at_type,
	             @ams,@svcspec,@gpu,@td,@notes,@rp,@sc)`

	basePath := fmt.Sprintf("/tmf-api/AiM/v4/aiModel/%s", id)
	_, err := r.db.ExecContext(ctx, q,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("category", in.Category), sql.Named("description", in.Description),
		sql.Named("end_date", in.EndDate), sql.Named("has_started", in.HasStarted),
		sql.Named("is_bundle", in.IsBundle), sql.Named("is_service_enabled", in.IsServiceEnabled),
		sql.Named("is_stateful", in.IsStateful), sql.Named("name", in.Name),
		sql.Named("service_date", in.ServiceDate), sql.Named("service_type", in.ServiceType),
		sql.Named("start_date", in.StartDate), sql.Named("start_mode", in.StartMode),
		sql.Named("state", in.State),
		sql.Named("base_type", in.AtBaseType), sql.Named("schema_location", in.AtSchemaLocation),
		sql.Named("at_type", in.AtType),
		sql.Named("ams", amsRaw), sql.Named("svcspec", svcspecRaw),
		sql.Named("gpu", gpuRaw), sql.Named("td", tdRaw),
		sql.Named("notes", noteRaw), sql.Named("rp", rpRaw), sql.Named("sc", scRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiModelRepo) Update(ctx context.Context, id string, in models.AiModelUpdate) (*models.AiModel, error) {
	existing, err := r.Get(ctx, id)
	if err != nil || existing == nil {
		return existing, err
	}
	if in.Category != nil {
		existing.Category = in.Category
	}
	if in.Description != nil {
		existing.Description = in.Description
	}
	if in.EndDate != nil {
		existing.EndDate = in.EndDate
	}
	if in.HasStarted != nil {
		existing.HasStarted = in.HasStarted
	}
	if in.IsBundle != nil {
		existing.IsBundle = in.IsBundle
	}
	if in.IsServiceEnabled != nil {
		existing.IsServiceEnabled = in.IsServiceEnabled
	}
	if in.IsStateful != nil {
		existing.IsStateful = in.IsStateful
	}
	if in.Name != nil {
		existing.Name = *in.Name
	}
	if in.ServiceDate != nil {
		existing.ServiceDate = in.ServiceDate
	}
	if in.ServiceType != nil {
		existing.ServiceType = in.ServiceType
	}
	if in.StartDate != nil {
		existing.StartDate = in.StartDate
	}
	if in.StartMode != nil {
		existing.StartMode = in.StartMode
	}
	if in.State != nil {
		existing.State = in.State
	}
	if in.AiModelSpecification != nil {
		existing.AiModelSpecification = in.AiModelSpecification
	}
	if in.ServiceSpecification != nil {
		existing.ServiceSpecification = in.ServiceSpecification
	}
	if in.GPU != nil {
		existing.GPU = in.GPU
	}
	if in.TrainingData != nil {
		existing.TrainingData = in.TrainingData
	}
	if in.Note != nil {
		existing.Note = in.Note
	}
	if in.RelatedParty != nil {
		existing.RelatedParty = in.RelatedParty
	}
	if in.ServiceCharacteristic != nil {
		existing.ServiceCharacteristic = in.ServiceCharacteristic
	}

	amsRaw, _ := marshalJSON(existing.AiModelSpecification)
	svcspecRaw, _ := marshalJSON(existing.ServiceSpecification)
	gpuRaw, _ := marshalJSON(existing.GPU)
	tdRaw, _ := marshalJSON(existing.TrainingData)
	noteRaw, _ := marshalJSON(existing.Note)
	rpRaw, _ := marshalJSON(existing.RelatedParty)
	scRaw, _ := marshalJSON(existing.ServiceCharacteristic)

	q := `UPDATE AiModel SET
	        category=@category, description=@description, end_date=@end_date,
	        has_started=@has_started, is_bundle=@is_bundle, is_service_enabled=@is_service_enabled,
	        is_stateful=@is_stateful, name=@name, service_date=@service_date,
	        service_type=@service_type, start_date=@start_date, start_mode=@start_mode, state=@state,
	        base_type=@base_type, schema_location=@schema_location, at_type=@at_type,
	        ai_model_specification=@ams, service_specification=@svcspec,
	        gpu=@gpu, training_data=@td, notes=@notes, related_party=@rp, service_characteristic=@sc
	      WHERE id=@id`
	_, err = r.db.ExecContext(ctx, q,
		sql.Named("category", existing.Category), sql.Named("description", existing.Description),
		sql.Named("end_date", existing.EndDate), sql.Named("has_started", existing.HasStarted),
		sql.Named("is_bundle", existing.IsBundle), sql.Named("is_service_enabled", existing.IsServiceEnabled),
		sql.Named("is_stateful", existing.IsStateful), sql.Named("name", existing.Name),
		sql.Named("service_date", existing.ServiceDate), sql.Named("service_type", existing.ServiceType),
		sql.Named("start_date", existing.StartDate), sql.Named("start_mode", existing.StartMode),
		sql.Named("state", existing.State),
		sql.Named("base_type", existing.AtBaseType), sql.Named("schema_location", existing.AtSchemaLocation),
		sql.Named("at_type", existing.AtType),
		sql.Named("ams", amsRaw), sql.Named("svcspec", svcspecRaw),
		sql.Named("gpu", gpuRaw), sql.Named("td", tdRaw),
		sql.Named("notes", noteRaw), sql.Named("rp", rpRaw), sql.Named("sc", scRaw),
		sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AiModelRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM AiModel WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func expandAiModel(r *models.AiModel) {
	unmarshalJSON(r.AiModelSpecificationRaw, &r.AiModelSpecification)
	unmarshalJSON(r.ServiceSpecificationRaw, &r.ServiceSpecification)
	unmarshalJSON(r.GPURaw, &r.GPU)
	unmarshalJSON(r.TrainingDataRaw, &r.TrainingData)
	unmarshalJSON(r.NoteRaw, &r.Note)
	unmarshalJSON(r.RelatedPartyRaw, &r.RelatedParty)
	unmarshalJSON(r.ServiceCharacteristicRaw, &r.ServiceCharacteristic)
}
