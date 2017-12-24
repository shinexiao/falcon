package core

import (
	"falcon/src/session"
	"falcon/src/cache"
)

type FalconContext interface {
	CacheManager() cache.CacheManager                        //获取缓存管理
	ReusableSessionManager() *session.ReusableSessionManager //获取session管理器
}
