package repository

import (
	"context"
	coreRepository "golang-microservices-boilerplate/pkg/core/repository"
	"golang-microservices-boilerplate/services/patient-service/internal/entity"

	// Remove unused core imports if not needed directly by this interface's methods
	// "golang-microservices-boilerplate/pkg/core/types"

	"github.com/google/uuid"
)

// PatientRepository defines the interface for patient data persistence.
// It includes methods specific to Patient entities beyond the base repository operations.
// Implementations are expected to embed coreRepository.GormBaseRepository[entity.Patient]
// to provide standard CRUD methods (Save, Update, Delete, FindAll, etc.).
type PatientRepository interface {
	coreRepository.BaseRepository[entity.Patient]

	// AddMedicalRecord adds a medical record to a specific patient.
	AddMedicalRecord(ctx context.Context, patientID uuid.UUID, record *entity.MedicalRecord) error

	// --- Methods provided by embedding GormBaseRepository ---
	// Create(ctx context.Context, patient *entity.Patient) error
	// Update(ctx context.Context, patient *entity.Patient) error
	// Delete(ctx context.Context, id uuid.UUID, hardDelete bool) error
	// FindByID(ctx context.Context, id uuid.UUID) (*entity.Patient, error)
	// FindAll(ctx context.Context, opts coreTypes.FilterOptions) (*coreTypes.PaginationResult[entity.Patient], error)
	// FindWithFilter(ctx context.Context, filter map[string]interface{}, opts coreTypes.FilterOptions) (*coreTypes.PaginationResult[entity.Patient], error)
	// Count(ctx context.Context, filter map[string]interface{}) (int64, error)
	// Transaction(ctx context.Context, fn func(txRepo repository.BaseRepository[entity.Patient]) error) error
	// CreateMany(ctx context.Context, entities []*entity.Patient) error
	// UpdateMany(ctx context.Context, entities []*entity.Patient) error
	// DeleteMany(ctx context.Context, ids []uuid.UUID, hardDelete bool) error
}

// Note: Depending on the chosen database technology (SQL, NoSQL),
// the exact method signatures might vary slightly.
// Error handling should be consistent (e.g., return custom error types
// like ErrPatientNotFound).
