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

const (
	DIR_TEMPLATES = "tmpl"
	DIR_PAGES     = "pages"
	HTML_BASE     = "base.html"
	HTML_INDEX    = "index.html"
	HTML_PAGE     = "page.html"
)

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
		err := t.ExecuteTemplate(w, HTML_BASE, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func indexHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Index handler called\n")
		pages := make([]string, 0)
		filepath.Walk(DIR_PAGES, makeIndexer(&pages))
		err := t.ExecuteTemplate(w, HTML_BASE, pages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func makeTemplate(name string) *template.Template {
	tmpl_name := path.Join(DIR_TEMPLATES, name)
	tmpl_base := path.Join(DIR_TEMPLATES, HTML_BASE)
	return template.Must(template.ParseFiles(tmpl_name, tmpl_base))
}

func main() {
	cache := core.EmptyMemoryCache()

	indexTemplate := makeTemplate(HTML_INDEX)
	pageTemplate := makeTemplate(HTML_PAGE)

	http.HandleFunc("/page/", pageHandler(pageTemplate, cache))
	http.HandleFunc("/", indexHandler(indexTemplate))

	http.ListenAndServe(":8888", nil)
}
