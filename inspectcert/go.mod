module htdvisser.dev/exp/inspectcert

go 1.18

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.15.0
	github.com/spf13/pflag v1.0.5
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20230401152009-0799cd76c762
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.8.0 // indirect
)
