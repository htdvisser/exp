module htdvisser.dev/exp/inspectcert

go 1.20

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.15.0
	github.com/spf13/pflag v1.0.5
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20231004203327-512bce5a3dac
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
