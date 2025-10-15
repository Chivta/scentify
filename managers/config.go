package managers

import (
	"encoding/json"
	"os"
)

type Config struct {
	OpenAiAPIKey	string `json:"openai_api_key"`
	PropmtId		string `json:"prompt_id"`
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