package usecase

import (
	"context"
	"time"

	coreUseCase "golang-microservices-boilerplate/pkg/core/usecase"
	pb "golang-microservices-boilerplate/proto/appointment-service" // Import proto
	"golang-microservices-boilerplate/services/appointment-service/internal/entity"

	"github.com/google/uuid"
	// We might need staff service client/interface here to check availability
	// staffUseCase "golang-microservices-boilerplate/services/staff-service/internal/usecase"
)

// AppointmentUseCase defines the interface for appointment management operations.
type AppointmentUseCase interface {
	// ScheduleAppointment creates a new appointment.
	coreUseCase.BaseUseCase[entity.Appointment, pb.ScheduleAppointmentRequest, pb.RescheduleAppointmentRequest]
	// It needs to check doctor availability before scheduling.
	ScheduleAppointment(ctx context.Context, req *pb.ScheduleAppointmentRequest) (*entity.Appointment, error)

	// GetAppointmentDetails retrieves details of a specific appointment by ID.
	GetAppointmentDetails(ctx context.Context, appointmentID uuid.UUID) (*entity.Appointment, error)

	// UpdateAppointmentStatus updates the status of an appointment (e.g., Confirmed, Cancelled).
	UpdateAppointmentStatus(ctx context.Context, appointmentID uuid.UUID, req *pb.UpdateAppointmentStatusRequest) (*entity.Appointment, error)

	// RescheduleAppointment changes the time or duration of an existing appointment.
	// It needs to re-check doctor availability.
	RescheduleAppointment(ctx context.Context, appointmentID uuid.UUID, req *pb.RescheduleAppointmentRequest) (*entity.Appointment, error)

	// CancelAppointment cancels an existing appointment.
	CancelAppointment(ctx context.Context, appointmentID uuid.UUID) error

	// GetAppointmentsForPatient retrieves all appointments for a given patient.
	GetAppointmentsForPatient(ctx context.Context, patientID uuid.UUID) ([]*entity.Appointment, error)

	// GetAppointmentsForDoctor retrieves appointments for a doctor within a time range.
	GetAppointmentsForDoctor(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) ([]*entity.Appointment, error)

	// CheckDoctorAvailability checks if a specific time slot is free for a doctor.
	// Returns true if available, false otherwise.
	CheckDoctorAvailability(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) (bool, error)
}
