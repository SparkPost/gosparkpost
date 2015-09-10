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

func New(cfg *config.Config) (api API, err error) {
	api.Config = cfg
	api.Client, err = rest.New(api.Config.ApiBase)
	if err != nil {
		return
	}

	return
}

/*
How do we want to be able to use this API?
- import the api library
- call api.New()
- api library loads config from (file, ...)
- api object is returned
- api.Templates.Create(...) Just Works
*/
