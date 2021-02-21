module htdvisser.dev/exp/nats/cmd/stickyrouter

go 1.16

replace htdvisser.dev/exp/natsconfig => ../../../natsconfig

require (
	github.com/go-redis/redis/v8 v8.6.0
	github.com/nats-io/nats.go v1.10.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/natsconfig v0.0.0-00010101000000-000000000000
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redisconfig v0.8.5
)
