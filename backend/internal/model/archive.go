package model

// ArchiveMonth 归档月份统计。
type ArchiveMonth struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}
