module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.13.0
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20211029165221-6e7872819dc8 // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20210930055331-09e40ccb5157
)
