package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"math/big"
)

// Post представляет структуру записи в блоге
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// Хранилище постов в памяти
var (
	posts   = make(map[int]*Post)
	nextID  = 1
	postsMu sync.RWMutex
)

func init() {
	// Добавляем тестовые посты
	posts[1] = &Post{
		ID:        1,
		Title:     "Добро пожаловать!",
		Content:   "Это первый пост в нашем блоге. Здесь будет интересный контент!",
		Author:    "Admin",
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
	posts[2] = &Post{
		ID:        2,
		Title:     "Go - отличный язык для бэкенда",
		Content:   "Go обеспечивает высокую производительность, простоту разработки и отличную поддержку конкурентности.",
		Author:    "Developer",
		CreatedAt: time.Now().Add(-12 * time.Hour),
	}
	nextID = 3
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func factorial(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

// Обработчик для получения всех постов
func getPosts(w http.ResponseWriter, r *http.Request) {
	factorial(2000)
	w.Write([]byte("ok"))

	if r.Method == "OPTIONS" {
		enableCORS(w)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	postsMu.RLock()
	defer postsMu.RUnlock()

	allPosts := make([]*Post, 0, len(posts))
	for _, post := range posts {
		allPosts = append(allPosts, post)
	}

	json.NewEncoder(w).Encode(allPosts)
}

// Обработчик для получения конкретного поста
func getPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		enableCORS(w)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/api/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postsMu.RLock()
	post, exists := posts[id]
	postsMu.RUnlock()

	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

// Обработчик для создания нового поста
func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		enableCORS(w)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postsMu.Lock()
	post.ID = nextID
	nextID++
	post.CreatedAt = time.Now()
	posts[post.ID] = &post
	postsMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// Обработчик для обновления поста
func updatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		enableCORS(w)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var updatedPost Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postsMu.Lock()
	post, exists := posts[id]
	if !exists {
		postsMu.Unlock()
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if updatedPost.Title != "" {
		post.Title = updatedPost.Title
	}
	if updatedPost.Content != "" {
		post.Content = updatedPost.Content
	}
	if updatedPost.Author != "" {
		post.Author = updatedPost.Author
	}
	postsMu.Unlock()

	json.NewEncoder(w).Encode(post)
}

// Обработчик для удаления поста
func deletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		enableCORS(w)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postsMu.Lock()
	_, exists := posts[id]
	if !exists {
		postsMu.Unlock()
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	delete(posts, id)
	postsMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// Health check
func healthCheck(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "ok",
		"message": "Blog API is running",
		"version": "1.0.0",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Регистрируем обработчики для API
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "OPTIONS" {
			getPosts(w, r)
		} else if r.Method == "POST" {
			createPost(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET", "OPTIONS":
			getPost(w, r)
		case "PUT":
			updatePost(w, r)
		case "DELETE":
			deletePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Раздача статических файлов из директории frontend
	// ВАЖНО: поместите файлы index.html, styles.css, app.js в директорию frontend/
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Printf("🚀 Сервер запущен на http://localhost%s\n", port)
	fmt.Println("📝 API endpoints:")
	fmt.Println("   GET    /health        - Health check")
	fmt.Println("   GET    /api/posts     - Получить все посты")
	fmt.Println("   POST   /api/posts     - Создать новый пост")
	fmt.Println("   GET    /api/posts/:id - Получить пост по ID")
	fmt.Println("   PUT    /api/posts/:id - Обновить пост")
	fmt.Println("   DELETE /api/posts/:id - Удалить пост")
	fmt.Println("📂 Статические файлы: /")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
