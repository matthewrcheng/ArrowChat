package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// handler that serves a specified HTML template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// handles HTTP requests by rendering the specified HTML template
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse the template file once
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// execute parsed template, writing output to response writer
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "Addr of Application.")
	flag.Parse()

	// create a new room
	r := newRoom()

	// serve the home page and chat room
	http.Handle("/", &templateHandler{filename: "home.html"})
	http.Handle("/chat", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	// start the web server on the specified address
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
