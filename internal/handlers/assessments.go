package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vellalasantosh/wound_iq_api/internal/db"
	"github.com/vellalasantosh/wound_iq_api/internal/models"
	"github.com/vellalasantosh/wound_iq_api/internal/util"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *db.DB
}

func ListAssessments(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, pageSize := util.GetPaginationParams(c)

		patientIDStr := c.Query("patient_id")
		clinicianIDStr := c.Query("clinician_id")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		var filter models.AssessmentFilter

		if patientIDStr != "" {
			if v, err := strconv.Atoi(patientIDStr); err == nil {
				filter.PatientID = &v
			}
		}

		if clinicianIDStr != "" {
			if v, err := strconv.Atoi(clinicianIDStr); err == nil {
				filter.ClinicianID = &v
			}
		}

		if startDate != "" {
			filter.StartDate = &startDate
		}

		if endDate != "" {
			filter.EndDate = &endDate
		}

		database := &db.DB{DB: dbConn}

		assessments, err := database.GetAssessments(
			c.Request.Context(),
			filter,
			page,
			pageSize,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":      assessments,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

func CreateAssessment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var a models.Assessment
		if err := c.ShouldBindJSON(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call stored function add_full_assessment with all parameters (46)
		// Match the order from your SQL function exactly.
		query := `
            SELECT add_full_assessment(
                $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
                $11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
                $21,$22,$23,$24,$25,$26,$27,$28,$29,$30,
                $31,$32,$33,$34,$35,$36,$37,$38,$39,$40,
                $41,$42,$43,$44,$45,$46
            )
        `
		var newID int
		// prepare args in same sequence as function signature
		args := []interface{}{
			a.ClinicianID, a.PatientID, a.Location, a.Etiology, a.DepthOfInjury, a.Stage, a.Chronicity, a.HealingStatus, a.ReturnToClinic,
			a.LocalizedSymptoms, a.SystemicSymptoms, a.PainPresent, a.PainScore, a.CultureResults, a.Antibiotic,
			a.GranulationPercent, a.EpithelialPercent, a.SloughPercent, a.EscharPercent, a.NecroticPercent, a.Debridement,
			a.BloodPressure, a.Temperature, a.Pulse, a.RespirationRate, a.OxygenSaturation,
			a.WLength, a.WWidth, a.WDepth, a.WTunneling, a.Undermining, a.Edges, a.SkinCondition, a.Edema, a.Blister,
			a.ExudateType, a.ExudateAmount, a.Odor,
			a.PrimaryDressing, a.SecondaryDressing, a.TertiaryDressing, a.Frequency, a.Supplies, a.Orders,
		}

		// Ensure length
		if len(args) != 46 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("expected 46 args, got %d", len(args))})
			return
		}

		err := db.QueryRow(query, args...).Scan(&newID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		a.AssessmentID = newID
		if a.Date.IsZero() {
			a.Date = time.Now().UTC()
		}
		c.JSON(http.StatusCreated, gin.H{"data": a})
	}
}

func GetAssessment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var a models.Assessment
		var dt time.Time
		err := db.QueryRow("SELECT assessment_id, clinician_id, patient_id, date, location, etiology, depth_of_injury, stage, chronicity, healing_status, return_to_clinic FROM assessment WHERE assessment_id=$1", id).
			Scan(&a.AssessmentID, &a.ClinicianID, &a.PatientID, &dt, &a.Location, &a.Etiology, &a.DepthOfInjury, &a.Stage, &a.Chronicity, &a.HealingStatus, &a.ReturnToClinic)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		a.Date = dt.UTC()
		c.JSON(http.StatusOK, gin.H{"data": a})
	}
}

func UpdateAssessment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var a models.Assessment
		if err := c.ShouldBindJSON(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := db.Exec("UPDATE assessment SET clinician_id=$1, patient_id=$2, location=$3, etiology=$4, depth_of_injury=$5, stage=$6, chronicity=$7, healing_status=$8, return_to_clinic=$9 WHERE assessment_id=$10",
			a.ClinicianID, a.PatientID, a.Location, a.Etiology, a.DepthOfInjury, a.Stage, a.Chronicity, a.HealingStatus, a.ReturnToClinic, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": a})
	}
}

func DeleteAssessment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := db.Exec("DELETE FROM assessment WHERE assessment_id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func GetAssessmentFull(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		rows, err := db.Query("SELECT assessment_id, assessment_date, patient_name, clinician_name, location, stage, healing_status, pain_score, granulation_percent, w_length, w_width FROM get_assessment_full($1)", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		if rows.Next() {
			var out struct {
				AssessmentID   int     `json:"assessment_id"`
				AssessmentDate string  `json:"assessment_date"`
				PatientName    string  `json:"patient_name"`
				ClinicianName  string  `json:"clinician_name"`
				Location       string  `json:"location"`
				Stage          string  `json:"stage"`
				HealingStatus  string  `json:"healing_status"`
				PainScore      string  `json:"pain_score"`
				Granulation    int     `json:"granulation_percent"`
				WLength        float64 `json:"w_length"`
				WWidth         float64 `json:"w_width"`
			}
			var dt time.Time
			if err := rows.Scan(&out.AssessmentID, &dt, &out.PatientName, &out.ClinicianName, &out.Location, &out.Stage, &out.HealingStatus, &out.PainScore, &out.Granulation, &out.WLength, &out.WWidth); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			out.AssessmentDate = dt.UTC().Format(time.RFC3339)
			c.JSON(http.StatusOK, gin.H{"data": out})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}
