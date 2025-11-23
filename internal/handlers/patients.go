package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"wound_iq_api/internal/models"
	"wound_iq_api/internal/util"

	"github.com/gin-gonic/gin"
)

func ListPatients(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, pageSize := util.GetPaginationParams(c)
		offset := (page - 1) * pageSize

		rows, err := db.Query("SELECT patient_id, full_name, date_of_birth, gender, medical_record_number FROM get_all_patients() LIMIT $1 OFFSET $2", pageSize, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var list []models.Patient
		for rows.Next() {
			var p models.Patient
			var dob time.Time
			if err := rows.Scan(&p.PatientID, &p.FullName, &dob, &p.Gender, &p.MedicalRecordNumber); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			p.DateOfBirth = dob.UTC()
			list = append(list, p)
		}
		c.JSON(http.StatusOK, gin.H{"data": list, "page": page, "page_size": pageSize})
	}
}

func CreatePatient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p models.Patient
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var newID int
		err := db.QueryRow(
			"SELECT add_patient($1, $2, $3, $4)",
			p.FullName, p.DateOfBirth, p.Gender, p.MedicalRecordNumber,
		).Scan(&newID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		p.PatientID = newID
		c.JSON(http.StatusCreated, gin.H{"data": p})
	}
}

func GetPatient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var p models.Patient
		var dob time.Time
		err := db.QueryRow("SELECT patient_id, full_name, date_of_birth, gender, medical_record_number FROM patient WHERE patient_id=$1", id).
			Scan(&p.PatientID, &p.FullName, &dob, &p.Gender, &p.MedicalRecordNumber)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		p.DateOfBirth = dob.UTC()
		c.JSON(http.StatusOK, gin.H{"data": p})
	}
}

func UpdatePatient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var p models.Patient
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := db.Exec("UPDATE patient SET full_name=$1, date_of_birth=$2, gender=$3, medical_record_number=$4 WHERE patient_id=$5",
			p.FullName, p.DateOfBirth, p.Gender, p.MedicalRecordNumber, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": p})
	}
}

func DeletePatient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := db.Exec("DELETE FROM patient WHERE patient_id=$1", id)
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

func GetPatientHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		rows, err := db.Query("SELECT assessment_id, assessment_date, location, stage, healing_status FROM get_patient_wound_history($1)", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		type hist struct {
			AssessmentID   int    `json:"assessment_id"`
			AssessmentDate string `json:"assessment_date"`
			Location       string `json:"location"`
			Stage          string `json:"stage"`
			HealingStatus  string `json:"healing_status"`
		}
		var out []hist
		for rows.Next() {
			var h hist
			var dt time.Time
			if err := rows.Scan(&h.AssessmentID, &dt, &h.Location, &h.Stage, &h.HealingStatus); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			h.AssessmentDate = dt.UTC().Format(time.RFC3339)
			out = append(out, h)
		}
		c.JSON(http.StatusOK, gin.H{"data": out})
	}
}
