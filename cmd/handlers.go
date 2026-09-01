package main

import (
	"fmt"
	"go-weebsite/web"
	"html/template"
	"net/http"
	"strings"
)

func (app *App) ParsePage(pageFile string) (*template.Template, error) {
	return template.ParseFS(web.WebFS,
		"templates/base.html",
		pageFile,
		"templates/components/*.html",
	)
}
func (app *App) Render(w http.ResponseWriter, pageFile string, data any) {
	// 1. Dynamically build the template tree for this specific page request
	tmpl, err := app.ParsePage(pageFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template compilation error: %s", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Render error: %s", err), http.StatusInternalServerError)
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("[HTTP] %s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func (app *App) Custom404(next *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := next.Handler(r)
		if pattern == "" || (pattern == "GET /" && r.URL.Path != "/") {
			w.WriteHeader(http.StatusNotFound)
			app.Render(w, "templates/404.html", PageData{Title: "404 Not Found"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title:       "Home",
		Description: "Software developer from India. I build things, read books, and meditate.",
		Page:        "home",
	}
	app.Render(w, "templates/home.html", pageData)
}

func (app *App) WorkHandler(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title:       "Work",
		Description: "Career log and professional engineering experience.",
		Page:        "work",
	}
	app.Render(w, "templates/work.html", pageData)
}

// BLOG ROUTING
func (app *App) BlogListHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	app.Render(w, "templates/blogs.html", map[string]any{
		"Title": "blogs", "Description": "Thoughts, notes and writing from Arjun.",
		"Page": "blogs", "Posts": filterContent(BlogPosts, tag), "CurrentTag": tag,
	})
}

// LIBRARY ROUTING
func (app *App) LibraryListHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	app.Render(w, "templates/library.html", map[string]any{
		"Title": "library", "Description": "Books and media that shaped my thinking.",
		"Page": "library", "Media": filterContent(LibraryItems, tag), "CurrentTag": tag,
	})
}

func (app *App) BlogPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		slug = strings.TrimPrefix(r.URL.Path, "/blog/")
	}
	for _, post := range BlogPosts {
		if post.Slug == slug {
			// Cast post.BodyHTML to template.HTML so Go knows it's safe execution data
			app.Render(w, "templates/post.html", map[string]any{
				"Title":       post.Title,
				"Description": post.Summary,
				"Date":        post.Date,
				"ReadingTime": post.ReadingTime,
				"Page":        "blogs",
				"BodyHTML":    template.HTML(post.BodyHTML),
			})
			return
		}
	}
	http.NotFound(w, r)
}

func (app *App) LibraryPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		slug = strings.TrimPrefix(r.URL.Path, "/library/")
	}
	for _, item := range LibraryItems {
		if item.Slug == slug {
			// Cast item.BodyHTML to template.HTML here too
			app.Render(w, "templates/post.html", map[string]any{
				"Title":       item.Title,
				"Description": item.Summary,
				"Date":        item.Date,
				"ReadingTime": item.ReadingTime,
				"Page":        "library",
				"BodyHTML":    template.HTML(item.BodyHTML),
			})
			return
		}
	}
	http.NotFound(w, r)
}
