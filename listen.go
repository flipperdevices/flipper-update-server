package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleReindex(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.FormValue("key") != cfg.Github.GithubToken {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("wrong key"))
		return
	}

	go func() {
		err := regenDirectory()
		if err != nil {
			log.Println("Regen", err)
		}
	}()

	w.Write([]byte("ok"))
}

func handleDirectory(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(latestDirectory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
