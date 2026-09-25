package constants

// DeadlineType 期限类型枚举。
const (
	DeadlineTypeHearing  = "hearing"
	DeadlineTypeAppeal   = "appeal"
	DeadlineTypeEvidence = "evidence"
	DeadlineTypeFiling   = "filing"
	DeadlineTypeOther    = "other"
)

// DeadlineTypeValues 全部期限类型值。
var DeadlineTypeValues = []string{DeadlineTypeHearing, DeadlineTypeAppeal, DeadlineTypeEvidence, DeadlineTypeFiling, DeadlineTypeOther}

// DeadlineStatus 期限状态枚举。
const (
	DeadlineStatusPending   = "pending"
	DeadlineStatusCompleted = "completed"
)

// DeadlineStatusValues 全部期限状态值。
var DeadlineStatusValues = []string{DeadlineStatusPending, DeadlineStatusCompleted}

// DeadlineView 期限中心视图枚举。
const (
	DeadlineViewUpcoming  = "upcoming"
	DeadlineViewOverdue   = "overdue"
	DeadlineViewCompleted = "completed"
)

// DeadlineViewValues 全部期限中心视图值。
var DeadlineViewValues = []string{DeadlineViewUpcoming, DeadlineViewOverdue, DeadlineViewCompleted}

// IsValidDeadlineType 校验期限类型。
func IsValidDeadlineType(s string) bool {
	for _, v := range DeadlineTypeValues {
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
