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

// AddMonthsEnd は、months だけ先の月の **最終日** を返します。
//   - 基準日の day が対象月に存在しない場合はその月の末日に丸める
//   - 例: 2025-01-31 +1 month → 2025-02-28
func AddMonthsEnd(t time.Time, months int) time.Time {
	year, month, _ := t.Date()
	loc := t.Location()

	// (year, month+months, 1) を作る → 先に月だけ足す
	startOfTargetMonth := time.Date(year, month+time.Month(months), 1,
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)

	// 翌月の 1 日を取得し、そこから -1 日で「対象月の最終日」になる
	lastDayOfTarget := startOfTargetMonth.AddDate(0, 1, -1)
	return lastDayOfTarget
}

// AddMonthsPreferSameOrEnd は、
//  1. 対象月に同じ day があればそれを返す
//  2. 無ければその月の最終日を返す
func AddMonthsPreferSameOrEnd(t time.Time, months int) time.Time {
	year, month, day := t.Date()
	loc := t.Location()

	targetMonth := month + time.Month(months)
	targetYear := year + int(targetMonth-1)/12
	targetMonth = (targetMonth-1)%12 + 1

	// 試しに同じ day を作ってみる（存在しなければ次月へロールオーバー）
	candidate := time.Date(targetYear, targetMonth, day,
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)

	if candidate.Month() == targetMonth {
		// 同じ day が作れた → そのまま返す
		return candidate
	}
	// 作れなかった → 月末に丸める
	return AddMonthsEnd(t, months)
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
