package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/achristie/go-blueprints/meander"
)

func main() {
	meander.APIKey = os.Getenv("PLACES_APIKEY")
	http.HandleFunc("/journeys", cors(func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	}))

	http.HandleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
		q := &meander.Query{
			Journey: strings.Split(r.URL.Query().Get("journey"), "|"),
		}
		var err error
		q.Lat, err = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q.Lng, err = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q.Radius, err = strconv.Atoi(r.URL.Query().Get("radius"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q.CostRangeStr = r.URL.Query().Get("cost")
		places := q.Run()
		respond(w, r, places)
	}))
	http.ListenAndServe(":8081", http.DefaultServeMux)
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Origin", "*")
		f(w, r)
	}
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}
