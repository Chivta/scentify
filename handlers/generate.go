package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"scentify/managers"
)

type GenerateHandler struct {
	Generator *managers.ScentGenerator
	ImageSearcher *managers.ImageSearcher
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

	description := string(bodyByte[:256])
	log.Println(description)

	notes, err := h.Generator.GenerateNotes(description)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(notes)

	notesWithImages, err := h.ImageSearcher.GetQueryImageLinks(notes)
	if err!=nil{
		log.Println(err)
		return
	}

	content := make([]map[string]string, len(notes))

	for i, noteImage := range notesWithImages {
		if noteImage.Link == ""{
			noteImage.Link="static/images/no_image.jpg"
		}

		content[i] = map[string]string{
			"image": noteImage.Link,
			"note":  noteImage.Query,
		}
	}

	contentJson, err := json.Marshal(content)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(contentJson)
}
