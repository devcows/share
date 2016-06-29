package api

import (
	"github.com/devcows/share/lib"
)

type AddResponse struct {
	UpnpOpened   bool       `json:"upnp_opened"`
	Status       bool       `json:"status"`
	ErrorMessage string     `json:"error_message"`
	Server       lib.Server `json:"server"`
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
