package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Todo は1つのタスクを表す構造体です
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	todos  = make(map[int]*Todo)
	nextID = 1
	mutex  sync.RWMutex
)

// Handler はすべてのAPIリクエストを処理します
func Handler(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーの設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// プリフライトリクエストの処理
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// パスの解析
	path := strings.TrimPrefix(r.URL.Path, "/api/todos")
	path = strings.TrimPrefix(path, "/")
	pathParts := strings.Split(path, "/")

	// ルーティング
	switch {
	case path == "" && r.Method == "GET":
		getTodos(w, r)
	case path == "" && r.Method == "POST":
		createTodo(w, r)
	case len(pathParts) == 2 && pathParts[1] == "toggle" && r.Method == "POST":
		toggleTodo(w, r, pathParts[0])
	case len(pathParts) == 1 && r.Method == "DELETE":
		deleteTodo(w, r, pathParts[0])
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// getTodos は全てのToDoを取得します
func getTodos(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	todoList := make([]*Todo, 0, len(todos))
	for _, todo := range todos {
		todoList = append(todoList, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

// createTodo は新しいToDoを作成します
func createTodo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	todo := &Todo{
		ID:        nextID,
		Title:     input.Title,
		Completed: false,
	}
	todos[nextID] = todo
	nextID++
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// toggleTodo はToDoの完了状態を切り替えます
func toggleTodo(w http.ResponseWriter, r *http.Request, idStr string) {
	id := 0
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if todo, exists := todos[id]; exists {
		todo.Completed = !todo.Completed
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	} else {
		http.Error(w, "Todo not found", http.StatusNotFound)
	}
}

// deleteTodo はToDoを削除します
func deleteTodo(w http.ResponseWriter, r *http.Request, idStr string) {
	id := 0
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := todos[id]; exists {
		delete(todos, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Todo not found", http.StatusNotFound)
	}
} 