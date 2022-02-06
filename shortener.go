package shortener

import (
	"flag"
	"github.com/SukiEva/shortener/config"
	"github.com/SukiEva/shortener/storage"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	configPath = flag.String("config", "./config/config.json", "config path")
)

type shorten struct {
	conf   *config.Config
	store  storage.Store
	client *gin.Engine
}

func New() (*shorten, error) {
	conf, err := config.Read(*configPath)
	if err != nil {
		log.Println("Config Read Error: ", err)
		return nil, err
	}
	store, err := storage.NewStore(conf.Redis.Addr, conf.Redis.Password, conf.Redis.Db)
	if err != nil {
		log.Println("Redis Connect Error: ", err)
		return nil, err
	}
	return &shorten{
		conf:   conf,
		store:  store,
		client: gin.Default(),
	}, nil
}
