package entity

import (
	"time"

	coreEntity "golang-microservices-boilerplate/pkg/core/entity"

	"github.com/google/uuid"
)

// --- Lookup Entities (Using Name as PK) ---

// StaffRole represents a role a staff member can have.
type StaffRole struct {
	// Using Name as the primary key
	Name        string `gorm:"primaryKey;uniqueIndex;not null"` // Role name (e.g., "Doctor")
	Description string // Optional description
}

// TableName specifies the table name for StaffRole.
func (StaffRole) TableName() string {
	return "staff_roles"
}

// StaffStatus represents the employment status of a staff member.
type StaffStatus struct {
	// Using Name as the primary key
	Name        string `gorm:"primaryKey;uniqueIndex;not null"` // Status name (e.g., "Active")
	Description string // Optional description
}

// TableName specifies the table name for StaffStatus.
func (StaffStatus) TableName() string {
	return "staff_statuses"
}

// TaskStatus represents the status of a task.
type TaskStatus struct {
	// Using Name as the primary key
	Name        string `gorm:"primaryKey;uniqueIndex;not null"` // Status name (e.g., "Pending")
	Description string // Optional description
}

// TableName specifies the table name for TaskStatus.
func (TaskStatus) TableName() string {
	return "task_statuses"
}

// --- Main Staff Entity ---

// Staff represents a hospital staff member.
type Staff struct {
	coreEntity.BaseEntity
	FirstName   string    `gorm:"not null"`    // Staff member's first name
	LastName    string    `gorm:"not null"`    // Staff member's last name
	DateOfBirth time.Time `gorm:"type:date"`   // Staff member's date of birth
	PhoneNumber string    `gorm:"uniqueIndex"` // Contact phone number
	Address     string    // Staff member's address

	// Foreign keys are now strings referencing the Name field of lookup tables
	RoleID string `gorm:"not null;index"` // Foreign key referencing StaffRole(Name)
	// Role field for GORM relation. Use preload in queries: db.Preload("Role").Find(&staff)
	Role StaffRole `gorm:"foreignKey:RoleID;references:Name"`

	StatusID string `gorm:"not null;index"` // Foreign key referencing StaffStatus(Name)
	// Status field for GORM relation. Use preload in queries: db.Preload("Status").Find(&staff)
	Status StaffStatus `gorm:"foreignKey:StatusID;references:Name"`

	// --- STI Fields (Single Table Inheritance) ---
	Specialization string `gorm:"column:specialization"` // Specific to Doctors
	NurseType      string `gorm:"column:nurse_type"`     // Specific to Nurses

	// --- Relations ---
	// Staff workload/tasks are accessed via their schedule entries
	Schedule []ScheduleEntry `gorm:"foreignKey:StaffID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Workload []Task removed - Tasks are linked via ScheduleEntry
}

// TableName specifies the table name for Staff.
func (Staff) TableName() string {
	return "staff"
}

// --- Related Entities (Schedule, Task) ---

// ScheduleEntry represents the link between a staff member and a specific task in their schedule.
// It no longer has its own ID or timestamps, acting as a join structure.
type ScheduleEntry struct {
	// coreEntity.BaseEntity removed
	StaffID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"` // Foreign key to Staff (part of composite PK)
	// StartTime, EndTime, IsBooked removed

	// Link to Task (One ScheduleEntry -> One Task, mandatory)
	// Removed uniqueIndex to allow flexibility if a task could be in multiple schedule entries (though primary flow is Staff -> Schedule -> Task)
	TaskID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"` // Foreign key to Task (part of composite PK)
	Task   Task      `gorm:"foreignKey:TaskID"`                   // Task association (non-pointer)
}

// TableName specifies the table name for ScheduleEntry.
func (ScheduleEntry) TableName() string {
	return "staff_schedule_entries"
}

// Task represents a task assigned via a schedule entry.
type Task struct {
	coreEntity.BaseEntity
	// StaffID removed - Task is linked via ScheduleEntry
	Title       string    `gorm:"not null"` // Title of the task
	Description string    // Description of the task
	Priority    int       // Priority level
	StartTime   time.Time `gorm:"not null"` // Task start time (replaces TimeSlot)
	EndTime     time.Time `gorm:"not null"` // Task end time (replaces TimeSlot)
	// Date field removed

	// Link to TaskStatus (using string ID)
	StatusID string     `gorm:"not null;index"`                      // Renamed: Foreign key referencing TaskStatus(Name)
	Status   TaskStatus `gorm:"foreignKey:StatusID;references:Name"` // Renamed: GORM relation

}

// TableName specifies the table name for Task.
func (Task) TableName() string {
	return "staff_tasks"
}

// --- Constructors and Methods ---

// NewStaff creates a new Staff instance.
func NewStaff(firstName, lastName, phoneNumber, address string, dob time.Time, roleID, statusID string, specialization, nurseType string) *Staff { // roleID, statusID are now string
	staff := &Staff{
		FirstName:      firstName,
		LastName:       lastName,
		DateOfBirth:    dob,
		PhoneNumber:    phoneNumber,
		Address:        address,
		RoleID:         roleID,   // String
		StatusID:       statusID, // String
		Specialization: specialization,
		NurseType:      nurseType,
		Schedule:       []ScheduleEntry{},
		// Workload removed
	}
	return staff
}

// UpdateDetails updates staff member details.
func (s *Staff) UpdateDetails(firstName, lastName, phoneNumber, address string, dob time.Time, specialization, nurseType string) {
	changed := false
	if firstName != "" && s.FirstName != firstName {
		s.FirstName = firstName
		changed = true
	}
	if lastName != "" && s.LastName != lastName {
		s.LastName = lastName
		changed = true
	}
	if phoneNumber != "" && s.PhoneNumber != phoneNumber {
		s.PhoneNumber = phoneNumber
		changed = true
	}
	if address != "" && s.Address != address {
		s.Address = address
		changed = true
	}
	if !dob.IsZero() && s.DateOfBirth != dob {
		s.DateOfBirth = dob
		changed = true
	}

	// Only update role-specific fields if they are relevant
	if specialization != "" && s.Specialization != specialization {
		s.Specialization = specialization
		changed = true
	}
	if nurseType != "" && s.NurseType != nurseType {
		s.NurseType = nurseType
		changed = true
	}

	if changed {
		s.UpdatedAt = time.Now()
	}
}

// SetStatus updates the staff member's status ID.
func (s *Staff) SetStatus(statusID string) { // Expects string
	if s.StatusID != statusID {
		s.StatusID = statusID
		s.UpdatedAt = time.Now()
	}
}

// AddScheduleEntry adds a new entry to the staff's schedule.
// Assumes the entry comes with a valid TaskID already set.
// Since ScheduleEntry is now a simple join table, this just links Staff and Task.
func (s *Staff) AddScheduleEntry(entry ScheduleEntry) {
	if entry.StaffID == uuid.Nil {
		entry.StaffID = s.ID
	}
	// Ensure entry.TaskID is not uuid.Nil if Task is mandatory
	if entry.TaskID == uuid.Nil {
		// Handle error: A schedule entry must have an associated task
		// This might involve returning an error or logging.
		// For now, just appending, but this indicates an invalid state.
	}
	s.Schedule = append(s.Schedule, entry)
	// s.UpdatedAt = time.Now() // Staff record itself isn't directly changed by adding a schedule link
}

// AssignTask method removed from Staff entity.
