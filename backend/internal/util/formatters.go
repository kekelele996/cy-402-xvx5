package util

import (
	"fmt"
	"time"

	"cylawcase/internal/constants"
)

// formatters.go 同时提供日期格式化、状态文本、案件类型文本、账单类型文本，多个 handler/service 直接引用。

// FormatDateTime 格式化日期时间。
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// FormatDate 格式化日期。
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// CaseStatusText 案件状态文本。
func CaseStatusText(status string) string {
	switch status {
	case constants.CaseStatusFiled:
		return "已立案"
	case constants.CaseStatusInvestigating:
		return "调查取证"
	case constants.CaseStatusHearing:
		return "庭审中"
	case constants.CaseStatusClosed:
		return "已结案"
	case constants.CaseStatusArchived:
		return "已归档"
	default:
		return status
	}
}

// CaseTypeText 案件类型文本。
func CaseTypeText(t string) string {
	switch t {
	case constants.CaseTypeCivil:
		return "民事"
	case constants.CaseTypeCriminal:
		return "刑事"
	case constants.CaseTypeAdmin:
		return "行政"
	case constants.CaseTypeCommercial:
		return "商事"
	case constants.CaseTypeLabor:
		return "劳动"
	default:
		return t
	}
}

// BillingTypeText 费用类型文本。
func BillingTypeText(t string) string {
	switch t {
	case constants.BillingTypeAttorneyFee:
		return "律师费"
	case constants.BillingTypeCourtFee:
		return "诉讼费"
	case constants.BillingTypeTravelFee:
		return "差旅费"
	case constants.BillingTypeOther:
		return "其他"
	default:
		return t
	}
}

// BillingStatusText 账单状态文本。
func BillingStatusText(s string) string {
	switch s {
	case constants.BillingStatusPending:
		return "待支付"
	case constants.BillingStatusPaid:
		return "已支付"
	case constants.BillingStatusInvoiced:
		return "已开票"
	case constants.BillingStatusVoid:
		return "已作废"
	default:
		return s
	}
}

// DocumentTypeText 文档类型文本。
func DocumentTypeText(t string) string {
	switch t {
	case constants.DocTypeComplaint:
		return "起诉状"
	case constants.DocTypeDefense:
		return "答辩状"
	case constants.DocTypeEvidence:
		return "证据"
	case constants.DocTypeJudgment:
		return "判决书"
	case constants.DocTypeContract:
		return "合同"
	default:
		return "其他"
	}
}

// FormatMoney 格式化金额。
func FormatMoney(v float64) string {
	return fmt.Sprintf("¥%.2f", v)
}

// DeadlineTypeText 期限类型文本。
func DeadlineTypeText(t string) string {
	switch t {
	case constants.DeadlineTypeHearing:
		return "开庭"
	case constants.DeadlineTypeAppeal:
		return "上诉"
	case constants.DeadlineTypeEvidence:
		return "举证"
	case constants.DeadlineTypeTrial:
		return "审理"
	case constants.DeadlineTypeExecution:
		return "执行"
	case constants.DeadlineTypeLimitation:
		return "诉讼时效"
	case constants.DeadlineTypeReply:
		return "答辩"
	case constants.DeadlineTypeRegistration:
		return "立案登记"
	default:
		return "其他"
	}
}

// DeadlineStatusText 期限状态文本。
func DeadlineStatusText(s string) string {
	switch s {
	case constants.DeadlineStatusPending:
		return "未完成"
	case constants.DeadlineStatusCompleted:
		return "已完成"
	default:
		return s
	}
}

// DeadlineViewText 期限中心视图文本。
func DeadlineViewText(v string) string {
	switch v {
	case constants.DeadlineViewUpcoming:
		return "即将到期"
	case constants.DeadlineViewOverdue:
		return "已逾期"
	case constants.DeadlineViewCompleted:
		return "已完成"
	default:
		return v
	}
}

// DueDayText 根据到期时间计算剩余/逾期天数文案。
// 未完成：今天到期返回"今天到期"，未来返回"剩余 N 天"，过去返回"已逾期 N 天"；
// 已完成：返回"已完成"。
func DueDayText(dueAt time.Time, status string, now time.Time) string {
	if status == constants.DeadlineStatusCompleted {
		return "已完成"
	}
	days := DaysBetween(now, dueAt)
	switch {
	case days == 0:
		return "今天到期"
	case days > 0:
		return fmt.Sprintf("剩余 %d 天", days)
	default:
		return fmt.Sprintf("已逾期 %d 天", -days)
	}
}

// DaysBetween 按自然日计算 b - a（截断到天，时区一致）。
func DaysBetween(a, b time.Time) int {
	pa := time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, a.Location())
	pb := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, b.Location())
	return int(pb.Sub(pa).Hours() / 24)
}
