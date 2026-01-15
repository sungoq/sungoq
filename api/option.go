package api

type OptionFunc func(*API) *API

func WithAddr(addr string) OptionFunc {
	return func(api *API) *API {
		api.addr = addr
		return api
	}
}
