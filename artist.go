package main

import (
	"log"
	"net/url"
)

type Artist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
	Path string `json:"path"`
}

func (a *Artist) IsUrl() bool {
	u, err := url.Parse(a.Url)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
