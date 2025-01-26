package transport

import (
	"encoding/json"
	"github.com/sureshdsk/todo-goland-api/internal/todo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TodoItem struct {
	Item string `json:"item"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(todoSvc *todo.Service) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		todoItems, err := todoSvc.GetAll()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(todoItems)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		var t TodoItem
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = todoSvc.Add(t.Item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	})

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		results, err := todoSvc.Search(query)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	mux.HandleFunc("PATCH /todo/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("Invalid task ID:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req UpdateStatusRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println("Invalid request body:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = todoSvc.UpdateStatus(id, req.Status)
		if err != nil {
			log.Println("Failed to update status:", err)
			if strings.Contains(err.Error(), "task not found") {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("DELETE /todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("Invalid task ID:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = todoSvc.Delete(id)
		if err != nil {
			log.Println("Failed to delete task:", err)
			if strings.Contains(err.Error(), "task not found") {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	return &Server{mux: mux}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux)
}
