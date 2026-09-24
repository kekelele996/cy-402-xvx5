package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/dto"
	"cylawcase/internal/middleware"
	"cylawcase/internal/service"
	"cylawcase/internal/util"

	"github.com/gin-gonic/gin"
)

// DeadlineHandler 案件期限 HTTP 处理器。
type DeadlineHandler struct {
	svc    *service.DeadlineService
	logger *slog.Logger
}

// NewDeadlineHandler 构造期限处理器。
func NewDeadlineHandler(svc *service.DeadlineService, logger *slog.Logger) *DeadlineHandler {
	return &DeadlineHandler{svc: svc, logger: logger}
}

// CenterList 期限中心列表（即将到期/已逾期/已完成，可按责任人筛选）。
func (h *DeadlineHandler) CenterList(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	var f dto.DeadlineCenterQuery
	_ = c.ShouldBindQuery(&f)
	if f.View == "" {
		f.View = constants.DeadlineViewUpcoming
	}
	list, total, err := h.svc.CenterList(q.Page, q.PageSize, f.View, f.AssigneeID, f.CaseID)
	if err != nil {
		h.wrapError(c, err, "Deadline center list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// ListByCase 某案件的期限列表。
func (h *DeadlineHandler) ListByCase(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline list by case: invalid case id")
		return
	}
	list, err := h.svc.ListByCase(caseID)
	if err != nil {
		h.wrapError(c, err, "Deadline list by case failed")
		return
	}
	OK(c, list)
}

// Create 登记期限。
func (h *DeadlineHandler) Create(c *gin.Context) {
	var req dto.DeadlineCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline create: "+err.Error())
		return
	}
	dueAt, err := parseDeadlineDueAt(req.DueAt)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline[due_at="+req.DueAt+"] create: invalid datetime, use YYYY-MM-DD HH:mm")
		return
	}
	d, err := h.svc.Create(req.CaseID, req.AssigneeID, req.Type, req.Name, dueAt)
	if err != nil {
		h.wrapError(c, err, "Deadline[case_id="+strconv.FormatUint(req.CaseID, 10)+",name="+req.Name+"] create failed")
		return
	}
	OKWithMessage(c, constants.MsgDeadlineCreated, d)
}

// Update 修改未完成期限。
func (h *DeadlineHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline[id] update: invalid id")
		return
	}
	var req dto.DeadlineUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline update: "+err.Error())
		return
	}
	dueAt, err := parseDeadlineDueAt(req.DueAt)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline[due_at="+req.DueAt+"] update: invalid datetime, use YYYY-MM-DD HH:mm")
		return
	}
	d, err := h.svc.Update(id, req.AssigneeID, req.Type, req.Name, dueAt)
	if err != nil {
		h.wrapError(c, err, "Deadline[id="+strconv.FormatUint(id, 10)+"] update failed")
		return
	}
	OKWithMessage(c, constants.MsgDeadlineUpdated, d)
}

// Complete 标记完成（处理人取自登录态）。
func (h *DeadlineHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline[id] complete: invalid id")
		return
	}
	d, err := h.svc.Complete(id, middleware.GetUserID(c))
	if err != nil {
		h.wrapError(c, err, "Deadline[id="+strconv.FormatUint(id, 10)+"] complete failed")
		return
	}
	OKWithMessage(c, constants.MsgDeadlineCompleted, d)
}

// parseDeadlineDueAt 兼容 "2006-01-02 15:04"、"2006-01-02 15:04:05"、RFC3339 与仅日期。
func parseDeadlineDueAt(v string) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02"}
	var lastErr error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, v, time.Local)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func (h *DeadlineHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("deadline handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("deadline handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
