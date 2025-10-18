package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"scentify/managers"
	"slices"
)

type GenerateHandler struct {
	Generator *managers.ScentGenerator
	ImageSearcher *managers.ImageSearcher
}

type Request struct{
	Description string 	`json:"desc"`
	Language string		`json:"lang"`
	NoteAmount int		`json:"amnt"`
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

	var request Request
	err = json.Unmarshal(bodyByte, &request)
	if err != nil {
		log.Println("Error unmarshalling:",err)
		log.Println(string(bodyByte))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !requestValid(request){
		log.Println("Invalid request:",request)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(request)


	notes, err := h.Generator.GenerateNotes(request.Description,request.Language,request.NoteAmount)
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

var (
	allowedLanguages = []string{"en","ua"}
	maxNodeAmount = 10
	descriptionLimit = 256
)

func requestValid(request Request) bool {
	if len(request.Description)>descriptionLimit{
		request.Description = request.Description[:descriptionLimit]
	}

	if !slices.Contains(allowedLanguages, request.Language){
		return false
	}
	if request.NoteAmount > maxNodeAmount{
		return false
	}

	return true
}