package main

import (
	"os"
	"log"
	"flag"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/nnaka2992/go_web_app/chapter1/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"

	/*
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	*/
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "Application Port")
	flag.Parse()

	_ = godotenv.Load("test.env")

	gomniauth.SetSecurityKey(os.Getenv("SECURITY_KEY"))
	gomniauth.WithProviders(
		facebook.New(os.Getenv("FACEBOOK_KEY"),	os.Getenv("FACEBOOK_SECRETKEY"), "http://localhost:8080/auth/callback/facebook"),
		github.New(os.Getenv("GITHUB_KEY"),	os.Getenv("GITHUB_SECRETKEY"), "http://localhost:8080/auth/callback/github"),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRETKEY"), "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	if err:= http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

