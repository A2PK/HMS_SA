package repository

import (
	"context"
	"errors"

	"golang-microservices-boilerplate/pkg/core/repository"
	"golang-microservices-boilerplate/services/patient-service/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ensure GormPatientRepository implements PatientRepository
var _ PatientRepository = (*GormPatientRepository)(nil)

type GormPatientRepository struct {
	*repository.GormBaseRepository[entity.Patient] // Embed the generic GORM repository
}

// NewGormPatientRepository creates a new GORM patient repository
func NewGormPatientRepository(db *gorm.DB) *GormPatientRepository {
	//nolint:all // Suppress T does not satisfy Entity constraint due to pointer receivers
	baseRepo := repository.NewGormBaseRepository[entity.Patient](db)
	return &GormPatientRepository{
		GormBaseRepository: baseRepo,
	}
}

// FindByID overrides or enhances the base FindByID if needed (e.g., to preload associations)
func (r *GormPatientRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Patient, error) {
	patient := &entity.Patient{}
	// Use Preload to automatically load the MedicalHistory association
	// Use Clauses(clause.Associations) to load all defined associations
	err := r.DB.WithContext(ctx).Preload("MedicalHistory").First(patient, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Use the error message expected by the BaseRepository for consistency
			return nil, errors.New("entity not found")
		}
		return nil, err
	}
	return patient, nil
}

// AddMedicalRecord implements the specific method for adding a medical record.
// This assumes MedicalRecord has its own BaseEntity for ID/timestamps.
func (r *GormPatientRepository) AddMedicalRecord(ctx context.Context, patientID uuid.UUID, record *entity.MedicalRecord) error {
	// Ensure the record has the correct PatientID
	if record.PatientID != patientID {
		// This case might indicate a logic error elsewhere, but we enforce consistency here.
		record.PatientID = patientID
	}

	// Create the medical record. GORM's BeforeCreate hook on MedicalRecord's BaseEntity
	// should handle setting its own UUID.
	// This operation should ideally be within a transaction managed by the use case if multiple steps are involved.
	return r.DB.WithContext(ctx).Create(record).Error
}

// You can override other base methods like FindAll if specific preloading or filtering is always needed for Patients
/*
func (r *GormPatientRepository) FindAll(ctx context.Context, opts coreTypes.FilterOptions) (*coreTypes.PaginationResult[entity.Patient], error) {
	// Custom logic before or after calling the base method
	// Add specific preloads
	r.DB = r.DB.Preload("MedicalHistory")
	return r.GormBaseRepository.FindAll(ctx, opts)
}
*/
