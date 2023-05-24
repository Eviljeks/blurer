package clock

import "time"

var Now = time.Now

func GetCurrentTS() int64 {
	return Now().UnixMicro()
}
