package usecase

import (
	"context"

	pb "golang-microservices-boilerplate/proto/patient-service" // Import proto
	"golang-microservices-boilerplate/services/patient-service/internal/entity"

	"github.com/google/uuid"
)

// PatientUseCase defines the interface for patient management operations.
type PatientUseCase interface {
	// RegisterPatient registers a new patient in the system.
	RegisterPatient(ctx context.Context, req *pb.RegisterPatientRequest) (*entity.Patient, error)

	// GetPatientDetails retrieves the details of a specific patient by ID.
	GetPatientDetails(ctx context.Context, patientID uuid.UUID) (*entity.Patient, error)

	// UpdatePatientDetails updates the information for an existing patient.
	UpdatePatientDetails(ctx context.Context, patientID uuid.UUID, req *pb.UpdatePatientDetailsRequest) (*entity.Patient, error)

	// AddMedicalRecord adds a new medical record to a patient's history.
	AddMedicalRecord(ctx context.Context, patientID uuid.UUID, req *pb.AddMedicalRecordRequest) error

	// GetPatientMedicalHistory retrieves the medical history for a specific patient.
	GetPatientMedicalHistory(ctx context.Context, patientID uuid.UUID) ([]entity.MedicalRecord, error)
}
