module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.13.0
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20211030113544-8213ee4fc10e
)
