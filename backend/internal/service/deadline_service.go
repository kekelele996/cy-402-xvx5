package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
	"cylawcase/internal/util"
)

// DeadlineService 案件期限业务逻辑。
type DeadlineService struct {
	repo     *repository.DeadlineRepository
	caseRepo *repository.CaseRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewDeadlineService 构造期限服务。
func NewDeadlineService(repo *repository.DeadlineRepository, caseRepo *repository.CaseRepository,
	userRepo *repository.UserRepository, logger *slog.Logger) *DeadlineService {
	return &DeadlineService{repo: repo, caseRepo: caseRepo, userRepo: userRepo, logger: logger}
}

// Create 登记期限。已结案/归档案件只能补录过去日期；同一案件同一天同名称只保留一条。
func (s *DeadlineService) Create(caseID, assigneeID uint64, deadlineType, name string, dueAt time.Time) (*model.Deadline, error) {
	if !constants.IsValidDeadlineType(deadlineType) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[type="+deadlineType+"] create: invalid type")
	}
	name = truncateSpace(name)
	if name == "" {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[name] create: name must not be blank")
	}
	if dueAt.IsZero() {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[due_at] create: due_at required")
	}
	c, err := s.caseRepo.FindByID(caseID)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[case_id=%d] create: case not found", caseID)
	}
	if _, err := s.userRepo.FindByID(assigneeID); err != nil {
		return nil, util.Wrap(err, "Deadline[assignee_id=%d] create: assignee not found", assigneeID)
	}
	if err := s.checkCaseClosedPastOnly(c.Status, dueAt, "create"); err != nil {
		return nil, err
	}
	dueDate := dueAt.Format("2006-01-02")
	if err := s.ensureUnique(caseID, dueDate, name, 0); err != nil {
		return nil, err
	}
	d := &model.Deadline{
		CaseID:     caseID,
		Type:       deadlineType,
		Name:       name,
		DueAt:      dueAt,
		DueDate:    dueDate,
		AssigneeID: assigneeID,
		Status:     constants.DeadlineStatusPending,
	}
	if err := s.repo.Create(d); err != nil {
		s.logger.Error(constants.LogDeadlineCreateFailed, "error", err.Error())
		return nil, util.Wrap(err, "Deadline[case_id=%d, name=%s] create failed", caseID, name)
	}
	s.logger.Info(constants.LogDeadlineCreateSuccess, "deadline_id", d.ID, "case_id", caseID, "due_at", d.DueAt)
	return d, nil
}

// Update 修改未完成期限。已完成项不可再改；结案/归档案件下同样只能改为过去日期。
func (s *DeadlineService) Update(id, assigneeID uint64, deadlineType, name string, dueAt time.Time) (*model.Deadline, error) {
	d, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] update find failed", id)
	}
	if d.Status != constants.DeadlineStatusPending {
		s.logger.Warn(constants.LogDeadlineUpdateFailed, "deadline_id", id, "status", d.Status)
		return nil, util.NewAppError(constants.CodeDeadlineStatusConflict,
			"Deadline[id="+u64(id)+"] update failed: completed deadline can not modify")
	}
	if !constants.IsValidDeadlineType(deadlineType) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[type="+deadlineType+"] update: invalid type")
	}
	name = truncateSpace(name)
	if name == "" {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[name] update: name must not be blank")
	}
	if dueAt.IsZero() {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[due_at] update: due_at required")
	}
	c, err := s.caseRepo.FindByID(d.CaseID)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] update: case not found", id)
	}
	if _, err := s.userRepo.FindByID(assigneeID); err != nil {
		return nil, util.Wrap(err, "Deadline[assignee_id=%d] update: assignee not found", assigneeID)
	}
	if err := s.checkCaseClosedPastOnly(c.Status, dueAt, "update"); err != nil {
		return nil, err
	}
	dueDate := dueAt.Format("2006-01-02")
	if err := s.ensureUnique(d.CaseID, dueDate, name, id); err != nil {
		return nil, err
	}
	d.Type = deadlineType
	d.Name = name
	d.DueAt = dueAt
	d.DueDate = dueDate
	d.AssigneeID = assigneeID
	if err := s.repo.Update(d); err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] update save failed", id)
	}
	s.logger.Info(constants.LogDeadlineUpdateSuccess, "deadline_id", d.ID, "case_id", d.CaseID)
	return d, nil
}

// Complete 标记完成，保留处理人与完成时间；已完成项幂等拒绝。
func (s *DeadlineService) Complete(id, operatorID uint64) (*model.Deadline, error) {
	d, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] complete find failed", id)
	}
	if d.Status == constants.DeadlineStatusCompleted {
		s.logger.Warn(constants.LogDeadlineCompleteFailed, "deadline_id", id, "status", d.Status)
		return nil, util.NewAppError(constants.CodeDeadlineStatusConflict,
			"Deadline[id="+u64(id)+"] complete failed: already completed by user_id="+u64ptrValue(d.CompletedByID))
	}
	now := time.Now()
	d.Status = constants.DeadlineStatusCompleted
	d.CompletedByID = &operatorID
	d.CompletedAt = &now
	if err := s.repo.Update(d); err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] complete save failed", id)
	}
	s.logger.Info(constants.LogDeadlineCompleteSuccess, "deadline_id", d.ID, "completed_by", operatorID)
	return d, nil
}

// CenterList 期限中心：按即将到期/已逾期/已完成视图查看，可按责任人筛选。
func (s *DeadlineService) CenterList(page, pageSize int, view string, assigneeID, caseID uint64) ([]model.DeadlineViewItem, int64, error) {
	if view == "" {
		view = constants.DeadlineViewUpcoming
	}
	if !constants.IsValidDeadlineView(view) {
		return nil, 0, util.NewAppError(constants.CodeValidationFailed, "Deadline[view="+view+"] list: invalid view")
	}
	list, total, err := s.repo.CenterList(page, pageSize, view, assigneeID, caseID)
	if err != nil {
		return nil, 0, util.Wrap(err, "Deadline center list[view=%s] failed", view)
	}
	s.logger.Info(constants.LogDeadlineList, "view", view, "assignee_id", assigneeID, "total", total)
	return list, total, nil
}

// ListByCase 查询某案件全部期限。
func (s *DeadlineService) ListByCase(caseID uint64) ([]model.DeadlineViewItem, error) {
	if _, err := s.caseRepo.FindByID(caseID); err != nil {
		return nil, util.Wrap(err, "Deadline[case_id=%d] list: case not found", caseID)
	}
	return s.repo.ListByCaseView(caseID)
}

// checkCaseClosedPastOnly 已结案/归档案件只能登记或修改为过去日期。
func (s *DeadlineService) checkCaseClosedPastOnly(caseStatus string, dueAt time.Time, action string) error {
	if caseStatus != constants.CaseStatusClosed && caseStatus != constants.CaseStatusArchived {
		return nil
	}
	if util.DaysBetween(time.Now(), dueAt) > 0 {
		s.logger.Warn(constants.LogDeadlineCreateFailed, "case_status", caseStatus, "due_at", dueAt)
		return util.NewAppError(constants.CodeDeadlinePastOnly,
			fmt.Sprintf("Deadline[case_status=%s] %s failed: closed/archived case only allows past date", caseStatus, action))
	}
	return nil
}

// ensureUnique 同一案件同一天同名称只保留一条。
func (s *DeadlineService) ensureUnique(caseID uint64, dueDate, name string, excludeID uint64) error {
	exist, err := s.repo.FindByCaseDayName(caseID, dueDate, name, excludeID)
	if err == nil && exist != nil {
		return util.NewAppError(constants.CodeDeadlineDuplicate,
			"Deadline[case_id="+u64(caseID)+", due_date="+dueDate+", name="+name+"] duplicate: only one record allowed")
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return util.Wrap(err, "Deadline[case_id=%d] uniqueness check failed", caseID)
	}
	return nil
}

func truncateSpace(v string) string {
	return strings.TrimSpace(v)
}

func u64ptrValue(v *uint64) string {
	if v == nil {
		return "0"
	}
	return u64(*v)
}
