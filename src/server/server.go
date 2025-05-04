package server

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jnsoft/htmxgo/src/models"
)

type Templates struct {
	templates *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Templates) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("../static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	templates := NewTemplates()

	data := models.IndexModel{
		Count:       0,
		CurrentDate: time.Now().Format("2006-01-02"),
		CurrentTime: time.Now().Format("15:04:05"),
	}

	contacts := models.NewData()

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleHome(w, r, templates, &data)
	})

	mux.HandleFunc("/increment", func(w http.ResponseWriter, r *http.Request) {
		HandleIncrement(w, r, templates, &data)
	})

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		HandleContacts(w, r, templates, &data)
	})

	return mux

}

func HandleHome(w http.ResponseWriter, r *http.Request, templates *Templates, data *models.IndexModel) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err := templates.Render(w, "index.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
	}

}

func HandleIncrement(w http.ResponseWriter, r *http.Request, templates *Templates, data *models.IndexModel) {
	if r.Method == http.MethodPost {
		data.Count++
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", data.Count)))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleContacts(w http.ResponseWriter, r *http.Request, templates *Templates, data *models.IndexModel) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err := templates.Render(w, "contacts.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
	}
}
