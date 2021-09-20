package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"chapter7/meander"
)

func main() {
	var (
		envfile = flag.String("env", "meander.env", "environment file's rellative path")
	)
	env, err := readEnv(*envfile)
	if err != nil {
		log.Fatal("Failed to read environment file: ", err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	meander.APIKey = env["key"]
	http.HandleFunc("/journeys", cors(func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	}))
	http.HandleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
		q := &meander.Query{
			Journey: strings.Split(r.URL.Query().Get("journey"), "|"),
		}
		q.Lat, _ = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		q.Lng, _ = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
		q.Radius, _ = strconv.Atoi(r.URL.Query().Get("radius"))
		q.CostRangeStr = r.URL.Query().Get("cost")
		places := q.Run()
		respond(w, r, places)

	}))
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}

func readEnv(filename string) (map[string]string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("Could not open file[%s]: %s", filename, err)
	}
	reader := bufio.NewReaderSize(f, 1024*2)

	env := make(map[string]string)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to read lines from file[%s]: %s", filename, err)
		}
		res := bytes.Split(line, []byte("="))
		env[string(res[0])] = string(res[1])
	}
	return env, nil
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}
