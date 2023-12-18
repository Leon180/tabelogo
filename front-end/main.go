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

	config, err := LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// handleize function: render
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderFrontEnd(w, "main.gohtml", config)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		renderFrontEnd(w, "login.gohtml", config)
	})
	http.HandleFunc("/regist", func(w http.ResponseWriter, r *http.Request) {
		renderFrontEnd(w, "regist.gohtml", config)
	})
	// build connection and start listen...
	fmt.Println("Starting front end server and listen on port: " + frontEndPort)
	err = http.ListenAndServe(frontEndPort, nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templatesFS embed.FS

func renderFrontEnd(w http.ResponseWriter, t string, config Config) {

	partials, err := ReadFilesName("templates/partials")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFS(templatesFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerURL  string
		WebsiteURL string
	}

	data.BrokerURL = config.BrokerURLDeployment
	data.WebsiteURL = config.WebsiteURLDeployment

	if err := tmpl.Execute(w, data); err != nil {
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
