module htdvisser.dev/exp/inspectcert

go 1.16

replace htdvisser.dev/exp/tlsconfig => ../tlsconfig

require (
	github.com/fatih/color v1.10.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073 // indirect
	htdvisser.dev/exp/pflagenv v1.0.0
	htdvisser.dev/exp/tlsconfig v0.0.0-20210222072022-ad7a9e56ebee
)
