package shortener

import (
	"flag"
	"github.com/SukiEva/shortener/config"
	"github.com/SukiEva/shortener/storage"
	"github.com/gin-gonic/gin"
	"github.com/marksalpeter/token/v2"
	"log"
	"net/http"
	"strings"
)

var (
	configPath = flag.String("config", "./config/config.json", "config path")
)

func (s *shorten) Serve() {
	log.Println("Start Server...")
	s.router.POST("/api/v1/generate", generate(s))
	s.router.POST("/api/v1/expire", expire(s))
	s.router.GET("/:token", redirect(s))
	s.router.GET("/", home(s))
	http.ListenAndServe(s.conf.Server.Port, s.router)
	defer s.store.Close()
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
		router: gin.Default(),
	}, nil
}

func redirect(s *shorten) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.Param("token")
		log.Printf("Redirect request from %s (%s), token: %s\n", c.Request.RemoteAddr, c.Request.Referer(), t)
		if s.store.Check(t) {
			http.NotFound(c.Writer, c.Request)
			return
		}
		url, err := s.store.Get(t)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("%s: redirected to %s\n", c.Request.URL, url)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func generate(s *shorten) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Generate request from %s (%s)", c.Request.RemoteAddr, c.Request.Referer())
		jsonObj := generated{}
		if err := c.ShouldBindJSON(&jsonObj); err != nil {
			log.Println("Api Generate Error: ", err)
			c.JSON(http.StatusOK, gin.H{
				"code":  "0",
				"error": err.Error(),
			})
			return
		}
		t := s.random()
		if jsonObj.Exp == 0 {
			jsonObj.Exp = s.conf.Exp
		}
		if _, err := s.store.Set(t, jsonObj.Url, jsonObj.Exp); err != nil {
			log.Println("Api Generate Error: ", err)
			c.JSON(http.StatusOK, gin.H{
				"code":  "0",
				"error": err.Error(),
			})
			return
		}
		log.Printf("Generate %s to %s\n", jsonObj.Url, strings.Join([]string{s.conf.Server.Host, s.conf.Server.Port, "/"}, ""))
		c.JSON(http.StatusOK, gin.H{
			"code":  "200",
			"token": t,
			"url":   jsonObj.Url,
		})
	}
}

func expire(s *shorten) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Expire request from %s (%s)", c.Request.RemoteAddr, c.Request.Referer())
		jsonObj := expired{}
		if err := c.ShouldBindJSON(&jsonObj); err != nil {
			log.Println("Api Expire Error: ", err)
			c.JSON(http.StatusOK, gin.H{
				"code":  "0",
				"error": err.Error(),
			})
			return
		}
		if err := s.store.Expire(jsonObj.Token, jsonObj.Exp); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  "0",
				"error": err.Error(),
			})
			return
		}
		log.Printf("Expire %s to %s\n", jsonObj.Token, jsonObj.Exp)
		c.JSON(http.StatusOK, gin.H{
			"code":  "200",
			"token": jsonObj.Token,
			"exp":   jsonObj.Exp,
		})
	}
}

func home(s *shorten) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "https://github.com/SukiEva/shortener")
	}
}

func (s *shorten) random() string {
	t := token.New().Encode()
	for !s.store.Check(t) {
		t = token.New().Encode()
	}
	return t
}
