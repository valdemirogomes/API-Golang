package dto

type SearchParams struct {
	Limit  *int64
	Offset *int64
	Min    *float64
	Max    *float64
	Sort   *string
	Title  *string
	Name   *string
}

type Metadata struct {
	Total        int64 `json:"total"`
	Limit        int64 `json:"limit"`
	Offset       int64 `json:"offset"`
	TotalEntries int64 `json:"total_entries"`
}
