package dto

type NewsFeed struct {
	Id        int    `json:"id"`
	Published int    `json:"published"`
	Type      int    `json:"type"`
	Message   string `json:"message"`
}
