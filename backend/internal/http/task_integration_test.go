//go:build integration

package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type wantStatus string

func assertTask(t *testing.T, got task.Task, want wantTask, status wantStatus) {
	t.Helper()

	if status == "" {
		status = "todo"
	}
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
	if string(got.Status) != string(status) {
		t.Errorf("status = %q, want %q", got.Status, status)
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
				assertTask(t, got, *tt.want, "")
			}
		})
	}
}

func createTestTask(t *testing.T, baseURL, title string) task.Task {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"title": title})
	resp, err := http.Post(baseURL+"/api/tasks", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}
	defer resp.Body.Close()
	var created task.Task
	json.NewDecoder(resp.Body).Decode(&created)
	return created
}

func TestGetTask(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	created := createTestTask(t, ts.URL, "Test task")

	t.Run("existing", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + "/api/tasks/" + fmt.Sprint(created.ID))
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("not found", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + "/api/tasks/99999")
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})
}

func TestUpdateTask(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	created := createTestTask(t, ts.URL, "Original")

	t.Run("valid", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{"title": "Updated", "status": "done"})
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/tasks/"+fmt.Sprint(created.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		var got task.Task
		json.NewDecoder(resp.Body).Decode(&got)
		assertTask(t, got, wantTask{title: "Updated", priority: "medium"}, "done")
	})

	t.Run("not found", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{"title": "X"})
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/tasks/99999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{"title": "X", "status": "invalid"})
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/tasks/"+fmt.Sprint(created.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})
}

func TestDeleteTask(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	created := createTestTask(t, ts.URL, "To delete")

	t.Run("existing", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/tasks/"+fmt.Sprint(created.ID), nil)
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/tasks/99999", nil)
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})
}
