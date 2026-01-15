package main

import (
	"github.com/gowok/gowok"
	"github.com/sungoq/sungoq/api"
	"github.com/sungoq/sungoq/service"
	_ "github.com/sungoq/sungoq/service"
)

func main() {
	gowok.
		Configures(
			service.Configure,
			api.Configure,
		).
		Run()
}
