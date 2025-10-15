package managers

import (
	"encoding/json"
	"os"
)

type Config struct {
	OpenAiAPIKey	string `json:"openai_api_key"`
	PropmtId		string `json:"prompt_id"`
	SerpApiKey		string `json:"serp_api_key"`
}

func GetConfig() (*Config, error){
	data, err := os.ReadFile("config.json")
	if err != nil{
		return nil, err
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err!=nil{
		return nil, err
	}

	return &cfg, nil
}