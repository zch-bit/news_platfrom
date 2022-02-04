package models

import "time"

type Response struct {
	Status       string `json:"status"`
	TotalResults int    `json:"total_results"`
	Articles     []News `json:"articles"`
}

type News struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
	FullContent string    `json:"full_content"`
}
