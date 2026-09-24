package constants

// DeadlineType 期限类型枚举。
const (
	DeadlineTypeHearing      = "hearing"
	DeadlineTypeAppeal       = "appeal"
	DeadlineTypeEvidence     = "evidence"
	DeadlineTypeTrial        = "trial"
	DeadlineTypeExecution    = "execution"
	DeadlineTypeLimitation   = "limitation"
	DeadlineTypeReply        = "reply"
	DeadlineTypeRegistration = "registration"
	DeadlineTypeOther        = "other"
)

// DeadlineTypeValues 全部期限类型值。
var DeadlineTypeValues = []string{
	DeadlineTypeHearing, DeadlineTypeAppeal, DeadlineTypeEvidence, DeadlineTypeTrial,
	DeadlineTypeExecution, DeadlineTypeLimitation, DeadlineTypeReply, DeadlineTypeRegistration,
	DeadlineTypeOther,
}

// DeadlineStatus 期限状态枚举（pending 未完成 / completed 已完成）。
const (
	DeadlineStatusPending   = "pending"
	DeadlineStatusCompleted = "completed"
)

// DeadlineStatusValues 全部期限状态值。
var DeadlineStatusValues = []string{DeadlineStatusPending, DeadlineStatusCompleted}

// 期限中心视图：upcoming 即将到期 / overdue 已逾期 / completed 已完成。
const (
	DeadlineViewUpcoming  = "upcoming"
	DeadlineViewOverdue   = "overdue"
	DeadlineViewCompleted = "completed"
)

// DeadlineViewValues 全部期限中心视图值。
var DeadlineViewValues = []string{DeadlineViewUpcoming, DeadlineViewOverdue, DeadlineViewCompleted}

// UpcomingWindowDays 即将到期窗口：未完成且截止时间在未来 N 天内。
const UpcomingWindowDays = 7

// IsValidDeadlineType 校验期限类型。
func IsValidDeadlineType(s string) bool {
	for _, v := range DeadlineTypeValues {
		if v == s {
			return true
		}
	}
	return false
}

// IsValidDeadlineStatus 校验期限状态。
func IsValidDeadlineStatus(s string) bool {
	for _, v := range DeadlineStatusValues {
		if v == s {
			return true
		}
	}
	return false
}

// IsValidDeadlineView 校验期限中心视图。
func IsValidDeadlineView(s string) bool {
	for _, v := range DeadlineViewValues {
		if v == s {
			return true
		}
	}
	return false
}
