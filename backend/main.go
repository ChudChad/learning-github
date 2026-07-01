package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Todo est une tâche de la liste.
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// store garde les todos en mémoire (pas de base de données pour rester simple).
type store struct {
	mu     sync.Mutex
	todos  []Todo
	nextID int
}

func newStore() *store {
	return &store{nextID: 1}
}

func (s *store) list() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Retourne une copie pour éviter les accès concurrents.
	out := make([]Todo, len(s.todos))
	copy(out, s.todos)
	return out
}

func (s *store) add(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Todo{ID: s.nextID, Title: title, Done: false}
	s.nextID++
	s.todos = append(s.todos, t)
	return t
}

// update modifie le titre et/ou l'état "done" d'un todo. Retourne false si l'id n'existe pas.
func (s *store) update(id int, title *string, done *bool) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.todos {
		if s.todos[i].ID == id {
			if title != nil {
				s.todos[i].Title = *title
			}
			if done != nil {
				s.todos[i].Done = *done
			}
			return s.todos[i], true
		}
	}
	return Todo{}, false
}

func (s *store) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.todos {
		if s.todos[i].ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return true
		}
	}
	return false
}

// server rassemble les dépendances des handlers.
type server struct {
	store *store
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// handleTodos gère /api/todos : GET (lister) et POST (créer).
func (s *server) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.store.list())
	case http.MethodPost:
		var body struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "corps JSON invalide", http.StatusBadRequest)
			return
		}
		title := strings.TrimSpace(body.Title)
		if title == "" {
			http.Error(w, "le titre est requis", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, s.store.add(title))
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// handleTodo gère /api/todos/{id} : PUT (modifier) et DELETE (supprimer).
func (s *server) handleTodo(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var body struct {
			Title *string `json:"title"`
			Done  *bool   `json:"done"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "corps JSON invalide", http.StatusBadRequest)
			return
		}
		todo, ok := s.store.update(id, body.Title, body.Done)
		if !ok {
			http.Error(w, "todo introuvable", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, todo)
	case http.MethodDelete:
		if !s.store.delete(id) {
			http.Error(w, "todo introuvable", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// cors autorise le frontend SvelteKit (autre origine) à appeler l'API.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	srv := &server{store: newStore()}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/todos", srv.handleTodos)
	mux.HandleFunc("/api/todos/", srv.handleTodo)

	addr := ":8080"
	log.Printf("Backend ToDo démarré sur http://localhost%s", addr)
	if err := http.ListenAndServe(addr, cors(mux)); err != nil {
		log.Fatal(err)
	}
}
}