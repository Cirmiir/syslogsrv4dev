package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	_ "net/http/pprof"

	rice "github.com/GeertJohan/go.rice"
	store "github.com/cirmiir/syslogsrv4dev/storages"
)

type filterModel struct {
	Title    string
	Binding  string
	JSONData string
}

type page struct {
	Title    string
	Body     []byte
	Filters  []filterModel
	Settings string
}

type apiResponse struct {
	Data       []store.LogEvent
	TotalCount int
}

type searchCriteria struct {
	store.Criteria
	StartDate time.Time
	EndDate   time.Time
	Facility  []int
	Severity  []int
	AppName   []string
}

func runWeblistener(wg *sync.WaitGroup) {
	defer wg.Done()

	appBox, _ := rice.FindBox("./app/build")

	httpServer := http.NewServeMux()

	httpServer.HandleFunc("/api/", webapi)
	httpServer.HandleFunc("/", serveAppHandler(appBox))
	httpServer.Handle("/static/", http.FileServer(appBox.HTTPBox()))
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(cfg.WebPort), httpServer)
}

func serveAppHandler(app *rice.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := app.Open("index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		filters := []store.Filter{store.Facility, store.Application, store.Severity}
		models := make([]filterModel, 0)

		for _, filter := range filters {
			data, _ := storage.FetchFilter(filter)
			arr, _ := json.Marshal(data)
			models = append(models,
				filterModel{
					Title:    filter.GetTitle(),
					Binding:  filter.GetBinding(),
					JSONData: string(arr),
				})
		}

		conf, _ := json.Marshal(cfg)

		p := &page{
			Title:    "Log Server",
			Filters:  models,
			Settings: string(conf),
		}

		file, _ := indexFile.Stat()
		length := file.Size()

		content := make([]byte, length)
		_, err = indexFile.Read(content)
		if err != nil {
			return
		}
		t, _ := template.New("index").Parse(string(content))
		t.Execute(w, p)
	}
}

func webapi(w http.ResponseWriter, r *http.Request) {

	criteria := searchCriteria{}
	err := json.NewDecoder(r.Body).Decode(&criteria)

	if err != nil {
		log.Println(err)
	}

	start := time.Date(criteria.StartDate.Year(), criteria.StartDate.Month(), criteria.StartDate.Day(), 0, 0, 0, 1, time.UTC)
	end := time.Date(criteria.EndDate.Year(), criteria.EndDate.Month(), criteria.EndDate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Add(time.Duration(-1))

	cr := store.FetchCriteria{
		Criteria:  criteria.Criteria,
		StartDate: start,
		EndDate:   end,
		Filters: map[store.Filter]interface{}{
			store.Facility:    criteria.Facility,
			store.Severity:    criteria.Severity,
			store.Application: criteria.AppName,
		},
	}

	logs, totalCount := storage.FetchEvents(cr, criteria.Page, criteria.Limit)

	data, err := json.Marshal(apiResponse{Data: logs, TotalCount: totalCount})

	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")

	w.Write(data)
}
