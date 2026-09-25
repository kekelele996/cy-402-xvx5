package model

import "time"

// Deadline 案件期限实体：同一案件同一天同名称只留一条（联合唯一索引）。
type Deadline struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	CaseID       uint64     `gorm:"not null;index;uniqueIndex:uk_deadline_case_date_name" json:"case_id"`
	DeadlineType string     `gorm:"size:30;not null;default:other" json:"deadline_type"`
	Name         string     `gorm:"size:200;not null;uniqueIndex:uk_deadline_case_date_name" json:"name"`
	DueDate      time.Time  `gorm:"type:date;not null;uniqueIndex:uk_deadline_case_date_name" json:"due_date"`
	OwnerID      uint64     `gorm:"not null;index" json:"owner_id"`
	Status       string     `gorm:"size:30;not null;default:pending;index" json:"status"`
	CompletedBy  *uint64    `json:"completed_by"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TableName 指定表名。
func (Deadline) TableName() string { return "deadlines" }

// DeadlineWithCase 期限中心视图：附带案号、案件标题与剩余/逾期天数。
type DeadlineWithCase struct {
	Deadline
	CaseNo    string `json:"case_no"`
	CaseTitle string `json:"case_title"`
	Days      int    `json:"days"` // 正数=剩余天数，0=今天到期，负数=逾期天数
}
