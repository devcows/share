package api

import (
	"../lib"
)

type AddResponse struct {
	Path         string   `json:"path"`
	Status       bool     `json:"status"`
	ListIps      []string `json:"list_ips"`
	UpnpOpened   bool     `json:"upnp_opened"`
	ErrorMessage string   `json:"error_message"`
}

type RmResponse struct {
	Status       bool   `json:"status"`
	ErrorMessage string `json:"error_message"`
}

type PsResponse struct {
	Status       bool         `json:"status"`
	ErrorMessage string       `json:"error_message"`
	Servers      []lib.Server `json:"servers"`
}
