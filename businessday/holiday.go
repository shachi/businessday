package businessday

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// グローバル変数
var (
	holidayMap    map[string]string     // "2006-01-02" → 名称
	closedWeekday map[time.Weekday]bool // 営業しない曜日集合
	defaultClosed = []time.Weekday{time.Saturday, time.Sunday}
)

// init でデフォルトの休日（土・日）を設定しておく
func init() {
	closedWeekday = make(map[time.Weekday]bool, len(defaultClosed))
	for _, d := range defaultClosed {
		closedWeekday[d] = true
	}
}

// 祝日のロード（Shift‑JIS → UTF‑8）
func LoadJapaneseHolidays(csvPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	sjisReader := transform.NewReader(bufio.NewReader(f), japanese.ShiftJIS.NewDecoder())
	r := csv.NewReader(sjisReader)
	r.Comma = ','
	r.TrimLeadingSpace = true

	tmp := make(map[string]string)

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(rec) < 2 {
			continue // 壊れた行はスキップ
		}

		dateStr := strings.TrimSpace(rec[0])
		name := strings.TrimSpace(rec[1])

		t, err := time.Parse("2006/1/2", dateStr)
		if err != nil {
			continue
		}
		key := t.Format("2006-01-02")
		tmp[key] = name
	}
	holidayMap = tmp
	return nil
}

// 休日曜日（営業しない曜日）のロード
//
//	フォーマット: weekday,備考   （weekday は 0‑6 の整数）
func LoadClosedWeekdays(csvPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = ','
	r.TrimLeadingSpace = true

	tmp := make(map[time.Weekday]bool)

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(rec) < 1 {
			continue
		}

		// ヘッダー行やコメントは数値に変換できなければスキップ
		num, err := strconv.Atoi(strings.TrimSpace(rec[0]))
		if err != nil {
			continue
		}
		if num < 0 || num > 6 {
			continue // 範囲外は無視
		}
		tmp[time.Weekday(num)] = true
	}

	// 上書き（必要なら Add したい場合は merge を実装）
	closedWeekday = tmp
	return nil
}

// プログラムから直接設定したいときのヘルパー
func SetClosedWeekdays(days []time.Weekday) {
	tmp := make(map[time.Weekday]bool, len(days))
	for _, d := range days {
		tmp[d] = true
	}
	closedWeekday = tmp
}

// 判定ロジック（休日＝土日＋CSV の祝日）
func IsHoliday(t time.Time) bool {
	if closedWeekday[t.Weekday()] { // 営業しない曜日か？
		return true
	}
	key := t.Format("2006-01-02")
	_, ok := holidayMap[key]
	return ok
}

// 週末判定は内部で使わなくても OK（外部からも利用できる）
// （※以前の実装を残すだけのラッパー）
func isWeekend(t time.Time) bool {
	return closedWeekday[t.Weekday()]
}
