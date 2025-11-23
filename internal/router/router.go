package router

import (
	"database/sql"

	"github.com/vellalasantosh/wound_iq_api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		p := v1.Group("/patients")
		{
			p.GET("", handlers.ListPatients(db))
			p.POST("", handlers.CreatePatient(db))
			p.GET("/:id", handlers.GetPatient(db))
			p.PUT("/:id", handlers.UpdatePatient(db))
			p.DELETE("/:id", handlers.DeletePatient(db))
			p.GET("/:id/history", handlers.GetPatientHistory(db))
		}

		c := v1.Group("/clinicians")
		{
			c.GET("", handlers.ListClinicians(db))
			c.POST("", handlers.CreateClinician(db))
			c.GET("/:id", handlers.GetClinician(db))
			c.PUT("/:id", handlers.UpdateClinician(db))
			c.DELETE("/:id", handlers.DeleteClinician(db))
		}

		a := v1.Group("/assessments")
		{
			a.GET("", handlers.ListAssessments(db))
			a.POST("", handlers.CreateAssessment(db))
			a.GET("/:id", handlers.GetAssessment(db))
			a.PUT("/:id", handlers.UpdateAssessment(db))
			a.DELETE("/:id", handlers.DeleteAssessment(db))
			a.GET("/:id/full", handlers.GetAssessmentFull(db))
		}
	}

	return r
}
