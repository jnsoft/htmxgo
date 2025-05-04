package server

import (
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

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleHome(w, r, templates, &data)
	})

	return mux

}

func HandleHome(w http.ResponseWriter, r *http.Request, templates *Templates, data *models.IndexModel) {

	data.Count++
	data.CurrentDate = time.Now().Format("2006-01-02")
	data.CurrentTime = time.Now().Format("15:04:05")

	if err := templates.Render(w, "index.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
	}
}
