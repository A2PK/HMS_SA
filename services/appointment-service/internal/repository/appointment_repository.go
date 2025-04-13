package repository

import (
	"context"
	"time"

	// Remove unused core imports
	// coreTypes "golang-microservices-boilerplate/pkg/core/types"
	coreRepository "golang-microservices-boilerplate/pkg/core/repository"
	"golang-microservices-boilerplate/services/appointment-service/internal/entity"

	"github.com/google/uuid"
)

// AppointmentRepository defines the interface for appointment data persistence.
// Implementations are expected to embed coreRepository.GormBaseRepository[entity.Appointment].
type AppointmentRepository interface {
	coreRepository.BaseRepository[entity.Appointment]
	// FindByPatientID retrieves all appointments for a specific patient.
	FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Appointment, error)

	// FindByDoctorID retrieves all appointments for a specific doctor within a time range.
	FindByDoctorID(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) ([]*entity.Appointment, error)

	// FindByDateRange retrieves all appointments within a specific date range.
	FindByDateRange(ctx context.Context, startTime time.Time, endTime time.Time) ([]*entity.Appointment, error)

	// CheckDoctorAvailability checks if a doctor has conflicting appointments at a given time.
	CheckDoctorAvailability(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) (bool, error)

	// --- Methods provided by embedding GormBaseRepository ---
	// Create(ctx context.Context, appointment *entity.Appointment) error
	// FindByID(ctx context.Context, id uuid.UUID) (*entity.Appointment, error)
	// Update(ctx context.Context, appointment *entity.Appointment) error
	// Delete(ctx context.Context, id uuid.UUID, hardDelete bool) error
	// FindAll(ctx context.Context, opts coreTypes.FilterOptions) (*coreTypes.PaginationResult[entity.Appointment], error)
	// FindWithFilter(ctx context.Context, filter map[string]interface{}, opts coreTypes.FilterOptions) (*coreTypes.PaginationResult[entity.Appointment], error)
	// Count(ctx context.Context, filter map[string]interface{}) (int64, error)
	// Transaction(ctx context.Context, fn func(txRepo repository.BaseRepository[entity.Appointment]) error) error
	// CreateMany(ctx context.Context, entities []*entity.Appointment) error
	// UpdateMany(ctx context.Context, entities []*entity.Appointment) error
	// DeleteMany(ctx context.Context, ids []uuid.UUID, hardDelete bool) error
}

// Note: Error handling (e.g., ErrAppointmentNotFound, ErrTimeSlotUnavailable) is crucial.
