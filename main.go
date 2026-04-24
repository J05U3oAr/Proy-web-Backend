package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"series-tracker/database"
	"series-tracker/handlers"

	"github.com/rs/cors"
)

func main() {
	// Initialize SQLite database
	database.Init(resolveDatabasePath())

	mux := http.NewServeMux()

	// Route: GET /series and POST /series
	mux.HandleFunc("/series", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAllSeries(w, r)
		case http.MethodPost:
			handlers.CreateSeries(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"method_not_allowed","message":"Method %s not allowed on /series"}`, r.Method)
		}
	})

	// Route: GET /series/:id, PUT /series/:id, DELETE /series/:id
	mux.HandleFunc("/series/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetSeriesByID(w, r)
		case http.MethodPut:
			handlers.UpdateSeries(w, r)
		case http.MethodDelete:
			handlers.DeleteSeries(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"method_not_allowed","message":"Method %s not allowed on /series/:id"}`, r.Method)
		}
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok","service":"series-tracker-api"}`)
	})

	// CORS configuration
	// CORS (Cross-Origin Resource Sharing) is a browser security policy that blocks
	// requests from a different origin (domain/port) than the server.
	// We configure it here to allow the frontend (running on a different port) to call our API.
	c := cors.New(cors.Options{
		AllowedOrigins: resolveAllowedOrigins(),
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(mux)

	port := resolvePort()
	log.Printf("Series Tracker API running on port %s", port)
	log.Printf("Endpoints:")
	log.Printf("   GET    /series        - list all series")
	log.Printf("   POST   /series        - create a series (201)")
	log.Printf("   GET    /series/:id    - get one series")
	log.Printf("   PUT    /series/:id    - update a series")
	log.Printf("   DELETE /series/:id    - delete a series (204)")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}

func resolveDatabasePath() string {
	if customPath := strings.TrimSpace(os.Getenv("DATABASE_PATH")); customPath != "" {
		return customPath
	}

	if volumeMount := strings.TrimSpace(os.Getenv("RAILWAY_VOLUME_MOUNT_PATH")); volumeMount != "" {
		return volumeMount + "/series.db"
	}

	if volumeMount := strings.TrimSpace(os.Getenv("RAILWAY_VOLUME_MOUNT")); volumeMount != "" {
		return volumeMount + "/series.db"
	}

	paths := []string{
		"./series.db",
		"../series.db",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "./series.db"
}

func resolvePort() string {
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		return port
	}

	return "8080"
}

func resolveAllowedOrigins() []string {
	origins := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if origins == "" {
		return []string{"*"}
	}

	parts := strings.Split(origins, ",")
	allowed := make([]string, 0, len(parts))

	for _, origin := range parts {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed = append(allowed, origin)
		}
	}

	if len(allowed) == 0 {
		return []string{"*"}
	}

	return allowed
}
