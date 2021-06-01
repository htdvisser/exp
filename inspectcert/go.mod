module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.12.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20210511113859-b0526f3d8744 // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20210430064303-1ac598279b4b
)
