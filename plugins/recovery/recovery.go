package recovery

import "github.com/abingzo/bups/common/plugin"

type Recovery struct {

}

func New() plugin.Plugin {
	return &Recovery{}
}

func (r *Recovery) Start(args []string) {
	panic("implement me")
}

func (r *Recovery) Caller(single plugin.Single) {
	panic("implement me")
}

func (r *Recovery) GetName() string {
	return "recovery"
}

func (r *Recovery) GetType() plugin.Type {
	panic("implement me")
}

func (r *Recovery) GetSupport() []uint32 {
	panic("implement me")
}

func (r *Recovery) SetSource(source *plugin.Source) {
	panic("implement me")
}

