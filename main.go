package main

import (
	"log"
	"net/http"
	"os"
	"scentify/handlers"
	"scentify/managers"
)

func main() {
	cfg, err := managers.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	generator := managers.GetScentGenerator(cfg)
	imageSearcher := managers.GetImageSearcher(cfg)

	http.Handle("/generate", &handlers.GenerateHandler{
		Generator: generator,
		ImageSearcher: imageSearcher,
	})

	log.SetOutput(os.Stdout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
