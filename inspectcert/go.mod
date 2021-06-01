module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.12.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20210531225629-47163c9f4e4f // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20210513141452-488c777cf587
)
