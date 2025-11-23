package models

import "time"

type Patient struct {
    PatientID           int       `json:"patient_id,omitempty"`
    FullName            string    `json:"full_name" binding:"required"`
    DateOfBirth         time.Time `json:"date_of_birth" binding:"required"`
    Gender              string    `json:"gender" binding:"required"`
    MedicalRecordNumber string    `json:"medical_record_number" binding:"required"`
}
