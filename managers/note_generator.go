package managers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/httputil"
	"strings"
)

type ScentGenerator struct {
	cfg *Config
}

func GetScentGenerator(cfg *Config) *ScentGenerator {
	return &ScentGenerator{cfg: cfg}
}

type APIResponse struct {
	Output []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func (sgen *ScentGenerator) fetchAPIResponse(description string, variables map[string]string) ([]byte, error) {
	payload := map[string]any{
		"model": "gpt-5-nano",
		"reasoning": map[string]string{
			"effort": "low",
		},
		"prompt": map[string]any{
			"id":        sgen.cfg.PropmtId,
			"variables": variables,
		},

		"input": description,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sgen.cfg.OpenAiAPIKey)

	// dump, err := httputil.DumpRequest(req, true) // 'true' includes body
	// if err != nil {
	//     return "",err
	// }

	// log.Println(string(dump))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println(string(body))
		return nil, fmt.Errorf("bad response")
	}

	return body, nil
}

func (sgen *ScentGenerator) getApiResponse(description string, variables map[string]string) (APIResponse, error) {
	var response APIResponse

	body, err := sgen.fetchAPIResponse(description, variables)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (sgen *ScentGenerator) GenerateNotes(description string, variables map[string]any) ([]string, error) {
	variablesStr := make(map[string]string,len(variables))
	for k,v := range variables{
		switch v.(type){
		case int:
			variablesStr[k]=fmt.Sprintf("%d",v)
		case string:
			variablesStr[k]=fmt.Sprintf("%s",v)
		}
	}

	response, err := sgen.getApiResponse(description, variablesStr)

	if err != nil {
		return nil, err
	}

	var scents string

	for _, output := range response.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					scents = content.Text
				}
			}
		}

	}

	return strings.Split(scents, ","), nil
}
