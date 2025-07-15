package entities

import "time"

type Post struct {
	ID			int			`json:"id"`
	Title		string		`json:"json"`
	Content		string		`json:"content"`
	Slug		string		`json:"slug"`
	AuthorID	string		`json:"author_id"`
	Author		*User		`json:"author,omitempty"`
	Published	bool		`json:"published"`
	PublishedAt	*time.Time	`json:"published_at,omitempty"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

