package server

import (
	"fmt"

	"github.com/alexmeuer/juke/pkg/openapi"
)

func Serve(host string, port uint16) error {
	r := openapi.NewRouter(openapi.ApiHandleFunctions{
		RoomsAPI: &RoomsAPI{},
	})
	return r.Run(fmt.Sprintf("%s:%d", host, port))
}
