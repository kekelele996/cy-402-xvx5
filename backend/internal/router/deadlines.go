package router

import (
	"cylawcase/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerDeadlineRoutes 案件期限路由。
func (r *Router) registerDeadlineRoutes(g *gin.RouterGroup) {
	deadlines := g.Group("/deadlines")
	deadlines.Use(middleware.AuthRequired(r.cfg))
	deadlines.GET("", r.deadline.ListCenter)
	deadlines.POST("", r.deadline.Create)
	deadlines.GET("/by-case/:id", r.deadline.ListByCase)
	deadlines.PUT("/:id", r.deadline.Update)
	deadlines.POST("/:id/complete", r.deadline.Complete)
}
