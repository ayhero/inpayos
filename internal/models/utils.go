package models

import "time"

// getCurrentTimeMillis 获取当前毫秒时间戳
func getCurrentTimeMillis() int64 {
	return time.Now().UnixMilli()
}

// GetCurrentTimeMillis 导出的获取当前毫秒时间戳函数
func GetCurrentTimeMillis() int64 {
	return getCurrentTimeMillis()
}
