package main

import (
	"log"
	"net/url"
)

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
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
