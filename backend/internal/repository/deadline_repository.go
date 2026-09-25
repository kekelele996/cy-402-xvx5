package repository

import (
	"errors"
	"fmt"
	"time"

	"cylawcase/internal/model"

	"gorm.io/gorm"
)

// DeadlineRepository 案件期限仓储。
type DeadlineRepository struct {
	db *gorm.DB
}

// NewDeadlineRepository 构造案件期限仓储。
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

// Update 更新期限。
func (r *DeadlineRepository) Update(d *model.Deadline) error {
	if err := r.db.Save(d).Error; err != nil {
		return fmt.Errorf("update deadline: %w", err)
	}
	return nil
}

// ListByCase 查询某案件的期限，按截止日期升序。
func (r *DeadlineRepository) ListByCase(caseID uint64) ([]model.Deadline, error) {
	var list []model.Deadline
	if err := r.db.Where("case_id = ?", caseID).Order("due_date ASC, id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list deadlines by case: %w", err)
	}
	return list, nil
}

// ExistsByCaseDateName 判断同案件同天同名期限是否已存在，excludeID 用于更新时排除自身。
// due_date 为 date 类型，统一格式化为 YYYY-MM-DD 后按 ::date 比较，避免会话时区造成偏差。
func (r *DeadlineRepository) ExistsByCaseDateName(caseID uint64, dueDate time.Time, name string, excludeID uint64) (bool, error) {
	var n int64
	q := r.db.Model(&model.Deadline{}).
		Where("case_id = ? AND due_date = ?::date AND name = ?", caseID, dueDate.Format("2006-01-02"), name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&n).Error; err != nil {
		return false, fmt.Errorf("count deadline by case/date/name: %w", err)
	}
	return n > 0, nil
}

// ListCenter 期限中心分页查询：按状态与截止日期范围筛选（关联案件取案号与标题），可按责任人过滤。
func (r *DeadlineRepository) ListCenter(status string, dueFrom, dueTo *time.Time, ownerID uint64, page, pageSize int) ([]model.DeadlineWithCase, int64, error) {
	apply := func(q *gorm.DB) *gorm.DB {
		if status != "" {
			q = q.Where("deadlines.status = ?", status)
		}
		if dueFrom != nil {
			q = q.Where("deadlines.due_date >= ?::date", dueFrom.Format("2006-01-02"))
		}
		if dueTo != nil {
			q = q.Where("deadlines.due_date < ?::date", dueTo.Format("2006-01-02"))
		}
		if ownerID > 0 {
			q = q.Where("deadlines.owner_id = ?", ownerID)
		}
		return q
	}
	var total int64
	if err := apply(r.db.Model(&model.Deadline{})).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count deadlines: %w", err)
	}
	order := "deadlines.due_date ASC, deadlines.id ASC"
	if status == "completed" {
		order = "deadlines.completed_at DESC, deadlines.id DESC"
	}
	var list []model.DeadlineWithCase
	q := apply(r.db.Model(&model.Deadline{}).
		Select("deadlines.*, cases.case_no AS case_no, cases.title AS case_title").
		Joins("JOIN cases ON cases.id = deadlines.case_id"))
	if err := q.Order(order).Offset((page - 1) * pageSize).Limit(pageSize).Scan(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list deadlines center: %w", err)
	}
	return list, total, nil
}
