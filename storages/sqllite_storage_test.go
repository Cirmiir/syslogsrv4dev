package storages

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	unittestDb = "unittest.db"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

type testBody func(*testing.T)

type testCase struct {
	Name string
	Body testBody
}

var tests = []testCase{
	{Name: "NewEventTest", Body: newLog},
	{Name: "FilterIncludeFacilityTest", Body: func(t *testing.T) {
		filter(t, Facility, []int{2}, 1)
	}},
	{Name: "FilterExcludeFacilityTest", Body: func(t *testing.T) {
		filter(t, Facility, []int{1}, 0)
	}},
	{Name: "FilterExcludeSeverityTest", Body: func(t *testing.T) {
		filter(t, Severity, []int{2}, 0)
	}},
	{Name: "FilterIncludeSeverityTest", Body: func(t *testing.T) {
		filter(t, Severity, []int{1}, 1)
	}},
	{Name: "FilterIncludeAppNameTest", Body: func(t *testing.T) {
		filter(t, Application, []string{"UnitTestApp"}, 1)
	}},
	{Name: "FilterExcludeAppNameTest", Body: func(t *testing.T) {
		filter(t, Application, []string{"UnitTestApp1"}, 0)
	}},
	{Name: "DateFilterOneDay", Body: func(t *testing.T) {
		startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		filterDate(t, startDate, endDate, 1)
	}},
	{Name: "DateFilterOutOfRange", Body: func(t *testing.T) {
		startDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		filterDate(t, startDate, endDate, 0)
	}},
	{Name: "DateFilterStartDate", Body: func(t *testing.T) {
		startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		filterDate(t, startDate, endDate, 1)
	}},
	{Name: "DateFilterEndDate", Body: func(t *testing.T) {
		startDate := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		filterDate(t, startDate, endDate, 1)
	}},
}

func TestStorage(t *testing.T) {
	for _, tc := range tests {
		testSetup()
		tc := tc // capture range variable
		t.Run(tc.Name, tc.Body)
	}
}

func newLog(t *testing.T) {
	db := NewSqlLite3Storage("unittest.db")

	if db == nil {
		t.Error("Database cannot be initialized.")
	}

	appName, message, host, sev, facility := "UnitTestApp", "Test Message", "localhost", 1, 2
	date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	db.SaveEvent(LogEvent{
		AppName:  appName,
		Host:     host,
		Severity: sev,
		Facility: facility,
		Message:  message,
		Date:     date,
	})

	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	logs, cnt := db.FetchEvents(FetchCriteria{StartDate: startDate, EndDate: endDate}, 0, 1)

	if cnt != 1 {
		t.Fatalf("The wrong count of events.")
	}

	log := logs[0]

	assertEqual(t, log.AppName, appName)
	assertEqual(t, log.Facility, facility)
	assertEqual(t, log.Severity, sev)
	assertEqual(t, log.Date, date)
	assertEqual(t, log.Host, host)
	assertEqual(t, log.Message, message)
}

func filter(t *testing.T, filter Filter, values interface{}, expectedCount int) {
	db := NewSqlLite3Storage("unittest.db")

	if db == nil {
		t.Error("Database cannot be initialized.")
	}

	appName, message, host, sev, facility := "UnitTestApp", "Test Message", "localhost", 1, 2
	date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	db.SaveEvent(LogEvent{
		AppName:  appName,
		Host:     host,
		Severity: sev,
		Facility: facility,
		Message:  message,
		Date:     date,
	})

	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	logs, cnt := db.FetchEvents(FetchCriteria{StartDate: startDate, EndDate: endDate, Filters: map[Filter]interface{}{filter: values}}, 0, 1)

	if cnt != expectedCount || len(logs) != expectedCount {
		t.Fatalf("The wrong count of events. Filter doesn't work. Actual %v. Expected %v", cnt, 1)
	}
}

func filterDate(t *testing.T, startDate time.Time, endDate time.Time, expectedCount int) {
	db := NewSqlLite3Storage("unittest.db")

	if db == nil {
		t.Error("Database cannot be initialized.")
	}

	appName, message, host, sev, facility := "UnitTestApp", "Test Message", "localhost", 1, 2
	date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	db.SaveEvent(LogEvent{
		AppName:  appName,
		Host:     host,
		Severity: sev,
		Facility: facility,
		Message:  message,
		Date:     date,
	})

	logs, cnt := db.FetchEvents(FetchCriteria{StartDate: startDate, EndDate: endDate}, 0, 1)

	if cnt != expectedCount || len(logs) != expectedCount {
		t.Fatalf("The wrong count of events. Filter doesn't work. Actual %v. Expected %v", cnt, expectedCount)
	}
}

func testSetup() bool {
	db, err := sql.Open("sqlite3", unittestDb)

	if err != nil {
		panic("Error")
	}

	_, err = db.Exec(`DELETE FROM logs`)

	return err == nil

}

func setup() bool {
	if _, err := os.Stat(unittestDb); err == nil {
		if err = os.Remove(unittestDb); err != nil {
			return false
		}
	}
	return true
}

func shutdown() {
}

func TestMain(m *testing.M) {
	if !setup() {
		os.Exit(1)
	}
	code := m.Run()
	shutdown()
	os.Exit(code)
}
