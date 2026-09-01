package main

import (
	"bytes"
	"fmt"
	"go-weebsite/web"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

func (app *App) ExportStatic(outDir string) error {
	// 1. Ensure output directory exists
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	// 2. Copy embedded static assets (web/static -> dist/static)
	staticFS, err := fs.Sub(web.WebFS, "static")
	if err != nil {
		return err
	}

	err = fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		targetPath := filepath.Join(outDir, "static", path)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}
		data, err := fs.ReadFile(staticFS, path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, 0644)
	})
	if err != nil {
		return fmt.Errorf("copy static assets: %w", err)
	}

	// Helper to render and write an HTML page
	writePage := func(pagePath string, tmplFile string, data any) error {
		tmpl, err := app.ParsePage(tmplFile)
		if err != nil {
			return fmt.Errorf("parse %s: %w", tmplFile, err)
		}
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "base.html", data); err != nil {
			return fmt.Errorf("render %s: %w", tmplFile, err)
		}
		dest := filepath.Join(outDir, pagePath)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		return os.WriteFile(dest, buf.Bytes(), 0644)
	}

	// 3. Render Home page
	if err := writePage("index.html", "templates/home.html", PageData{
		Title:       "Home",
		Description: "Software developer from India. I build things, read books, and meditate.",
		Page:        "home",
	}); err != nil {
		return err
	}

	// 4. Render Blogs list
	if err := writePage("blogs/index.html", "templates/blogs.html", map[string]any{
		"Title":       "blogs",
		"Description": "Thoughts, notes and writing from Arjun.",
		"Page":        "blogs",
		"Posts":       BlogPosts,
		"CurrentTag":  "",
	}); err != nil {
		return err
	}

	// 5. Render individual Blog posts
	for _, post := range BlogPosts {
		postDest := filepath.Join("blog", post.Slug, "index.html")
		if err := writePage(postDest, "templates/post.html", map[string]any{
			"Title":       post.Title,
			"Description": post.Summary,
			"Date":        post.Date,
			"ReadingTime": post.ReadingTime,
			"Page":        "blogs",
			"BodyHTML":    template.HTML(post.BodyHTML),
		}); err != nil {
			return err
		}
	}

	// 6. Render Library list
	if err := writePage("library/index.html", "templates/library.html", map[string]any{
		"Title":       "library",
		"Description": "Books and media that shaped my thinking.",
		"Page":        "library",
		"Media":       LibraryItems,
		"CurrentTag":  "",
	}); err != nil {
		return err
	}

	// 7. Render individual Library items
	for _, item := range LibraryItems {
		itemDest := filepath.Join("library", item.Slug, "index.html")
		if err := writePage(itemDest, "templates/post.html", map[string]any{
			"Title":       item.Title,
			"Description": item.Summary,
			"Date":        item.Date,
			"ReadingTime": item.ReadingTime,
			"Page":        "library",
			"BodyHTML":    template.HTML(item.BodyHTML),
		}); err != nil {
			return err
		}
	}

	// 8. Render 404 page
	if err := writePage("404.html", "templates/404.html", PageData{
		Title:       "404 Not Found",
		Description: "Page not found",
		Page:        "",
	}); err != nil {
		return err
	}

	// 9. Write Cloudflare Pages _redirects for clean URLs (no trailing slash needed)
	redirects := "/blogs /blogs/index.html 200\n"
	redirects += "/library /library/index.html 200\n"
	for _, post := range BlogPosts {
		redirects += fmt.Sprintf("/blog/%s /blog/%s/index.html 200\n", post.Slug, post.Slug)
	}
	for _, item := range LibraryItems {
		redirects += fmt.Sprintf("/library/%s /library/%s/index.html 200\n", item.Slug, item.Slug)
	}
	redirects += "/* /404.html 404\n"
	if err := os.WriteFile(filepath.Join(outDir, "_redirects"), []byte(redirects), 0644); err != nil {
		return fmt.Errorf("write _redirects: %w", err)
	}

	fmt.Printf("[Success] Static website exported to '%s'\n", outDir)
	return nil
}
