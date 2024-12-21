package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func TestServerUp(t *testing.T) {
	done := make(chan bool)
	// Start the RunServer function in a separate goroutine
	go func() {
		main()
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)
	open("http://localhost:8080")

	// Wait for manual confirmation
	<-done
}

func TestCSSMimeType(t *testing.T) {
	// Set the flag for the config file

	// Start the RunServer function in a separate goroutine
	go func() {
		main()
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	// Make an actual HTTP request to the server
	resp, err := http.Get("http://localhost:8080/static/css/styles.css")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "text/css"
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, expected) {
		t.Errorf("handler returned wrong content type: got [%v] want [%v]", contentType, expected)
	}
}

func TestDisplayHandler(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	conf.PhotoDir = tempDir

	// Create a test photo directory with some images
	photoDir := filepath.Join(tempDir, "test_album")
	err := os.Mkdir(photoDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create some test image files
	for i := 1; i <= 3; i++ {
		file, err := os.Create(filepath.Join(photoDir, strconv.Itoa(i)+".jpg"))
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	// Create a request to the display handler
	req, err := http.NewRequest("GET", "/display/test_album", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(displayHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body for expected content
	expected := "Image 1 of 3"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
