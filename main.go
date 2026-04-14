package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type FileInfo struct {
	Path string
	Size int64
	Ext  string
	Hash string
}

type ExtStats struct {
	Count int
	Size  int64
}

func collectFiles(root string) ([]string, error) {
	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // skips files we can't read
		}
		if info.IsDir() {
			return nil // skip directories
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close() // defer ensures the file is closed even if errors occur

	h := md5.New()                           // creates a new md5 hash
	if _, err := io.Copy(h, f); err != nil { // copies the file content to the hash
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil // returns the md5 hash as a string
}

func processFiles(paths []string) []FileInfo {
	results := make(chan FileInfo, len(paths)) // make is used for creating channels
	var wg sync.WaitGroup

	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done() // defer ensures the WaitGroup counter is decremented even if errors occur

			info, err := os.Stat(path)
			if err != nil {
				return
			}

			hash, _ := hashFile(path)

			results <- FileInfo{
				Path: path,
				Size: info.Size(),
				Ext:  filepath.Ext(path),
				Hash: hash,
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var files []FileInfo
	for f := range results {
		files = append(files, f)
	}
	return files
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
	)
	switch {
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func analyze(files []FileInfo) {
	var totalSize int64
	extMap := make(map[string]*ExtStats)
	hashMap := make(map[string][]FileInfo)

	for _, f := range files {
		totalSize += f.Size

		if _, ok := extMap[f.Ext]; !ok {
			extMap[f.Ext] = &ExtStats{}
		}
		extMap[f.Ext].Count++
		extMap[f.Ext].Size += f.Size

		if f.Hash != "" {
			hashMap[f.Hash] = append(hashMap[f.Hash], f)
		}
	}

	fmt.Printf("Found %d files (%s)\n\n", len(files), formatSize(totalSize))

	//sort extension
	type extEntry struct {
		Ext   string
		Stats *ExtStats
	}

	var exts []extEntry
	for ext, stats := range extMap {
		name := ext
		if name == "" {
			name = "(none)"
		}
		exts = append(exts, extEntry{name, stats})

	}
	sort.Slice(exts, func(i, j int) bool {
		return exts[i].Stats.Count > exts[j].Stats.Count
	})

	fmt.Println("By Extension:")
	for _, e := range exts {
		fmt.Printf("%-8s: %4d files %s\n", e.Ext, e.Stats.Count, formatSize(e.Stats.Size))
	}

	// find duplicate
	fmt.Println("\nDuplicates:")
	found := false
	for _, group := range hashMap {
		if len(group) < 2 {
			continue
		}
		found = true
		name := filepath.Base(group[0].Path)
		fmt.Printf(" %s (%d copies, %s each)\n", name, len(group), formatSize(group[0].Size))
	}

	if !found {
		fmt.Println(" No Duplicates found")
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: filescan <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	fmt.Printf("Scanning: %s\n", root)

	paths, err := collectFiles(root)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(paths) == 0 {
		fmt.Println("No Files found")
		return
	}

	files := processFiles(paths)
	analyze(files)

}
