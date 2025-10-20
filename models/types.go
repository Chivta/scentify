package models

type GenerateRequest struct {
	Description 	string	`json:"description"`
	NoteAmount  	int		`json:"noteAmount"`
	Silliness   	int		`json:"silliness"`
	GenerateImages 	bool	`json:"generateImages"`
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