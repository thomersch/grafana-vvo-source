package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/thomersch/grafana-vvo-source/source"
)

type responseTable struct {
	Columns []map[string]interface{} `json:"columns"`
	Rows    [][]interface{}          `json:"rows"`
	Typ     string                   `json:"type"`
}

type requestData struct {
	Targets []map[string]string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// This is not a good handling. But the Simple JSON plugin gives up if there are too many results.
		var sts []string
		for _, st := range source.Stations {
			if st.Place == "Dresden" {
				sts = append(sts, fmt.Sprintf("%s", st.Station))
			}
		}
		json.NewEncoder(w).Encode(sts)
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		var req requestData
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			return
		}
		r.Body.Close()

		station, ok := req.Targets[0]["target"]
		if !ok {
			log.Println("no target specified")
			return
		}
		tms, err := source.NextTimes(station, "Dresden")
		if err != nil {
			log.Println(err)
			return
		}

		var rows [][]interface{}
		for _, tm := range tms {
			rows = append(rows, []interface{}{tm.Line, tm.In})
		}

		t := responseTable{
			Columns: []map[string]interface{}{
				{"text": "HST"},
				{"text": "min", "type": "int"},
			},
			Rows: rows,
			Typ:  "table",
		}

		rtls := []responseTable{t}
		json.NewEncoder(w).Encode(&rtls)
	})

	http.ListenAndServe(":8999", nil)
}
