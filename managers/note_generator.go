package managers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"log"
	"net/http/httputil"
)

type NoteGenerator struct {
	cfg *Config 
}

func GetNoteGenerator(cfg *Config) (*NoteGenerator){
	return &NoteGenerator{cfg: cfg}
}


type APIResponse struct {
    Output []struct {
		Content []struct {
			Text string		`json:"text"`
		} 				`json:"content"`
	} 				`json:"output"`
}

func (ngen *NoteGenerator) Generate(description string) (string, error){
	payload := map[string]interface{}{
		"model": "gpt-5-nano",
		"reasoning": map[string]string{
			"effort": "low",
		},
		"prompt": map[string]string{
			"id": ngen.cfg.PropmtId,
		},
		"input": description,
	}

	jsonData, err := json.Marshal(payload)
	if err!=nil{
		return "",err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", nil
    }
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Authorization","Bearer " + ngen.cfg.OpenAiAPIKey)

	dump, err := httputil.DumpRequest(req, true) // 'true' includes body
    if err != nil {
        return "",err
    }

	log.Println(string(dump))


	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()


	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		return "",err
	}
	
	log.Println(string(body))

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("bad response")
	}

	var response APIResponse

	err = json.Unmarshal(body,&response)
	if err != nil{
		return "", err
	}
	log.Println(response)

	return response.Output[1].Content[0].Text,nil
}	