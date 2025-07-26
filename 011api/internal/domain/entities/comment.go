package entities

import (
	"fmt"
	"strings"
	"time"
)

type Comment struct {
	ID			int			`json:"id"`
	PostID		int			`json:"post_id"`
	UserID		int			`json:"user_id"`
	Content		string		`json:"content"`
	User		*User		`json:"user,omitempty"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type CommentInput struct {
	Content	string	`json:"comment" validate:"required,min=1,max=1000"`
}

func (c *Comment) Validate() error {
	if c.Content == "" {
		return fmt.Errorf("content is required")
	}
		
	if len(c.Content) > 1000 {
		return fmt.Errorf("content must be at most 1000 characters")
	}

	c.Content = strings.TrimSpace(c.Content)

	return nil
}