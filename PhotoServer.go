package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for the server
type Config struct {
	PhotoDir string `yaml:"photoDir"`
	Port     string `yaml:"port"`
}

var conf Config

// Supported image file extensions
var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// Define the Chapter type globally
type Chapter struct {
	ID     string
	Images []Image
}

// Define the Image type globally
type Image struct {
	Index        int
	IndexPlusOne int
	FileName     string
}

func main() {
	// Load configuration
	loadConfig()

	// Set up the HTTP server
	setupServer()

	// Start the server
	log.Printf("Serving photos from: %s", conf.PhotoDir)
	log.Printf("Listening on port: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}

// loadConfig reads and parses the configuration file
func loadConfig() {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	if conf.PhotoDir == "" {
		log.Fatal("Please specify 'photoDir' in config.yaml")
	}
}

// setupServer configures the HTTP routes and handlers
func setupServer() {
	http.Handle("/photos/", http.StripPrefix("/photos", http.FileServer(http.Dir(conf.PhotoDir))))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(filepath.Join(".", "static")))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/display/", displayHandler)
}

// indexHandler handles the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	dirs, err := os.ReadDir(conf.PhotoDir)
	if err != nil {
		http.Error(w, "Directory not found", http.StatusNotFound)
		return
	}

	// Sort directories by name
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	// Prepare data for the template
	type Photo struct {
		Name string
	}
	var photos []Photo
	for _, dir := range dirs {
		if dir.IsDir() {
			thumbnailPath := filepath.Join(conf.PhotoDir, dir.Name(), "1.jpg")
			if _, err := os.Stat(thumbnailPath); err == nil {
				photos = append(photos, Photo{Name: dir.Name()})
			} else {
				log.Printf("No thumbnail found for %s", dir.Name())
			}
		}
	}

	// Render the template
	renderTemplate(w, "templates/index.html", struct{ Photos []Photo }{Photos: photos})
}

// displayHandler handles the display page
func displayHandler(w http.ResponseWriter, r *http.Request) {
	dirName := strings.Trim(strings.TrimPrefix(r.URL.Path, "/display/"), "/")
	files, err := os.ReadDir(filepath.Join(conf.PhotoDir, dirName))
	if err != nil {
		http.Error(w, "Directory not found", http.StatusNotFound)
		return
	}

	// Sort files with .txt files prioritized
	sortFiles(files)

	// Prepare chapters and images
	chapters, totalImages := prepareChapters(dirName, files)

	// Render the template
	renderTemplate(w, "templates/display.html", struct {
		AlbumName   string
		Chapters    []Chapter
		TotalImages int
	}{AlbumName: dirName, Chapters: chapters, TotalImages: totalImages})
}

// sortFiles sorts files numerically or alphabetically, prioritizing .txt files
func sortFiles(files []os.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		nameI, nameJ := files[i].Name(), files[j].Name()
		extI := strings.ToLower(filepath.Ext(nameI)) // Remove unused extJ

		// Extract numeric prefixes if present
		numI, errI := strconv.Atoi(strings.TrimSuffix(nameI, filepath.Ext(nameI)))
		numJ, errJ := strconv.Atoi(strings.TrimSuffix(nameJ, filepath.Ext(nameJ)))

		// Prioritize .txt files if numeric prefixes are the same
		if errI == nil && errJ == nil && numI == numJ {
			return extI == ".txt"
		}

		// Compare numerically if both have numeric prefixes
		if errI == nil && errJ == nil {
			return numI < numJ
		}

		// Prioritize .txt files if names are the same
		if nameI == nameJ {
			return extI == ".txt"
		}

		// Default to alphabetical comparison
		return nameI < nameJ
	})
}

// prepareChapters organizes files into chapters and calculates total images, including subfolders as chapters
func prepareChapters(dirName string, files []os.DirEntry) ([]Chapter, int) {
	var chapters []Chapter
	var currentChapter Chapter
	var subfolders []os.DirEntry // Array to store subfolders
	imageIndex := 0
	firstChapterCreated := false

	// Process files and collect subfolders
	for _, file := range files {
		if file.IsDir() {
			subfolders = append(subfolders, file) // Collect subfolders for the next loop
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".txt" {
			if len(currentChapter.Images) > 0 {
				chapters = append(chapters, currentChapter)
			}
			currentChapter = createChapterFromTxt(dirName, file)
			firstChapterCreated = true
		} else if imageExtensions[ext] {
			if !firstChapterCreated {
				currentChapter = Chapter{ID: "1"}
				firstChapterCreated = true
			}
			currentChapter.Images = append(currentChapter.Images, Image{
				Index:        imageIndex,
				IndexPlusOne: imageIndex + 1,
				FileName:     file.Name(),
			})
			imageIndex++
		}
	}

	if len(currentChapter.Images) > 0 {
		chapters = append(chapters, currentChapter)
	}

	// Process subfolders
	for _, folder := range subfolders {
		subDirName := filepath.Join(dirName, folder.Name())
		subFiles, err := os.ReadDir(filepath.Join(conf.PhotoDir, subDirName))
		if err != nil {
			log.Printf("Error reading subfolder %s: %v", subDirName, err)
			continue
		}

		// Create a chapter for the subfolder
		subChapter := Chapter{ID: folder.Name()}
		for _, subFile := range subFiles {
			ext := strings.ToLower(filepath.Ext(subFile.Name()))
			if imageExtensions[ext] {
				subChapter.Images = append(subChapter.Images, Image{
					Index:        imageIndex,
					IndexPlusOne: imageIndex + 1,
					FileName:     filepath.Join(folder.Name(), subFile.Name()),
				})
				imageIndex++
			}
		}

		// Add the subfolder chapter to the list
		if len(subChapter.Images) > 0 {
			chapters = append(chapters, subChapter)
		}
	}

	return chapters, imageIndex
}

// createChapterFromTxt creates a chapter from a .txt file
func createChapterFromTxt(dirName string, file os.DirEntry) Chapter {
	txtFilePath := filepath.Join(conf.PhotoDir, dirName, file.Name())
	content, err := os.ReadFile(txtFilePath)
	if err != nil {
		log.Printf("Error reading chapter file %s: %v", txtFilePath, err)
		return Chapter{ID: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))}
	}
	return Chapter{ID: strings.TrimSpace(string(content))}
}

// renderTemplate parses and executes a template
func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}
