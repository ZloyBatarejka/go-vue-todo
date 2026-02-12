package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"goTodo/backend/models"
)

// TodoRepository определяет интерфейс для работы с задачами
// Использование интерфейса позволяет легко тестировать и менять реализацию
type TodoRepository interface {
	Create(todo *models.Todo, userID int64) error
	GetAllByUserID(userID int64) ([]*models.Todo, error)
	GetByIDForUser(id int64, userID int64) (*models.Todo, error)
	DeleteForUser(id int64, userID int64) error
}

// todoRepository реализует TodoRepository
type todoRepository struct {
	db *sql.DB
}

// NewTodoRepository создает новый экземпляр репозитория
func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

// Create создает новую задачу в БД для конкретного пользователя.
// ID генерируется самой БД через DEFAULT/IDENTITY у колонки todos.id.
func (r *todoRepository) Create(todo *models.Todo, userID int64) error {
	// Дату создания задаём на бэкенде (входящее значение игнорируем)
	todo.Date = time.Now().UTC().Format(time.RFC3339)

	query := `
		INSERT INTO todos (value, date, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, value, date
	`

	err := r.db.QueryRow(query, todo.Value, todo.Date, userID).Scan(
		&todo.ID,
		&todo.Value,
		&todo.Date,
	)

	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

// GetAllByUserID получает все задачи текущего пользователя.
func (r *todoRepository) GetAllByUserID(userID int64) ([]*models.Todo, error) {
	query := `SELECT id, value, date FROM todos WHERE user_id = $1 ORDER BY id DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		todo := &models.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Value, &todo.Date); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todos: %w", err)
	}

	return todos, nil
}

// GetByIDForUser получает задачу по ID, только если она принадлежит пользователю.
func (r *todoRepository) GetByIDForUser(id int64, userID int64) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `SELECT id, value, date FROM todos WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&todo.ID,
		&todo.Value,
		&todo.Date,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("todo with id %d not found: %w", id, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return todo, nil
}

// DeleteForUser удаляет задачу по ID только в рамках текущего пользователя.
func (r *todoRepository) DeleteForUser(id int64, userID int64) error {
	query := `DELETE FROM todos WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found: %w", id, sql.ErrNoRows)
	}

	return nil
}
