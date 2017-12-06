package cache

type CacheManager interface {
	//设置缓存
	Set(key, value string, expireTime int64)
	//获取缓存
	Get(key string) string
}
