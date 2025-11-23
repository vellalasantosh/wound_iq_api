package models

type AssessmentFilter struct {
	PatientID   *int
	ClinicianID *int
	StartDate   *string
	EndDate     *string
}
