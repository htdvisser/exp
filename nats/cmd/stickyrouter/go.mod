module htdvisser.dev/exp/nats/cmd/stickyrouter

go 1.20

replace htdvisser.dev/exp/natsconfig => ../../../natsconfig

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/nats-io/nats.go v1.30.2
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.3.0
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/natsconfig v0.0.0-20231004203327-512bce5a3dac
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redisconfig v0.8.11
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.4.1 // indirect
	github.com/nats-io/nkeys v0.4.5 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/time v0.3.0 // indirect
)
