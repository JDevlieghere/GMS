package core

import (
	"html/template"
)

type Page struct {
	Title     string
	Body      template.HTML
	Timestamp string
}

type FilePage struct {
	Title     string
	Filename  string
	Timestamp string
}
