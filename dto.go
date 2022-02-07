package shortener

import (
	"github.com/SukiEva/shortener/config"
	"github.com/SukiEva/shortener/storage"
	"github.com/gin-gonic/gin"
)

type shorten struct {
	conf   *config.Config
	store  storage.Store
	router *gin.Engine
}

type generated struct {
	Url string `json:"url" binding:"required,url"`
	Exp int    `json:"exp" binding:"-"`
}

type expired struct {
	Token string `json:"token" binding:"required"`
	Exp   int    `json:"exp" binding:"required"`
}
