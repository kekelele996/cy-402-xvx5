package dto

// DeadlineCreateRequest 登记期限请求。
type DeadlineCreateRequest struct {
	CaseID     uint64 `json:"case_id" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=hearing appeal evidence trial execution limitation reply registration other"`
	Name       string `json:"name" binding:"required,max=200"`
	DueAt      string `json:"due_at" binding:"required"`
	AssigneeID uint64 `json:"assignee_id" binding:"required"`
}

// DeadlineUpdateRequest 修改未完成期限请求。
type DeadlineUpdateRequest struct {
	Type       string `json:"type" binding:"required,oneof=hearing appeal evidence trial execution limitation reply registration other"`
	Name       string `json:"name" binding:"required,max=200"`
	DueAt      string `json:"due_at" binding:"required"`
	AssigneeID uint64 `json:"assignee_id" binding:"required"`
}

// DeadlineCenterQuery 期限中心查询参数。
type DeadlineCenterQuery struct {
	View       string `form:"view" binding:"omitempty,oneof=upcoming overdue completed"`
	AssigneeID uint64 `form:"assignee_id"`
	CaseID     uint64 `form:"case_id"`
}
