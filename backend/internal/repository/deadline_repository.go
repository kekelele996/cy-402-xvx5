package repository

import (
	"errors"
	"fmt"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"

	"gorm.io/gorm"
)

// DeadlineRepository 案件期限仓储。
type DeadlineRepository struct {
	db *gorm.DB
}

// NewDeadlineRepository 构造期限仓储。
func NewDeadlineRepository(db *gorm.DB) *DeadlineRepository {
	return &DeadlineRepository{db: db}
}

// Create 创建期限。
func (r *DeadlineRepository) Create(d *model.Deadline) error {
	if err := r.db.Create(d).Error; err != nil {
		return fmt.Errorf("create deadline: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询期限。
func (r *DeadlineRepository) FindByID(id uint64) (*model.Deadline, error) {
	var d model.Deadline
	if err := r.db.First(&d, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find deadline by id: %w", err)
	}
	return &d, nil
}

// FindByCaseDayName 校验同一案件同一天同名称唯一，excludeID 用于修改时排除自身。
func (r *DeadlineRepository) FindByCaseDayName(caseID uint64, dueDate, name string, excludeID uint64) (*model.Deadline, error) {
	var d model.Deadline
	q := r.db.Where("case_id = ? AND due_date = ? AND name = ?", caseID, dueDate, name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find deadline by case/day/name: %w", err)
	}
	return &d, nil
}

// Update 更新期限。
func (r *DeadlineRepository) Update(d *model.Deadline) error {
	if err := r.db.Save(d).Error; err != nil {
		return fmt.Errorf("update deadline: %w", err)
	}
	return nil
}

// deadlineViewSelect 期限中心列表字段：期限本体 + 案号/标题/状态 + 责任人/完成人姓名。
const deadlineViewSelect = "deadlines.*, cases.case_no AS case_no, cases.title AS case_title, " +
	"cases.status AS case_status, COALESCE(NULLIF(u.real_name, ''), u.username) AS assignee_name, " +
	"COALESCE(NULLIF(cu.real_name, ''), cu.username) AS completer_name"

func deadlineViewQuery(db *gorm.DB) *gorm.DB {
	return db.Table("deadlines").
		Select(deadlineViewSelect).
		Joins("LEFT JOIN cases ON cases.id = deadlines.case_id").
		Joins("LEFT JOIN users u ON u.id = deadlines.assignee_id").
		Joins("LEFT JOIN users cu ON cu.id = deadlines.completed_by_id")
}

// viewWhere 视图过滤：upcoming 未完成且在窗口内（含今天）、overdue 未完成且已过期、completed 已完成。
func viewWhere(q *gorm.DB, view string, today time.Time) *gorm.DB {
	windowEnd := today.AddDate(0, 0, constants.UpcomingWindowDays)
	switch view {
	case constants.DeadlineViewUpcoming:
		return q.Where("deadlines.status = ? AND deadlines.due_date >= ? AND deadlines.due_date <= ?",
			constants.DeadlineStatusPending, today.Format("2006-01-02"), windowEnd.Format("2006-01-02"))
	case constants.DeadlineViewOverdue:
		return q.Where("deadlines.status = ? AND deadlines.due_date < ?",
			constants.DeadlineStatusPending, today.Format("2006-01-02"))
	case constants.DeadlineViewCompleted:
		return q.Where("deadlines.status = ?", constants.DeadlineStatusCompleted)
	default:
		return q
	}
}

// CenterList 期限中心分页查询，支持视图/责任人/案件筛选。
func (r *DeadlineRepository) CenterList(page, pageSize int, view string, assigneeID, caseID uint64) ([]model.DeadlineViewItem, int64, error) {
	today := time.Now()
	base := r.db.Table("deadlines")
	base = viewWhere(base, view, today)
	if assigneeID > 0 {
		base = base.Where("deadlines.assignee_id = ?", assigneeID)
	}
	if caseID > 0 {
		base = base.Where("deadlines.case_id = ?", caseID)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count deadlines: %w", err)
	}

	q := deadlineViewQuery(r.db)
	q = viewWhere(q, view, today)
	if assigneeID > 0 {
		q = q.Where("deadlines.assignee_id = ?", assigneeID)
	}
	if caseID > 0 {
		q = q.Where("deadlines.case_id = ?", caseID)
	}
	switch view {
	case "completed":
		q = q.Order("deadlines.completed_at DESC NULLS LAST, deadlines.id DESC")
	default:
		q = q.Order("deadlines.due_at ASC, deadlines.id ASC")
	}

	var list []model.DeadlineViewItem
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Scan(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list deadline center: %w", err)
	}
	return list, total, nil
}

// ListByCaseView 查询某案件的全部期限（含责任人/完成人姓名）。
func (r *DeadlineRepository) ListByCaseView(caseID uint64) ([]model.DeadlineViewItem, error) {
	var list []model.DeadlineViewItem
	if err := deadlineViewQuery(r.db).
		Where("deadlines.case_id = ?", caseID).
		Order("deadlines.status ASC, deadlines.due_at ASC, deadlines.id DESC").
		Scan(&list).Error; err != nil {
		return nil, fmt.Errorf("list deadlines by case: %w", err)
	}
	return list, nil
}
