package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/httputil"
	"strconv"
	"strings"
	. "scentify/models"
	. "scentify/config"
)

type ScentGenerator struct {
	cfg *Config
}

func GetScentGenerator(cfg *Config) *ScentGenerator {
	return &ScentGenerator{cfg: cfg}
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

var (
	sillinessLevels = map[int]string{
		1:"All notes must be real perfume ingredients (e.g., jasmine, amber, cedarwood).",
		2:"Notes should be mostly real (75% real, 25% invented scents).",
		3:"Notes should be half real, half invented (50/50 mix).",
		4:"Notes should be mostly imaginary scents (25% real, 75% invented).",
		5:"All notes must be invented scents (e.g., “midnight rain,” “old library”).",
		6:"Notes should be mostly invented scents (75% invented, 25% abstract vibes).",
		7:"Notes should be half invented, half vibe-like concepts (50/50).",
		8:"Notes should mostly be vibes (75% vibes, 25% invented scents)",
		9:"All notes must be vibes only (e.g. “melancholy,” “sunset regret”).",
		10:"All notes must be fully symbolic or surreal vibes — pure emotion or concept, not related to scent at all.",
	}
)

func formatVariables(variables map[string]any) (map[string]string){
	variablesFormated := make(map[string]string,len(variables))

	if sillines,ok := variables["silliness"]; ok {
		sillinesInt, ok := sillines.(int)
		if !ok {
			log.Println("Sillines is not an int. silliness:",sillines)
		}else{
			variables["silliness"] = sillinessLevels[sillinesInt]
		}
	}else{
		log.Println("Did not found silliness")
	}

	for k,v := range variables{
		switch vConv:=v.(type){
		case int:
			variablesFormated[k]=strconv.Itoa(vConv)
		case string:
			variablesFormated[k]=vConv
		}
	}

	return variablesFormated
}

func (sgen *ScentGenerator) GenerateNotes(description string, variables map[string]any) ([]string, error) {
	variablesFormated := formatVariables(variables)

	log.Println("Variables formated:",variablesFormated)

	response, err := sgen.getApiResponse(description, variablesFormated)

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
