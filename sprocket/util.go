package sprocket

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-zoo/bone"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

// GetTemplates ...
func GetTemplates(viewFS *embed.FS, patterns []string, filenames []string) (*template.Template, error) {
	templateFS, err := template.ParseFS(viewFS, patterns...)
	if err != nil {
		return nil, err
	}
	return templateFS.New("base").ParseFiles(filenames...)
}

// TransactionQuery ...
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

// TxTransactionQuery ...
func TxTransactionQuery(tx *sqlx.Tx, query string, args ...interface{}) (sql.Result, error) {
	sqlResult, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return sqlResult, nil
}

// RowScanWrap ...
func RowScanWrap(r *sql.Row, dest ...interface{}) error {
	return r.Scan(dest...)
}

// RowsScanWrap ...
func RowsScanWrap(r *sql.Rows, dest ...interface{}) error {
	return r.Scan(dest...)
}

// JustJSONMarshal ...
func JustJSONMarshal(v interface{}) string {
	result, err := json.Marshal(v)
	if err != nil {
		log.Panic(err)
	}
	return string(result)
}

// InterfaceToString ...
func InterfaceToString(v interface{}) string {
	if v != nil {
		return v.(string)
	}
	return ""
}

// InterfaceArrayToStringArray ...
func InterfaceArrayToStringArray(t []interface{}) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = fmt.Sprint(v)
	}
	return s
}

// ParseToDatabaseDate a helper method to parse string to time
// layout example `02/01/2006`
func ParseToDatabaseDate(layout string, v string) interface{} {
	parsedDate, err := time.Parse(layout, strings.TrimSpace(v))
	if err != nil {
		return nil
	}
	return parsedDate
}

// ExtractQueryParamAsInt ...
func ExtractQueryParamAsInt(r *http.Request, name string, defaultValue int) int {
	s := bone.GetQuery(r, name)
	if len(s) > 0 {
		if tmp, err := strconv.Atoi(s[0]); err == nil {
			return tmp
		}
	}
	return defaultValue
}

// ExtractQueryParamAsString ...
func ExtractQueryParamAsString(r *http.Request, name string, defaultValue string) string {
	s := bone.GetQuery(r, name)
	if len(s) > 0 {
		return s[0]
	}
	return defaultValue
}

// MetaValueStringExtractor ...
func MetaValueStringExtractor(v interface{}) string {
	if v != nil {
		return v.(string)
	}
	return ""
}

// LoadYAML used to load config.yml into the program
func LoadYAML(filename string, configuration interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(configuration)
	if err != nil {
		return err
	}
	return nil
}

// RespondOkayJSON ...
func RespondOkayJSON(w http.ResponseWriter, payload interface{}) {
	RespondStatusCodeWithJSON(w, http.StatusOK, payload)
}

// RespondInternalServerError ...
func RespondInternalServerError(w http.ResponseWriter, payload interface{}) {
	RespondStatusCodeWithJSON(w, http.StatusInternalServerError, payload)
}

// RespondNotImplementedJSON ...
func RespondNotImplementedJSON(w http.ResponseWriter, payload interface{}) {
	RespondStatusCodeWithJSON(w, http.StatusNotImplemented, payload)
}

// RespondNotFoundJSON ...
func RespondNotFoundJSON(w http.ResponseWriter, payload interface{}) {
	RespondStatusCodeWithJSON(w, http.StatusNotFound, payload)
}

// RespondBadRequestJSON ...
func RespondBadRequestJSON(w http.ResponseWriter, payload interface{}) {
	RespondStatusCodeWithJSON(w, http.StatusBadRequest, payload)
}

// RespondStatusCodeWithJSON ...
func RespondStatusCodeWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		crashPayload := []byte(`{
			"success": false,
			"result": null,
			"errors": ["` + fmt.Sprintf("%s", err.Error()) + `"]
		}`)
		w.Write(crashPayload)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(response)
}

// GenerateMySQLTotalCount this is used to generate total records base on the given SQL query for pagination
func GenerateMySQLTotalCount(sql string) string {
	upperSQL := strings.ToUpper(sql)
	upperSQL = strings.ReplaceAll(upperSQL, "\n", " ")
	indexOfFrom := strings.Index(upperSQL, " FROM ")
	if indexOfFrom != -1 {
		return "SELECT COUNT(*) FROM " + strings.ReplaceAll(sql[indexOfFrom+6:], "\n", " ")
	}
	return ""
}
