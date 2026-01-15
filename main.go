package main

import (
	"github.com/sungoq/sungoq/api"
	"github.com/sungoq/sungoq/service"
)

func main() {

	svc, err := service.New()
	if err != nil {
		panic(err)
	}

	api, err := api.New(
		api.WithService(svc),
	)
	if err != nil {
		panic(err)
	}

	err = api.Start()
	if err != nil {
		panic(err)
	}

}
