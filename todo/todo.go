package todo

import (
	"sync"
	"time"
)

// Todo は1つのタスクを表す構造体です
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoList はTodoを管理するための構造体です
type TodoList struct {
	mutex sync.RWMutex
	todos map[int]*Todo
	nextID int
}

// NewTodoList は新しいTodoListを作成します
func NewTodoList() *TodoList {
	return &TodoList{
		todos: make(map[int]*Todo),
		nextID: 1,
	}
}

// Add は新しいTodoを追加します
func (l *TodoList) Add(title string) *Todo {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	todo := &Todo{
		ID:        l.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	
	l.todos[l.nextID] = todo
	l.nextID++
	
	return todo
}

// GetAll は全てのTodoを取得します
func (l *TodoList) GetAll() []*Todo {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	todos := make([]*Todo, 0, len(l.todos))
	for _, todo := range l.todos {
		todos = append(todos, todo)
	}
	return todos
}

// Toggle はTodoの完了状態を切り替えます
func (l *TodoList) Toggle(id int) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if todo, exists := l.todos[id]; exists {
		todo.Completed = !todo.Completed
		return true
	}
	return false
}

// Delete は指定されたIDのTodoを削除します
func (l *TodoList) Delete(id int) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, exists := l.todos[id]; exists {
		delete(l.todos, id)
		return true
	}
	return false
} 