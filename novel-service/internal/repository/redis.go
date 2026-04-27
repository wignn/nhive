package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/novelhive/novel-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		ttl:    15 * time.Minute,
	}
}

func (c *RedisCache) GetNovel(slug string) (*domain.Novel, error) {
	key := fmt.Sprintf("novel:%s", slug)
	data, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	var novel domain.Novel
	if err := json.Unmarshal(data, &novel); err != nil {
		return nil, err
	}
	return &novel, nil
}

func (c *RedisCache) SetNovel(novel *domain.Novel) error {
	key := fmt.Sprintf("novel:%s", novel.Slug)
	data, err := json.Marshal(novel)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), key, data, c.ttl).Err()
}

func (c *RedisCache) InvalidateNovel(slug string) {
	c.client.Del(context.Background(), fmt.Sprintf("novel:%s", slug))
}

func (c *RedisCache) GetChapter(novelSlug string, number int) (*domain.Chapter, error) {
	key := fmt.Sprintf("chapter:%s:%d", novelSlug, number)
	data, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	var chapter domain.Chapter
	if err := json.Unmarshal(data, &chapter); err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (c *RedisCache) SetChapter(novelSlug string, chapter *domain.Chapter) error {
	key := fmt.Sprintf("chapter:%s:%d", novelSlug, chapter.Number)
	data, err := json.Marshal(chapter)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), key, data, c.ttl).Err()
}

func (c *RedisCache) InvalidateChapter(novelSlug string, number int) {
	c.client.Del(context.Background(), fmt.Sprintf("chapter:%s:%d", novelSlug, number))
}
