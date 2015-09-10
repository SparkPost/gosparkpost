package api

import (
	"bitbucket.org/yargevad/go-sparkpost/config"
	"github.com/yargevad/rest"
)

type API struct {
	Config *config.Config
	Client *rest.Client
}

type Error struct {
	Message     string `json:"message"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Part        string `json:"part,omitempty"`
	Line        int    `json:"line,omitempty"`
}

func (api *API) Init(cfg *config.Config) (err error) {
	if api == nil {
		// handle the case where we call Init on something
		// that api.API is embedded in
		api = &API{}
	}
	api.Config = cfg
	api.Client, err = rest.New(api.Config.ApiBase)
	if err != nil {
		return
	}

	return
}

/*
How do we want to be able to use this API?
- import library for API we want to use (Templates, ...)
- declare Templates object, call Init on it
	- loads config from (file, ...)
- Templates.Create(...) uses config loaded by Init
*/
