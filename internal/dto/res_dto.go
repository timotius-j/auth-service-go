package dto

type Response[T any] struct {
	Data    T             `                   json:"data,omitempty"`
	Paging  *PageMetadata `                   json:"paging,omitempty"`
	Error   *Error        `                   json:"error,omitempty"`
	Success bool          `binding:"required" json:"success"`
}

type Error struct {
	Message string `json:"message,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

type HealthCheckResponse struct {
	Message string `binding:"required" json:"message"`
}
