package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

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

// NewDeadlineHandler 构造案件期限处理器。
func NewDeadlineHandler(svc *service.DeadlineService, logger *slog.Logger) *DeadlineHandler {
	return &DeadlineHandler{svc: svc, logger: logger}
}

// Create 登记期限。
func (h *DeadlineHandler) Create(c *gin.Context) {
	var req dto.DeadlineCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline create: "+err.Error())
		return
	}
	dueDate, err := dto.ParseDueDate(req.DueDate)
	if err != nil || dueDate == nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline create: invalid due_date")
		return
	}
	d, err := h.svc.Create(req.CaseID, req.DeadlineType, req.Name, *dueDate, req.OwnerID)
	if err != nil {
		h.wrapError(c, err, "Deadline[case_id="+strconv.FormatUint(req.CaseID, 10)+"] create failed")
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
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline[id="+strconv.FormatUint(id, 10)+"] update: "+err.Error())
		return
	}
	dueDate, err := dto.ParseDueDate(req.DueDate)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline[id="+strconv.FormatUint(id, 10)+"] update: invalid due_date")
		return
	}
	d, err := h.svc.Update(id, req.DeadlineType, req.Name, dueDate, req.OwnerID)
	if err != nil {
		h.wrapError(c, err, "Deadline update failed")
		return
	}
	OKWithMessage(c, constants.MsgDeadlineUpdated, d)
}

// Complete 标记完成，记录处理人与时间。
func (h *DeadlineHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline[id] complete: invalid id")
		return
	}
	d, err := h.svc.Complete(id, middleware.GetUserID(c))
	if err != nil {
		h.wrapError(c, err, "Deadline complete failed")
		return
	}
	OKWithMessage(c, constants.MsgDeadlineCompleted, d)
}

// ListByCase 按案件查看期限。
func (h *DeadlineHandler) ListByCase(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Deadline list: invalid case id")
		return
	}
	list, err := h.svc.ListByCase(caseID)
	if err != nil {
		h.wrapError(c, err, "Deadline list by case failed")
		return
	}
	OK(c, list)
}

// ListCenter 期限中心：view=upcoming/overdue/completed，可按责任人筛选。
func (h *DeadlineHandler) ListCenter(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	ownerID, _ := strconv.ParseUint(c.Query("owner_id"), 10, 64)
	view := c.DefaultQuery("view", constants.DeadlineViewUpcoming)
	list, total, err := h.svc.ListCenter(view, ownerID, q.Page, q.PageSize)
	if err != nil {
		h.wrapError(c, err, "Deadline center list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
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
