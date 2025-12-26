package task

import "time"

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type CreateRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

type UpdateRequest struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Status      *Status   `json:"status,omitempty"`
	Priority    *Priority `json:"priority,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

func (r CreateRequest) Validate() error {
	if r.Title == "" {
		return ErrValidation("title is required")
	}
	if r.Priority != "" && !isValidPriority(r.Priority) {
		return ErrValidation("invalid priority: must be low, medium, or high")
	}
	return nil
}

func (r UpdateRequest) Validate() error {
	if r.Status != nil && !isValidStatus(*r.Status) {
		return ErrValidation("invalid status: must be todo, in_progress, or done")
	}
	if r.Priority != nil && !isValidPriority(*r.Priority) {
		return ErrValidation("invalid priority: must be low, medium, or high")
	}
	return nil
}

func isValidStatus(s Status) bool {
	return s == StatusTodo || s == StatusInProgress || s == StatusDone
}

func isValidPriority(p Priority) bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}
