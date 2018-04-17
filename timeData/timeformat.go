package timeData

import (
	"fmt"
	"time"
)

func GetNowTimeFormat(layout string) string {
	//const layout = "2006-01-02 15:04:05"

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

//DBのIDとして現在日時のフォーマットを変換
// func GetDBIDFormat() string {

// }
