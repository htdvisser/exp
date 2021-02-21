module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.10.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20210220050731-9a76102bfb43 // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20210219065638-51d25998b4f6
)
