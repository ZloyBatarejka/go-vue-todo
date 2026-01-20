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
	Create(todo *models.Todo) error
	GetAll() ([]*models.Todo, error)
	GetByID(id int64) (*models.Todo, error)
	Delete(id int64) error
}

// todoRepository реализует TodoRepository
type todoRepository struct {
	db *sql.DB
}

// NewTodoRepository создает новый экземпляр репозитория
func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

// Create создает новую задачу в БД
// ID генерируется автоматически на бэкенде через последовательность PostgreSQL
// Если последовательности нет, создаем её при первом использовании
func (r *todoRepository) Create(todo *models.Todo) error {
	// Сначала проверяем и создаем последовательность если её нет
	_, err := r.db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS todos_id_seq;
		SELECT setval('todos_id_seq', COALESCE((SELECT MAX(id) FROM todos), 0) + 1, false);
	`)
	if err != nil {
		return fmt.Errorf("failed to setup sequence: %w", err)
	}

	// Дату создания задаём на бэкенде (входящее значение игнорируем)
	todo.Date = time.Now().UTC().Format(time.RFC3339)

	// Используем nextval для генерации ID
	query := `
		INSERT INTO todos (id, value, date) 
		VALUES (nextval('todos_id_seq'), $1, $2) 
		RETURNING id, value, date
	`

	err = r.db.QueryRow(query, todo.Value, todo.Date).Scan(
		&todo.ID,
		&todo.Value,
		&todo.Date,
	)

	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

// GetAll получает все задачи из БД
func (r *todoRepository) GetAll() ([]*models.Todo, error) {
	query := `SELECT id, value, date FROM todos ORDER BY id DESC`

	rows, err := r.db.Query(query)
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

// GetByID получает задачу по ID
func (r *todoRepository) GetByID(id int64) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `SELECT id, value, date FROM todos WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Value,
		&todo.Date,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("todo with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return todo, nil
}

// Delete удаляет задачу по ID
func (r *todoRepository) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}
