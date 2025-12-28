//go:build integration

package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	lofamhttp "github.com/stadtaev/lofam/backend/internal/http"
	"github.com/stadtaev/lofam/backend/internal/sqlite"
	"github.com/stadtaev/lofam/backend/internal/task"
)

type wantTask struct {
	title       string
	description string
	priority    string
	dueDate     string
}

func assertTask(t *testing.T, got task.Task, want wantTask) {
	t.Helper()

	if got.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if got.Title != want.title {
		t.Errorf("title = %q, want %q", got.Title, want.title)
	}
	if got.Description != want.description {
		t.Errorf("description = %q, want %q", got.Description, want.description)
	}
	if string(got.Priority) != want.priority {
		t.Errorf("priority = %q, want %q", got.Priority, want.priority)
	}
	if got.Status != "todo" {
		t.Errorf("status = %q, want %q", got.Status, "todo")
	}
	if got.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}

	assertDueDate(t, got.DueDate, want.dueDate)
}

func assertDueDate(t *testing.T, got *time.Time, want string) {
	t.Helper()

	if want == "" {
		if got != nil {
			t.Errorf("dueDate = %v, want nil", got)
		}
		return
	}
	if got == nil {
		t.Errorf("dueDate = nil, want %q", want)
		return
	}
	if got.Format(time.RFC3339) != want {
		t.Errorf("dueDate = %q, want %q", got.Format(time.RFC3339), want)
	}
}

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	db, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	store := sqlite.NewTaskStore(db)
	service := task.NewService(store)
	server := lofamhttp.NewServer(service)

	return httptest.NewServer(server.Router())
}

func TestCreateTask(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
		want       *wantTask // nil for error cases
	}{
		{
			name:       "valid task with title only",
			body:       map[string]any{"title": "Buy groceries"},
			wantStatus: http.StatusCreated,
			want:       &wantTask{title: "Buy groceries", priority: "medium"},
		},
		{
			name: "valid task with all fields",
			body: map[string]any{
				"title":       "Complete project",
				"description": "Finish the backend API",
				"priority":    "high",
				"dueDate":     "2025-12-31T00:00:00Z",
			},
			wantStatus: http.StatusCreated,
			want: &wantTask{
				title:       "Complete project",
				description: "Finish the backend API",
				priority:    "high",
				dueDate:     "2025-12-31T00:00:00Z",
			},
		},
		{
			name:       "missing title",
			body:       map[string]any{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty title",
			body:       map[string]any{"title": ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid priority",
			body:       map[string]any{"title": "Test task", "priority": "invalid"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			resp, err := http.Post(ts.URL+"/api/tasks", "application/json", bytes.NewReader(body))
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.want != nil {
				var got task.Task
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				assertTask(t, got, *tt.want)
			}
		})
	}
}
