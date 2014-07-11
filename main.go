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
	"regexp"
	"time"
)

const (
	PAGES_DIR     = "pages"
	TEMPLATES_DIR = "tmpl"
)

type (
	Page struct {
		Body      template.HTML
		Timestamp time.Time
	}
)

var (
	validPath = regexp.MustCompile("^/(page|p)/([a-zA-Z0-9]+)$")
	addr      = flag.Bool("addr", false, "find open address and print to final-port.txt")
	templates = template.Must(template.ParseFiles(path.Join(TEMPLATES_DIR, "view.html")))
	pages     = make(map[string]*Page)
)

func loadPage(slug string) (*Page, error) {
	filename := slug + ".md"
	filepath := path.Join(PAGES_DIR, filename)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	finfo, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	body := template.HTML(blackfriday.MarkdownCommon(content))
	page := &Page{Body: body, Timestamp: finfo.ModTime()}
	pages[slug] = page
	return page, nil
}

func fetchPage(slug string) (*Page, error) {
	page, ok := pages[slug]
	if !ok {
		return loadPage(slug)
	}
	return page, nil
}

func pageHandler(w http.ResponseWriter, r *http.Request, slug string) {
	p, err := fetchPage(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "view", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/page/", makeHandler(pageHandler))

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
