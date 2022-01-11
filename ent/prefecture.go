// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/nagokos/connefut_backend/ent/prefecture"
)

// Prefecture is the model entity for the Prefecture schema.
type Prefecture struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Prefecture) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case prefecture.FieldID, prefecture.FieldName:
			values[i] = new(sql.NullString)
		case prefecture.FieldCreatedAt, prefecture.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Prefecture", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Prefecture fields.
func (pr *Prefecture) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case prefecture.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				pr.ID = value.String
			}
		case prefecture.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pr.CreatedAt = value.Time
			}
		case prefecture.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pr.UpdatedAt = value.Time
			}
		case prefecture.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pr.Name = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Prefecture.
// Note that you need to call Prefecture.Unwrap() before calling this method if this Prefecture
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Prefecture) Update() *PrefectureUpdateOne {
	return (&PrefectureClient{config: pr.config}).UpdateOne(pr)
}

// Unwrap unwraps the Prefecture entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *Prefecture) Unwrap() *Prefecture {
	tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("ent: Prefecture is not a transactional entity")
	}
	pr.config.driver = tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Prefecture) String() string {
	var builder strings.Builder
	builder.WriteString("Prefecture(")
	builder.WriteString(fmt.Sprintf("id=%v", pr.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(pr.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(pr.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(pr.Name)
	builder.WriteByte(')')
	return builder.String()
}

// Prefectures is a parsable slice of Prefecture.
type Prefectures []*Prefecture

func (pr Prefectures) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}