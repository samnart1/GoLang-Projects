package entities

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

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

type PostInput struct {
	Title	string	`json:"title" validate:"required,min=1,max=255"`
	Content	string	`json:"content" validate:"required,min=1"`
	Slug	string	`json:"slug,omitempty" validate:"omitempty,min=1,max=255"`
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(p.Title) > 255 {
		return fmt.Errorf("title must be less than 256 characters")
	}
	
	if p.Content == "" {
		return fmt.Errorf("content is required")
	}

	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)

	if p.Slug == "" {
		p.Slug = generateSlug(p.Content)
	} else {
		p.Slug = generateSlug(p.Slug)
	}

	return nil
}

func generateSlug(text string) string {
	slug := strings.ToLower(text)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.TrimSuffix(slug, "-")
	}

	return slug
}

func (p *Post) isPublished() bool {
	return p.Published && p.PublishedAt != nil
}