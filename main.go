package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"./todo"
)

// グローバル変数の定義
var (
	tmpl     *template.Template
	todoList *todo.TodoList
)

func init() {
	// テンプレートとTodoListの初期化
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	todoList = todo.NewTodoList()
}

func main() {
	// 静的ファイルの提供
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// APIエンドポイントの設定
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/todos", handleTodos)
	http.HandleFunc("/api/todos/toggle/", handleToggleTodo)
	http.HandleFunc("/api/todos/delete/", handleDeleteTodo)

	// サーバーの起動
	fmt.Println("サーバーを起動します。http://localhost:8080 にアクセスしてください。")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// インデックスページのハンドラー
func handleIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Todos []*todo.Todo
	}{
		Todos: todoList.GetAll(),
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ToDoのCRUD APIハンドラー
func handleTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(todoList.GetAll())

	case http.MethodPost:
		var input struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todo := todoList.Add(input.Title)
		json.NewEncoder(w).Encode(todo)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ToDoの完了状態を切り替えるハンドラー
func handleToggleTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.URL.Path[len("/api/todos/toggle/"):])
	if success := todoList.Toggle(id); !success {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ToDoを削除するハンドラー
func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.URL.Path[len("/api/todos/delete/"):])
	if success := todoList.Delete(id); !success {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
} 