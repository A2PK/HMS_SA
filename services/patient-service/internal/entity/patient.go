package entity

import (
	"time"

	coreEntity "golang-microservices-boilerplate/pkg/core/entity"

	"github.com/google/uuid"
)

// Patient represents a patient in the system.
type Patient struct {
	coreEntity.BaseEntity                 // Embedded base entity
	FirstName             string          `json:"first_name" gorm:"not null"`                                                                 // Patient's first name
	LastName              string          `json:"last_name" gorm:"not null"`                                                                  // Patient's last name
	DateOfBirth           time.Time       `json:"date_of_birth" gorm:"type:date"`                                                             // Patient's date of birth
	Gender                string          `json:"gender"`                                                                                     // Patient's gender
	PhoneNumber           string          `json:"phone_number" gorm:"uniqueIndex"`                                                            // Phone number, email, etc.
	Address               string          `json:"address"`                                                                                    // Patient's address
	MedicalHistory        []MedicalRecord `json:"medical_history" gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // List of medical records (GORM relation)
	// Removed ID, CreatedAt, UpdatedAt as they are in BaseEntity
}

// TableName overrides the default table name for the Patient model.
func (Patient) TableName() string {
	return "patients"
}

// MedicalRecord represents a single entry in a patient's medical history.
type MedicalRecord struct {
	coreEntity.BaseEntity           // Embedded base entity for MedicalRecord
	PatientID             uuid.UUID `json:"patient_id" gorm:"type:uuid;not null;index"` // Foreign key to Patient
	Date                  time.Time `json:"date" gorm:"not null"`                       // Date of the record
	StaffID               uuid.UUID `json:"doctor_id" gorm:"type:uuid;not null;index"`  // ID of the doctor who made the record (Could be uuid.UUID if Staff uses UUID)
	Diagnosis             string    `json:"diagnosis" gorm:"not null"`                  // Diagnosis details
	Treatment             string    `json:"treatment" gorm:"not null"`                  // Treatment provided
	Notes                 string    `json:"notes" gorm:"type:text"`                     // Additional notes
	// Removed RecordID
}

// TableName overrides the default table name for the MedicalRecord model.
func (MedicalRecord) TableName() string {
	return "medical_records"
}

// NewPatient creates a new Patient instance.
func NewPatient(firstName, lastName, gender, phoneNumber, address string, dob time.Time) *Patient {
	patient := &Patient{
		FirstName:      firstName,
		LastName:       lastName,
		DateOfBirth:    dob,
		Gender:         gender,
		PhoneNumber:    phoneNumber,
		Address:        address,
		MedicalHistory: []MedicalRecord{}, // Initialize as empty slice
	}
	// BaseEntity fields (ID, CreatedAt, UpdatedAt) are typically handled by GORM hooks or the repository
	return patient
}

// AddMedicalRecord adds a new record to the patient's medical history.
// Note: This method might change if using GORM associations directly.
func (p *Patient) AddMedicalRecord(record MedicalRecord) {
	// Ensure the record has the PatientID set if not done elsewhere
	if record.PatientID == uuid.Nil {
		record.PatientID = p.ID
	}
	p.MedicalHistory = append(p.MedicalHistory, record)
	p.UpdatedAt = time.Now() // Manually update timestamp if not relying solely on GORM hooks
}

// UpdateDetails updates non-identifying patient details.
func (p *Patient) UpdateDetails(firstName, lastName, gender, phoneNumber, address string, dob time.Time) {
	changed := false
	if firstName != "" && p.FirstName != firstName {
		p.FirstName = firstName
		changed = true
	}
	if lastName != "" && p.LastName != lastName {
		p.LastName = lastName
		changed = true
	}
	if gender != "" && p.Gender != gender {
		p.Gender = gender
		changed = true
	}
	if phoneNumber != "" && p.PhoneNumber != phoneNumber {
		p.PhoneNumber = phoneNumber
		changed = true
	}
	if address != "" && p.Address != address {
		p.Address = address
		changed = true
	}
	if !dob.IsZero() && p.DateOfBirth != dob {
		p.DateOfBirth = dob
		changed = true
	}

	if changed {
		p.UpdatedAt = time.Now() // Manually update timestamp if not relying solely on GORM hooks
	}
}
