// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: proto/patient-service/patient.proto

package patient_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Patient struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // UUID as string
	FirstName      string                 `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName       string                 `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	DateOfBirth    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=date_of_birth,json=dateOfBirth,proto3" json:"date_of_birth,omitempty"`
	Gender         string                 `protobuf:"bytes,5,opt,name=gender,proto3" json:"gender,omitempty"`
	PhoneNumber    string                 `protobuf:"bytes,6,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Address        string                 `protobuf:"bytes,7,opt,name=address,proto3" json:"address,omitempty"`
	MedicalHistory []*MedicalRecord       `protobuf:"bytes,8,rep,name=medical_history,json=medicalHistory,proto3" json:"medical_history,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt      *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Patient) Reset() {
	*x = Patient{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Patient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Patient) ProtoMessage() {}

func (x *Patient) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Patient.ProtoReflect.Descriptor instead.
func (*Patient) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{0}
}

func (x *Patient) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Patient) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *Patient) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *Patient) GetDateOfBirth() *timestamppb.Timestamp {
	if x != nil {
		return x.DateOfBirth
	}
	return nil
}

func (x *Patient) GetGender() string {
	if x != nil {
		return x.Gender
	}
	return ""
}

func (x *Patient) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *Patient) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Patient) GetMedicalHistory() []*MedicalRecord {
	if x != nil {
		return x.MedicalHistory
	}
	return nil
}

func (x *Patient) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Patient) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type MedicalRecord struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                                // UUID as string
	PatientId     string                 `protobuf:"bytes,2,opt,name=patient_id,json=patientId,proto3" json:"patient_id,omitempty"` // UUID as string
	Date          *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=date,proto3" json:"date,omitempty"`
	StaffId       string                 `protobuf:"bytes,4,opt,name=staff_id,json=staffId,proto3" json:"staff_id,omitempty"`
	Diagnosis     string                 `protobuf:"bytes,5,opt,name=diagnosis,proto3" json:"diagnosis,omitempty"`
	Treatment     string                 `protobuf:"bytes,6,opt,name=treatment,proto3" json:"treatment,omitempty"`
	Notes         string                 `protobuf:"bytes,7,opt,name=notes,proto3" json:"notes,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MedicalRecord) Reset() {
	*x = MedicalRecord{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MedicalRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MedicalRecord) ProtoMessage() {}

func (x *MedicalRecord) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MedicalRecord.ProtoReflect.Descriptor instead.
func (*MedicalRecord) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{1}
}

func (x *MedicalRecord) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MedicalRecord) GetPatientId() string {
	if x != nil {
		return x.PatientId
	}
	return ""
}

func (x *MedicalRecord) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *MedicalRecord) GetStaffId() string {
	if x != nil {
		return x.StaffId
	}
	return ""
}

func (x *MedicalRecord) GetDiagnosis() string {
	if x != nil {
		return x.Diagnosis
	}
	return ""
}

func (x *MedicalRecord) GetTreatment() string {
	if x != nil {
		return x.Treatment
	}
	return ""
}

func (x *MedicalRecord) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

func (x *MedicalRecord) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *MedicalRecord) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// Request for RegisterPatient
type RegisterPatientRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FirstName     string                 `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName      string                 `protobuf:"bytes,2,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Gender        string                 `protobuf:"bytes,3,opt,name=gender,proto3" json:"gender,omitempty"`
	PhoneNumber   string                 `protobuf:"bytes,4,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Address       string                 `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	DateOfBirth   *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=date_of_birth,json=dateOfBirth,proto3" json:"date_of_birth,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterPatientRequest) Reset() {
	*x = RegisterPatientRequest{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterPatientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterPatientRequest) ProtoMessage() {}

func (x *RegisterPatientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterPatientRequest.ProtoReflect.Descriptor instead.
func (*RegisterPatientRequest) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterPatientRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *RegisterPatientRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *RegisterPatientRequest) GetGender() string {
	if x != nil {
		return x.Gender
	}
	return ""
}

func (x *RegisterPatientRequest) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *RegisterPatientRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *RegisterPatientRequest) GetDateOfBirth() *timestamppb.Timestamp {
	if x != nil {
		return x.DateOfBirth
	}
	return nil
}

// Response for RegisterPatient (returns the created patient)
type RegisterPatientResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Patient       *Patient               `protobuf:"bytes,1,opt,name=patient,proto3" json:"patient,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterPatientResponse) Reset() {
	*x = RegisterPatientResponse{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterPatientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterPatientResponse) ProtoMessage() {}

func (x *RegisterPatientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterPatientResponse.ProtoReflect.Descriptor instead.
func (*RegisterPatientResponse) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterPatientResponse) GetPatient() *Patient {
	if x != nil {
		return x.Patient
	}
	return nil
}

// Request for GetPatientDetails
type GetPatientDetailsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PatientId     string                 `protobuf:"bytes,1,opt,name=patient_id,json=patientId,proto3" json:"patient_id,omitempty"` // UUID as string
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPatientDetailsRequest) Reset() {
	*x = GetPatientDetailsRequest{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPatientDetailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPatientDetailsRequest) ProtoMessage() {}

func (x *GetPatientDetailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPatientDetailsRequest.ProtoReflect.Descriptor instead.
func (*GetPatientDetailsRequest) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{4}
}

func (x *GetPatientDetailsRequest) GetPatientId() string {
	if x != nil {
		return x.PatientId
	}
	return ""
}

// Response for GetPatientDetails
type GetPatientDetailsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Patient       *Patient               `protobuf:"bytes,1,opt,name=patient,proto3" json:"patient,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPatientDetailsResponse) Reset() {
	*x = GetPatientDetailsResponse{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPatientDetailsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPatientDetailsResponse) ProtoMessage() {}

func (x *GetPatientDetailsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPatientDetailsResponse.ProtoReflect.Descriptor instead.
func (*GetPatientDetailsResponse) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{5}
}

func (x *GetPatientDetailsResponse) GetPatient() *Patient {
	if x != nil {
		return x.Patient
	}
	return nil
}

// Request for UpdatePatientDetails
type UpdatePatientDetailsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PatientId     string                 `protobuf:"bytes,1,opt,name=patient_id,json=patientId,proto3" json:"patient_id,omitempty"` // UUID as string
	FirstName     string                 `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"` // Use optional or field masks for partial updates
	LastName      string                 `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Gender        string                 `protobuf:"bytes,4,opt,name=gender,proto3" json:"gender,omitempty"`
	PhoneNumber   string                 `protobuf:"bytes,5,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Address       string                 `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
	DateOfBirth   *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=date_of_birth,json=dateOfBirth,proto3" json:"date_of_birth,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdatePatientDetailsRequest) Reset() {
	*x = UpdatePatientDetailsRequest{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePatientDetailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePatientDetailsRequest) ProtoMessage() {}

func (x *UpdatePatientDetailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePatientDetailsRequest.ProtoReflect.Descriptor instead.
func (*UpdatePatientDetailsRequest) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{6}
}

func (x *UpdatePatientDetailsRequest) GetPatientId() string {
	if x != nil {
		return x.PatientId
	}
	return ""
}

func (x *UpdatePatientDetailsRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *UpdatePatientDetailsRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *UpdatePatientDetailsRequest) GetGender() string {
	if x != nil {
		return x.Gender
	}
	return ""
}

func (x *UpdatePatientDetailsRequest) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *UpdatePatientDetailsRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *UpdatePatientDetailsRequest) GetDateOfBirth() *timestamppb.Timestamp {
	if x != nil {
		return x.DateOfBirth
	}
	return nil
}

// Response for UpdatePatientDetails
type UpdatePatientDetailsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Patient       *Patient               `protobuf:"bytes,1,opt,name=patient,proto3" json:"patient,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdatePatientDetailsResponse) Reset() {
	*x = UpdatePatientDetailsResponse{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePatientDetailsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePatientDetailsResponse) ProtoMessage() {}

func (x *UpdatePatientDetailsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePatientDetailsResponse.ProtoReflect.Descriptor instead.
func (*UpdatePatientDetailsResponse) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{7}
}

func (x *UpdatePatientDetailsResponse) GetPatient() *Patient {
	if x != nil {
		return x.Patient
	}
	return nil
}

// Request for AddMedicalRecord
type AddMedicalRecordRequest struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	PatientId string                 `protobuf:"bytes,1,opt,name=patient_id,json=patientId,proto3" json:"patient_id,omitempty"` // UUID as string
	// Pass the data needed to create a new medical record
	Date          *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	StaffId       string                 `protobuf:"bytes,3,opt,name=staff_id,json=staffId,proto3" json:"staff_id,omitempty"`
	Diagnosis     string                 `protobuf:"bytes,4,opt,name=diagnosis,proto3" json:"diagnosis,omitempty"`
	Treatment     string                 `protobuf:"bytes,5,opt,name=treatment,proto3" json:"treatment,omitempty"`
	Notes         string                 `protobuf:"bytes,6,opt,name=notes,proto3" json:"notes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddMedicalRecordRequest) Reset() {
	*x = AddMedicalRecordRequest{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddMedicalRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMedicalRecordRequest) ProtoMessage() {}

func (x *AddMedicalRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMedicalRecordRequest.ProtoReflect.Descriptor instead.
func (*AddMedicalRecordRequest) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{8}
}

func (x *AddMedicalRecordRequest) GetPatientId() string {
	if x != nil {
		return x.PatientId
	}
	return ""
}

func (x *AddMedicalRecordRequest) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *AddMedicalRecordRequest) GetStaffId() string {
	if x != nil {
		return x.StaffId
	}
	return ""
}

func (x *AddMedicalRecordRequest) GetDiagnosis() string {
	if x != nil {
		return x.Diagnosis
	}
	return ""
}

func (x *AddMedicalRecordRequest) GetTreatment() string {
	if x != nil {
		return x.Treatment
	}
	return ""
}

func (x *AddMedicalRecordRequest) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

// Request for GetPatientMedicalHistory
type GetPatientMedicalHistoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PatientId     string                 `protobuf:"bytes,1,opt,name=patient_id,json=patientId,proto3" json:"patient_id,omitempty"` // UUID as string
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPatientMedicalHistoryRequest) Reset() {
	*x = GetPatientMedicalHistoryRequest{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPatientMedicalHistoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPatientMedicalHistoryRequest) ProtoMessage() {}

func (x *GetPatientMedicalHistoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPatientMedicalHistoryRequest.ProtoReflect.Descriptor instead.
func (*GetPatientMedicalHistoryRequest) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{9}
}

func (x *GetPatientMedicalHistoryRequest) GetPatientId() string {
	if x != nil {
		return x.PatientId
	}
	return ""
}

// Response for GetPatientMedicalHistory
type GetPatientMedicalHistoryResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	MedicalHistory []*MedicalRecord       `protobuf:"bytes,1,rep,name=medical_history,json=medicalHistory,proto3" json:"medical_history,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *GetPatientMedicalHistoryResponse) Reset() {
	*x = GetPatientMedicalHistoryResponse{}
	mi := &file_proto_patient_service_patient_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPatientMedicalHistoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPatientMedicalHistoryResponse) ProtoMessage() {}

func (x *GetPatientMedicalHistoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_patient_service_patient_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPatientMedicalHistoryResponse.ProtoReflect.Descriptor instead.
func (*GetPatientMedicalHistoryResponse) Descriptor() ([]byte, []int) {
	return file_proto_patient_service_patient_proto_rawDescGZIP(), []int{10}
}

func (x *GetPatientMedicalHistoryResponse) GetMedicalHistory() []*MedicalRecord {
	if x != nil {
		return x.MedicalHistory
	}
	return nil
}

var File_proto_patient_service_patient_proto protoreflect.FileDescriptor

const file_proto_patient_service_patient_proto_rawDesc = "" +
	"\n" +
	"#proto/patient-service/patient.proto\x12\x0epatientservice\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1bgoogle/protobuf/empty.proto\"\xa8\x03\n" +
	"\aPatient\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1d\n" +
	"\n" +
	"first_name\x18\x02 \x01(\tR\tfirstName\x12\x1b\n" +
	"\tlast_name\x18\x03 \x01(\tR\blastName\x12>\n" +
	"\rdate_of_birth\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\vdateOfBirth\x12\x16\n" +
	"\x06gender\x18\x05 \x01(\tR\x06gender\x12!\n" +
	"\fphone_number\x18\x06 \x01(\tR\vphoneNumber\x12\x18\n" +
	"\aaddress\x18\a \x01(\tR\aaddress\x12F\n" +
	"\x0fmedical_history\x18\b \x03(\v2\x1d.patientservice.MedicalRecordR\x0emedicalHistory\x129\n" +
	"\n" +
	"created_at\x18\t \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\n" +
	" \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"\xd1\x02\n" +
	"\rMedicalRecord\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1d\n" +
	"\n" +
	"patient_id\x18\x02 \x01(\tR\tpatientId\x12.\n" +
	"\x04date\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x04date\x12\x19\n" +
	"\bstaff_id\x18\x04 \x01(\tR\astaffId\x12\x1c\n" +
	"\tdiagnosis\x18\x05 \x01(\tR\tdiagnosis\x12\x1c\n" +
	"\ttreatment\x18\x06 \x01(\tR\ttreatment\x12\x14\n" +
	"\x05notes\x18\a \x01(\tR\x05notes\x129\n" +
	"\n" +
	"created_at\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\t \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"\xe9\x01\n" +
	"\x16RegisterPatientRequest\x12\x1d\n" +
	"\n" +
	"first_name\x18\x01 \x01(\tR\tfirstName\x12\x1b\n" +
	"\tlast_name\x18\x02 \x01(\tR\blastName\x12\x16\n" +
	"\x06gender\x18\x03 \x01(\tR\x06gender\x12!\n" +
	"\fphone_number\x18\x04 \x01(\tR\vphoneNumber\x12\x18\n" +
	"\aaddress\x18\x05 \x01(\tR\aaddress\x12>\n" +
	"\rdate_of_birth\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\vdateOfBirth\"L\n" +
	"\x17RegisterPatientResponse\x121\n" +
	"\apatient\x18\x01 \x01(\v2\x17.patientservice.PatientR\apatient\"9\n" +
	"\x18GetPatientDetailsRequest\x12\x1d\n" +
	"\n" +
	"patient_id\x18\x01 \x01(\tR\tpatientId\"N\n" +
	"\x19GetPatientDetailsResponse\x121\n" +
	"\apatient\x18\x01 \x01(\v2\x17.patientservice.PatientR\apatient\"\x8d\x02\n" +
	"\x1bUpdatePatientDetailsRequest\x12\x1d\n" +
	"\n" +
	"patient_id\x18\x01 \x01(\tR\tpatientId\x12\x1d\n" +
	"\n" +
	"first_name\x18\x02 \x01(\tR\tfirstName\x12\x1b\n" +
	"\tlast_name\x18\x03 \x01(\tR\blastName\x12\x16\n" +
	"\x06gender\x18\x04 \x01(\tR\x06gender\x12!\n" +
	"\fphone_number\x18\x05 \x01(\tR\vphoneNumber\x12\x18\n" +
	"\aaddress\x18\x06 \x01(\tR\aaddress\x12>\n" +
	"\rdate_of_birth\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\vdateOfBirth\"Q\n" +
	"\x1cUpdatePatientDetailsResponse\x121\n" +
	"\apatient\x18\x01 \x01(\v2\x17.patientservice.PatientR\apatient\"\xd5\x01\n" +
	"\x17AddMedicalRecordRequest\x12\x1d\n" +
	"\n" +
	"patient_id\x18\x01 \x01(\tR\tpatientId\x12.\n" +
	"\x04date\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x04date\x12\x19\n" +
	"\bstaff_id\x18\x03 \x01(\tR\astaffId\x12\x1c\n" +
	"\tdiagnosis\x18\x04 \x01(\tR\tdiagnosis\x12\x1c\n" +
	"\ttreatment\x18\x05 \x01(\tR\ttreatment\x12\x14\n" +
	"\x05notes\x18\x06 \x01(\tR\x05notes\"@\n" +
	"\x1fGetPatientMedicalHistoryRequest\x12\x1d\n" +
	"\n" +
	"patient_id\x18\x01 \x01(\tR\tpatientId\"j\n" +
	" GetPatientMedicalHistoryResponse\x12F\n" +
	"\x0fmedical_history\x18\x01 \x03(\v2\x1d.patientservice.MedicalRecordR\x0emedicalHistory2\xa5\x04\n" +
	"\x0ePatientService\x12b\n" +
	"\x0fRegisterPatient\x12&.patientservice.RegisterPatientRequest\x1a'.patientservice.RegisterPatientResponse\x12h\n" +
	"\x11GetPatientDetails\x12(.patientservice.GetPatientDetailsRequest\x1a).patientservice.GetPatientDetailsResponse\x12q\n" +
	"\x14UpdatePatientDetails\x12+.patientservice.UpdatePatientDetailsRequest\x1a,.patientservice.UpdatePatientDetailsResponse\x12S\n" +
	"\x10AddMedicalRecord\x12'.patientservice.AddMedicalRecordRequest\x1a\x16.google.protobuf.Empty\x12}\n" +
	"\x18GetPatientMedicalHistory\x12/.patientservice.GetPatientMedicalHistoryRequest\x1a0.patientservice.GetPatientMedicalHistoryResponseB8Z6golang-microservices-boilerplate/proto/patient-serviceb\x06proto3"

var (
	file_proto_patient_service_patient_proto_rawDescOnce sync.Once
	file_proto_patient_service_patient_proto_rawDescData []byte
)

func file_proto_patient_service_patient_proto_rawDescGZIP() []byte {
	file_proto_patient_service_patient_proto_rawDescOnce.Do(func() {
		file_proto_patient_service_patient_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_patient_service_patient_proto_rawDesc), len(file_proto_patient_service_patient_proto_rawDesc)))
	})
	return file_proto_patient_service_patient_proto_rawDescData
}

var file_proto_patient_service_patient_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_patient_service_patient_proto_goTypes = []any{
	(*Patient)(nil),                          // 0: patientservice.Patient
	(*MedicalRecord)(nil),                    // 1: patientservice.MedicalRecord
	(*RegisterPatientRequest)(nil),           // 2: patientservice.RegisterPatientRequest
	(*RegisterPatientResponse)(nil),          // 3: patientservice.RegisterPatientResponse
	(*GetPatientDetailsRequest)(nil),         // 4: patientservice.GetPatientDetailsRequest
	(*GetPatientDetailsResponse)(nil),        // 5: patientservice.GetPatientDetailsResponse
	(*UpdatePatientDetailsRequest)(nil),      // 6: patientservice.UpdatePatientDetailsRequest
	(*UpdatePatientDetailsResponse)(nil),     // 7: patientservice.UpdatePatientDetailsResponse
	(*AddMedicalRecordRequest)(nil),          // 8: patientservice.AddMedicalRecordRequest
	(*GetPatientMedicalHistoryRequest)(nil),  // 9: patientservice.GetPatientMedicalHistoryRequest
	(*GetPatientMedicalHistoryResponse)(nil), // 10: patientservice.GetPatientMedicalHistoryResponse
	(*timestamppb.Timestamp)(nil),            // 11: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),                    // 12: google.protobuf.Empty
}
var file_proto_patient_service_patient_proto_depIdxs = []int32{
	11, // 0: patientservice.Patient.date_of_birth:type_name -> google.protobuf.Timestamp
	1,  // 1: patientservice.Patient.medical_history:type_name -> patientservice.MedicalRecord
	11, // 2: patientservice.Patient.created_at:type_name -> google.protobuf.Timestamp
	11, // 3: patientservice.Patient.updated_at:type_name -> google.protobuf.Timestamp
	11, // 4: patientservice.MedicalRecord.date:type_name -> google.protobuf.Timestamp
	11, // 5: patientservice.MedicalRecord.created_at:type_name -> google.protobuf.Timestamp
	11, // 6: patientservice.MedicalRecord.updated_at:type_name -> google.protobuf.Timestamp
	11, // 7: patientservice.RegisterPatientRequest.date_of_birth:type_name -> google.protobuf.Timestamp
	0,  // 8: patientservice.RegisterPatientResponse.patient:type_name -> patientservice.Patient
	0,  // 9: patientservice.GetPatientDetailsResponse.patient:type_name -> patientservice.Patient
	11, // 10: patientservice.UpdatePatientDetailsRequest.date_of_birth:type_name -> google.protobuf.Timestamp
	0,  // 11: patientservice.UpdatePatientDetailsResponse.patient:type_name -> patientservice.Patient
	11, // 12: patientservice.AddMedicalRecordRequest.date:type_name -> google.protobuf.Timestamp
	1,  // 13: patientservice.GetPatientMedicalHistoryResponse.medical_history:type_name -> patientservice.MedicalRecord
	2,  // 14: patientservice.PatientService.RegisterPatient:input_type -> patientservice.RegisterPatientRequest
	4,  // 15: patientservice.PatientService.GetPatientDetails:input_type -> patientservice.GetPatientDetailsRequest
	6,  // 16: patientservice.PatientService.UpdatePatientDetails:input_type -> patientservice.UpdatePatientDetailsRequest
	8,  // 17: patientservice.PatientService.AddMedicalRecord:input_type -> patientservice.AddMedicalRecordRequest
	9,  // 18: patientservice.PatientService.GetPatientMedicalHistory:input_type -> patientservice.GetPatientMedicalHistoryRequest
	3,  // 19: patientservice.PatientService.RegisterPatient:output_type -> patientservice.RegisterPatientResponse
	5,  // 20: patientservice.PatientService.GetPatientDetails:output_type -> patientservice.GetPatientDetailsResponse
	7,  // 21: patientservice.PatientService.UpdatePatientDetails:output_type -> patientservice.UpdatePatientDetailsResponse
	12, // 22: patientservice.PatientService.AddMedicalRecord:output_type -> google.protobuf.Empty
	10, // 23: patientservice.PatientService.GetPatientMedicalHistory:output_type -> patientservice.GetPatientMedicalHistoryResponse
	19, // [19:24] is the sub-list for method output_type
	14, // [14:19] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_proto_patient_service_patient_proto_init() }
func file_proto_patient_service_patient_proto_init() {
	if File_proto_patient_service_patient_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_patient_service_patient_proto_rawDesc), len(file_proto_patient_service_patient_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_patient_service_patient_proto_goTypes,
		DependencyIndexes: file_proto_patient_service_patient_proto_depIdxs,
		MessageInfos:      file_proto_patient_service_patient_proto_msgTypes,
	}.Build()
	File_proto_patient_service_patient_proto = out.File
	file_proto_patient_service_patient_proto_goTypes = nil
	file_proto_patient_service_patient_proto_depIdxs = nil
}
