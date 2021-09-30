module htdvisser.dev/exp/nats/cmd/stickyrouter

go 1.16

replace htdvisser.dev/exp/natsconfig => ../../../natsconfig

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-redis/redis/v8 v8.11.3
	github.com/nats-io/nats.go v1.12.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/natsconfig v0.0.0-20210913052840-278314a2111d
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redisconfig v0.8.11
)
