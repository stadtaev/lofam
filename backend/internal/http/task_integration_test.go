//go:build integration

package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	lofamhttp "github.com/stadtaev/lofam/backend/internal/http"
	"github.com/stadtaev/lofam/backend/internal/sqlite"
	"github.com/stadtaev/lofam/backend/internal/task"
)

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
		wantTitle  string
	}{
		{
			name: "valid task with title only",
			body: map[string]any{
				"title": "Buy groceries",
			},
			wantStatus: http.StatusCreated,
			wantTitle:  "Buy groceries",
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
			wantTitle:  "Complete project",
		},
		{
			name:       "missing title",
			body:       map[string]any{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty title",
			body: map[string]any{
				"title": "",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid priority",
			body: map[string]any{
				"title":    "Test task",
				"priority": "invalid",
			},
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
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusCreated {
				var created task.Task
				if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if created.Title != tt.wantTitle {
					t.Errorf("got title %q, want %q", created.Title, tt.wantTitle)
				}

				if created.ID == 0 {
					t.Error("expected non-zero ID")
				}

				if created.Status != "todo" {
					t.Errorf("got status %q, want %q", created.Status, "todo")
				}

				if created.CreatedAt.IsZero() {
					t.Error("expected non-zero CreatedAt")
				}
			}
		})
	}
}
