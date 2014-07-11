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
		Timestamp string
	}
	Stat struct {
		Cached int
	}
)

var (
	validPath = regexp.MustCompile("^/(page|stat)/([a-zA-Z0-9]+)?$")
	addr      = flag.Bool("addr", false, "find open address and print to final-port.txt")
	templates = template.Must(template.ParseFiles(
		path.Join(TEMPLATES_DIR, "view.html"),
		path.Join(TEMPLATES_DIR, "stat.html")))
	pages = make(map[string]*Page)
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
	page := &Page{Body: body, Timestamp: timestamp}
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

func pageHandler(w http.ResponseWriter, r *http.Request, param []string) {
	p, err := fetchPage(param[2])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderPage(w, "view", p)
}

func statHandler(w http.ResponseWriter, r *http.Request, param []string) {
	err := templates.ExecuteTemplate(w, "stat.html", &Stat{Cached: len(pages)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, []string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/page/", makeHandler(pageHandler))
	http.HandleFunc("/stat/", makeHandler(statHandler))

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
