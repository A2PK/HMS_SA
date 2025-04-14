package usecase

import (
	"context"
	"errors"
	"time"

	pb "golang-microservices-boilerplate/proto/patient-service" // Alias for generated proto types

	coreLogger "golang-microservices-boilerplate/pkg/core/logger"
	coreRepository "golang-microservices-boilerplate/pkg/core/repository" // Added alias for core repo
	coreTypes "golang-microservices-boilerplate/pkg/core/types"           // Import coreTypes
	coreUseCase "golang-microservices-boilerplate/pkg/core/usecase"

	// coreDTO "golang-microservices-boilerplate/pkg/core/dto" // Not used currently

	"golang-microservices-boilerplate/services/patient-service/internal/entity"
	"golang-microservices-boilerplate/services/patient-service/internal/repository"

	"github.com/google/uuid"
)

// Ensure implementation satisfies the specific interface AND the generated gRPC server interface (implicitly)
var _ PatientUseCase = (*patientUseCase)(nil)

// patientUseCase implements the PatientUseCase interface.
// It embeds BaseUseCaseImpl for standard CRUD operations if applicable,
// but Patient management often requires custom logic beyond simple DTO mapping.
type patientUseCase struct {
	// We embed the core use case, but might not use its generic Create/Update directly
	// if our logic (like RegisterPatient) is significantly different.
	// Let's define the specific Create/Update DTOs as the proto messages.
	*coreUseCase.BaseUseCaseImpl[entity.Patient, pb.RegisterPatientRequest, pb.UpdatePatientDetailsRequest]

	patientRepo repository.PatientRepository // Specific repository interface
	logger      coreLogger.Logger
}

// NewPatientUseCase creates a new instance of patientUseCase.
func NewPatientUseCase(repo repository.PatientRepository, logger coreLogger.Logger) PatientUseCase {
	// Instantiate the embedded base use case
	// Type assertion needs to target the core repository interface
	baseUseCase := coreUseCase.NewBaseUseCase[entity.Patient, pb.RegisterPatientRequest, pb.UpdatePatientDetailsRequest](
		repo.(coreRepository.BaseRepository[entity.Patient]),
		logger,
	)

	return &patientUseCase{
		BaseUseCaseImpl: baseUseCase,
		patientRepo:     repo,
		logger:          logger,
	}
}

// RegisterPatient handles the registration logic, potentially bypassing the generic BaseUseCase.Create.
func (uc *patientUseCase) RegisterPatient(ctx context.Context, req *pb.RegisterPatientRequest) (*entity.Patient, error) {
	uc.logger.Info("Registering new patient", "firstName", req.FirstName, "lastName", req.LastName)

	// Manual validation (or use coreDTO.Validate if applicable to proto messages)
	if req.FirstName == "" || req.LastName == "" || req.PhoneNumber == "" {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "missing required patient information")
	}
	var dob time.Time
	if req.DateOfBirth != nil {
		dob = req.DateOfBirth.AsTime()
	}

	// Create entity using domain constructor
	patient := entity.NewPatient(req.FirstName, req.LastName, req.Gender, req.PhoneNumber, req.Address, dob)

	// Use the Save method from the embedded BaseUseCaseImpl's Repository (Explicit access)
	err := uc.BaseUseCaseImpl.Repository.Create(ctx, patient)
	if err != nil {
		uc.logger.Error("Failed to save patient", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to register patient")
	}

	uc.logger.Info("Patient registered successfully", "patientID", patient.ID.String())
	// The patient object now has the ID assigned by BeforeCreate/Save
	return patient, nil
}

// GetPatientDetails retrieves patient details by ID.
// Can potentially leverage the embedded BaseUseCase.GetByID if no custom logic is needed.
func (uc *patientUseCase) GetPatientDetails(ctx context.Context, patientID uuid.UUID) (*entity.Patient, error) {
	uc.logger.Info("Getting patient details", "patientID", patientID.String())
	if patientID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid patient ID")
	}

	// Use the specific repository's FindByID (which might have custom preloading)
	patient, err := uc.patientRepo.FindByID(ctx, patientID)
	if err != nil {
		if errors.Is(err, errors.New("entity not found")) { // Match repo error
			uc.logger.Warn("Patient not found", "patientID", patientID.String())
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "patient not found")
		}
		uc.logger.Error("Failed to get patient details", "patientID", patientID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to retrieve patient details")
	}
	return patient, nil
}

// UpdatePatientDetails handles update logic.
// Could use BaseUseCase.Update if mapping from pb.UpdatePatientDetailsRequest works directly.
// Here, we implement custom logic using the entity method.
func (uc *patientUseCase) UpdatePatientDetails(ctx context.Context, patientID uuid.UUID, req *pb.UpdatePatientDetailsRequest) (*entity.Patient, error) {
	uc.logger.Info("Updating patient details", "patientID", patientID.String())
	if patientID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid patient ID")
	}

	// Validate req? (e.g., ensure at least one field is being updated)
	// Use coreDTO.Validate(req) if validation tags are added to proto messages?

	// Fetch existing patient
	patient, err := uc.patientRepo.FindByID(ctx, patientID)
	if err != nil {
		if errors.Is(err, errors.New("entity not found")) {
			uc.logger.Warn("Patient not found for update", "patientID", patientID.String())
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "patient not found")
		}
		uc.logger.Error("Failed to find patient for update", "patientID", patientID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to find patient for update")
	}

	// Apply updates using the entity method
	var dob time.Time
	if req.DateOfBirth != nil {
		dob = req.DateOfBirth.AsTime()
	}
	patient.UpdateDetails(req.FirstName, req.LastName, req.Gender, req.PhoneNumber, req.Address, dob)

	// Save the updated entity using the embedded repository's Update (Explicit access)
	err = uc.BaseUseCaseImpl.Repository.Update(ctx, patient)
	if err != nil {
		uc.logger.Error("Failed to update patient", "patientID", patientID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to update patient details")
	}

	uc.logger.Info("Patient details updated successfully", "patientID", patientID.String())
	return patient, nil
}

// AddMedicalRecord adds a medical record.
func (uc *patientUseCase) AddMedicalRecord(ctx context.Context, patientID uuid.UUID, req *pb.AddMedicalRecordRequest) error {
	uc.logger.Info("Adding medical record", "patientID", patientID.String())
	if patientID == uuid.Nil {
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid patient ID")
	}
	// Basic validation for request data
	if req.StaffId == "" || req.Diagnosis == "" || req.Date == nil {
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "missing required medical record information")
	}

	// Parse StaffID string to UUID
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil {
		uc.logger.Error("Invalid StaffID format in AddMedicalRecord request", "staffIdString", req.StaffId, "error", err)
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid staff ID format")
	}

	// Convert proto request data to entity.MedicalRecord
	// Note: MedicalRecord ID is generated by the repository/BeforeCreate hook
	record := &entity.MedicalRecord{
		PatientID: patientID,
		Date:      req.Date.AsTime(),
		StaffID:   staffID, // Use parsed UUID
		Diagnosis: req.Diagnosis,
		Treatment: req.Treatment,
		Notes:     req.Notes,
	}

	// Use the dedicated method from the specific repository interface
	err = uc.patientRepo.AddMedicalRecord(ctx, patientID, record)
	if err != nil {
		uc.logger.Error("Failed to add medical record", "patientID", patientID.String(), "error", err)
		// Check for specific repo errors?
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to add medical record")
	}
	uc.logger.Info("Medical record added successfully", "patientID", patientID.String(), "recordID", record.ID.String())
	return nil
}

// GetPatientMedicalHistory retrieves medical history.
func (uc *patientUseCase) GetPatientMedicalHistory(ctx context.Context, patientID uuid.UUID) ([]entity.MedicalRecord, error) {
	uc.logger.Info("Getting medical history", "patientID", patientID.String())
	if patientID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid patient ID")
	}

	// Use GetPatientDetails which uses repo's FindByID with preloading
	patient, err := uc.GetPatientDetails(ctx, patientID)
	if err != nil {
		// Error already handled and logged in GetPatientDetails
		return nil, err
	}

	if patient.MedicalHistory == nil {
		return []entity.MedicalRecord{}, nil // Return empty slice, not nil
	}
	return patient.MedicalHistory, nil
}

// ListPatients retrieves a list of all patients.
// Returns a slice of pointers to patients.
func (uc *patientUseCase) ListPatients(ctx context.Context) ([]*entity.Patient, error) {
	uc.logger.Info("Listing all patients")

	// Call the embedded repository's FindAll method.
	// Use coreTypes.FilterOptions.
	filterOpts := coreTypes.FilterOptions{} // Use coreTypes
	paginationResult, err := uc.patientRepo.FindAll(ctx, filterOpts)

	if err != nil {
		uc.logger.Error("Failed to list patients", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to retrieve patients")
	}

	// Check if the result or Items field is nil.
	if paginationResult == nil || paginationResult.Items == nil {
		return []*entity.Patient{}, nil // Return empty slice of pointers
	}

	// Return the slice of pointers directly.
	return paginationResult.Items, nil
}
