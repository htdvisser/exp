module htdvisser.dev/exp/rhmanage

go 1.14

replace htdvisser.dev/exp/redis => ../redis

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20200615200032-f1bc736245b1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redis v0.0.0-20200902192356-4836516ca15c
)
