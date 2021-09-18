package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "Address of endpoint")
		mongo = flag.String("mongo", "localhost", "Address of MongoDB")
	)
	flag.Parse()
	log.Println("Connect to MongoDB", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("Failed to connect MongoDB: ", err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))

	log.Println("Start hosting Web Server", *addr)
	http.ListenAndServe(*addr, mux)
	log.Println("Stopping Web Server...")
}

// withAPIKey returns http.Handlefunc which is its argument
// when the api key is correct. Otherwise, the function returns
// statusunauthorized error.
func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "API Key is invalid")
			return
		}
		fn(w, r)
	}
}

// isValidapikey validate whther api key is correct or not
func isValidAPIKey(key string) bool {
	return key == "abc123"
}

// withData sets the contents of "ballots" to vars map.
// then it return http.Handlefunc which is passed as its argument.
func withData(d *mgo.Session, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDB := d.Copy()
		defer thisDB.Close()
		SetVar(r, "db", thisDB.DB("ballots"))
		f(w, r)
	}
}

// withVars manages context of vars map's context.
func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

// withCORS add header to ResponseWriter to allows clients
// to access API
// You have better use "https://github.com/fasterness/cors" for
// production
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
