package dto

import "time"

// DeadlineCreateRequest 登记期限请求。
type DeadlineCreateRequest struct {
	CaseID       uint64 `json:"case_id" binding:"required"`
	DeadlineType string `json:"deadline_type" binding:"required,oneof=hearing appeal evidence filing other"`
	Name         string `json:"name" binding:"required,max=200"`
	DueDate      string `json:"due_date" binding:"required"`
	OwnerID      uint64 `json:"owner_id" binding:"required"`
}

// DeadlineUpdateRequest 修改期限请求（仅未完成项可修改）。
type DeadlineUpdateRequest struct {
	DeadlineType string `json:"deadline_type" binding:"omitempty,oneof=hearing appeal evidence filing other"`
	Name         string `json:"name" binding:"omitempty,max=200"`
	DueDate      string `json:"due_date"`
	OwnerID      uint64 `json:"owner_id"`
}

// ParseDueDate 解析截止日期字符串（YYYY-MM-DD）。
func ParseDueDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
