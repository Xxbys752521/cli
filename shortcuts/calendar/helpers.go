// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package calendar

import (
	"fmt"
	"time"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

const (
	PrimaryCalendarIDStr = "primary"
)

// resolvePrimaryCalendarID calls the calendars.primary API to get the actual
// primary calendar ID. Some private deployments do not accept the literal
// string "primary" as a calendar_id, so this function resolves it to the real
// calendar ID (e.g. "feishu.cn_xxx@group.calendar.feishu.cn").
func resolvePrimaryCalendarID(runtime *common.RuntimeContext) (string, error) {
	data, err := runtime.CallAPI("POST", "/open-apis/calendar/v4/calendars/primary", nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to resolve primary calendar: %w", err)
	}
	calendars, _ := data["calendars"].([]interface{})
	if len(calendars) == 0 {
		return "", output.Errorf(output.ExitAPI, "no_primary_calendar", "no primary calendar found")
	}
	first, _ := calendars[0].(map[string]interface{})
	cal, _ := first["calendar"].(map[string]interface{})
	if cal == nil {
		cal = first
	}
	calID, _ := cal["calendar_id"].(string)
	if calID == "" {
		return "", output.Errorf(output.ExitAPI, "no_primary_calendar", "primary calendar has no calendar_id")
	}
	return calID, nil
}

// resolveStartEnd returns (startInput, endInput) from flags with defaults.
// --start defaults to today's date, --end defaults to start date (will be resolved to end-of-day by caller).
func resolveStartEnd(runtime *common.RuntimeContext) (string, string) {
	startInput := runtime.Str("start")
	if startInput == "" {
		startInput = time.Now().Format("2006-01-02")
	}
	endInput := runtime.Str("end")
	if endInput == "" {
		endInput = startInput
	}
	return startInput, endInput
}
