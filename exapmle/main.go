package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shachi/businessday"
)

func main() {
	if err := businessday.LoadJapaneseHolidays("../syukujitsu.csv"); err != nil {
		log.Fatalf("祝日ロード失敗: %v", err)
	}

	base, _ := time.Parse("2006-01-02", "2025-08-13")
	fmt.Println("基準日 :", base.Format("2006-01-02 (Mon)"))

	// 同じ day を保つ月足し
	nextMonthSameDay := businessday.AddMonthsSameDay(base, 1)
	fmt.Println("AddMonthsSameDay +1:", nextMonthSameDay.Format("2006-01-02 (Mon)"))

	// 営業日で「次の営業日」→土・日・祝日を除外
	nextBiz := businessday.NextBusinessDay(base, 1)
	fmt.Println("NextBusinessDay (+1):", nextBiz.Format("2006-01-02 (Mon)"))

	// 「前の営業日」例
	prevBiz := businessday.PrevBusinessDay(base, 1)
	fmt.Println("PrevBusinessDay (-1):", prevBiz.Format("2006-01-02 (Mon)"))

	// 祝日かどうか確認
	h := time.Date(2025, 8, 11, 0, 0, 0, 0, time.UTC)
	fmt.Printf("%s は祝日か? %v\n", h.Format("2006-01-02"), businessday.IsHoliday(h))
}
