module htdvisser.dev/exp/nats

go 1.14

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/pflagenv => ../pflagenv

replace htdvisser.dev/exp/redis => ../redis

replace htdvisser.dev/exp/tls => ../tls

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-redis/redis/v8 v8.2.3
	github.com/kr/pretty v0.1.0 // indirect
	github.com/nats-io/jwt v1.0.1 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/nkeys v0.2.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	golang.org/x/sync v0.0.0-20200930132711-30421366ff76
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redis v0.0.0-20200902192356-4836516ca15c
	htdvisser.dev/exp/tls v0.0.0-20200902192356-4836516ca15c
)
