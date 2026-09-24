package model

import "time"

// Deadline 案件期限实体。
type Deadline struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	CaseID        uint64     `gorm:"not null;index:idx_deadline_case_name_day,unique,priority:1" json:"case_id"`
	Type          string     `gorm:"size:30;not null;default:other" json:"type"`
	Name          string     `gorm:"size:200;not null;index:idx_deadline_case_name_day,unique,priority:2" json:"name"`
	DueAt         time.Time  `gorm:"not null;index" json:"due_at"`
	DueDate       string     `gorm:"size:10;not null;index:idx_deadline_case_name_day,unique,priority:3" json:"due_date"`
	AssigneeID    uint64     `gorm:"not null;index" json:"assignee_id"`
	Status        string     `gorm:"size:30;not null;default:pending;index" json:"status"`
	CompletedByID *uint64    `gorm:"index" json:"completed_by_id"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// TableName 指定表名。
func (Deadline) TableName() string { return "deadlines" }

// DeadlineViewItem 期限中心列表项，附带案件与责任人展示信息。
type DeadlineViewItem struct {
	Deadline
	CaseNo        string `json:"case_no"`
	CaseTitle     string `json:"case_title"`
	CaseStatus    string `json:"case_status"`
	AssigneeName  string `json:"assignee_name"`
	CompleterName string `json:"completer_name"`
}
