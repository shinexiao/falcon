package server

import (
	"cache"
	"session"
)

type FalconServer struct {
	cacheManager           cache.CacheManager
	reusableSessionManager *session.ReusableSessionManager
}

func NewFalconServer() *FalconServer {

	return &FalconServer{
		cacheManager:           cache.NewRedisCacheManager(),
		reusableSessionManager: session.NewReusableSessionManager(),
	}
}

/**
 * 获取缓存管理
 */
func (me *FalconServer) CacheManager() cache.CacheManager {

	return me.cacheManager
}

/**
 * 获取session管理器
 */
func (me *FalconServer) ReusableSessionManager() *session.ReusableSessionManager {

	return me.reusableSessionManager
}
