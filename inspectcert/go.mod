module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.10.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20210331175145-43e1dd70ce54 // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20210326073310-b9d6065f847d
)
