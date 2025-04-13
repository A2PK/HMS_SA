package controller

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "golang-microservices-boilerplate/proto/appointment-service"
	"golang-microservices-boilerplate/services/appointment-service/internal/entity"
)

// Mapper defines the interface for mapping between gRPC proto messages and internal types.
type Mapper interface {
	EntityToProto(apt *entity.Appointment) (*pb.Appointment, error)
	EntitiesToProto(apts []*entity.Appointment) ([]*pb.Appointment, error)
	EntityStatusToProto(status entity.AppointmentStatus) (pb.AppointmentStatus, error)
	ProtoStatusToEntity(status pb.AppointmentStatus) (entity.AppointmentStatus, error)
}

// Ensure AppointmentMapper implements Mapper interface.
var _ Mapper = (*AppointmentMapper)(nil)

// AppointmentMapper handles mapping between gRPC proto messages and internal types.
type AppointmentMapper struct{}

// NewAppointmentMapper creates a new instance of AppointmentMapper.
func NewAppointmentMapper() *AppointmentMapper { // Return concrete type for instantiation
	return &AppointmentMapper{}
}

// EntityToProto converts an entity.Appointment to a proto.Appointment.
func (m *AppointmentMapper) EntityToProto(apt *entity.Appointment) (*pb.Appointment, error) {
	if apt == nil {
		return nil, errors.New("cannot map nil entity to proto")
	}
	statusProto, err := m.EntityStatusToProto(apt.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to map status: %w", err)
	}

	return &pb.Appointment{
		Id:              apt.ID.String(),
		PatientId:       apt.PatientID.String(),
		DoctorId:        apt.DoctorID.String(),
		AppointmentTime: timestamppb.New(apt.AppointmentTime),
		Duration:        durationpb.New(apt.Duration),
		Reason:          apt.Reason,
		Status:          statusProto,
		Notes:           apt.Notes,
		CreatedAt:       timestamppb.New(apt.CreatedAt),
		UpdatedAt:       timestamppb.New(apt.UpdatedAt),
	}, nil
}

// EntitiesToProto converts a slice of entity.Appointment to a slice of proto.Appointment.
func (m *AppointmentMapper) EntitiesToProto(apts []*entity.Appointment) ([]*pb.Appointment, error) {
	protos := make([]*pb.Appointment, 0, len(apts))
	for _, apt := range apts {
		proto, err := m.EntityToProto(apt)
		if err != nil {
			// Log or handle individual mapping errors? For now, skip.
			continue
		}
		protos = append(protos, proto)
	}
	return protos, nil
}

// EntityStatusToProto converts entity.AppointmentStatus to pb.AppointmentStatus.
func (m *AppointmentMapper) EntityStatusToProto(status entity.AppointmentStatus) (pb.AppointmentStatus, error) {
	switch status {
	case entity.Scheduled:
		return pb.AppointmentStatus_SCHEDULED, nil
	case entity.Confirmed:
		return pb.AppointmentStatus_CONFIRMED, nil
	case entity.Cancelled:
		return pb.AppointmentStatus_CANCELLED, nil
	case entity.Completed:
		return pb.AppointmentStatus_COMPLETED, nil
	case entity.NoShow:
		return pb.AppointmentStatus_NO_SHOW, nil
	default:
		return pb.AppointmentStatus_APPOINTMENT_STATUS_UNSPECIFIED, fmt.Errorf("unknown entity status: %s", status)
	}
}

// ProtoStatusToEntity converts pb.AppointmentStatus to entity.AppointmentStatus.
func (m *AppointmentMapper) ProtoStatusToEntity(status pb.AppointmentStatus) (entity.AppointmentStatus, error) {
	switch status {
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
	case pb.AppointmentStatus_APPOINTMENT_STATUS_UNSPECIFIED:
		return "", errors.New("unspecified status cannot be mapped to entity")
	default:
		return "", fmt.Errorf("unknown proto status: %s", status.String())
	}
}
