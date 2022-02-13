module htdvisser.dev/exp/nats/cmd/stickyrouter

go 1.16

replace htdvisser.dev/exp/natsconfig => ../../../natsconfig

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/klauspost/compress v1.14.2 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/nats-server/v2 v2.7.2 // indirect
	github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/natsconfig v0.0.0-20220106142016-6857a7f82179
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redisconfig v0.8.11
)
