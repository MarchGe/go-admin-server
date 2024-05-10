package utils

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"log/slog"
	"os"
)

var _cache *bigcache.BigCache

type CloseFunc func()

var ErrCacheNotFound = bigcache.ErrEntryNotFound

const _cleanWindow = 0
const _cacheMaxEntrySize = 10000 // 10000 entries配合默认的单条entry大小500Byte，缓存空间最大5M
func InitializeCache() CloseFunc {
	c := bigcache.DefaultConfig(0)
	c.CleanWindow = _cleanWindow
	c.MaxEntriesInWindow = _cacheMaxEntrySize
	var err error
	_cache, err = bigcache.New(context.Background(), c)
	if err != nil {
		slog.Error("create cache error", slog.Any("err", err))
		os.Exit(1)
	}
	return func() {
		if err = _cache.Close(); err != nil {
			slog.Error("close bigcache error,", slog.Any("err", err))
		}
	}
}

func GetCache() *bigcache.BigCache {
	return _cache
}
