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
