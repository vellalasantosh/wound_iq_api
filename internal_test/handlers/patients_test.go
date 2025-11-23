package handlers_test

import (
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vellalasantosh/wound_iq_api/internal/handlers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// AnyTime allows matching any time value in sqlmock expectations.
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreatePatient_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT add_patient").
		WithArgs("John Doe", AnyTime{}, "Male", "MRN123").
		WillReturnRows(sqlmock.NewRows([]string{"add_patient"}).AddRow(42))

	router := gin.Default()
	router.POST("/v1/patients", handlers.CreatePatient(db))

	req := httptest.NewRequest("POST", "/v1/patients", strings.NewReader(`{"full_name":"John Doe","date_of_birth":"1990-01-01T00:00:00Z","gender":"Male","medical_record_number":"MRN123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
