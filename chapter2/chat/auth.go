package main

import (
	"net/http"
	"strings"
	"log"
	"fmt"

	// "github.com/nnaka2992/go_web_app/chapter1/trace"
	"github.com/stretchr/gomniauth"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// Not Autholized
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// Unhandled Error
		panic(err.Error())
	} else {
		// Autholized, Call wrapped handler
		h.next.ServeHTTP(w, r)
	}
}
// Helper to use authHandler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler wait third party login process
// Path must be like: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider);
		if err != nil {
			log.Fatalln("Faled to get Authentification Provider: ", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil);
		if err != nil {
			log.Fatalln("Faled to call GetBeginAuthURL: ", provider, "-", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Action %s is not supported", action)
	}
}
