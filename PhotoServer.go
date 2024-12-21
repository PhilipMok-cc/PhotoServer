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

type Config struct {
	PhotoDir string `yaml:"photoDir"`
	Port     string `yaml:"port"`
}

var conf Config

func main() {
	// Load configuration from the specified config file
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	// Use the configured photo directory
	photoDir := conf.PhotoDir
	if photoDir == "" {
		log.Fatal("Please specify 'photoDir' in config.yaml")
	}

	// Set up the HTTP server
	// Serve photos and static files
	http.Handle("/photos/", http.StripPrefix("/photos", http.FileServer(http.Dir(photoDir))))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(filepath.Join(".", "static")))))

	// Handle the index and display pages
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/display/", displayHandler)

	// Start the server
	log.Printf("Serving photos from: %s", photoDir)
	log.Printf("Listening on port: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Read all directories in the configured photo directory
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
			// Check if the directory contains a 1.jpg file
			thumbnailPath := filepath.Join(conf.PhotoDir, dir.Name(), "1.jpg")
			if _, err := os.Stat(thumbnailPath); err == nil {
				photos = append(photos, Photo{Name: dir.Name()})
			} else {
				log.Default().Printf("No thumbnail found for %s", dir.Name())
			}
		}
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, struct{ Photos []Photo }{Photos: photos})
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the directory name from the URL path
	dirName := strings.TrimPrefix(r.URL.Path, "/display/")
	dirName = strings.TrimSuffix(dirName, "/")

	// Get all files in the directory, using the configured photo directory
	files, err := os.ReadDir(filepath.Join(conf.PhotoDir, dirName))
	if err != nil {
		http.Error(w, "Directory not found", http.StatusNotFound)
		return
	}

	// Sort files by filename numerically
	sort.Slice(files, func(i, j int) bool {
		numI, _ := strconv.Atoi(filepath.Base(files[i].Name()[:len(files[i].Name())-4])) // Remove ".jpg" and convert to integer
		numJ, _ := strconv.Atoi(filepath.Base(files[j].Name()[:len(files[j].Name())-4])) // Remove ".jpg" and convert to integer

		return numI < numJ
	})

	// Prepare data for the template
	type Image struct {
		Index        int
		IndexPlusOne int
		FileName     string
	}
	var images []Image
	for i, file := range files {
		if file.IsDir() {
			continue
		}

		// Check if the file is a jpg
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			images = append(images, Image{Index: i, IndexPlusOne: i + 1, FileName: file.Name()})
		}
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/display.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, struct {
		AlbumName   string
		TotalImages int
		Images      []Image
	}{AlbumName: dirName, TotalImages: len(images), Images: images})
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
