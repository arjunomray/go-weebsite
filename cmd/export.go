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
	if err := writePage("index.html", "templates/home.html", PageData{Title: "Home"}); err != nil {
		return err
	}

	// 4. Render Blogs list
	if err := writePage("blogs/index.html", "templates/blogs.html", map[string]any{
		"Title":      "blogs",
		"Posts":      BlogPosts,
		"CurrentTag": "",
	}); err != nil {
		return err
	}

	// 5. Render individual Blog posts
	for _, post := range BlogPosts {
		postDest := filepath.Join("blog", post.Slug, "index.html")
		if err := writePage(postDest, "templates/post.html", map[string]any{
			"Title":    post.Title,
			"Date":     post.Date,
			"BodyHTML": template.HTML(post.BodyHTML),
		}); err != nil {
			return err
		}
	}

	// 6. Render Library list
	if err := writePage("library/index.html", "templates/library.html", map[string]any{
		"Title":      "library",
		"Media":      LibraryItems,
		"CurrentTag": "",
	}); err != nil {
		return err
	}

	// 7. Render individual Library items
	for _, item := range LibraryItems {
		itemDest := filepath.Join("library", item.Slug, "index.html")
		if err := writePage(itemDest, "templates/post.html", map[string]any{
			"Title":    item.Title,
			"Date":     item.Date,
			"BodyHTML": template.HTML(item.BodyHTML),
		}); err != nil {
			return err
		}
	}

	// 8. Render 404 page
	if err := writePage("404.html", "templates/404.html", PageData{Title: "404 Not Found"}); err != nil {
		return err
	}

	fmt.Printf("[Success] Static website exported to '%s'\n", outDir)
	return nil
}
