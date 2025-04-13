package entity

import (
	"fmt"
	"time"

	coreEntity "golang-microservices-boilerplate/pkg/core/entity"

	"github.com/google/uuid"
)

// AppointmentStatus defines the status of an appointment.
type AppointmentStatus string

const (
	Scheduled AppointmentStatus = "Scheduled"
	Confirmed AppointmentStatus = "Confirmed"
	Cancelled AppointmentStatus = "Cancelled"
	Completed AppointmentStatus = "Completed"
	NoShow    AppointmentStatus = "NoShow"
)

// Appointment represents a scheduled appointment between a patient and a doctor.
type Appointment struct {
	coreEntity.BaseEntity
	PatientID       uuid.UUID         `json:"patient_id" gorm:"type:uuid;index"` // ID of the patient
	DoctorID        uuid.UUID         `json:"doctor_id" gorm:"type:uuid;index"`  // ID of the doctor (a staff member)
	AppointmentTime time.Time         `json:"appointment_time" gorm:"index"`     // Date and time of the appointment
	Duration        time.Duration     `json:"duration"`                          // Duration of the appointment
	Reason          string            `json:"reason"`                            // Reason for the appointment
	Status          AppointmentStatus `json:"status"`                            // Current status of the appointment
	Notes           string            `json:"notes,omitempty"`                   // Optional notes about the appointment
	Place           string            `json:"place,omitempty" gorm:"index"`      // Added: Location/Place of the appointment
	// Removed ID, CreatedAt, UpdatedAt

	// Optional: Embed or reference full patient/doctor details if needed frequently,
	// but often just IDs are sufficient for the core appointment logic.
	// Patient *patientEntity.Patient `json:"patient,omitempty" gorm:"foreignKey:PatientID"` // Define patientEntity alias if needed
	// Doctor  *staffEntity.Staff    `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`   // Define staffEntity alias if needed
}

// DoctorAvailability struct might not be needed here if availability is primarily managed
// by the Staff service and checked via the client interface.
// If kept, it should also embed BaseEntity and use UUIDs.
/*
type DoctorAvailability struct {
	coreEntity.BaseEntity
	DoctorID  uuid.UUID `json:"doctor_id" gorm:"type:uuid;index"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsBooked  bool      `json:"is_booked"`
}
*/

// NewAppointment creates a new Appointment instance.
func NewAppointment(patientID, doctorID uuid.UUID, reason, place string, appointmentTime time.Time, duration time.Duration) *Appointment {
	apt := &Appointment{
		PatientID:       patientID,
		DoctorID:        doctorID,
		AppointmentTime: appointmentTime,
		Duration:        duration,
		Reason:          reason,
		Place:           place,
		Status:          Scheduled, // Initial status
	}
	// BaseEntity fields handled by GORM/Repository
	return apt
}

// SetStatus updates the appointment status, checking for valid transitions.
func (a *Appointment) SetStatus(newStatus AppointmentStatus) error {
	// Basic validation: Prevent setting back from terminal states
	if (a.Status == Completed || a.Status == Cancelled) && a.Status != newStatus {
		return fmt.Errorf("cannot change status of a %s appointment to %s", a.Status, newStatus)
	}

	if a.Status != newStatus {
		a.Status = newStatus
		a.UpdatedAt = time.Now()
	}
	return nil
}

// Reschedule updates the appointment time and optionally duration and place.
func (a *Appointment) Reschedule(newTime time.Time, newDuration *time.Duration, newPlace *string) error {
	// Check if appointment is in a state that allows rescheduling
	if a.Status != Scheduled && a.Status != Confirmed {
		return fmt.Errorf("cannot reschedule appointment with status %s", a.Status)
	}

	changed := false
	if !newTime.IsZero() && a.AppointmentTime != newTime {
		a.AppointmentTime = newTime
		changed = true
	}
	if newDuration != nil && a.Duration != *newDuration {
		a.Duration = *newDuration
		changed = true
	}
	if newPlace != nil && a.Place != *newPlace {
		a.Place = *newPlace
		changed = true
	}

	if changed {
		// Reset status to Scheduled after reschedule? Or keep Confirmed if it was Confirmed?
		// Let's keep Confirmed if it was, otherwise set to Scheduled.
		if a.Status != Confirmed {
			a.Status = Scheduled
		}
		a.UpdatedAt = time.Now()
	}
	return nil
}
