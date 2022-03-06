module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.13.0
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20220213111631-ce84b5198ac1
)
