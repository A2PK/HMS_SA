package controller

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "golang-microservices-boilerplate/proto/patient-service"
	"golang-microservices-boilerplate/services/patient-service/internal/entity"
)

// Mapper defines the interface for mapping between gRPC proto messages and internal types.
type Mapper interface {
	EntityToProto(patient *entity.Patient) (*pb.Patient, error)
	PatientsToProto(patients []*entity.Patient) ([]*pb.Patient, error)
	MedicalRecordsToProto(records []entity.MedicalRecord) ([]*pb.MedicalRecord, error)
	MedicalRecordToProto(record *entity.MedicalRecord) (*pb.MedicalRecord, error)
}

// Ensure PatientMapper implements Mapper interface.
var _ Mapper = (*PatientMapper)(nil)

// PatientMapper handles mapping between gRPC proto messages and internal types.
type PatientMapper struct{}

// NewPatientMapper creates a new instance of PatientMapper.
func NewPatientMapper() *PatientMapper { // Return concrete type for instantiation
	return &PatientMapper{}
}

// EntityToProto converts an entity.Patient to a proto.Patient.
func (m *PatientMapper) EntityToProto(patient *entity.Patient) (*pb.Patient, error) {
	if patient == nil {
		return nil, errors.New("cannot map nil patient entity to proto")
	}

	medicalHistoryProto, err := m.MedicalRecordsToProto(patient.MedicalHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to map medical history: %w", err)
	}

	return &pb.Patient{
		Id:             patient.ID.String(),
		FirstName:      patient.FirstName,
		LastName:       patient.LastName,
		DateOfBirth:    timestamppb.New(patient.DateOfBirth),
		Gender:         patient.Gender,
		PhoneNumber:    patient.PhoneNumber,
		Address:        patient.Address,
		MedicalHistory: medicalHistoryProto,
		CreatedAt:      timestamppb.New(patient.CreatedAt),
		UpdatedAt:      timestamppb.New(patient.UpdatedAt),
	}, nil
}

// MedicalRecordToProto converts an entity.MedicalRecord to a proto.MedicalRecord.
func (m *PatientMapper) MedicalRecordToProto(record *entity.MedicalRecord) (*pb.MedicalRecord, error) {
	if record == nil {
		return nil, errors.New("cannot map nil medical record entity to proto")
	}
	return &pb.MedicalRecord{
		Id:        record.ID.String(),
		PatientId: record.PatientID.String(),
		Date:      timestamppb.New(record.Date),
		StaffId:   record.StaffID.String(),
		Diagnosis: record.Diagnosis,
		Treatment: record.Treatment,
		Notes:     record.Notes,
		CreatedAt: timestamppb.New(record.CreatedAt),
		UpdatedAt: timestamppb.New(record.UpdatedAt),
	}, nil
}

// MedicalRecordsToProto converts a slice of entity.MedicalRecord to a slice of proto.MedicalRecord.
func (m *PatientMapper) MedicalRecordsToProto(records []entity.MedicalRecord) ([]*pb.MedicalRecord, error) {
	protos := make([]*pb.MedicalRecord, 0, len(records))
	for i := range records { // Iterate by index to pass pointer of element
		proto, err := m.MedicalRecordToProto(&records[i])
		if err != nil {
			// Log or handle individual mapping errors? For now, skip.
			continue
		}
		protos = append(protos, proto)
	}
	return protos, nil
}

// PatientsToProto converts a slice of *entity.Patient to a slice of *pb.Patient.
func (m *PatientMapper) PatientsToProto(patients []*entity.Patient) ([]*pb.Patient, error) {
	protos := make([]*pb.Patient, 0, len(patients))
	for _, patientEntity := range patients {
		proto, err := m.EntityToProto(patientEntity) // Reuse single entity mapping
		if err != nil {
			// Log or handle individual mapping errors? For now, skip.
			// You might want to return an error here or collect errors.
			fmt.Printf("Warning: skipping patient mapping due to error: %v\n", err) // Simple logging
			continue
		}
		protos = append(protos, proto)
	}
	// Currently, errors during individual mapping are only logged/skipped.
	// Return nil error unless a different strategy is needed.
	return protos, nil
}
