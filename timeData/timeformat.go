package timeData

import (
	"fmt"
	"time"
)

func GetNowTimeFormat() string {
	// timeNow := time.Now()
	// const layout = "2006-01-02 15:04:05"
	// fmt.Println("timeNow: " + timeNow.Format(layout))
	//日付フォーマット指定
	const layout = "2006-01-02 15:04:05"

	now := time.Now()
	fmt.Println(now.Format(time.RFC3339))

	nowUTC := now.UTC()
	fmt.Println(nowUTC.Format(time.RFC3339))

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	nowJST := nowUTC.In(jst)
	fmt.Println(nowJST.Format(layout))
	return nowJST.Format(layout)
	//return timeNow.Format(layout)
}

//中断処理について、現在時刻と比べ、差が一番小さいものをだす
