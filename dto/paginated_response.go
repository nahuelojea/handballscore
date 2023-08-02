package dto

type PaginatedResponse struct {
	Items        interface{} `json:"items"`
	TotalRecords int64       `json:"total_records"`
	TotalPages   int         `json:"total_pages"`
	CurrentPage  int         `json:"current_page"`
	PageSize     int         `json:"page_size"`
}

func NewPaginatedResponse(items interface{}, totalCount int64, totalPages, currentPage, pageSize int) *PaginatedResponse {
	return &PaginatedResponse{
		Items:        items,
		TotalRecords: totalCount,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PageSize:     pageSize,
	}
}
