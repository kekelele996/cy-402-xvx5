package router

import (
	"cylawcase/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerDeadlineRoutes 案件期限路由。
func (r *Router) registerDeadlineRoutes(g *gin.RouterGroup) {
	deadlines := g.Group("/deadlines")
	deadlines.Use(middleware.AuthRequired(r.cfg))
	deadlines.GET("", r.deadline.CenterList)
	deadlines.POST("", r.deadline.Create)
	deadlines.PUT("/:id", r.deadline.Update)
	deadlines.POST("/:id/complete", r.deadline.Complete)

	// 按案件查询期限：/cases/:id/deadlines
	caseDeadlines := g.Group("/cases")
	caseDeadlines.Use(middleware.AuthRequired(r.cfg))
	caseDeadlines.GET("/:id/deadlines", r.deadline.ListByCase)
}
