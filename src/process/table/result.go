package table

type PageResult struct {
	TotalDiffTime int64       `json:"total_diff_time"`
	List          []*ItemData `json:"list"`
	PageIndex     int         `json:"page_index"`
	PageSize      int         `json:"page_size"`
	TotalSize     int         `json:"total_size"`
}
