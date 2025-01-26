package todo_test

import (
	"context"
	"github.com/sureshdsk/todo-goland-api/internal/db"
	"github.com/sureshdsk/todo-goland-api/internal/todo"
	"reflect"
	"testing"
)

type MockDBManager struct {
	items []db.Item
}

func (m *MockDBManager) InsertItem(ctx context.Context, item db.Item) error {
	m.items = append(m.items, item)
	return nil
}

func (m *MockDBManager) GetAllItems(ctx context.Context) ([]db.Item, error) {
	return m.items, nil
}

func (m *MockDBManager) UpdateItemStatus(ctx context.Context, id int, status string) error {
	if id >= len(m.items) {
		return db.ErrTaskNotFound
	}
	m.items[id].Status = status
	return nil
}

func (m *MockDBManager) DeleteItem(ctx context.Context, id int) error {
	if id >= len(m.items) {
		return db.ErrTaskNotFound
	}
	m.items = append(m.items[:id], m.items[id+1:]...)
	return nil
}

func TestService_Search(t *testing.T) {
	tests := []struct {
		name     string
		todos    []string
		query    string
		expected []db.Item
	}{
		{
			name:     "given a todo of shop, and a search for sh should return shop",
			todos:    []string{"shop"},
			query:    "sh",
			expected: []db.Item{{Task: "shop", Status: "TO_BE_STARTED"}},
		},
		{
			name:     "still return shop, if case does not match",
			todos:    []string{"Shopping"},
			query:    "sh",
			expected: []db.Item{{Task: "Shopping", Status: "TO_BE_STARTED"}},
		},
		{
			name:     "spaces",
			todos:    []string{"go Shopping"},
			query:    "go",
			expected: []db.Item{{Task: "go Shopping", Status: "TO_BE_STARTED"}},
		},
		{
			name:     "space at the start of the word",
			todos:    []string{" go Shopping"},
			query:    "go",
			expected: []db.Item{{Task: " go Shopping", Status: "TO_BE_STARTED"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &MockDBManager{}
			svc := todo.NewService(d)
			for _, toAdd := range tt.todos {
				err := svc.Add(toAdd)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := svc.Search(tt.query)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Search() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestService_UpdateStatus(t *testing.T) {
	tests := []struct {
		name        string
		setupItems  []string
		updateID    int
		newStatus   string
		expectError bool
	}{
		{
			name:       "successfully update status",
			setupItems: []string{"test task"},
			updateID:   0,
			newStatus:  "COMPLETED",
		},
		{
			name:        "fail when task not found",
			setupItems:  []string{"test task"},
			updateID:    1,
			newStatus:   "COMPLETED",
			expectError: true,
		},
		{
			name:        "fail when status is empty",
			setupItems:  []string{"test task"},
			updateID:    0,
			newStatus:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &MockDBManager{}
			svc := todo.NewService(d)
			
			// Setup items
			for _, item := range tt.setupItems {
				err := svc.Add(item)
				if err != nil {
					t.Fatal(err)
				}
			}

			err := svc.UpdateStatus(tt.updateID, tt.newStatus)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Verify status was updated
				items, err := svc.GetAll()
				if err != nil {
					t.Fatal(err)
				}
				if items[tt.updateID].Status != tt.newStatus {
					t.Errorf("status not updated, got %s, want %s", items[tt.updateID].Status, tt.newStatus)
				}
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		setupItems  []string
		deleteID    int
		expectError bool
	}{
		{
			name:       "successfully delete task",
			setupItems: []string{"test task"},
			deleteID:   0,
		},
		{
			name:        "fail when task not found",
			setupItems:  []string{"test task"},
			deleteID:    1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &MockDBManager{}
			svc := todo.NewService(d)
			
			// Setup items
			for _, item := range tt.setupItems {
				err := svc.Add(item)
				if err != nil {
					t.Fatal(err)
				}
			}

			initialCount := len(d.items)
			err := svc.Delete(tt.deleteID)
			
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				// Verify no items were deleted
				if len(d.items) != initialCount {
					t.Errorf("items were deleted when they shouldn't have been")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Verify item was deleted
				if len(d.items) != initialCount-1 {
					t.Errorf("item was not deleted")
				}
			}
		})
	}
}
