package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/pilu/go-base62"
	"time"
)

const (
	UrlIdKey           = "next.url.id"
	ShortLinkKey       = "short.link:%s:url"
	UrlHashKey         = "url.hash:%s:url"
	ShortLinkDetailKey = "short.link:%s:detail"
)

type RedisCli struct {
	Cli *redis.Client
}

type UrlDetail struct {
	Url                 string        `json:"url"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

func (r RedisCli) Shorten(url string, exp int64) (string, error) {
	hashVal := toSha1(url)
	d, err := r.Cli.Get(fmt.Sprintf(UrlHashKey, hashVal)).Result()
	if err != nil {

	} else {
		if d == "{}" {
		} else {
			return d, nil
		}
	}

	err = r.Cli.Incr(UrlIdKey).Err()
	if err != nil {
		return "", err
	}

	id, err := r.Cli.Get(UrlIdKey).Int64()
	if err != nil {
		return "", err
	}
	eid := base62.Encode(int(id))

	err = r.Cli.Set(fmt.Sprintf(ShortLinkKey, eid), url, time.Minute*time.Duration(exp)).Err()

	if err != nil {
		return "", err
	}

	err = r.Cli.Set(fmt.Sprintf(UrlHashKey, hashVal), eid, time.Minute*time.Duration(exp)).Err()

	detail, err := json.Marshal(&UrlDetail{Url: url, CreatedAt: time.Now().String(), ExpirationInMinutes: time.Duration(exp)})

	if err != nil {
		return "", err
	}

	err = r.Cli.Set(fmt.Sprintf(ShortLinkDetailKey, eid), detail, time.Minute*time.Duration(exp)).Err()

	if err != nil {
		return "", err
	}
	return eid, nil
}

func (r RedisCli) ShortenLinkInfo(eid string) (interface{}, error) {
	detail, err := r.Cli.Get(fmt.Sprintf(ShortLinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{Code: 404, Err: errors.New("unknown short URL")}
	} else if err != nil {
		return "", err
	} else {
		return detail, nil
	}
}

func (r RedisCli) Unshorten(eid string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortLinkKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{Code: 404, Err: errors.New("unknown short URL")}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}

func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisCli{Cli: client}
}

func toSha1(s string) string {
	sha := sha1.New()
	return string(sha.Sum([]byte(s)))
}
