package main

import "time"

func GetCurrentTimeMills() int64 {
	curTime := time.Now()
	return curTime.UnixMilli()
}