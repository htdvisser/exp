module htdvisser.dev/exp/nats/cmd/stickyrouter

go 1.17

replace htdvisser.dev/exp/natsconfig => ../../../natsconfig

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/nats-io/nats.go v1.16.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	htdvisser.dev/exp/clicontext v1.1.0
	htdvisser.dev/exp/natsconfig v0.0.0-20220602072606-c615a5df9963
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/redisconfig v0.8.11
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/nats-io/jwt/v2 v2.2.1-0.20220113022732-58e87895b296 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20220728211354-c7608f3a8462 // indirect
	golang.org/x/sys v0.0.0-20220731174439-a90be440212d // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
