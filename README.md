# businessday
営業日の一月後や一月前の計算。前の日や次の日も

// 基本的な日付操作
func AddMonthsSameDay(t time.Time, months int) time.Time          // 従来実装（端数は次月へロールオーバー）
func AddMonthsEnd(t time.Time, months int) time.Time             // 目的の「翌月末」実装
func AddMonthsPreferSameOrEnd(t time.Time, months int) time.Time // 同日があれば同じ、無ければ月末

// 営業日計算（休日＝土・日＋祝日＋任意設定曜日）
func NextBusinessDay(t time.Time, n int) time.Time
func PrevBusinessDay(t time.Time, n int) time.Time

// 祝日／休日判定
func LoadJapaneseHolidays(csvPath string) error      // Shift‑JIS の CSV をロード
func IsHoliday(t time.Time) bool                    // 祝日 or 設定された休日曜日か？
func SetClosedWeekdays(days []time.Weekday)         // デフォルト（土・日）を上書き
func LoadClosedWeekdays(csvPath string) error       // CSV で任意の休日曜日をロード

exampleのmain.goに大体書いた
