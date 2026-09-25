package service

import (
	"log/slog"
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

// NewDeadlineService 构造案件期限服务。
func NewDeadlineService(repo *repository.DeadlineRepository, caseRepo *repository.CaseRepository,
	userRepo *repository.UserRepository, logger *slog.Logger) *DeadlineService {
	return &DeadlineService{repo: repo, caseRepo: caseRepo, userRepo: userRepo, logger: logger}
}

// Create 登记期限。已结案/归档案件只能补录过去日期；同案件同天同名称只留一条。
func (s *DeadlineService) Create(caseID uint64, deadlineType, name string, dueDate time.Time, ownerID uint64) (*model.Deadline, error) {
	cs, err := s.caseRepo.FindByID(caseID)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[case_id=%d] create: case not found", caseID)
	}
	if !constants.IsValidDeadlineType(deadlineType) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[deadline_type="+deadlineType+"] create: invalid type")
	}
	if _, err := s.userRepo.FindByID(ownerID); err != nil {
		return nil, util.Wrap(err, "Deadline[owner_id=%d] create: owner not found", ownerID)
	}
	today := todayBiz()
	if err := checkClosedCaseDueDate(cs.Status, dueDate, today); err != nil {
		return nil, err
	}
	exists, err := s.repo.ExistsByCaseDateName(caseID, dueDate, name, 0)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[case_id=%d] create duplicate check failed", caseID)
	}
	if exists {
		return nil, util.NewAppError(constants.CodeDeadlineDuplicate, "Deadline[case_id="+u64(caseID)+"] create: "+constants.MsgDeadlineDuplicate)
	}
	d := &model.Deadline{
		CaseID:       caseID,
		DeadlineType: deadlineType,
		Name:         name,
		DueDate:      dueDate,
		OwnerID:      ownerID,
		Status:       constants.DeadlineStatusPending,
	}
	if err := s.repo.Create(d); err != nil {
		s.logger.Error(constants.LogDeadlineCreateFailed, "error", err.Error(), "case_id", caseID)
		return nil, util.Wrap(err, "Deadline[case_id=%d] create save failed", caseID)
	}
	s.logger.Info(constants.LogDeadlineCreateSuccess, "deadline_id", d.ID, "case_id", caseID)
	return d, nil
}

// Update 修改未完成期限；已完成项不可修改。
func (s *DeadlineService) Update(id uint64, deadlineType, name string, dueDate *time.Time, ownerID uint64) (*model.Deadline, error) {
	d, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] update find failed", id)
	}
	if d.Status == constants.DeadlineStatusCompleted {
		return nil, util.NewAppError(constants.CodeDeadlineStatusConflict, "Deadline[id="+u64(id)+"] update: "+constants.MsgDeadlineAlreadyDone)
	}
	cs, err := s.caseRepo.FindByID(d.CaseID)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] update: case not found", id)
	}
	if deadlineType != "" {
		if !constants.IsValidDeadlineType(deadlineType) {
			return nil, util.NewAppError(constants.CodeValidationFailed, "Deadline[deadline_type="+deadlineType+"] update: invalid type")
		}
		d.DeadlineType = deadlineType
	}
	if name != "" {
		d.Name = name
	}
	if dueDate != nil {
		d.DueDate = *dueDate
	}
	if ownerID > 0 {
		if _, err := s.userRepo.FindByID(ownerID); err != nil {
			return nil, util.Wrap(err, "Deadline[owner_id=%d] update: owner not found", ownerID)
		}
		d.OwnerID = ownerID
	}
	today := todayBiz()
	if err := checkClosedCaseDueDate(cs.Status, d.DueDate, today); err != nil {
		return nil, err
	}
	exists, err := s.repo.ExistsByCaseDateName(d.CaseID, d.DueDate, d.Name, d.ID)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] update duplicate check failed", id)
	}
	if exists {
		return nil, util.NewAppError(constants.CodeDeadlineDuplicate, "Deadline[id="+u64(id)+"] update: "+constants.MsgDeadlineDuplicate)
	}
	if err := s.repo.Update(d); err != nil {
		s.logger.Error(constants.LogDeadlineUpdateFailed, "error", err.Error(), "deadline_id", id)
		return nil, util.Wrap(err, "Deadline[id=%d] update save failed", id)
	}
	s.logger.Info(constants.LogDeadlineUpdateSuccess, "deadline_id", id)
	return d, nil
}

// Complete 标记完成，保留处理人与完成时间；已完成项不可重复标记。
func (s *DeadlineService) Complete(id, operatorID uint64) (*model.Deadline, error) {
	d, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Deadline[id=%d] complete find failed", id)
	}
	if d.Status == constants.DeadlineStatusCompleted {
		return nil, util.NewAppError(constants.CodeDeadlineStatusConflict, "Deadline[id="+u64(id)+"] complete: already completed")
	}
	now := time.Now()
	d.Status = constants.DeadlineStatusCompleted
	d.CompletedBy = &operatorID
	d.CompletedAt = &now
	if err := s.repo.Update(d); err != nil {
		s.logger.Error(constants.LogDeadlineCompleteFailed, "error", err.Error(), "deadline_id", id)
		return nil, util.Wrap(err, "Deadline[id=%d] complete save failed", id)
	}
	s.logger.Info(constants.LogDeadlineCompleteSuccess, "deadline_id", id, "operator_id", operatorID)
	return d, nil
}

// ListByCase 按案件查看期限。
func (s *DeadlineService) ListByCase(caseID uint64) ([]model.Deadline, error) {
	return s.repo.ListByCase(caseID)
}

// ListCenter 期限中心查询：按视图（upcoming/overdue/completed）与责任人筛选，返回剩余/逾期天数。
func (s *DeadlineService) ListCenter(view string, ownerID uint64, page, pageSize int) ([]model.DeadlineWithCase, int64, error) {
	if !constants.IsValidDeadlineView(view) {
		return nil, 0, util.NewAppError(constants.CodeValidationFailed, "Deadline center list: invalid view "+view)
	}
	today := todayBiz()
	status, dueFrom, dueTo := deadlineFilterForView(view, today)
	list, total, err := s.repo.ListCenter(status, dueFrom, dueTo, ownerID, page, pageSize)
	if err != nil {
		return nil, 0, util.Wrap(err, "Deadline center[view=%s] list failed", view)
	}
	for i := range list {
		list[i].Days = daysBetween(list[i].DueDate, today)
	}
	return list, total, nil
}

// checkClosedCaseDueDate 已结案或归档案件只能补录过去日期。
func checkClosedCaseDueDate(caseStatus string, dueDate, today time.Time) error {
	if caseStatus == constants.CaseStatusClosed || caseStatus == constants.CaseStatusArchived {
		if !dueDate.Before(today) {
			return util.NewAppError(constants.CodeValidationFailed, "Deadline[case_status="+caseStatus+"] save: "+constants.MsgDeadlineCaseClosed)
		}
	}
	return nil
}

// deadlineFilterForView 期限中心视图 -> 状态与截止日期边界。
// upcoming=未完成且未到期（含今天）；overdue=未完成且早于今天；completed=已完成。
func deadlineFilterForView(view string, today time.Time) (status string, dueFrom, dueTo *time.Time) {
	switch view {
	case constants.DeadlineViewUpcoming:
		return constants.DeadlineStatusPending, &today, nil
	case constants.DeadlineViewOverdue:
		return constants.DeadlineStatusPending, nil, &today
	default:
		return constants.DeadlineStatusCompleted, nil, nil
	}
}

// daysBetween 计算截止日期相对今天的天数：正数剩余、0 今天、负数逾期。
func daysBetween(dueDate, today time.Time) int {
	due := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, time.UTC)
	base := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	return int(due.Sub(base).Hours() / 24)
}

// todayBiz 返回业务当天零点。系统面向国内律所、数据库会话时区为 Asia/Shanghai，
// 固定按东八区取自然日，避免凌晨时段跨天导致逾期/剩余天数偏差。
func todayBiz() time.Time {
	loc := time.FixedZone("Asia/Shanghai", 8*60*60)
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}
