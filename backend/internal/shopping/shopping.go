package shopping

import "time"

type Item struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateRequest struct {
	Title string `json:"title"`
}

func (r CreateRequest) Validate() error {
	if r.Title == "" {
		return ErrValidation("title is required")
	}
	return nil
}
