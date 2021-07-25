package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func handleReindex(c *gin.Context) {
	if c.PostForm("key") != cfg.Github.GithubToken {
		c.String(http.StatusForbidden, "wrong key")
		return
	}

	go func() {
		err := regenDirectory()
		if err != nil {
			log.Println("Regen", err)
		}
	}()

	c.String(http.StatusOK, "ok")
}

func serveDirectory(c *gin.Context) {
	c.JSON(http.StatusOK, latestDirectory)
}
