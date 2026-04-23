package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"series-tracker/database"
	"series-tracker/models"
)

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a standardized JSON error response
func writeError(w http.ResponseWriter, status int, errType, message string) {
	writeJSON(w, status, models.ErrorResponse{
		Error:   errType,
		Message: message,
	})
}

// GetAllSeries handles GET /series
func GetAllSeries(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, title, genre, status, episodes, description, image_url, created_at, updated_at
		FROM series
		ORDER BY created_at DESC
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", "Failed to retrieve series")
		return
	}
	defer rows.Close()

	seriesList := []models.Series{}
	for rows.Next() {
		var s models.Series
		err := rows.Scan(&s.ID, &s.Title, &s.Genre, &s.Status, &s.Episodes, &s.Description, &s.ImageURL, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "scan_error", "Failed to parse series data")
			return
		}
		seriesList = append(seriesList, s)
	}

	writeJSON(w, http.StatusOK, seriesList)
}

// GetSeriesByID handles GET /series/:id
func GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/series/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "The ID must be a valid integer")
		return
	}

	var s models.Series
	err = database.DB.QueryRow(`
		SELECT id, title, genre, status, episodes, description, image_url, created_at, updated_at
		FROM series WHERE id = ?
	`, id).Scan(&s.ID, &s.Title, &s.Genre, &s.Status, &s.Episodes, &s.Description, &s.ImageURL, &s.CreatedAt, &s.UpdatedAt)

	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "not_found", "Series with that ID does not exist")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", "Failed to retrieve series")
		return
	}

	writeJSON(w, http.StatusOK, s)
}

// CreateSeries handles POST /series — returns 201 Created
func CreateSeries(w http.ResponseWriter, r *http.Request) {
	var input models.SeriesInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON")
		return
	}

	validationErrors := validateSeriesInput(input)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":   "validation_error",
			"message": "One or more fields are invalid",
			"fields":  validationErrors,
		})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO series (title, genre, status, episodes, description, image_url)
		VALUES (?, ?, ?, ?, ?, ?)
	`, input.Title, input.Genre, input.Status, input.Episodes, input.Description, input.ImageURL)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", "Failed to create series")
		return
	}

	newID, _ := result.LastInsertId()

	var s models.Series
	database.DB.QueryRow(`
		SELECT id, title, genre, status, episodes, description, image_url, created_at, updated_at
		FROM series WHERE id = ?
	`, newID).Scan(&s.ID, &s.Title, &s.Genre, &s.Status, &s.Episodes, &s.Description, &s.ImageURL, &s.CreatedAt, &s.UpdatedAt)

	writeJSON(w, http.StatusCreated, s)
}

// UpdateSeries handles PUT /series/:id — returns 200 OK with updated resource
func UpdateSeries(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/series/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "The ID must be a valid integer")
		return
	}

	var exists int
	database.DB.QueryRow("SELECT COUNT(*) FROM series WHERE id = ?", id).Scan(&exists)
	if exists == 0 {
		writeError(w, http.StatusNotFound, "not_found", "Series with that ID does not exist")
		return
	}

	var input models.SeriesInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON")
		return
	}

	validationErrors := validateSeriesInput(input)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":   "validation_error",
			"message": "One or more fields are invalid",
			"fields":  validationErrors,
		})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE series
		SET title = ?, genre = ?, status = ?, episodes = ?, description = ?, image_url = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, input.Title, input.Genre, input.Status, input.Episodes, input.Description, input.ImageURL, id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", "Failed to update series")
		return
	}

	var s models.Series
	database.DB.QueryRow(`
		SELECT id, title, genre, status, episodes, description, image_url, created_at, updated_at
		FROM series WHERE id = ?
	`, id).Scan(&s.ID, &s.Title, &s.Genre, &s.Status, &s.Episodes, &s.Description, &s.ImageURL, &s.CreatedAt, &s.UpdatedAt)

	writeJSON(w, http.StatusOK, s)
}

// DeleteSeries handles DELETE /series/:id — returns 204 No Content
func DeleteSeries(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/series/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "The ID must be a valid integer")
		return
	}

	var exists int
	database.DB.QueryRow("SELECT COUNT(*) FROM series WHERE id = ?", id).Scan(&exists)
	if exists == 0 {
		writeError(w, http.StatusNotFound, "not_found", "Series with that ID does not exist")
		return
	}

	_, err = database.DB.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", "Failed to delete series")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// validateSeriesInput returns a map of field errors
func validateSeriesInput(input models.SeriesInput) map[string]string {
	errors := map[string]string{}

	if strings.TrimSpace(input.Title) == "" {
		errors["title"] = "Title is required and cannot be empty"
	}

	validStatuses := map[string]bool{
		"watching":      true,
		"completed":     true,
		"plan_to_watch": true,
		"dropped":       true,
	}
	if input.Status != "" && !validStatuses[input.Status] {
		errors["status"] = "Status must be one of: watching, completed, plan_to_watch, dropped"
	}

	if input.Episodes < 0 {
		errors["episodes"] = "Episodes cannot be negative"
	}

	return errors
}

// extractID parses an integer ID from a URL path after a given prefix
func extractID(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.TrimSuffix(idStr, "/")
	return strconv.Atoi(idStr)
}
