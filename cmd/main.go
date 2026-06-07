package main

import (
	"context"
	"fmt"
	"go-weebsite/web"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	errorChan := make(chan error)
	signalChan := make(chan os.Signal, 1)
	staticFS, err := fs.Sub(web.WebFS, "static")
	BlogPosts, err = InitContentDatabase("content/blog")
	if err != nil {
		fmt.Printf("[Warning] Could not initialize blog posts database: %v\n", err)
	} else {
		fmt.Printf("[Success] Loaded %d active blog posts into memory cache\n", len(BlogPosts))
	}

	LibraryItems, err = InitContentDatabase("content/library")
	if err != nil {
		fmt.Printf("[Warning] Could not initialize library database: %v\n", err)
	} else {
		fmt.Printf("[Success] Loaded %d active library items into memory cache\n", len(LibraryItems))
	}

	if err != nil {
		log.Fatal("Error Unable to load templates: ", err)
	}
	// templatesFS, err := fs.Sub(web.WebFS, "web/templates")

	if err != nil {
		log.Fatal("Error Unable to load templates: ", err)
	}
	t, err := template.ParseFS(web.WebFS, "templates/*.html", "templates/components/*.html")
	if err != nil {
		log.Fatal("Error Unable to load templates: ", err)
	}
	app := App{
		templates: t,
	}

	handler := http.ServeMux{}
	wrapped := LoggingMiddleware(app.Custom404(&handler))

	handler.Handle("GET /static/{s...}", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))
	handler.HandleFunc("GET /health", app.HealthHandler)
	handler.HandleFunc("GET /", app.HomeHandler)
	handler.HandleFunc("GET /blogs", app.BlogListHandler)
	handler.HandleFunc("GET /blog/", app.BlogPostHandler) // Matches /blog/<slug>

	// Register Library Paths
	handler.HandleFunc("GET /library", app.LibraryListHandler)
	handler.HandleFunc("GET /library/", app.LibraryPostHandler)

	s := &http.Server{
		Addr:           ":8000",
		Handler:        wrapped,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			} else {
				errorChan <- err
			}
		}
	}()

	signal.Notify(signalChan, os.Interrupt)

	select {
	case err := <-errorChan:
		log.Fatal("Server Error", err)
	case <-signalChan:
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatal("Failed to close the server gracefully: ", err)
		}
	}
}
