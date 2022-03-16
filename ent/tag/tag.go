// Code generated by entc, DO NOT EDIT.

package tag

import (
	"time"
)

const (
	// Label holds the string label denoting the tag type in the database.
	Label = "tag"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeRecruitmentTags holds the string denoting the recruitment_tags edge name in mutations.
	EdgeRecruitmentTags = "recruitment_tags"
	// Table holds the table name of the tag in the database.
	Table = "tags"
	// RecruitmentTagsTable is the table that holds the recruitment_tags relation/edge.
	RecruitmentTagsTable = "recruitment_tags"
	// RecruitmentTagsInverseTable is the table name for the RecruitmentTag entity.
	// It exists in this package in order to avoid circular dependency with the "recruitmenttag" package.
	RecruitmentTagsInverseTable = "recruitment_tags"
	// RecruitmentTagsColumn is the table column denoting the recruitment_tags relation/edge.
	RecruitmentTagsColumn = "tag_id"
)

// Columns holds all SQL columns for tag fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldName,
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