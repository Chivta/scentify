package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
)

type GenerateHandler struct {
}

func (h *GenerateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// if len(bodyByte) == 0{
	// 	return
	// }

	body := string(bodyByte[:128])
	log.Println(body)

	n := rand.IntN(7)
	content := make([]map[string]string, n)

	for i := 0; i < n; i++ {
		content[i] = map[string]string{
			"image" : "static/images/salt_vinegar.png",
			"title" : "Lady of the Night Flowerg",
		}
	}

	contentJson, err := json.Marshal(content)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(contentJson)
}

func main() {
	http.Handle("/generate", &GenerateHandler{})

	log.SetOutput(os.Stdout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
