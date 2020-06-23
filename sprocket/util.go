package sprocket

import (
	"database/sql"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
	"github.com/jmoiron/sqlx"
)

//GetTemplate ...
func GetTemplate(viewbox *rice.Box, base string, page string) (*template.Template, error) {
	base, err := viewbox.String(base)
	if err != nil {
		return nil, err
	}

	content, err := viewbox.String(page)
	if err != nil {
		return nil, err
	}

	x, err := template.New("base").Parse(base)
	if err != nil {
		return nil, err
	}

	x.New("content").Parse(content)
	if err != nil {
		return nil, err
	}
	return x, nil
}

//GetTemplates ...
func GetTemplates(viewbox *rice.Box, filenames []string) (*template.Template, error) {
	var x *template.Template
	for i := 0; i < len(filenames); i++ {
		tmp, err := viewbox.String(filenames[i])
		if err != nil {
			return nil, err
		}
		if i == 0 {
			x, err = template.New("base").Parse(tmp)
		} else {
			x.New("content").Parse(tmp)
		}
	}
	return x, nil
}

//TransactionQuery ...
func TransactionQuery(db *sqlx.DB, query string, args ...interface{}) (sql.Result, error) {
	err := db.Ping()
	if err == nil {
		tx, err := db.Beginx()
		if err == nil {
			sqlResult, err := tx.Exec(query, args...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			tx.Commit()
			return sqlResult, nil
		}
	}
	return nil, err
}

//TxTransactionQuery ...
func TxTransactionQuery(tx *sqlx.Tx, query string, args ...interface{}) (sql.Result, error) {
	sqlResult, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return sqlResult, nil
}

//RowScanWrap ...
func RowScanWrap(r *sql.Row, dest ...interface{}) error {
	return r.Scan(dest...)
}

//RowsScanWrap ...
func RowsScanWrap(r *sql.Rows, dest ...interface{}) error {
	return r.Scan(dest...)
}
