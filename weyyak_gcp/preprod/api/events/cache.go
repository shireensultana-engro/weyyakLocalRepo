package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type KeyValue struct {
	Key, Value string
}

func CacheSetKey(c *gin.Context) {
	//TODO - move to middleware
	config := c.MustGet("CONFIG").(Config)
	if !config.UseCache {
		c.JSON(http.StatusOK, gin.H{"msg": "Cache is not enabled"})
		return
	}

	var keyval KeyValue
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	ctx := c.MustGet("CONTEXT").(context.Context)
	err := c.ShouldBindJSON(&keyval)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Info().Str("key", keyval.Key).Str("value", keyval.Value)

	errSet := rdb.Set(ctx, keyval.Key, keyval.Value, 0).Err()
	if errSet != nil {
		log.Error().Msg(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"value": keyval.Value})
}

func CacheGetKey(c *gin.Context) {
	//TODO - move to middleware
	config := c.MustGet("CONFIG").(Config)
	if !config.UseCache {
		c.JSON(http.StatusOK, gin.H{"msg": "Cache is not enabled"})
		return
	}

	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	ctx := c.MustGet("CONTEXT").(context.Context)
	val, err := rdb.Get(ctx, c.Param("key")).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": val})
}

func CacheRemoveKey(c *gin.Context) {
	//TODO - move to middleware
	config := c.MustGet("CONFIG").(Config)
	if !config.UseCache {
		c.JSON(http.StatusOK, gin.H{"msg": "Cache is not enabled"})
		return
	}

	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	ctx := c.MustGet("CONTEXT").(context.Context)

	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, c.Param("key"), 10).Result()
		if err != nil {
			panic(err)
		}

		rdb.Del(ctx, keys[0])

		n += len(keys)
		if cursor == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": n,
		"key":   c.Param("key"),
	})
}
