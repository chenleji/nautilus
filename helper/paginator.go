package helper

import "math"

// page search util
type Paginator struct {
	*PageRequest
	TotalPages    int         `json:"totalPages"`
	TotalElements int         `json:"totalElements"`
	Content       interface{} `json:"content"`
}

// page begin from 1
type PageRequest struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
}

func (p Paginator) GetOffset(r *PageRequest) int {
	if r.CurrentPage != 0 {
		return (r.CurrentPage - 1) * r.PageSize
	}
	return 0
}

func (p Paginator) GetTotalPages(r *PageRequest, totalElements int) int {
	totalPages := math.Ceil(float64(totalElements) / float64(r.PageSize))
	return int(totalPages)
}
