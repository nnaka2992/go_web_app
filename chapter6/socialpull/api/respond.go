package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

// decodeBody encodes request to json.
// To modify Decode method to others such as binary,
// function can decode other format
func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// encodeBody encodes response to json.
// To modify Encode method to others such as binary,
// function can encode to other format
func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// respond write response with status code and contents.
func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		encodeBody(w, r, data)
	}
}

// func respondErr write response as an error.
func respondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {
	respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}

// respondHTTPErr wraps respondErr to return respondErr as HTTP 
func respondHTTPErr(w http.ResponseWriter, r *http.Request, status int) {
	respondErr(w, r, status, http.StatusText(status))
}
