// Code generated by entc, DO NOT EDIT.

package applicant

import (
	"fmt"
	"time"
)

const (
	// Label holds the string label denoting the applicant type in the database.
	Label = "applicant"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldManagementStatus holds the string denoting the management_status field in the database.
	FieldManagementStatus = "management_status"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldRecruitmentID holds the string denoting the recruitment_id field in the database.
	FieldRecruitmentID = "recruitment_id"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeRecruitment holds the string denoting the recruitment edge name in mutations.
	EdgeRecruitment = "recruitment"
	// Table holds the table name of the applicant in the database.
	Table = "applicants"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "applicants"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_id"
	// RecruitmentTable is the table that holds the recruitment relation/edge.
	RecruitmentTable = "applicants"
	// RecruitmentInverseTable is the table name for the Recruitment entity.
	// It exists in this package in order to avoid circular dependency with the "recruitment" package.
	RecruitmentInverseTable = "recruitments"
	// RecruitmentColumn is the table column denoting the recruitment relation/edge.
	RecruitmentColumn = "recruitment_id"
)

// Columns holds all SQL columns for applicant fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldManagementStatus,
	FieldUserID,
	FieldRecruitmentID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() string
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
)

// ManagementStatus defines the type for the "management_status" enum field.
type ManagementStatus string

// ManagementStatusBacklog is the default value of the ManagementStatus enum.
const DefaultManagementStatus = ManagementStatusBacklog

// ManagementStatus values.
const (
	ManagementStatusBacklog  ManagementStatus = "backlog"
	ManagementStatusChecked  ManagementStatus = "checked"
	ManagementStatusAccepted ManagementStatus = "accepted"
	ManagementStatusRejected ManagementStatus = "rejected"
)

func (ms ManagementStatus) String() string {
	return string(ms)
}

// ManagementStatusValidator is a validator for the "management_status" field enum values. It is called by the builders before save.
func ManagementStatusValidator(ms ManagementStatus) error {
	switch ms {
	case ManagementStatusBacklog, ManagementStatusChecked, ManagementStatusAccepted, ManagementStatusRejected:
		return nil
	default:
		return fmt.Errorf("applicant: invalid enum value for management_status field: %q", ms)
	}
}
