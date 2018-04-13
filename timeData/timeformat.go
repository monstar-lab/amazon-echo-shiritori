package timeData

import (
	"fmt"
	"time"
)

func GetNowTimeFormat() string {
	time := time.Now()

	return fmt.Sprintf("%d/%d/%d %d:%d:%d\n", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second())
}

//中断処理について、現在時刻と比べ、差が一番小さいものをだす
