package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmf915-api/internal/models"
)

type AlarmRepo struct{ db *sqlx.DB }

func NewAlarmRepo(db *sqlx.DB) *AlarmRepo { return &AlarmRepo{db: db} }

func (r *AlarmRepo) List(ctx context.Context, offset, limit int) ([]models.Alarm, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Alarm").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}
	rows := []models.Alarm{}
	q := `SELECT id,href,ack_state,ack_system_id,ack_user_id,alarm_changed_time,alarm_cleared_time,
	             alarm_details,alarm_escalation,alarm_raised_time,alarm_reporting_time,
	             alarmed_object_type,clear_system_id,clear_user_id,external_alarm_id,is_root_cause,
	             planned_outage_indicator,probable_cause,proposed_repaired_actions,reporting_system_id,
	             service_affecting,source_system_id,specific_problem,state,
	             base_type,schema_location,at_type,
	             alarm_type,alarmed_object,perceived_severity,crossed_threshold_information,
	             affected_service,comments,correlated_alarm,parent_alarm,places
	      FROM Alarm ORDER BY alarm_raised_time DESC
	      OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY`
	if err := r.db.SelectContext(ctx, &rows, q, sql.Named("offset", offset), sql.Named("limit", limit)); err != nil {
		return nil, 0, fmt.Errorf("select: %w", err)
	}
	return rows, total, nil
}

func (r *AlarmRepo) Get(ctx context.Context, id string) (*models.Alarm, error) {
	var row models.Alarm
	q := `SELECT id,href,ack_state,ack_system_id,ack_user_id,alarm_changed_time,alarm_cleared_time,
	             alarm_details,alarm_escalation,alarm_raised_time,alarm_reporting_time,
	             alarmed_object_type,clear_system_id,clear_user_id,external_alarm_id,is_root_cause,
	             planned_outage_indicator,probable_cause,proposed_repaired_actions,reporting_system_id,
	             service_affecting,source_system_id,specific_problem,state,
	             base_type,schema_location,at_type,
	             alarm_type,alarmed_object,perceived_severity,crossed_threshold_information,
	             affected_service,comments,correlated_alarm,parent_alarm,places
	      FROM Alarm WHERE id=@id`
	if err := r.db.GetContext(ctx, &row, q, sql.Named("id", id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	return &row, nil
}

func (r *AlarmRepo) Create(ctx context.Context, in models.AlarmCreate) (*models.Alarm, error) {
	id := uuid.NewString()
	atRaw, _ := toJSONStr(in.AlarmType)
	aoRaw, _ := toJSONStr(in.AlarmedObject)
	psRaw, _ := toJSONStr(in.PerceivedSeverity)

	q := `INSERT INTO Alarm
	      (id,href,ack_state,alarm_details,alarm_escalation,alarm_raised_time,alarmed_object_type,
	       external_alarm_id,is_root_cause,probable_cause,specific_problem,state,
	       base_type,schema_location,at_type,alarm_type,alarmed_object,perceived_severity)
	      VALUES(@id,@href,@as,@ad,@ae,@art,@aot,@eai,@irc,@pc,@sp,@state,@bt,@sl,@at,@altype,@ao,@ps)`
	basePath := fmt.Sprintf("/tmf-api/AiM/v4/alarm/%s", id)
	_, err := r.db.ExecContext(ctx, q,
		sql.Named("id", id), sql.Named("href", basePath),
		sql.Named("as", in.AckState), sql.Named("ad", in.AlarmDetails),
		sql.Named("ae", in.AlarmEscalation), sql.Named("art", in.AlarmRaisedTime),
		sql.Named("aot", in.AlarmedObjectType), sql.Named("eai", in.ExternalAlarmId),
		sql.Named("irc", in.IsRootCause), sql.Named("pc", in.ProbableCause),
		sql.Named("sp", in.SpecificProblem), sql.Named("state", in.State),
		sql.Named("bt", in.AtBaseType), sql.Named("sl", in.AtSchemaLocation), sql.Named("at", in.AtType),
		sql.Named("altype", atRaw), sql.Named("ao", aoRaw), sql.Named("ps", psRaw),
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AlarmRepo) Update(ctx context.Context, id string, in models.AlarmUpdate) (*models.Alarm, error) {
	e, err := r.Get(ctx, id)
	if err != nil || e == nil {
		return e, err
	}
	if in.AckState != nil {
		e.AckState = in.AckState
	}
	if in.AckSystemId != nil {
		e.AckSystemId = in.AckSystemId
	}
	if in.AckUserId != nil {
		e.AckUserId = in.AckUserId
	}
	if in.AlarmChangedTime != nil {
		e.AlarmChangedTime = in.AlarmChangedTime
	}
	if in.AlarmClearedTime != nil {
		e.AlarmClearedTime = in.AlarmClearedTime
	}
	if in.AlarmDetails != nil {
		e.AlarmDetails = in.AlarmDetails
	}
	if in.AlarmEscalation != nil {
		e.AlarmEscalation = in.AlarmEscalation
	}
	if in.PlannedOutageIndicator != nil {
		e.PlannedOutageIndicator = in.PlannedOutageIndicator
	}
	if in.ProbableCause != nil {
		e.ProbableCause = in.ProbableCause
	}
	if in.ProposedRepairedActions != nil {
		e.ProposedRepairedActions = in.ProposedRepairedActions
	}
	if in.ServiceAffecting != nil {
		e.ServiceAffecting = in.ServiceAffecting
	}
	if in.SpecificProblem != nil {
		e.SpecificProblem = in.SpecificProblem
	}
	if in.State != nil {
		e.State = in.State
	}

	q := `UPDATE Alarm SET
	        ack_state=@as,ack_system_id=@asid,ack_user_id=@auid,alarm_changed_time=@act,
	        alarm_cleared_time=@aclt,alarm_details=@ad,alarm_escalation=@ae,
	        planned_outage_indicator=@poi,probable_cause=@pc,proposed_repaired_actions=@pra,
	        service_affecting=@sa,specific_problem=@sp,state=@state
	      WHERE id=@id`
	_, err = r.db.ExecContext(ctx, q,
		sql.Named("as", e.AckState), sql.Named("asid", e.AckSystemId), sql.Named("auid", e.AckUserId),
		sql.Named("act", e.AlarmChangedTime), sql.Named("aclt", e.AlarmClearedTime),
		sql.Named("ad", e.AlarmDetails), sql.Named("ae", e.AlarmEscalation),
		sql.Named("poi", e.PlannedOutageIndicator), sql.Named("pc", e.ProbableCause),
		sql.Named("pra", e.ProposedRepairedActions), sql.Named("sa", e.ServiceAffecting),
		sql.Named("sp", e.SpecificProblem), sql.Named("state", e.State),
		sql.Named("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *AlarmRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM Alarm WHERE id=@id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func toJSONStr(v interface{}) (*string, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	s := string(b)
	return &s, nil
}
