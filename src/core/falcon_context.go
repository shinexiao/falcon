package core

import (
	"cache"
	"session"
	"ack"
)

type FalconContext interface {
	CacheManager() cache.CacheManager                        //获取缓存管理
	ReusableSessionManager() *session.ReusableSessionManager //获取session管理器
	AckTaskQueue() *ack.AckTaskQueue                         //获取ack任务管理器
}
