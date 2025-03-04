// Package api はVercelのサーバーレス関数としてToDoアプリのバックエンドを実装します。
// このパッケージでは、RESTful APIの原則に従ってToDo項目のCRUD操作を提供します。
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Todo は1つのタスクを表す構造体です。
// JSONタグを使用してJSONとの相互変換時のフィールド名を指定しています。
type Todo struct {
	ID        int    `json:"id"`        // タスクの一意識別子
	Title     string `json:"title"`     // タスクのタイトル
	Completed bool   `json:"completed"` // タスクの完了状態
}

// グローバル変数の定義
// 注意: サーバーレス環境では関数呼び出し間でデータが保持されない可能性があります
var (
	todos  = make(map[int]*Todo) // ToDoリストを保持するマップ
	nextID = 1                    // 次に割り当てるID
	mutex  sync.RWMutex          // 同時アクセスを制御するためのミューテックス
)

// Index はすべてのAPIリクエストを処理する関数です。
// URLパスに基づいて適切な処理関数にルーティングします。
func Index(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーの設定
	// クロスオリジンリクエストを許可するために必要です
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// プリフライトリクエストの処理
	// ブラウザが送信するOPTIONSリクエストに対応します
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// URLパスの解析
	// /api/todosプレフィックスを除去して実際のパスを取得します
	path := strings.TrimPrefix(r.URL.Path, "/api/todos")
	path = strings.TrimPrefix(path, "/")
	pathParts := strings.Split(path, "/")

	// リクエストのルーティング
	// パスとHTTPメソッドに基づいて適切なハンドラー関数を呼び出します
	switch {
	case path == "" && r.Method == "GET":
		getTodos(w, r)
	case path == "" && r.Method == "POST":
		// リクエストボディを確認して、actionフィールドがある場合はtoggleとして処理
		var input struct {
			Title  string `json:"title"`
			ID     int    `json:"id"`
			Action string `json:"action"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if input.Action == "toggle" {
			toggleTodo(w, r, input.ID)
		} else {
			createTodo(w, r, input.Title)
		}
	case path == "" && r.Method == "DELETE":
		var input struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		deleteTodo(w, r, input.ID)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// getTodos は全てのToDoを取得します。
// GET /api/todos に対応します。
func getTodos(w http.ResponseWriter, r *http.Request) {
	// 読み取り用のロックを取得
	mutex.RLock()
	defer mutex.RUnlock()

	// マップから配列に変換
	todoList := make([]*Todo, 0, len(todos))
	for _, todo := range todos {
		todoList = append(todoList, todo)
	}

	// JSONとしてレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

// createTodo は新しいToDoを作成します。
func createTodo(w http.ResponseWriter, r *http.Request, title string) {
	// 書き込み用のロックを取得
	mutex.Lock()
	// 新しいToDoを作成
	todo := &Todo{
		ID:        nextID,
		Title:     title,
		Completed: false,
	}
	todos[nextID] = todo
	nextID++
	mutex.Unlock()

	// 作成したToDoをJSONとして返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// toggleTodo はToDoの完了状態を切り替えます。
func toggleTodo(w http.ResponseWriter, r *http.Request, id int) {
	// 書き込み用のロックを取得
	mutex.Lock()
	defer mutex.Unlock()

	// ToDoの存在確認と状態の切り替え
	if todo, exists := todos[id]; exists {
		todo.Completed = !todo.Completed
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	} else {
		http.Error(w, "Todo not found", http.StatusNotFound)
	}
}

// deleteTodo はToDoを削除します。
func deleteTodo(w http.ResponseWriter, r *http.Request, id int) {
	// 書き込み用のロックを取得
	mutex.Lock()
	defer mutex.Unlock()

	// ToDoの存在確認と削除
	if _, exists := todos[id]; exists {
		delete(todos, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Todo not found", http.StatusNotFound)
	}
} 