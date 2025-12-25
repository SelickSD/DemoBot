package dto

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message,omitempty"`
	FinishReason string  `json:"finish_reason"`
}
