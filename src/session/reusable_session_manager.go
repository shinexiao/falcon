package session

import (
	"fmt"
	"cache"
	"connector/connection"
	"time"
	"crypto/md5"
	"encoding/hex"
)

const REUSABLE_SESSION_CACHE_KEY_PREFIX = "falcon:rs:"

/**
 * 可重用session管理器
 */
type ReusableSessionManager struct {
	ExpiredTime  int64
	CacheManager cache.CacheManager
}

/**
 * 要设置缓存管理者
 */
func NewReusableSessionManager() *ReusableSessionManager {

	//创建一个默认的缓存管理者

	return &ReusableSessionManager{
		CacheManager: cache.NewRedisCacheManager(),
	}
}

/**
 * 缓存session
 */
func (me *ReusableSessionManager) CacheSession(session *ReusableSession) bool {

	cacheKey := fmt.Sprintf("%s%s", REUSABLE_SESSION_CACHE_KEY_PREFIX, session.SessionId)

	me.CacheManager.Set(cacheKey, Encode(session.Context), me.ExpiredTime)

	return true
}

/**
 * 查询缓存的session
 */
func (me *ReusableSessionManager) QuerySession(sessionId string) *ReusableSession {

	cacheKey := fmt.Sprintf("%s%s", REUSABLE_SESSION_CACHE_KEY_PREFIX, sessionId)

	value := me.CacheManager.Get(cacheKey)

	if len(value) == 0 {
		return nil
	}

	return Decode(value)
}

/**
 * 生成快速重连session
 */
func (me *ReusableSessionManager) GenSession(context *connection.SessionContext) *ReusableSession {

	now := time.Now().Unix() * 1000 //以毫秒方式

	session := NewReusableSession()
	session.Context = context
	session.SessionId = me.md5(fmt.Sprintf("%s%d", context.DeviceId, now))
	session.ExpireTime = now + me.ExpiredTime

	return session
}

/**
 * 计算md5值
 */
func (me *ReusableSessionManager) md5(value string) string {
	hash := md5.New()
	hash.Write([]byte(value))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
