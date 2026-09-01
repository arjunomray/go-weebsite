package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

// Helper to filter items by tag slice array matching parameters
func filterContent(items []ContentItem, tag string) []ContentItem {
	if tag == "" {
		return items
	}
	var filtered []ContentItem
	for _, item := range items {
		for _, t := range item.Tags {
			if t == strings.ToLower(tag) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}

func InitContentDatabase(dirPath string) ([]ContentItem, error) {
	var items []ContentItem

	files, err := filepath.Glob(filepath.Join(dirPath, "*.md"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Split file by "---" boundaries to isolate metadata from body text
		parts := strings.SplitN(string(data), "---", 3)
		if len(parts) < 3 {
			continue
		}

		metaLines := strings.Split(strings.TrimSpace(parts[1]), "\n")
		markdownBody := parts[2]

		item := ContentItem{
			Slug: strings.TrimSuffix(filepath.Base(file), ".md"),
		}

		// Simple key-value text flag parser
		for _, line := range metaLines {
			kv := strings.SplitN(line, ":", 2)
			if len(kv) != 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])

			switch key {
			case "title":
				item.Title = val
			case "summary":
				item.Summary = val
			case "stage":
				item.Stage = strings.TrimSpace(strings.ToLower(val))
			case "status":
				item.Status = strings.TrimSpace(val)
			case "date":
				item.Date, _ = time.Parse("2006-01-02", val)
			case "tags":
				for _, tag := range strings.Split(val, ",") {
					item.Tags = append(item.Tags, strings.TrimSpace(strings.ToLower(tag)))
				}
			}
		}

		// Only publish items explicitly marked as published
		if item.Stage != "published" {
			continue
		}

		// Transpile raw Markdown bytes directly to HTML string structures
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(markdownBody), &buf); err == nil {
			item.BodyHTML = buf.String()
		}

		// Compute reading time: ~200 words per minute
		wordCount := len(strings.Fields(markdownBody))
		item.ReadingTime = wordCount / 200
		if item.ReadingTime < 1 {
			item.ReadingTime = 1
		}

		items = append(items, item)
	}

	// Sort items chronologically (newest first)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Date.Before(items[j].Date) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	return items, nil
}
