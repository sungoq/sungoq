package api

import "github.com/sungoq/sungoq/service"

type OptionFunc func(*API) *API

func WithAddr(addr string) OptionFunc {
	return func(api *API) *API {
		api.addr = addr
		return api
	}
}

func WithService(svc *service.Service) OptionFunc {
	return func(api *API) *API {
		api.service = svc
		return api
	}
}
