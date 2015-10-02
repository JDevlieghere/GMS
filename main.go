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
	"strings"
)

const TEMPLATES_DIR = "tmpl"
const PAGES_DIR = "pages"

func makeIndexer(pages *[]string) func(path string, f os.FileInfo, err error) error {
	return func(path string, f os.FileInfo, err error) error {
		name := f.Name()
		validFile := regexp.MustCompile("^([a-zA-Z0-9]+).(md|MD)$")
		ext := validFile.FindStringSubmatch(name)
		if ext != nil {
			log.Printf("Indexed page: %v\n", ext[1])
			*pages = append(*pages, ext[1])
		}
		return nil
	}
}

func pageHandler(t *template.Template, c core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile("^(/page|/)(/[a-zA-Z0-9]+)?$")
		params := validPath.FindStringSubmatch(r.URL.Path)
		if params == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		slug := strings.TrimLeft(params[2], "/")
		log.Printf("Page handler called with slug: %v\n", slug)
		p := c.GetPage(slug)
		if p == nil {
			http.NotFound(w, r)
			return
		}
		err := t.ExecuteTemplate(w, "base.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func indexHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Index handler called\n")
		pages := make([]string, 0)
		filepath.Walk(PAGES_DIR, makeIndexer(&pages))
		err := t.ExecuteTemplate(w, "base.html", pages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func makeTemplate(name string, base string) *template.Template {
	tmpl_name := path.Join(TEMPLATES_DIR, name)
	tmpl_base := path.Join(TEMPLATES_DIR, base)
	return template.Must(template.ParseFiles(tmpl_name, tmpl_base))
}

func main() {
	cache := core.EmptyMemoryCache()

	indexTemplate := makeTemplate("index.html", "base.html")
	pageTemplate := makeTemplate("page.html", "base.html")

	http.HandleFunc("/page/", pageHandler(pageTemplate, cache))
	http.HandleFunc("/", indexHandler(indexTemplate))

	http.ListenAndServe(":8888", nil)
}
