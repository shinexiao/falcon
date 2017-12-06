package common

import "strings"

const (
	MOBILE  int8 = 1 // "android", "ios"
	PC      int8 = 2 //"windows", "mac", "linux"
	WEB     int8 = 3 //"web", "h5"
	UNKNOWN int8 = -1
)

/**
 * 查找客户端分类
 */
func FindClientType(osName string) int8 {

	switch strings.ToLower(osName) {
	case "android", "ios":
		return MOBILE
	case "windows", "mac", "linux":
		return PC
	case "web", "h5":
		return WEB
	default:
		return UNKNOWN
	}
}
