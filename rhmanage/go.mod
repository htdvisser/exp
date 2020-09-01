module htdvisser.dev/exp/rhmanage

go 1.14

replace htdvisser.dev/exp/redis => ../redis

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/pflagenv => ../pflagenv

require (
	github.com/go-redis/redis/v8 v8.0.0-beta.9 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/exp v0.0.0-20200831210406-1ff542fc73e3 // indirect
	htdvisser.dev/exp/clicontext v1.0.0
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redis v0.0.0-20200729194012-b5e9825fa931
)
