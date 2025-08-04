package main

import (
	"fmt"
	"log"
	"net/http"

	"cliapp/internal/utils"
	"cliapp/pkg/imagebuilder"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/build", buildHandler)

	log.Println("Starting web UI on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Builder CI</title></head>
<body>
<h1>Builder CI</h1>
<form action="/build" method="post">
<label>Repository URL:</label><br>
<input type="text" name="repoURL" required><br><br>
<label>Image Name:</label><br>
<input type="text" name="imageName" required><br><br>
<label>Image Tag:</label><br>
<input type="text" name="imageTag" required><br><br>
<input type="submit" value="Start Build">
</form>
</body>
</html>`)
}

func buildHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	repoURL := r.FormValue("repoURL")
	imageName := r.FormValue("imageName")
	imageTag := r.FormValue("imageTag")

	if repoURL == "" || imageName == "" || imageTag == "" {
		http.Error(w, "repoURL, imageName and imageTag are required", http.StatusBadRequest)
		return
	}

	if err := utils.Clone(repoURL); err != nil {
		http.Error(w, "error cloning repository: "+err.Error(), http.StatusInternalServerError)
		return
	}

	repoName := utils.GetRepoNameFromURL(repoURL)
	language, _ := utils.DetectLanguage(repoName)
	imgBuilder := imagebuilder.NewImageBuilderImpl(language)
	if imgBuilder.SelectBuilder() == "" {
		http.Error(w, "unsupported language: "+language, http.StatusBadRequest)
		return
	}

	if err := imgBuilder.Build(imageName); err != nil {
		http.Error(w, "build failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := imgBuilder.TagImage(imageName, imageTag); err != nil {
		http.Error(w, "tagging failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Build and tagging complete! Image: %s:%s", imageName, imageTag)
}
