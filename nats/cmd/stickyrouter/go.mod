module htdvisser.dev/exp/nats/cmd/stickyrouter

go 1.20

replace htdvisser.dev/exp/natsconfig => ../../../natsconfig

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/nats-io/nats.go v1.32.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.5.0
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/natsconfig v0.0.0-20231206185358-cf15410f4841
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redisconfig v0.8.11
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
