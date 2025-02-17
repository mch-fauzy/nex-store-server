package dto

type PaginationResponse struct {
	TotalPages   int `json:"totalPages,omitempty"`
	CurrentPage  int `json:"currentPage,omitempty"`
	NextPage     int `json:"nextPage,omitempty"`
	PreviousPage int `json:"previousPage,omitempty"`
}
