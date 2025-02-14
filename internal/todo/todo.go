package todo

import (
	"context"
	"errors"
	"fmt"
	"github.com/sureshdsk/todo-goland-api/internal/db"
	"strings"
)

type DBManager interface {
	InsertItem(ctx context.Context, item db.Item) error
	GetAllItems(ctx context.Context) ([]db.Item, error)
	UpdateItemStatus(ctx context.Context, id int, status string) error
	DeleteItem(ctx context.Context, id int) error
}

type Service struct {
	db DBManager
}

func NewService(db DBManager) *Service {
	return &Service{
		db: db,
	}
}

func (svc *Service) Add(todo string) error {
	items, err := svc.GetAll()
	if err != nil {
		return fmt.Errorf("could not get items: %w", err)
	}
	for _, t := range items {
		if t.Task == todo {
			return errors.New("todo already exists")
		}
	}
	if err := svc.db.InsertItem(context.Background(), db.Item{
		Task:   todo,
		Status: "TO_BE_STARTED",
	}); err != nil {
		return fmt.Errorf("could not insert item: %w", err)
	}
	return nil
}

func (svc *Service) Search(query string) ([]db.Item, error) {
	items, err := svc.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get items: %w", err)
	}

	var results []db.Item
	query = strings.ToLower(query)
	for _, todo := range items {
		if strings.Contains(strings.ToLower(todo.Task), query) {
			results = append(results, todo)
		}
	}
	return results, nil
}

func (svc *Service) GetAll() ([]db.Item, error) {
	var results []db.Item
	items, err := svc.db.GetAllItems(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get all items: %w", err)
	}
	for _, item := range items {
		results = append(results, db.Item{
			ID:     item.ID,
			Task:   item.Task,
			Status: item.Status,
		})
	}
	return results, nil
}

func (svc *Service) UpdateStatus(id int, status string) error {
	if status == "" {
		return errors.New("status is required")
	}
	
	if err := svc.db.UpdateItemStatus(context.Background(), id, status); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

func (svc *Service) Delete(id int) error {
	if err := svc.db.DeleteItem(context.Background(), id); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}
