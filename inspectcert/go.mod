module htdvisser.dev/exp/inspectcert

go 1.17

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.13.0
	github.com/spf13/pflag v1.0.5
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20220410102725-442ae7b529aa
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20220502124256-b6088ccd6cba // indirect
)
