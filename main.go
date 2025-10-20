package main

import (
	"log"
	"net/http"
	"os"
	"scentify/handlers"
	"scentify/services"
	"scentify/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	generator := services.GetScentGenerator(cfg)
	imageSearcher := services.GetImageSearcher(cfg)

	http.Handle("/generate", &handlers.GenerateHandler{
		Generator: generator,
		ImageSearcher: imageSearcher,
	})

	log.SetOutput(os.Stdout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
