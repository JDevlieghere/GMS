package main

import (
	"GMS/core"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const TEMPLATES_DIR = "tmpl"
const PAGES_DIR = "pages"

func makeIndexer(pages *[]string) func(path string, f os.FileInfo, err error) error {
	return func(path string, f os.FileInfo, err error) error {
		name := f.Name()
		validFile := regexp.MustCompile("^([a-zA-Z0-9]+).(md|MD)$")
		ext := validFile.FindStringSubmatch(name)
		if ext != nil {
			log.Printf("Found page: %v\n", ext[1])
			*pages = append(*pages, ext[1])
		}
		return nil
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	validPath := regexp.MustCompile("^(/page|/)(/[a-zA-Z0-9]+)?$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	log.Printf("Page handler called with parameters: %v\n", m)
	if m == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	p := core.Page{}
	err := t.ExecuteTemplate(w, "base.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	log.Printf("Index handler called\n")
	pages := make([]string, 0)
	filepath.Walk(PAGES_DIR, makeIndexer(&pages))
	err := t.ExecuteTemplate(w, "base.html", pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *template.Template), t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, t)
	}
}

func makeTemplate(name string, base string) *template.Template {
	tmpl_name := path.Join(TEMPLATES_DIR, name)
	tmpl_base := path.Join(TEMPLATES_DIR, base)
	return template.Must(template.ParseFiles(tmpl_name, tmpl_base))
}

func main() {
	indexTemplate := makeTemplate("index.html", "base.html")
	pageTemplate := makeTemplate("page.html", "base.html")

	http.HandleFunc("/page/", makeHandler(pageHandler, pageTemplate))
	http.HandleFunc("/", makeHandler(indexHandler, indexTemplate))

	http.ListenAndServe(":8888", nil)
}
