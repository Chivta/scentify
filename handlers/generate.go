package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"scentify/managers"
	"strings"
)

type GenerateHandler struct {
	Generator *managers.NoteGenerator
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

	description := string(bodyByte[:128])
	log.Println(description)

	notes, err := h.Generator.Generate(description)
	if err != nil{
		log.Println(err)
		return
	}

	notesSlice := strings.Split(notes, ",")

	content := make([]map[string]string, len(notesSlice))

	for i,note := range notesSlice {
		content[i] = map[string]string{
			"image" : "static/images/salt_vinegar.png",
			"note" : note,
		}
	}

	contentJson, err := json.Marshal(content)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(contentJson)
}