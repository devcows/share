package api

import "github.com/tylerb/graceful"

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

type Server struct {
	Path    string   `json:"path"`
	ID      int      `json:"id"`
	ListIps []string `json:"list_ips"`
	Srv     *graceful.Server
}
