package usecase

import (
	"time"
	// "github.com/google/uuid" // No longer needed for this struct
)

// AvailableTimeSlot represents a time slot where a staff member (usually a doctor)
// is available for appointments, as determined by the staff service.
// This matches the structure returned by the staff service's GetDoctorAvailability RPC.
type AvailableTimeSlot struct {
	// EntryID, StaffID, IsBooked removed as they are not part of the core availability info returned.
	StartTime time.Time `json:"start_time"` // Start time of the availability block
	EndTime   time.Time `json:"end_time"`   // End time of the availability block
}
