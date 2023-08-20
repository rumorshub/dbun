package dbun

import (
	"github.com/roadrunner-server/endure/v2/dep"
)

const PluginName = "db.bun"

type Plugin struct {
	opener *BunOpener
}

func (p *Plugin) Init(opener SQLDBOpener) error {
	p.opener = NewOpener(opener)

	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*Opener)(nil), p.BunOpener),
	}
}

func (p *Plugin) BunOpener() *BunOpener {
	return p.opener
}
