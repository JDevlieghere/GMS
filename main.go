package main

import (
	"flag"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

const (
	PAGES_DIR     = "pages"
	TEMPLATES_DIR = "tmpl"
)

type (
	Page struct {
		Title     string
		Body      template.HTML
		Timestamp string
	}
)

var (
	validPath = regexp.MustCompile("^(/page|/)(/[a-zA-Z0-9]+)?$")
	validFile = regexp.MustCompile("^([a-zA-Z0-9]+).(md|MD)$")
	addr      = flag.Bool("addr", false, "find open address and print to final-port.txt")
	tmpl      = make(map[string]*template.Template)
	cache     = make(map[string]*Page)
	pages     = make([]string, 0)
)

func loadPage(slug string) (*Page, error) {
	filename := slug + ".md"
	filepath := path.Join(PAGES_DIR, filename)
	// File Content
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	// File Modification Date
	finfo, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	// Create Page from File
	body := template.HTML(blackfriday.MarkdownCommon(content))
	timestamp := finfo.ModTime().Format(time.ANSIC)
	page := &Page{Title: slug, Body: body, Timestamp: timestamp}
	cache[slug] = page
	return page, nil
}

func fetchPage(slug string) (*Page, error) {
	page, ok := cache[slug]
	if !ok {
		return loadPage(slug)
	}
	return page, nil
}

func pageHandler(w http.ResponseWriter, r *http.Request, param []string) {
	p, err := fetchPage(param[2])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = tmpl["view.html"].ExecuteTemplate(w, "base.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexPage(path string, f os.FileInfo, err error) error {
	name := f.Name()
	ext := validFile.FindStringSubmatch(name)
	if ext != nil {
		pages = append(pages, ext[1])
	}
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request, param []string) {
	pages = make([]string, 0)
	filepath.Walk(PAGES_DIR, indexPage)
	err := tmpl["index.html"].ExecuteTemplate(w, "base.html", pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, []string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		fn(w, r, m)
	}
}

func makeTemplate(name string, base string) {
	tmpl_name := path.Join(TEMPLATES_DIR, name)
	tmpl_base := path.Join(TEMPLATES_DIR, base)
	tmpl[name] = template.Must(template.ParseFiles(tmpl_name, tmpl_base))
}

func main() {
	makeTemplate("index.html", "base.html")
	makeTemplate("view.html", "base.html")

	flag.Parse()
	http.HandleFunc("/page/", makeHandler(pageHandler))
	http.HandleFunc("/", makeHandler(indexHandler))

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":80", nil)
}
