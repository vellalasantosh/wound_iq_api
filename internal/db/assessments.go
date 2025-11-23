package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"wound_iq_api/internal/models"
)

// IMPORTANT: this keeps database/sql import ACTIVE
type DB struct {
	DB *sql.DB
}

func (d *DB) GetAssessments(
	ctx context.Context,
	filter models.AssessmentFilter,
	page, pageSize int,
) ([]models.Assessment, error) {

	offset := (page - 1) * pageSize

	baseQuery := `
        SELECT 
            assessment_id, clinician_id, patient_id, date,
            location, etiology, depth_of_injury, stage,
            chronicity, healing_status, return_to_clinic
        FROM assessment
    `

	where := []string{}
	args := []any{}
	i := 1

	if filter.PatientID != nil {
		where = append(where, fmt.Sprintf("patient_id = $%d", i))
		args = append(args, *filter.PatientID)
		i++
	}

	if filter.ClinicianID != nil {
		where = append(where, fmt.Sprintf("clinician_id = $%d", i))
		args = append(args, *filter.ClinicianID)
		i++
	}

	if filter.StartDate != nil {
		where = append(where, fmt.Sprintf("date >= $%d", i))
		args = append(args, *filter.StartDate)
		i++
	}

	if filter.EndDate != nil {
		where = append(where, fmt.Sprintf("date <= $%d", i))
		args = append(args, *filter.EndDate)
		i++
	}

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY date DESC LIMIT %d OFFSET %d", pageSize, offset)

	rows, err := d.DB.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []models.Assessment

	for rows.Next() {
		var a models.Assessment
		err := rows.Scan(
			&a.AssessmentID, &a.ClinicianID, &a.PatientID, &a.Date,
			&a.Location, &a.Etiology, &a.DepthOfInjury, &a.Stage,
			&a.Chronicity, &a.HealingStatus, &a.ReturnToClinic,
		)
		if err != nil {
			return nil, err
		}

		assessments = append(assessments, a)
	}

	return assessments, nil
}
