package models

import "time"

type Assessment struct {
    AssessmentID    int       `json:"assessment_id,omitempty"`
    ClinicianID     int       `json:"clinician_id" binding:"required"`
    PatientID       int       `json:"patient_id" binding:"required"`
    Date            time.Time `json:"date,omitempty"`
    Location        string    `json:"location" binding:"required"`
    Etiology        string    `json:"etiology" binding:"required"`
    DepthOfInjury   string    `json:"depth_of_injury" binding:"required"`
    Stage           string    `json:"stage" binding:"required"`
    Chronicity      string    `json:"chronicity" binding:"required"`
    HealingStatus   string    `json:"healing_status" binding:"required"`
    ReturnToClinic  bool      `json:"return_to_clinic"`

    // infection_and_pain
    LocalizedSymptoms string `json:"localized_symptoms"`
    SystemicSymptoms  string `json:"systemic_symptoms"`
    PainPresent       string `json:"pain_present"`
    PainScore         string `json:"pain_score"`
    CultureResults    string `json:"culture_results"`
    Antibiotic        string `json:"antibiotic"`

    // tissue_status
    GranulationPercent int    `json:"granulation_percent"`
    EpithelialPercent  int    `json:"epithelial_percent"`
    SloughPercent      int    `json:"slough_percent"`
    EscharPercent      int    `json:"eschar_percent"`
    NecroticPercent    int    `json:"necrotic_percent"`
    Debridement        string `json:"debridement"`

    // vitals
    BloodPressure     string  `json:"blood_pressure"`
    Temperature       float64 `json:"temperature"`
    Pulse             int     `json:"pulse"`
    RespirationRate   int     `json:"respiration_rate"`
    OxygenSaturation  int     `json:"oxygen_saturation"`

    // wound_condition
    WLength      float64 `json:"w_length"`
    WWidth       float64 `json:"w_width"`
    WDepth       float64 `json:"w_depth"`
    WTunneling   bool    `json:"w_tunneling"`
    Undermining  bool    `json:"undermining"`
    Edges        string  `json:"edges"`
    SkinCondition string `json:"skin_condition"`
    Edema         string `json:"edema"`
    Blister       string `json:"blister"`

    // exudate
    ExudateType   string `json:"exudate_type"`
    ExudateAmount string `json:"exudate_amount"`
    Odor          string `json:"odor"`

    // treatment
    PrimaryDressing   string `json:"primary_dressing"`
    SecondaryDressing string `json:"secondary_dressing"`
    TertiaryDressing  string `json:"tertiary_dressing"`
    Frequency         string `json:"frequency"`
    Supplies          string `json:"supplies"`
    Orders            string `json:"orders"`
}
