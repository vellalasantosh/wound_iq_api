package models

type Clinician struct {
    ClinicianID int    `json:"clinician_id,omitempty"`
    FullName    string `json:"full_name" binding:"required"`
    Role        string `json:"role" binding:"required"`
    Department  string `json:"department" binding:"required"`
    ContactInfo string `json:"contact_info" binding:"required"`
    License     string `json:"license_number" binding:"required"`
}
