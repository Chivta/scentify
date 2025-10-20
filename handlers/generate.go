package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"scentify/managers"
)

type GenerateHandler struct {
	Generator     *managers.ScentGenerator
	ImageSearcher *managers.ImageSearcher
}

type Request struct {
	Description 	string	`json:"description"`
	NoteAmount  	int		`json:"noteAmount"`
	Silliness   	int		`json:"silliness"`
	GenerateImages 	bool	`json:"generateImages"`
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

	var request Request
	err = json.Unmarshal(bodyByte, &request)
	if err != nil {
		log.Println("Error unmarshalling:", err)
		log.Println(string(bodyByte))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !requestValid(request) {
		log.Println("Invalid request:", request)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(request)

	notes, err := h.Generator.GenerateNotes(request.Description, 
		map[string]any{
			"note_amount": request.NoteAmount, 
			"silliness": request.Silliness,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(notes)

	imageLinks := make([]string,len(notes))

	if request.GenerateImages{
		imageLinks, err = h.ImageSearcher.GetQueryImageLinks(notes)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}


	response := make([]map[string]string, len(notes))

	for i := range imageLinks {
		if imageLinks[i] == "" {
			imageLinks[i] = "static/images/no_image.png"
		}

		response[i] = map[string]string{
			"image": imageLinks[i],
			"note":  notes[i],
		}
	}

	contentJson, err := json.Marshal(response)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(contentJson)
}

var (
	minSillinessLevel = 1
	maxSillinessLevel = 10
	minNodeAmount     = 1
	maxNodeAmount     = 10
	descriptionLimit  = 256
)

func requestValid(request Request) bool {
	if len(request.Description) > descriptionLimit {
		request.Description = request.Description[:descriptionLimit]
	}else if len(request.Description) == 0 {
		return false
	}

	if request.NoteAmount > maxNodeAmount || request.NoteAmount < minNodeAmount {
		return false
	}

	if request.Silliness > maxSillinessLevel || request.Silliness < minSillinessLevel {
		return false
	}

	return true
}
