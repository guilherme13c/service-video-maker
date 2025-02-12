package model

type GenerateVideoRequest struct {
	Id   string `json:"id"`
	Data struct {
		BackgroundURL string   `json:"background_url"`
		ImageList     []string `json:"image_list"`
		VoiceURL      string   `json:"voice_url"`
		Transcripts   []struct {
			Words string  `json:"words"`
			Start float64 `json:"start"`
			End   float64 `json:"end"`
		} `json:"transcripts"`
	} `json:"data"`
	WebhookURL *string `json:"webhook_url,omitempty"`
}
