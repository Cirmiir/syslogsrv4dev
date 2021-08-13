package storages

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SqlLite3Storage struct {
	Db *sql.DB
}

func NewSqlLite3Storage(dbFile string) *SqlLite3Storage {

	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatal(err)
	}

	initScript := `create table if not exists Logs (host text, message text, severity int, facility int, date datetime, appname text)`
	_, err = db.Exec(initScript)
	if err != nil {
		log.Printf("%q: %s\n", err, initScript)
		return nil
	}

	return &SqlLite3Storage{Db: db}
}

func (storage *SqlLite3Storage) SaveEvent(event LogEvent) (bool, error) {
	insertScript := `INSERT INTO Logs(host, message, severity, facility, date, appname) VALUES(?,?,?,?,?,?)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := storage.Db.PrepareContext(ctx, insertScript)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, event.Host, event.Message, event.Severity, event.Facility, event.Date.Format(time.RFC3339), event.AppName)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return false, err
	}
	return true, nil

}

func escape(val interface{}) string {
	switch v := val.(type) {
	case string:
		return "'" + v + "'"
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

func addFilter(whereClause string, filter Filter, values interface{}) string {
	filterValues := make([]string, 0)
	switch reflect.TypeOf(values).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(values)

		length := s.Len()

		if length == 0 {
			return whereClause
		}

		for i := 0; i < length; i++ {
			filterValues = append(filterValues, escape(s.Index(i).Interface()))
		}

		if len(whereClause) == 0 {
			return filter.GetColumnForSQL() + " IN (" + strings.Join(filterValues, ",") + ")"
		}

		return whereClause + " AND " + filter.GetColumnForSQL() + " IN (" + strings.Join(filterValues, ",") + ")"
	default:
		if len(whereClause) == 0 {
			return filter.GetColumnForSQL() + " = " + escape(values)
		}

		return whereClause + " AND " + filter.GetColumnForSQL() + " = " + escape(values)
	}
}

func (storage *SqlLite3Storage) GenerateWhere(criteria FetchCriteria) string {
	res := ""
	for key, v := range criteria.Filters {
		res = addFilter(res, key, v)
	}
	return res
}

func (storage *SqlLite3Storage) FetchEvents(criteria FetchCriteria, page int, limit int) ([]LogEvent, int) {

	where := storage.GenerateWhere(criteria)

	criteria.StartDate = time.Date(criteria.StartDate.Year(), criteria.StartDate.Month(), criteria.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	criteria.EndDate = time.Date(criteria.EndDate.Year(), criteria.EndDate.Month(), criteria.EndDate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Add(time.Duration(-1))

	if len(where) != 0 {
		where = " AND " + where
	}

	selectScript := `
	SELECT host, message,severity, facility, date, appname
	FROM  Logs
	WHERE date >= ? AND date <= ?
	` + where + `
	ORDER By date DESC
	LIMIT ?, ?`

	rows, err := storage.Db.Query(selectScript, criteria.StartDate.Format(time.RFC3339), criteria.EndDate.Format(time.RFC3339), page*limit, limit)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	result := make([]LogEvent, 0)

	for rows.Next() {
		var (
			host     string
			message  string
			severity int
			facility int
			appname  string
			date     time.Time
		)
		if err := rows.Scan(&host, &message, &severity, &facility, &date, &appname); err != nil {
			log.Fatal(err)
		}

		result = append(result, LogEvent{Host: host, Message: message, Severity: severity, Facility: facility, Date: date, AppName: appname})
	}

	var count int
	row := storage.Db.QueryRow("SELECT COUNT(*) FROM Logs WHERE date >= ? AND date <= ?"+where, criteria.StartDate.Format(time.RFC3339), criteria.EndDate.Format(time.RFC3339))
	err = row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return result, count
}

func (storage *SqlLite3Storage) FetchFilter(filter Filter) ([]Lookup, error) {
	selectScript := `SELECT DISTINCT ` + filter.GetColumnForSQL() + ` FROM Logs`
	rows, err := storage.Db.Query(selectScript)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	result := make([]Lookup, 0)

	for rows.Next() {
		var (
			key string
		)
		if err := rows.Scan(&key); err != nil {
			log.Fatal(err)
		}
		id, err := strconv.Atoi(key)
		if err == nil {
			result = append(result, Lookup{Key: id, Name: key})
		} else {
			result = append(result, Lookup{Key: key, Name: key})
		}
	}

	return result, err
}
