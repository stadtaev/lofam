package domain

import "time"

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	ProjectID   int64        `json:"projectId"`
	DueDate     *time.Time   `json:"dueDate,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type CreateTaskRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	ProjectID   int64        `json:"projectId"`
	DueDate     *time.Time   `json:"dueDate,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string       `json:"title,omitempty"`
	Description *string       `json:"description,omitempty"`
	Status      *TaskStatus   `json:"status,omitempty"`
	Priority    *TaskPriority `json:"priority,omitempty"`
	DueDate     *time.Time    `json:"dueDate,omitempty"`
}

func (r CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return ErrValidation("title is required")
	}
	if r.ProjectID == 0 {
		return ErrValidation("projectId is required")
	}
	if r.Priority != "" && !isValidPriority(r.Priority) {
		return ErrValidation("invalid priority: must be low, medium, or high")
	}
	return nil
}

func (r UpdateTaskRequest) Validate() error {
	if r.Status != nil && !isValidStatus(*r.Status) {
		return ErrValidation("invalid status: must be todo, in_progress, or done")
	}
	if r.Priority != nil && !isValidPriority(*r.Priority) {
		return ErrValidation("invalid priority: must be low, medium, or high")
	}
	return nil
}

func isValidStatus(s TaskStatus) bool {
	return s == TaskStatusTodo || s == TaskStatusInProgress || s == TaskStatusDone
}

func isValidPriority(p TaskPriority) bool {
	return p == TaskPriorityLow || p == TaskPriorityMedium || p == TaskPriorityHigh
}
