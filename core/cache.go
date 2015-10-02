package core

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

const (
	PAGES_DIR = "pages"
)

type Cache interface {
	GetPage(slug string) *Page
}

type MemoryCache struct {
	Pages map[string]*Page
}

func EmptyMemoryCache() *MemoryCache {
	return &MemoryCache{Pages: make(map[string]*Page)}
}

func (cache *MemoryCache) GetPage(slug string) *Page {
	// Try Cache
	page, ok := cache.Pages[slug]
	if ok {
		log.Printf("Loading page from MemoryCache: %v\n", slug)
		return page
	}

	// Load Page
	page, err := loadPage(PAGES_DIR, slug)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Cache Page
	log.Printf("Storing page in MemoryCache: %v\n", slug)
	cache.Pages[slug] = page

	return page
}

func loadPage(root string, slug string) (*Page, error) {

	log.Printf("Loading page from file: %v\n", slug)

	filename := slug + ".md"
	filepath := path.Join(root, filename)

	// Read File Content
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Read File Modification Date
	finfo, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	// Create Page from File
	body := template.HTML(blackfriday.MarkdownCommon(content))
	timestamp := finfo.ModTime().Format(time.ANSIC)
	page := &Page{Title: slug, Body: body, Timestamp: timestamp}

	return page, nil
}
