package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
	"path/filepath"
	"time"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	folder := "/home/sebin/Documents/projects/message-broker/rmq/static/"
	folder = fmt.Sprintf("%s%s/", folder, currentTime.Format("2006-01-02"))
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// folder 
    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse input form
    err = r.ParseMultipartForm(10 << 20) // Limit file size
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Handle file uploads
    for _, key := range []string{"video", "image1", "image2",} {
        file, header, err := r.FormFile(key)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error retrieving the file for key %s: %v", key, err), http.StatusInternalServerError)
            return
        }
        defer file.Close()
	ext := filepath.Ext(header.Filename)

        // Save the file
        fileBytes, err := ioutil.ReadAll(file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
	err = ioutil.WriteFile(fmt.Sprintf("%s%s%s", folder,key,ext), fileBytes, 0644)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // Handle bounding box array (pseudo-code, implement according to your needs)
    bbox := r.FormValue("bbox")
	bboxFileName := fmt.Sprintf("%s%s", folder,"bbox.txt")
	err = ioutil.WriteFile(bboxFileName, []byte(bbox), 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


    // Process bbox...

    fmt.Fprintf(w, "Upload successful")
}

func main() {
    http.HandleFunc("/upload", uploadHandler)
    http.ListenAndServe(":8080", nil)
}
