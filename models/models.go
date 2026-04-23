package models

// Series represents a TV series in the tracker
type Series struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	Status      string `json:"status"` // "watching", "completed", "plan_to_watch", "dropped"
	Episodes    int    `json:"episodes"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// SeriesInput is used for creating/updating a series
type SeriesInput struct {
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	Status      string `json:"status"`
	Episodes    int    `json:"episodes"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// ErrorResponse is the standard error JSON response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse wraps a successful message
type SuccessResponse struct {
	Message string `json:"message"`
}
