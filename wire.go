//+build wireinject

package ignite

import (
	"github.com/google/wire"

	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/model"
	"github.com/go-ignite/ignite/server"
	"github.com/go-ignite/ignite/service"
	"github.com/go-ignite/ignite/state"
)

func Init() (*Ignite, error) {
	wire.Build(
		config.Set,
		model.InitHandler,
		service.Set,
		state.Set,
		server.Set,
		wire.Struct(new(Ignite), "*"),
	)

	return nil, nil
}
