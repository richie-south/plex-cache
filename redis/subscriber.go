package redisH

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"plexcache/models"
	"plexcache/utils"
	s "strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func SubscribeToExpired(rdb *redis.Client) *redis.PubSub {
	ctx := context.Background()
	plexExpirerKey := ":plex-expirer"
	location := "/cache"

	subscriber := rdb.PSubscribe(ctx, "__keyevent@0__:expired")
	go func() {
		defer subscriber.Close()

		for msg := range subscriber.Channel() {

			if !s.HasSuffix(msg.Payload, plexExpirerKey) {
				continue
			}

			dataKey := s.Split(msg.Payload, plexExpirerKey)[0]
			log.Println("Time to remove", dataKey)
			storedValue, err := rdb.Get(ctx, dataKey).Result()

			if err == redis.Nil {
				log.Println("redis.Nil", err)
				continue
			} else if err != nil {
				log.Println("Error retriving key to delete from redis", err)
				continue
			}

			var episodeCache models.EpisodeCache
			if err := json.Unmarshal([]byte(storedValue), &episodeCache); err != nil {
				continue
			}

			err = utils.RemoveFile(location + episodeCache.EpisodeFilePath)
			log.Println("Removed", episodeCache.EpisodeFilePath)

			if err != nil {
				log.Println("failed to remove file", err)
				continue
			}

			for _, srtPath := range episodeCache.SrtFilePaths {
				err = utils.RemoveFile(location + srtPath)
				log.Println("Removed srt", srtPath)

				if err != nil {
					log.Println("failed to remove srt", err)
					continue
				}
			}

			err = rdb.Del(ctx, dataKey).Err()

			if err != nil {
				log.Println("failed to remove dataKey from redis", err)
				continue
			}

		}
	}()

	return subscriber
}

func SaveEpisodeCacheToRedis(rdb *redis.Client, episodesToCache []models.EpisodeCache) error {
	ctx := context.Background()
	pipe := rdb.Pipeline()
	ttl := 20 * 24 * time.Hour
	for _, item := range episodesToCache {
		marshaled, err := json.Marshal(item)
		if err != nil {
			fmt.Println("could not marshal episode")
		}

		pipe.Set(ctx, item.RatingKey, marshaled, -1)
		pipe.Set(ctx, fmt.Sprintf("%s%s", item.RatingKey, ":plex-expirer"), "", ttl)
	}
	_, err := pipe.Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
