package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "golang-microservices-boilerplate/proto/appointment-service"
	staff_pb "golang-microservices-boilerplate/proto/staff-service"

	coreLogger "golang-microservices-boilerplate/pkg/core/logger"
	coreRepository "golang-microservices-boilerplate/pkg/core/repository"
	coreUseCase "golang-microservices-boilerplate/pkg/core/usecase"

	// coreDTO "golang-microservices-boilerplate/pkg/core/dto"

	"golang-microservices-boilerplate/services/appointment-service/internal/entity"
	"golang-microservices-boilerplate/services/appointment-service/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	// No longer importing anything from staff service internal
)

// StaffServiceClient defines the interface for communication with the staff service.
// This would typically be implemented by a gRPC client adapter.
type StaffServiceClient interface {
	// GetDoctorAvailability checks if a doctor has available slots in a given time range.
	// Returns AvailableTimeSlot used by this service.
	GetDoctorAvailability(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) ([]AvailableTimeSlot, error)

	// MarkSlotAsBooked tells the staff service to mark a specific time slot as booked.
	// MarkSlotAsBooked(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) error

	// MarkSlotAsAvailable tells the staff service to mark a specific time slot as available (e.g., after cancellation).
	// MarkSlotAsAvailable(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) error
}

// --- gRPC Client Adapter ---

// Ensure adapter implements the interface.
var _ StaffServiceClient = (*grpcStaffServiceClientAdapter)(nil)

// grpcStaffServiceClientAdapter implements StaffServiceClient using a gRPC client.
type grpcStaffServiceClientAdapter struct {
	client staff_pb.StaffServiceClient // The actual generated gRPC client
	logger coreLogger.Logger
}

// NewGrpcStaffServiceClientAdapter creates a new adapter.
func NewGrpcStaffServiceClientAdapter(client staff_pb.StaffServiceClient, logger coreLogger.Logger) StaffServiceClient {
	return &grpcStaffServiceClientAdapter{
		client: client,
		logger: logger,
	}
}

// GetDoctorAvailability implements the StaffServiceClient interface.
// Updated to return []AvailableTimeSlot
func (a *grpcStaffServiceClientAdapter) GetDoctorAvailability(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) ([]AvailableTimeSlot, error) {
	a.logger.Debug("Calling StaffService.GetDoctorAvailability via gRPC adapter", "doctorID", doctorID, "start", startTime, "end", endTime)

	protoReq := &staff_pb.GetDoctorAvailabilityRequest{
		DoctorId:  doctorID.String(),
		StartTime: timestamppb.New(startTime),
		EndTime:   timestamppb.New(endTime),
	}

	protoResp, err := a.client.GetDoctorAvailability(ctx, protoReq)
	if err != nil {
		// Basic gRPC error handling
		st, ok := status.FromError(err)
		if ok {
			a.logger.Error("StaffService.GetDoctorAvailability gRPC error", "code", st.Code(), "message", st.Message())
			// Map gRPC status codes to potential coreUseCase errors if needed
			if st.Code() == codes.NotFound {
				return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "doctor not found or not available in staff service")
			}
			// Add other mappings as necessary
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, fmt.Sprintf("staff service communication error: %s", st.Message()))
		}
		// Handle non-gRPC errors
		a.logger.Error("StaffService.GetDoctorAvailability non-gRPC error", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, fmt.Sprintf("failed to call staff service: %v", err))
	}

	if protoResp == nil || protoResp.AvailableSlots == nil {
		return []AvailableTimeSlot{}, nil // Return empty slice of correct type
	}

	// Map response proto (repeated TimeSlot) back to AvailableTimeSlot type
	entitySlots := make([]AvailableTimeSlot, 0, len(protoResp.AvailableSlots))
	for _, protoSlot := range protoResp.AvailableSlots {
		if protoSlot == nil {
			continue
		}
		// Basic validation on response data
		if protoSlot.StartTime == nil || protoSlot.EndTime == nil || !protoSlot.StartTime.IsValid() || !protoSlot.EndTime.IsValid() || protoSlot.StartTime.AsTime().After(protoSlot.EndTime.AsTime()) {
			a.logger.Warn("Received invalid time slot from StaffService", "start", protoSlot.StartTime, "end", protoSlot.EndTime)
			continue
		}

		// Map to AvailableTimeSlot
		entitySlots = append(entitySlots, AvailableTimeSlot{
			StartTime: protoSlot.StartTime.AsTime(),
			EndTime:   protoSlot.EndTime.AsTime(),
		})
	}

	a.logger.Debug("Received available time slots from StaffService", "count", len(entitySlots))
	return entitySlots, nil
}

// --- End gRPC Client Adapter ---

// Ensure implementation satisfies the interface
var _ AppointmentUseCase = (*appointmentUseCase)(nil)

// appointmentUseCase implements the AppointmentUseCase interface.
type appointmentUseCase struct {
	// Embed BaseUseCaseImpl with Appointment entity and relevant proto DTOs
	*coreUseCase.BaseUseCaseImpl[entity.Appointment, pb.ScheduleAppointmentRequest, pb.RescheduleAppointmentRequest] // Using Reschedule for UpdateDTO example

	appointmentRepo repository.AppointmentRepository // Specific repository interface
	staffClient     StaffServiceClient               // Client for staff service interaction
	logger          coreLogger.Logger
}

// NewAppointmentUseCase creates a new instance.
func NewAppointmentUseCase(
	repo repository.AppointmentRepository,
	staffClient StaffServiceClient,
	logger coreLogger.Logger,
) AppointmentUseCase {
	baseUseCase := coreUseCase.NewBaseUseCase[entity.Appointment, pb.ScheduleAppointmentRequest, pb.RescheduleAppointmentRequest](
		repo.(coreRepository.BaseRepository[entity.Appointment]),
		logger,
	)
	return &appointmentUseCase{
		BaseUseCaseImpl: baseUseCase,
		appointmentRepo: repo,
		staffClient:     staffClient,
		logger:          logger,
	}
}

// CheckDoctorAvailability checks if a doctor is available.
func (uc *appointmentUseCase) CheckDoctorAvailability(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) (bool, error) {
	uc.logger.Debug("Checking doctor availability", "doctorID", doctorID, "start", startTime, "end", endTime)
	if doctorID == uuid.Nil {
		return false, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid doctor ID")
	}
	if startTime.IsZero() || endTime.IsZero() || endTime.Before(startTime) {
		return false, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid time range")
	}

	// 1. Check Staff Service for general availability slots
	availableTimeSlots, err := uc.staffClient.GetDoctorAvailability(ctx, doctorID, startTime, endTime)
	if err != nil {
		// Check if the error is a UseCaseError indicating NotFound
		var ucErr *coreUseCase.UseCaseError
		if errors.As(err, &ucErr) && ucErr.Type == coreUseCase.ErrNotFound {
			uc.logger.Info("Doctor not found or unavailable via StaffService", "doctorID", doctorID)
			return false, nil // Doctor doesn't exist or has no availability reported by staff service
		}
		// For other errors from the client, log and return them
		uc.logger.Error("Failed to check doctor availability via staff service client", "doctorID", doctorID, "error", err)
		return false, err // Don't wrap internal error here, let the original propagate
	}

	if len(availableTimeSlots) == 0 {
		uc.logger.Info("Doctor has no available time slots reported by StaffService client", "doctorID", doctorID)
		return false, nil // Staff service reports no availability in this window
	}

	// 2. Check if the *requested* time range [startTime, endTime) fits within *any* of the available slots
	isCovered := false
	for _, slot := range availableTimeSlots {
		// Check if requested start is within or at the beginning of the slot AND requested end is within or at the end of the slot
		if !startTime.Before(slot.StartTime) && endTime.Before(slot.EndTime.Add(time.Nanosecond)) { // Use Add(time.Nanosecond) for inclusive end
			isCovered = true
			break
		}
	}
	if !isCovered {
		uc.logger.Info("Requested time slot does not fit within any available slot from StaffService", "doctorID", doctorID)
		return false, nil // Requested time doesn't fit within the general availability blocks
	}

	// 3. Double-check against appointments already booked in *this* service's repository.
	// This confirms the specific requested slot is free, even if the broader window was available.
	isLocallyFree, err := uc.appointmentRepo.CheckDoctorAvailability(ctx, doctorID, startTime, endTime)
	if err != nil {
		uc.logger.Error("Failed to check appointment repository for conflicts", "doctorID", doctorID, "error", err)
		return false, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to confirm appointment availability locally")
	}
	if !isLocallyFree {
		uc.logger.Info("Doctor has conflicting appointments in local repository for the specific slot", "doctorID", doctorID)
	}

	return isLocallyFree, nil
}

// ScheduleAppointment creates a new appointment.
func (uc *appointmentUseCase) ScheduleAppointment(ctx context.Context, req *pb.ScheduleAppointmentRequest) (*entity.Appointment, error) {
	uc.logger.Info("Scheduling appointment", "patientID", req.PatientId, "doctorID", req.DoctorId, "place", req.Place)
	// Validate IDs
	patientID, errP := uuid.Parse(req.PatientId)
	doctorID, errD := uuid.Parse(req.DoctorId)
	if errP != nil || errD != nil || req.Reason == "" || req.AppointmentTime == nil || !req.AppointmentTime.IsValid() || req.Duration == nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "missing or invalid required appointment information")
	}
	appointmentTime := req.AppointmentTime.AsTime()
	duration := req.Duration.AsDuration()
	if appointmentTime.IsZero() || duration <= 0 {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid appointment time or duration")
	}
	endTime := appointmentTime.Add(duration)

	// 1. Check Doctor Availability
	available, err := uc.CheckDoctorAvailability(ctx, doctorID, appointmentTime, endTime)
	if err != nil {
		// Error already logged and wrapped by CheckDoctorAvailability
		return nil, err
	}
	if !available {
		uc.logger.Warn("Doctor not available for requested slot", "doctorID", req.DoctorId, "time", appointmentTime)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrConflict, "doctor is not available at the requested time")
	}

	// 2. Create Appointment Entity - Pass Place
	appointment := entity.NewAppointment(patientID, doctorID, req.Reason, req.Place, appointmentTime, duration)

	// 3. Save Appointment locally using embedded base repo's CREATE method
	err = uc.BaseUseCaseImpl.Repository.Create(ctx, appointment)
	if err != nil {
		uc.logger.Error("Failed to create appointment", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to schedule appointment")
	}

	// 4. Optional: Notify Staff Service (Saga/Outbox pattern recommended)
	uc.logger.Info("Appointment scheduled successfully", "appointmentID", appointment.ID.String())
	return appointment, nil
}

// GetAppointmentDetails retrieves appointment details.
func (uc *appointmentUseCase) GetAppointmentDetails(ctx context.Context, appointmentID uuid.UUID) (*entity.Appointment, error) {
	uc.logger.Info("Getting appointment details", "appointmentID", appointmentID.String())
	if appointmentID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid appointment ID")
	}
	// Use embedded base GetByID which uses repo FindByID
	appointment, err := uc.BaseUseCaseImpl.Repository.FindByID(ctx, appointmentID)
	if err != nil {
		if errors.Is(err, errors.New("entity not found")) {
			uc.logger.Warn("Appointment not found", "appointmentID", appointmentID.String())
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "appointment not found")
		}
		uc.logger.Error("Failed to get appointment details", "appointmentID", appointmentID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to retrieve appointment details")
	}
	return appointment, nil
}

// UpdateAppointmentStatus updates the status.
func (uc *appointmentUseCase) UpdateAppointmentStatus(ctx context.Context, appointmentID uuid.UUID, req *pb.UpdateAppointmentStatusRequest) (*entity.Appointment, error) {
	uc.logger.Info("Updating appointment status", "appointmentID", appointmentID.String(), "newStatus", req.Status)
	if appointmentID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid appointment ID")
	}
	// Map and validate proto status
	newStatus, err := mapProtoToEntityStatus(req.Status)
	if err != nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, err.Error())
	}

	// Fetch appointment using base repo method
	apt, err := uc.BaseUseCaseImpl.Repository.FindByID(ctx, appointmentID)
	if err != nil {
		if errors.Is(err, errors.New("entity not found")) {
			uc.logger.Warn("Appointment not found for status update", "appointmentID", appointmentID.String())
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "appointment not found")
		}
		uc.logger.Error("Failed to find appointment for status update", "appointmentID", appointmentID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to find appointment")
	}

	// Apply status change using entity method and handle error
	if err := apt.SetStatus(newStatus); err != nil {
		// Handle potential invalid status transition from entity logic
		uc.logger.Warn("Invalid status transition attempted", "appointmentID", appointmentID.String(), "from", apt.Status, "to", newStatus, "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrConflict, err.Error())
	}

	// Save using base repo method's UPDATE
	err = uc.BaseUseCaseImpl.Repository.Update(ctx, apt)
	if err != nil {
		uc.logger.Error("Failed to update appointment status", "appointmentID", appointmentID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to update appointment status")
	}

	// Optional: Notify Staff Service if cancelled

	uc.logger.Info("Appointment status updated successfully", "appointmentID", appointmentID.String())
	return apt, nil
}

// RescheduleAppointment changes the time/duration/place.
func (uc *appointmentUseCase) RescheduleAppointment(ctx context.Context, appointmentID uuid.UUID, req *pb.RescheduleAppointmentRequest) (*entity.Appointment, error) {
	uc.logger.Info("Rescheduling appointment", "appointmentID", appointmentID.String(), "newPlace", req.Place)
	if appointmentID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid appointment ID")
	}
	if req.NewTime == nil || !req.NewTime.IsValid() {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "new appointment time cannot be empty or invalid")
	}
	newTime := req.NewTime.AsTime()
	if newTime.IsZero() || newTime.Before(time.Now()) { // Prevent rescheduling to the past
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "cannot reschedule appointment to the past or zero time")
	}

	// Fetch existing appointment using base repo method
	apt, err := uc.BaseUseCaseImpl.Repository.FindByID(ctx, appointmentID)
	if err != nil {
		if errors.Is(err, errors.New("entity not found")) {
			uc.logger.Warn("Appointment not found for reschedule", "appointmentID", appointmentID.String())
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "appointment not found")
		}
		uc.logger.Error("Failed to find appointment for reschedule", "appointmentID", appointmentID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to find appointment for reschedule")
	}

	// Duration handling
	currentDuration := apt.Duration
	var newDurationPtr *time.Duration
	if req.NewDuration != nil {
		d := req.NewDuration.AsDuration()
		if d <= 0 {
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid new duration provided")
		}
		currentDuration = d
		newDurationPtr = &d
	}
	newEndTime := newTime.Add(currentDuration)

	// Place handling (optional update)
	var newPlacePtr *string
	if req.Place != "" { // Check if place is provided in request
		newPlace := req.Place
		newPlacePtr = &newPlace
	}

	// Check availability for the new time slot
	available, err := uc.CheckDoctorAvailability(ctx, apt.DoctorID, newTime, newEndTime)
	if err != nil {
		uc.logger.Error("Failed availability check during reschedule", "appointmentID", appointmentID.String(), "error", err)
		// Don't wrap the error again if it's already a UseCaseError
		if _, ok := err.(*coreUseCase.UseCaseError); ok {
			return nil, err
		}
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to check doctor availability for reschedule")
	}
	if !available {
		uc.logger.Warn("Doctor not available for requested reschedule slot", "appointmentID", appointmentID.String(), "doctorID", apt.DoctorID, "newTime", newTime)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrConflict, "doctor is not available at the requested new time")
	}

	// Apply the changes using the entity method, including place
	if err := apt.Reschedule(newTime, newDurationPtr, newPlacePtr); err != nil {
		uc.logger.Warn("Failed to apply reschedule to entity", "appointmentID", appointmentID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrConflict, err.Error())
	}

	// Save changes using base repo Update
	err = uc.BaseUseCaseImpl.Repository.Update(ctx, apt)
	if err != nil {
		uc.logger.Error("Failed to update appointment after reschedule", "appointmentID", appointmentID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to save rescheduled appointment")
	}

	uc.logger.Info("Appointment rescheduled successfully", "appointmentID", apt.ID.String(), "newTime", newTime)
	return apt, nil
}

// CancelAppointment cancels an appointment.
func (uc *appointmentUseCase) CancelAppointment(ctx context.Context, appointmentID uuid.UUID) error {
	// Reuse UpdateAppointmentStatus logic
	// We need to construct the appropriate request proto
	// Note: UpdateAppointmentStatus now accepts the proto request directly
	updateReq := &pb.UpdateAppointmentStatusRequest{
		AppointmentId: appointmentID.String(), // This needs conversion if Update takes UUID
		Status:        pb.AppointmentStatus_CANCELLED,
	}
	// Assuming UpdateAppointmentStatus takes (ctx, uuid, *protoReq)
	_, err := uc.UpdateAppointmentStatus(ctx, appointmentID, updateReq)
	if err != nil {
		// Map internal errors to potentially simpler user-facing ones if needed
		if errors.Is(err, errors.New("appointment not found")) {
			return errors.New("appointment not found")
		}
		if errors.Is(err, errors.New("cannot update status of a completed or cancelled appointment")) {
			return errors.New("appointment already completed or cancelled")
		}
		return errors.New("failed to cancel appointment") // Generic fallback
	}
	return nil
}

// GetAppointmentsForPatient retrieves appointments by patient ID.
func (uc *appointmentUseCase) GetAppointmentsForPatient(ctx context.Context, patientID uuid.UUID) ([]*entity.Appointment, error) {
	if patientID == uuid.Nil {
		return nil, errors.New("invalid patient ID")
	}
	apts, err := uc.appointmentRepo.FindByPatientID(ctx, patientID)
	if err != nil {
		// uc.logger.Error("Failed to get appointments for patient", "patientID", patientID, "error", err)
		return nil, errors.New("failed to retrieve patient appointments")
	}
	return apts, nil
}

// GetAppointmentsForDoctor retrieves appointments by doctor ID and time range.
func (uc *appointmentUseCase) GetAppointmentsForDoctor(ctx context.Context, doctorID uuid.UUID, startTime time.Time, endTime time.Time) ([]*entity.Appointment, error) {
	if doctorID == uuid.Nil {
		return nil, errors.New("invalid doctor ID")
	}
	if startTime.IsZero() || endTime.IsZero() || endTime.Before(startTime) {
		return nil, errors.New("invalid time range")
	}
	apts, err := uc.appointmentRepo.FindByDoctorID(ctx, doctorID, startTime, endTime)
	if err != nil {
		// uc.logger.Error("Failed to get appointments for doctor", "doctorID", doctorID, "error", err)
		return nil, errors.New("failed to retrieve doctor appointments")
	}
	return apts, nil
}

// --- Helper Functions ---

// mapProtoToEntityStatus converts proto enum to entity enum
func mapProtoToEntityStatus(protoStatus pb.AppointmentStatus) (entity.AppointmentStatus, error) {
	switch protoStatus {
	case pb.AppointmentStatus_SCHEDULED:
		return entity.Scheduled, nil
	case pb.AppointmentStatus_CONFIRMED:
		return entity.Confirmed, nil
	case pb.AppointmentStatus_CANCELLED:
		return entity.Cancelled, nil
	case pb.AppointmentStatus_COMPLETED:
		return entity.Completed, nil
	case pb.AppointmentStatus_NO_SHOW:
		return entity.NoShow, nil
	default:
		return "", errors.New("invalid or unspecified appointment status")
	}
}
