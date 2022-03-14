// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ApplicantsColumns holds the columns for the "applicants" table.
	ApplicantsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "management_status", Type: field.TypeEnum, Enums: []string{"backlog", "checked", "accepted", "rejected"}, Default: "backlog"},
		{Name: "recruitment_id", Type: field.TypeString, Nullable: true},
		{Name: "user_id", Type: field.TypeString, Nullable: true},
	}
	// ApplicantsTable holds the schema information for the "applicants" table.
	ApplicantsTable = &schema.Table{
		Name:       "applicants",
		Columns:    ApplicantsColumns,
		PrimaryKey: []*schema.Column{ApplicantsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "applicants_recruitments_applicants",
				Columns:    []*schema.Column{ApplicantsColumns[4]},
				RefColumns: []*schema.Column{RecruitmentsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "applicants_users_applicants",
				Columns:    []*schema.Column{ApplicantsColumns[5]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "applicant_user_id_recruitment_id",
				Unique:  true,
				Columns: []*schema.Column{ApplicantsColumns[5], ApplicantsColumns[4]},
			},
			{
				Name:    "applicant_user_id",
				Unique:  false,
				Columns: []*schema.Column{ApplicantsColumns[5]},
			},
			{
				Name:    "applicant_recruitment_id",
				Unique:  false,
				Columns: []*schema.Column{ApplicantsColumns[4]},
			},
		},
	}
	// CompetitionsColumns holds the columns for the "competitions" table.
	CompetitionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// CompetitionsTable holds the schema information for the "competitions" table.
	CompetitionsTable = &schema.Table{
		Name:       "competitions",
		Columns:    CompetitionsColumns,
		PrimaryKey: []*schema.Column{CompetitionsColumns[0]},
	}
	// PrefecturesColumns holds the columns for the "prefectures" table.
	PrefecturesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
	}
	// PrefecturesTable holds the schema information for the "prefectures" table.
	PrefecturesTable = &schema.Table{
		Name:       "prefectures",
		Columns:    PrefecturesColumns,
		PrimaryKey: []*schema.Column{PrefecturesColumns[0]},
	}
	// RecruitmentsColumns holds the columns for the "recruitments" table.
	RecruitmentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "title", Type: field.TypeString, SchemaType: map[string]string{"postgres": "varchar(60)"}},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"unnecessary", "opponent", "individual", "member", "joining", "others"}, Default: "unnecessary"},
		{Name: "place", Type: field.TypeString, Nullable: true},
		{Name: "start_at", Type: field.TypeTime, Nullable: true},
		{Name: "content", Type: field.TypeString, Nullable: true, SchemaType: map[string]string{"postgres": "varchar(10000)"}},
		{Name: "location_lat", Type: field.TypeFloat64, Nullable: true},
		{Name: "location_lng", Type: field.TypeFloat64, Nullable: true},
		{Name: "capacity", Type: field.TypeInt, Nullable: true},
		{Name: "closing_at", Type: field.TypeTime, Nullable: true},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"draft", "published", "closed"}, Default: "draft"},
		{Name: "competition_id", Type: field.TypeString, Nullable: true},
		{Name: "prefecture_id", Type: field.TypeString, Nullable: true},
		{Name: "user_id", Type: field.TypeString, Nullable: true},
	}
	// RecruitmentsTable holds the schema information for the "recruitments" table.
	RecruitmentsTable = &schema.Table{
		Name:       "recruitments",
		Columns:    RecruitmentsColumns,
		PrimaryKey: []*schema.Column{RecruitmentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "recruitments_competitions_recruitments",
				Columns:    []*schema.Column{RecruitmentsColumns[13]},
				RefColumns: []*schema.Column{CompetitionsColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "recruitments_prefectures_recruitments",
				Columns:    []*schema.Column{RecruitmentsColumns[14]},
				RefColumns: []*schema.Column{PrefecturesColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "recruitments_users_recruitments",
				Columns:    []*schema.Column{RecruitmentsColumns[15]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "recruitment_user_id",
				Unique:  false,
				Columns: []*schema.Column{RecruitmentsColumns[15]},
			},
			{
				Name:    "recruitment_prefecture_id",
				Unique:  false,
				Columns: []*schema.Column{RecruitmentsColumns[14]},
			},
			{
				Name:    "recruitment_competition_id",
				Unique:  false,
				Columns: []*schema.Column{RecruitmentsColumns[13]},
			},
		},
	}
	// StocksColumns holds the columns for the "stocks" table.
	StocksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "recruitment_id", Type: field.TypeString, Nullable: true},
		{Name: "user_id", Type: field.TypeString, Nullable: true},
	}
	// StocksTable holds the schema information for the "stocks" table.
	StocksTable = &schema.Table{
		Name:       "stocks",
		Columns:    StocksColumns,
		PrimaryKey: []*schema.Column{StocksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "stocks_recruitments_stocks",
				Columns:    []*schema.Column{StocksColumns[3]},
				RefColumns: []*schema.Column{RecruitmentsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "stocks_users_stocks",
				Columns:    []*schema.Column{StocksColumns[4]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "stock_user_id_recruitment_id",
				Unique:  true,
				Columns: []*schema.Column{StocksColumns[4], StocksColumns[3]},
			},
			{
				Name:    "stock_user_id",
				Unique:  false,
				Columns: []*schema.Column{StocksColumns[4]},
			},
			{
				Name:    "stock_recruitment_id",
				Unique:  false,
				Columns: []*schema.Column{StocksColumns[3]},
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, SchemaType: map[string]string{"postgres": "varchar(50)"}},
		{Name: "email", Type: field.TypeString, Unique: true, SchemaType: map[string]string{"postgres": "varchar(100)"}},
		{Name: "role", Type: field.TypeEnum, Enums: []string{"admin", "general"}, Default: "general"},
		{Name: "avatar", Type: field.TypeString, Default: "https://abs.twimg.com/sticky/default_profile_images/default_profile.png"},
		{Name: "introduction", Type: field.TypeString, Nullable: true, SchemaType: map[string]string{"postgres": "varchar(4000)"}},
		{Name: "email_verification_status", Type: field.TypeEnum, Enums: []string{"unnecessary", "pending", "verified"}, Default: "pending"},
		{Name: "email_verification_token", Type: field.TypeString, Nullable: true},
		{Name: "email_verification_token_expires_at", Type: field.TypeTime, Nullable: true},
		{Name: "password_digest", Type: field.TypeString, Nullable: true},
		{Name: "last_sign_in_at", Type: field.TypeTime, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "user_email_verification_token",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[9]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ApplicantsTable,
		CompetitionsTable,
		PrefecturesTable,
		RecruitmentsTable,
		StocksTable,
		UsersTable,
	}
)

func init() {
	ApplicantsTable.ForeignKeys[0].RefTable = RecruitmentsTable
	ApplicantsTable.ForeignKeys[1].RefTable = UsersTable
	RecruitmentsTable.ForeignKeys[0].RefTable = CompetitionsTable
	RecruitmentsTable.ForeignKeys[1].RefTable = PrefecturesTable
	RecruitmentsTable.ForeignKeys[2].RefTable = UsersTable
	StocksTable.ForeignKeys[0].RefTable = RecruitmentsTable
	StocksTable.ForeignKeys[1].RefTable = UsersTable
}
