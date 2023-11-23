package main

import (
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Object struct {
	Str string
	Num int
}

func initRadisCache() *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": ":" + radisPort,
		},
	})
	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	if mycache == nil {
		panic("failed to create cache")
	}
	return mycache
}
