package main

import (		
		"database/sql"
		"encoding/json"
		"log"
		"net/http"
		"os"
		"github.com/gorilla/mux"
		_ "github.com/lib/pq"
)

type post struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
}

func main() {
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
				log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS posts (id SERIAL PRIMARY KEY, title TEXT, body TEXT)")
		if err != nil {
				log.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/api/posts", getPosts(db)).Methods("GET")
		router.HandleFunc("/api/posts", createPost(db)).Methods("POST")
		router.HandleFunc("/api/posts/{id}", getPost(db)).Methods("GET")
		router.HandleFunc("/api/posts/{id}", updatePost(db)).Methods("PUT")
		router.HandleFunc("/api/posts/{id}", deletePost(db)).Methods("DELETE")
		enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))
		log.Fatal(http.ListenAndServe(":8000", enhancedRouter))

}

func enableCORS(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				h.ServeHTTP(w, r)
		})
}

func jsonContentTypeMiddleware(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				h.ServeHTTP(w, r)
				})
}

func getPosts(db *sql.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
				rows, err := db.Query("SELECT * FROM posts")
				if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}
				defer rows.Close()

				posts := []post{}
				for rows.Next() {
						var p post
						if err := rows.Scan(&p.ID, &p.Title, &p.Body); err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
								return
						}
						posts = append(posts, p)
				}

				if err := rows.Err(); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}

				if err := json.NewEncoder(w).Encode(posts); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}
		}
}

func getPost(db *sql.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id := vars["id"]

				var p post
				row := db.QueryRow("SELECT * FROM posts WHERE id = $1", id)
				if err := row.Scan(&p.ID, &p.Title, &p.Body); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}

				if err := json.NewEncoder(w).Encode(p); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}
		}
}

func createPost(db *sql.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
				var p post
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
				}

				row := db.QueryRow("INSERT INTO posts (title, body) VALUES ($1, $2) RETURNING id", p.Title, p.Body)
				if err := row.Scan(&p.ID); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}

				if err := json.NewEncoder(w).Encode(p); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}
		}
}

func updatePost(db *sql.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id := vars["id"]

				var p post
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
				}

				_, err := db.Exec("UPDATE posts SET title = $1, body = $2 WHERE id = $3", p.Title, p.Body, id)
				if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}

				if err := json.NewEncoder(w).Encode(p); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}
		}
}

func deletePost(db *sql.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id := vars["id"]

				_, err := db.Exec("DELETE FROM posts WHERE id = $1", id)
				if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
				}

				w.WriteHeader(http.StatusNoContent)
		}
}



