package businessday

import (
	"time"
)

// 同じ day を保つ月加算
func AddMonthsSameDay(t time.Time, months int) time.Time {
	year, month, day := t.Date()
	loc := t.Location()

	newMonth := month + time.Month(months)
	newYear := year + int(newMonth-1)/12
	newMonth = (newMonth-1)%12 + 1

	candidate := time.Date(newYear, newMonth, day,
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond(), loc)

	if candidate.Month() != newMonth { // 例: 2月に31日は作れない
		// 先の月を取得し、同じ day に再設定
		rollover := candidate.AddDate(0, 0, -candidate.Day()+1) // 当月の1日へ
		candidate = time.Date(rollover.Year(), rollover.Month(),
			day, t.Hour(), t.Minute(), t.Second(),
			t.Nanosecond(), loc)
	}
	return candidate
}

// 営業日判定（土・日・祝日）
func IsWeekend(t time.Time) bool {
	wd := t.Weekday()
	return wd == time.Saturday || wd == time.Sunday
}

// 前後の営業日取得（n は正数）
func NextBusinessDay(t time.Time, n int) time.Time {
	if n < 0 {
		return PrevBusinessDay(t, -n)
	}
	cur := t
	for i := 0; i < n; {
		cur = cur.AddDate(0, 0, 1) // +1 day
		if !IsHoliday(cur) {
			i++
		}
	}
	return cur
}

func PrevBusinessDay(t time.Time, n int) time.Time {
	if n < 0 {
		return NextBusinessDay(t, -n)
	}
	cur := t
	for i := 0; i < n; {
		cur = cur.AddDate(0, 0, -1) // -1 day
		if !IsHoliday(cur) {
			i++
		}
	}
	return cur
}
