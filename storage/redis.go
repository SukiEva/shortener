package storage

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Store interface {
	Set(token, url string, exp int) (bool, error)
	Get(token string) (string, error)
	Expire(token string, exp int) error
	Delete(token string) error
	Close() error
}

func NewStore() (Store, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	if _, err := db.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return &rdb{
		db:  db,
		ctx: ctx,
	}, nil
}

type rdb struct {
	db  *redis.Client
	ctx context.Context
}

// Set Redis `SET token url exp NX`
func (r *rdb) Set(token, url string, exp int) (bool, error) {
	if exp < 0 {
		exp = 0
	}
	return r.db.SetNX(r.ctx, token, url, time.Hour*24*time.Duration(exp)).Result()
}

// Get Redis `GET token`
func (r *rdb) Get(token string) (string, error) {
	return r.db.Get(r.ctx, token).Result()
}

// Expire Refresh expire time for existing token
func (r *rdb) Expire(token string, exp int) error {
	if exp < 0 {
		exp = 0
	}
	_, err := r.db.Expire(r.ctx, token, time.Hour*24*time.Duration(exp)).Result()
	if err == redis.Nil {
		return errors.New("token is not exists")
	}
	return err
}

// Delete Redis `DEL token`
func (r *rdb) Delete(token string) error {
	_, err := r.db.Del(r.ctx, token).Result()
	if err == redis.Nil {
		return errors.New("token is not exists")
	}
	return err
}

// Close Flush data and close connection to redis
func (r *rdb) Close() error {
	if _, err := r.db.BgSave(r.ctx).Result(); err != nil {
		log.Println("Redis BgSave Error: ", err)
	}
	return r.db.Close()
}
