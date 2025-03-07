// Package handler はVercelのサーバーレス関数としてToDoアプリのバックエンドを実装します。
// このパッケージでは、RESTful APIの原則に従ってToDo項目のCRUD操作を提供します。
package handler

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Todo は1つのタスクを表す構造体です。
// JSONタグを使用してJSONとの相互変換時のフィールド名を指定しています。
type Todo struct {
	ID    int    `json:"id"`    // タスクの一意識別子
	Title string `json:"title"` // タスクのタイトル
}

// グローバル変数の定義
// 注意: サーバーレス環境では関数呼び出し間でデータが保持されない可能性があります
var (
	todos  = make(map[int]*Todo) // ToDoリストを保持するマップ
	nextID = 1                   // 次に割り当てるID
	mutex  sync.RWMutex         // 同時アクセスを制御するためのミューテックス
)

// Index はすべてのAPIリクエストを処理する関数です。
func Index(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーの設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// プリフライトリクエストの処理
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// URLパスの解析
	switch r.Method {
	case "GET":
		getTodos(w)
	case "POST":
		createTodo(w, r)
	case "DELETE":
		deleteTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getTodos は全てのToDoを取得します。
// GET /api/todos に対応します。
func getTodos(w http.ResponseWriter) {
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
		ID:    nextID,
		Title: input.Title,
	}
	todos[nextID] = todo
	nextID++
	mutex.Unlock()

	// 作成したToDoをJSONとして返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// deleteTodo はToDoを削除します。
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// ToDoの存在確認と削除
	if _, exists := todos[input.ID]; exists {
		delete(todos, input.ID)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Todo not found", http.StatusNotFound)
	}
} 