package ignite

import (
	"github.com/go-ignite/ignite/server"
	"github.com/go-ignite/ignite/state"
)

type Ignite struct {
	stateHandler *state.Handler
	server       *server.Server
}

func (i *Ignite) Start() error {
	i.stateHandler.Start()
	return i.server.Start()
}
