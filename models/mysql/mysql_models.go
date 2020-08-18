package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

//region MySQLNullTime

// MySQLNullTime is an alias for mysql.NullTime data type
type MySQLNullTime struct {
	mysql.NullTime
}

// MarshalJSON for MySQLNullTime
func (nt MySQLNullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// UnmarshalJSON for MySQLNullTime
func (nt *MySQLNullTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		nt.Valid = false
		return err
	}

	nt.Time = x
	nt.Valid = true
	return nil
}

// NullTime is an alias for mysql.NullTime data type
//type NullTime struct {
//	mysql.NullTime
//}

//endregion

//region MySQLNullString

// MySQLNullString ...
type MySQLNullString struct {
	sql.NullString
}

// MarshalJSON for MySQLNullString
func (nt MySQLNullString) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte(""), nil
	}
	return []byte(nt.String), nil
}

// UnmarshalJSON for MySQLNullString
func (nt *MySQLNullString) UnmarshalJSON(b []byte) error {
	s := string(b)
	nt.String = s
	nt.Valid = true
	return nil
}

//endregion

//region MySQLNullInt64

//MySQLNullInt64 ...
type MySQLNullInt64 struct {
	sql.NullInt64
}

//MarshalJSON for MySQLNullInt64
func (t MySQLNullInt64) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Int64)
}

//endregion

//region MySQLNullInt32

// MySQLNullInt32 ...
type MySQLNullInt32 struct {
	sql.NullInt32
}

// MarshalJSON for MySQLNullInt32
func (t MySQLNullInt32) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Int32)
}

//endregion
