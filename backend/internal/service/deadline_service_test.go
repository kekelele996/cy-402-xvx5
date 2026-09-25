package service

import (
	"testing"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/util"
)

func TestDeadlineFilterForView(t *testing.T) {
	today := time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)

	status, from, to := deadlineFilterForView(constants.DeadlineViewUpcoming, today)
	if status != constants.DeadlineStatusPending || from == nil || !from.Equal(today) || to != nil {
		t.Errorf("upcoming filter wrong: status=%s from=%v to=%v", status, from, to)
	}

	status, from, to = deadlineFilterForView(constants.DeadlineViewOverdue, today)
	if status != constants.DeadlineStatusPending || from != nil || to == nil || !to.Equal(today) {
		t.Errorf("overdue filter wrong: status=%s from=%v to=%v", status, from, to)
	}

	status, from, to = deadlineFilterForView(constants.DeadlineViewCompleted, today)
	if status != constants.DeadlineStatusCompleted || from != nil || to != nil {
		t.Errorf("completed filter wrong: status=%s from=%v to=%v", status, from, to)
	}
}

func TestDaysBetween(t *testing.T) {
	today := time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		due  time.Time
		want int
	}{
		{time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC), 3},
		{time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC), -5},
	}
	for _, tc := range cases {
		if got := daysBetween(tc.due, today); got != tc.want {
			t.Errorf("daysBetween(%v) = %d, want %d", tc.due, got, tc.want)
		}
	}
}

func TestCheckClosedCaseDueDate(t *testing.T) {
	today := time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)
	past := today.AddDate(0, 0, -1)
	future := today.AddDate(0, 0, 1)

	// 进行中案件：过去/未来日期都允许
	for _, status := range []string{constants.CaseStatusFiled, constants.CaseStatusInvestigating, constants.CaseStatusHearing} {
		if err := checkClosedCaseDueDate(status, future, today); err != nil {
			t.Errorf("open case %s should allow future date: %v", status, err)
		}
	}
	// 已结案/已归档：只允许过去日期
	for _, status := range []string{constants.CaseStatusClosed, constants.CaseStatusArchived} {
		if err := checkClosedCaseDueDate(status, past, today); err != nil {
			t.Errorf("closed case %s should allow past date: %v", status, err)
		}
		if err := checkClosedCaseDueDate(status, today, today); err == nil {
			t.Errorf("closed case %s should reject today", status)
		}
		if err := checkClosedCaseDueDate(status, future, today); err == nil {
			t.Errorf("closed case %s should reject future date", status)
		} else {
			appErr, ok := err.(*util.AppError)
			if !ok || appErr.Code != constants.CodeValidationFailed {
				t.Errorf("closed case %s error should be AppError(validation), got %v", status, err)
			}
		}
	}
}

func TestDeadlineValidators(t *testing.T) {
	if !constants.IsValidDeadlineType(constants.DeadlineTypeHearing) {
		t.Error("hearing should be valid deadline type")
	}
	if constants.IsValidDeadlineType("bogus") {
		t.Error("bogus should be invalid deadline type")
	}
	if !constants.IsValidDeadlineView(constants.DeadlineViewUpcoming) {
		t.Error("upcoming should be valid view")
	}
	if constants.IsValidDeadlineView("bogus") {
		t.Error("bogus should be invalid view")
	}
}
