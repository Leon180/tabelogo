package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	frontEndPort = ":8081"
)

func main() {
	// handleize function: render
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderFrontEnd(w, "content.gohtml")
	})

	// build connection and start listen...
	fmt.Println("Starting front end server and listen on port: " + frontEndPort)
	err := http.ListenAndServe(frontEndPort, nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templatesFS embed.FS

func renderFrontEnd(w http.ResponseWriter, t string) {

	partials, err := ReadFilesName("templates")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templatesFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// var data struct {
	// 	BrokerURL string
	// }

	// data.BrokerURL = os.Getenv("BROKER_URL")
	// data.BrokerURL = "http://localhost:8080"

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ReadFilesName(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, path+"/"+file.Name())
	}

	return fileNames, nil
}
