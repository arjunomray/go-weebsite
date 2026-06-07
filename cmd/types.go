package main

import (
	"html/template"
	"time"
)

type ContentItem struct {
	Slug     string // URL path tracker (e.g., "my-first-post")
	Title    string
	Date     time.Time
	Tags     []string // For tag filtering
	Summary  string   // Brief card snippet
	BodyHTML string   // The parsed HTML block
	Stage    string
}

// Global memory stores for blazing fast routing
var BlogPosts []ContentItem
var LibraryItems []ContentItem

// InitContentDatabase scans a folder and compiles all markdown files

type PageData struct {
	Title string
}

type App struct {
	templates *template.Template
}
