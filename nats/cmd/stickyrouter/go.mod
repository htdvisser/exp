module htdvisser.dev/exp/nats/cmd/stickyrouter

go 1.18

replace htdvisser.dev/exp/natsconfig => ../../../natsconfig

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/nats-io/nats.go v1.25.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.1.0
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/natsconfig v0.0.0-20230211100859-1e66560f78a2
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redisconfig v0.8.11
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/nats-io/jwt/v2 v2.2.1-0.20220113022732-58e87895b296 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
