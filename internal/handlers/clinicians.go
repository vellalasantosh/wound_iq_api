package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/vellalasantosh/wound_iq_api/internal/models"
	"github.com/vellalasantosh/wound_iq_api/internal/util"

	"github.com/gin-gonic/gin"
)

func ListClinicians(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, pageSize := util.GetPaginationParams(c)
		offset := (page - 1) * pageSize

		rows, err := db.Query("SELECT clinician_id, full_name, role, department, contact_info, license_number FROM clinician ORDER BY full_name LIMIT $1 OFFSET $2", pageSize, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var list []models.Clinician
		for rows.Next() {
			var cl models.Clinician
			if err := rows.Scan(&cl.ClinicianID, &cl.FullName, &cl.Role, &cl.Department, &cl.ContactInfo, &cl.License); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			list = append(list, cl)
		}
		c.JSON(http.StatusOK, gin.H{"data": list, "page": page, "page_size": pageSize})
	}
}

func CreateClinician(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cl models.Clinician
		if err := c.ShouldBindJSON(&cl); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var newID int
		err := db.QueryRow("INSERT INTO clinician (full_name, role, department, contact_info, license_number) VALUES ($1,$2,$3,$4,$5) RETURNING clinician_id", cl.FullName, cl.Role, cl.Department, cl.ContactInfo, cl.License).Scan(&newID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cl.ClinicianID = newID
		c.JSON(http.StatusCreated, gin.H{"data": cl})
	}
}

func GetClinician(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var cl models.Clinician
		err := db.QueryRow("SELECT clinician_id, full_name, role, department, contact_info, license_number FROM clinician WHERE clinician_id=$1", id).
			Scan(&cl.ClinicianID, &cl.FullName, &cl.Role, &cl.Department, &cl.ContactInfo, &cl.License)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": cl})
	}
}

func UpdateClinician(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var cl models.Clinician
		if err := c.ShouldBindJSON(&cl); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := db.Exec("UPDATE clinician SET full_name=$1, role=$2, department=$3, contact_info=$4, license_number=$5 WHERE clinician_id=$6", cl.FullName, cl.Role, cl.Department, cl.ContactInfo, cl.License, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": cl})
	}
}

func DeleteClinician(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := db.Exec("DELETE FROM clinician WHERE clinician_id=$1", id)
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
