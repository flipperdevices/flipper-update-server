package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
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

func serveLatest(c *gin.Context) {
	chID := c.Param("channel")
	var ch *channel
	for _, i := range latestDirectory.Channels {
		if i.ID == chID {
			ch = &i
		}
	}
	if ch == nil {
		c.String(http.StatusNotFound, "no such channel")
		return
	}
	if len(ch.Versions) == 0 {
		c.String(http.StatusNotFound, "no versions in this channel")
		return
	}

	ver := ch.Versions[len(ch.Versions)-1]
	target := strings.ReplaceAll(c.Param("target"), "-", "/")
	t := c.Param("type")

	for _, f := range ver.Files {
		if f.Type == t && f.Target == target {
			c.Redirect(http.StatusFound, f.URL)
			return
		}
	}

	c.String(http.StatusNotFound, "no such target or type")
	return
}