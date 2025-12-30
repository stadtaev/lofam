package note

import "time"

type Color string

const (
	ColorYellow Color = "yellow"
	ColorPink   Color = "pink"
	ColorGreen  Color = "green"
)

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Color     Color     `json:"color"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Color   Color  `json:"color"`
}

func (r CreateRequest) Validate() error {
	if r.Title == "" {
		return ErrValidation("title is required")
	}
	if !isValidColor(r.Color) {
		return ErrValidation("color must be yellow, pink, or green")
	}
	return nil
}

type UpdateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Color   Color  `json:"color"`
}

func (r UpdateRequest) Validate() error {
	if r.Title == "" {
		return ErrValidation("title is required")
	}
	if !isValidColor(r.Color) {
		return ErrValidation("color must be yellow, pink, or green")
	}
	return nil
}

func isValidColor(c Color) bool {
	return c == ColorYellow || c == ColorPink || c == ColorGreen
}
