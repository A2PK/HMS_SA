package repository

import (
	"context"
	"time"

	coreRepository "golang-microservices-boilerplate/pkg/core/repository"
	"golang-microservices-boilerplate/services/appointment-service/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ensure GormAppointmentRepository implements AppointmentRepository
var _ AppointmentRepository = (*GormAppointmentRepository)(nil)

type GormAppointmentRepository struct {
	*coreRepository.GormBaseRepository[entity.Appointment] // Embed the generic GORM repository
}

// NewGormAppointmentRepository creates a new GORM appointment repository
func NewGormAppointmentRepository(db *gorm.DB) *GormAppointmentRepository {
	//nolint:all // Suppress T does not satisfy Entity constraint due to pointer receivers
	baseRepo := coreRepository.NewGormBaseRepository[entity.Appointment](db)
	return &GormAppointmentRepository{
		GormBaseRepository: baseRepo,
	}
}

// --- Implementations for methods specific to AppointmentRepository ---

// FindByID could be overridden for specific preloading if needed
/*
func (r *GormAppointmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Appointment, error) {
    // Example: Preload patient and doctor info if needed often
    appointment := &entity.Appointment{}
    err := r.DB.WithContext(ctx).Preload("Patient").Preload("Doctor").First(appointment, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("entity not found")
        }
        return nil, err
    }
    return appointment, nil
}
*/

// FindByPatientID retrieves appointments for a specific patient.
func (r *GormAppointmentRepository) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Appointment, error) {
	var appointments []*entity.Appointment
	err := r.DB.WithContext(ctx).Where("patient_id = ?", patientID).Order("appointment_time DESC").Find(&appointments).Error
	return appointments, err
}

// FindByDoctorID retrieves appointments for a specific doctor within a time range.
func (r *GormAppointmentRepository) FindByDoctorID(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) ([]*entity.Appointment, error) {
	var appointments []*entity.Appointment
	err := r.DB.WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Where("appointment_time >= ?", startTime).
		Where("appointment_time < ?", endTime).
		Order("appointment_time ASC").
		Find(&appointments).Error
	return appointments, err
}

// FindByDateRange retrieves all appointments within a specific date range.
func (r *GormAppointmentRepository) FindByDateRange(ctx context.Context, startTime time.Time, endTime time.Time) ([]*entity.Appointment, error) {
	var appointments []*entity.Appointment
	err := r.DB.WithContext(ctx).
		Where("appointment_time >= ?", startTime).
		Where("appointment_time < ?", endTime).
		Order("appointment_time ASC").
		Find(&appointments).Error
	return appointments, err
}

// CheckDoctorAvailability checks for conflicting appointments for a doctor in the given time range.
// Returns true if NO conflicting appointments are found (slot is free), false otherwise.
func (r *GormAppointmentRepository) CheckDoctorAvailability(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&entity.Appointment{}).
		Where("doctor_id = ?", doctorID).
		Where("status != ?", entity.Cancelled).                                                // Don't count cancelled appointments as conflicts
		Where("appointment_time < ?", endTime).                                                // Existing appointment starts before the new one ends
		Where("TIMESTAMPADD(SECOND, duration / 1000000000, appointment_time) > ?", startTime). // Existing appointment ends after the new one starts (Requires adapting Duration to DB function)
		// Note: The above line is complex and DB-specific. GORM might not directly translate duration.
		// A simpler check might be less precise: Where("appointment_time BETWEEN ? AND ?", startTime, endTime)
		// Or query potentially overlapping appointments and check overlap in Go code.
		// Simplified check (less precise): check if any appointment starts within the slot.
		// Where("appointment_time >= ? AND appointment_time < ?", startTime, endTime).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil // True (available) if count is 0
}

// Override FindAll if necessary for specific preloading for Appointments
/*
func (r *GormAppointmentRepository) FindAll(ctx context.Context, opts coreTypes.FilterOptions) (*coreTypes.PaginationResult[entity.Appointment], error) {
    r.DB = r.DB.Preload("Patient").Preload("Doctor") // Example
    return r.GormBaseRepository.FindAll(ctx, opts)
}
*/
