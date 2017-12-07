package cache

/**
 * 实现基于redis的缓存管理
 */
type RedisCacheManager struct {
}

func NewRedisCacheManager() *RedisCacheManager {

	return &RedisCacheManager{}
}

//设置缓存
func (me *RedisCacheManager) Set(key, value string, expireTime int64) {

	//todo
}

//获取缓存
func (me *RedisCacheManager) Get(key string) string {

	//todo
	return ""
}
