// Code generated by entc, DO NOT EDIT.

package recruitment

import (
	"fmt"
	"time"
)

const (
	// Label holds the string label denoting the recruitment type in the database.
	Label = "recruitment"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldPlace holds the string denoting the place field in the database.
	FieldPlace = "place"
	// FieldStartAt holds the string denoting the start_at field in the database.
	FieldStartAt = "start_at"
	// FieldContent holds the string denoting the content field in the database.
	FieldContent = "content"
	// FieldLocationLat holds the string denoting the locationlat field in the database.
	FieldLocationLat = "location_lat"
	// FieldLocationLng holds the string denoting the locationlng field in the database.
	FieldLocationLng = "location_lng"
	// FieldCapacity holds the string denoting the capacity field in the database.
	FieldCapacity = "capacity"
	// FieldClosingAt holds the string denoting the closing_at field in the database.
	FieldClosingAt = "closing_at"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldPrefectureID holds the string denoting the prefecture_id field in the database.
	FieldPrefectureID = "prefecture_id"
	// FieldCompetitionID holds the string denoting the competition_id field in the database.
	FieldCompetitionID = "competition_id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// EdgeStocks holds the string denoting the stocks edge name in mutations.
	EdgeStocks = "stocks"
	// EdgeApplicants holds the string denoting the applicants edge name in mutations.
	EdgeApplicants = "applicants"
	// EdgeRecruitmentTags holds the string denoting the recruitment_tags edge name in mutations.
	EdgeRecruitmentTags = "recruitment_tags"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgePrefecture holds the string denoting the prefecture edge name in mutations.
	EdgePrefecture = "prefecture"
	// EdgeCompetition holds the string denoting the competition edge name in mutations.
	EdgeCompetition = "competition"
	// Table holds the table name of the recruitment in the database.
	Table = "recruitments"
	// StocksTable is the table that holds the stocks relation/edge.
	StocksTable = "stocks"
	// StocksInverseTable is the table name for the Stock entity.
	// It exists in this package in order to avoid circular dependency with the "stock" package.
	StocksInverseTable = "stocks"
	// StocksColumn is the table column denoting the stocks relation/edge.
	StocksColumn = "recruitment_id"
	// ApplicantsTable is the table that holds the applicants relation/edge.
	ApplicantsTable = "applicants"
	// ApplicantsInverseTable is the table name for the Applicant entity.
	// It exists in this package in order to avoid circular dependency with the "applicant" package.
	ApplicantsInverseTable = "applicants"
	// ApplicantsColumn is the table column denoting the applicants relation/edge.
	ApplicantsColumn = "recruitment_id"
	// RecruitmentTagsTable is the table that holds the recruitment_tags relation/edge.
	RecruitmentTagsTable = "recruitment_tags"
	// RecruitmentTagsInverseTable is the table name for the RecruitmentTag entity.
	// It exists in this package in order to avoid circular dependency with the "recruitmenttag" package.
	RecruitmentTagsInverseTable = "recruitment_tags"
	// RecruitmentTagsColumn is the table column denoting the recruitment_tags relation/edge.
	RecruitmentTagsColumn = "recruitment_id"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "recruitments"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_id"
	// PrefectureTable is the table that holds the prefecture relation/edge.
	PrefectureTable = "recruitments"
	// PrefectureInverseTable is the table name for the Prefecture entity.
	// It exists in this package in order to avoid circular dependency with the "prefecture" package.
	PrefectureInverseTable = "prefectures"
	// PrefectureColumn is the table column denoting the prefecture relation/edge.
	PrefectureColumn = "prefecture_id"
	// CompetitionTable is the table that holds the competition relation/edge.
	CompetitionTable = "recruitments"
	// CompetitionInverseTable is the table name for the Competition entity.
	// It exists in this package in order to avoid circular dependency with the "competition" package.
	CompetitionInverseTable = "competitions"
	// CompetitionColumn is the table column denoting the competition relation/edge.
	CompetitionColumn = "competition_id"
)

// Columns holds all SQL columns for recruitment fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldTitle,
	FieldType,
	FieldPlace,
	FieldStartAt,
	FieldContent,
	FieldLocationLat,
	FieldLocationLng,
	FieldCapacity,
	FieldClosingAt,
	FieldStatus,
	FieldPrefectureID,
	FieldCompetitionID,
	FieldUserID,
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
	// TitleValidator is a validator for the "title" field. It is called by the builders before save.
	TitleValidator func(string) error
	// ContentValidator is a validator for the "content" field. It is called by the builders before save.
	ContentValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() string
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
)

// Type defines the type for the "type" enum field.
type Type string

// TypeUnnecessary is the default value of the Type enum.
const DefaultType = TypeUnnecessary

// Type values.
const (
	TypeUnnecessary Type = "unnecessary"
	TypeOpponent    Type = "opponent"
	TypeIndividual  Type = "individual"
	TypeMember      Type = "member"
	TypeJoining     Type = "joining"
	TypeOthers      Type = "others"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeUnnecessary, TypeOpponent, TypeIndividual, TypeMember, TypeJoining, TypeOthers:
		return nil
	default:
		return fmt.Errorf("recruitment: invalid enum value for type field: %q", _type)
	}
}

// Status defines the type for the "status" enum field.
type Status string

// StatusDraft is the default value of the Status enum.
const DefaultStatus = StatusDraft

// Status values.
const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusClosed    Status = "closed"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusDraft, StatusPublished, StatusClosed:
		return nil
	default:
		return fmt.Errorf("recruitment: invalid enum value for status field: %q", s)
	}
}
