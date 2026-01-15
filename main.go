package main

import (
	"github.com/sungoq/sungoq/api"
	_ "github.com/sungoq/sungoq/service"
)

func main() {

	api, err := api.New()
	if err != nil {
		panic(err)
	}

	err = api.Start()
	if err != nil {
		panic(err)
	}

}
